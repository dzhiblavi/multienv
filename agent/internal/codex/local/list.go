package local

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"

	"multienv/common/codex/model"
	"multienv/common/log"
	"multienv/common/subprocess"
	"multienv/common/util"
)

type tmuxPane struct {
	Id    string `json:"id"`
	Tty   string `json:"tty"`
	Title string `json:"title"`
}

type codexFilesInfo struct {
	command      string
	tty          string
	sessionId    string
	ambiguousTty bool
	ambiguousId  bool
}

func ListSessions(ctx context.Context) (model.SessionsById, error) {
	return getAllCodexSessions(ctx)
}

func GetSessionToTmuxPaneMapping(ctx context.Context) (map[string]string, error) {
	filesInfo, err := getCodexFilesInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("get active Codex session files: %w", err)
	}

	result := make(map[string]string)
	filesInfo = codexFilesInfoForTmux(filesInfo)
	if len(filesInfo) == 0 {
		return result, nil
	}

	paneByTty, err := getTtyToTmuxPaneMapping(ctx)
	if err != nil {
		return nil, fmt.Errorf("get tmux panes: %w", err)
	}

	for _, info := range filesInfo {
		if pane, ok := paneByTty[info.tty]; ok {
			result[info.sessionId] = pane.Id
		}
	}

	return result, nil
}

func getAllCodexSessions(ctx context.Context) (model.SessionsById, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("get hostname: %w", err)
	}

	syncedAt := time.Now()
	result := make(model.SessionsById)
	guardianIds := make(map[string]bool)
	db, err := openCodexDatabase(ctx, "state")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("open Codex sessions database: %w", err)
	}

	if err == nil {
		defer db.Close()
		const query = `
			select id, coalesce(source, ''), coalesce(archived, 0), coalesce(cwd, ''),
				coalesce(first_user_message, ''),
				substr(replace(replace(replace(replace(
					case when history_mode = 'legacy'
						then coalesce(nullif(title, ''), nullif(name, ''), nullif(preview, ''), '')
						else coalesce(nullif(name, ''), nullif(title, ''), nullif(preview, ''), '')
					end,
					char(9), ' '), char(10), ' '), char(13), ' '), char(27), ''), 1, 100)
			from threads
		`
		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("query Codex sessions: %w", err)
		}

		defer rows.Close()

		for rows.Next() {
			var sessionId, source, workDir, firstMessage, title string
			var archived int64
			if err := rows.Scan(&sessionId, &source, &archived, &workDir, &firstMessage, &title); err != nil {
				return nil, fmt.Errorf("read Codex session: %w", err)
			}
			if isGuardianSessionSource(source) {
				guardianIds[sessionId] = true
				continue
			}
			state := model.NotLoaded
			if archived != 0 {
				state = model.Archived
			}
			result[sessionId] = model.Session{
				Hostname:     hostname,
				Id:           sessionId,
				WorkDir:      workDir,
				Title:        title,
				FirstMessage: firstMessage,
				State:        state,
				SyncedAt:     syncedAt,
			}
		}

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("read Codex sessions: %w", err)
		}
		if err := rows.Close(); err != nil {
			return nil, fmt.Errorf("close Codex sessions cursor: %w", err)
		}
	}

	filesInfo, err := getCodexFilesInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("get active Codex session files: %w", err)
	}
	// An open writer lock must not reintroduce a filtered Guardian thread.
	filesInfo = slices.DeleteFunc(filesInfo, func(info codexFilesInfo) bool {
		return guardianIds[info.sessionId]
	})
	threadsDB, err := openCodexDatabase(ctx, "thread_history")
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) || len(filesInfo) > 0 {
			return nil, fmt.Errorf("open Codex thread history database: %w", err)
		}
	} else {
		defer threadsDB.Close()
	}

	var paneByTty map[string]tmuxPane
	if len(codexFilesInfoForTmux(filesInfo)) > 0 {
		paneByTty, err = getTtyToTmuxPaneMapping(ctx)
		if err != nil {
			log.Warn(ctx, "failed to get tmux panes", "error", err)
		}
	}

	for _, info := range filesInfo {
		if !validCodexFilesInfo(info) {
			continue
		}

		var pane tmuxPane
		if info.tty != "" && !info.ambiguousTty && !info.ambiguousId {
			pane = paneByTty[info.tty]
		}

		session, exists := result[info.sessionId]
		if !exists {
			session = model.Session{
				Hostname: hostname,
				Id:       info.sessionId,
				State:    model.Unknown,
				SyncedAt: syncedAt,
			}
		}

		session.TmuxPaneId = pane.Id

		if pane.Id != "" && strings.HasPrefix(pane.Title, "[ ! ] Action Required") {
			session.State = model.ActionRequired
		} else {
			state, err := getCodexSessionState(ctx, threadsDB, info.sessionId)
			if err != nil {
				log.Warn(ctx, "failed to get active Codex session state",
					"session_id", info.sessionId, "pane_id", pane.Id, "error", err)
				state = model.Unknown
			}

			session.State = state
		}

		result[info.sessionId] = session
	}

	if threadsDB != nil {
		if err := fillCodexSessionMessages(ctx, threadsDB, result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// Keep the initial user request and the latest user message or final agent answer.
// Read history for inactive and archived sessions as well as loaded sessions.
func fillCodexSessionMessages(ctx context.Context, db *sql.DB, sessions model.SessionsById) error {
	const query = `
		with messages as (
			select thread_id, rollout_ordinal, item_type,
				case item_type
					when 'userMessage' then (
						select group_concat(json_extract(value, '$.text'), char(10))
						from json_each(item_json, '$.content')
						where json_extract(value, '$.type') = 'text'
					)
					else json_extract(item_json, '$.text')
				end as message
			from thread_items
			where item_type = 'userMessage'
				or (item_type = 'agentMessage'
					and coalesce(json_extract(item_json, '$.phase'), 'final_answer')
						= 'final_answer')
		), useful_messages as (
			select * from messages
			where trim(message, ' ' || char(9) || char(10) || char(13)) <> ''
		), bounds as (
			select thread_id,
				min(case when item_type = 'userMessage' then rollout_ordinal end) as first,
				max(rollout_ordinal) as last
			from useful_messages
			group by thread_id
		)
		select bounds.thread_id, coalesce(first.message, ''), coalesce(last.message, '')
		from bounds
		left join useful_messages first
			on first.thread_id = bounds.thread_id and first.rollout_ordinal = bounds.first
		left join useful_messages last
			on last.thread_id = bounds.thread_id and last.rollout_ordinal = bounds.last
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("query Codex session messages: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var sessionId, firstMessage, lastMessage string
		if err := rows.Scan(&sessionId, &firstMessage, &lastMessage); err != nil {
			return fmt.Errorf("read Codex session messages: %w", err)
		}
		session, ok := sessions[sessionId]
		if !ok {
			continue
		}
		if strings.TrimSpace(session.FirstMessage) == "" {
			session.FirstMessage = firstMessage
		}
		session.LastMessage = lastMessage
		sessions[sessionId] = session
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("read Codex session messages: %w", err)
	}
	if err := rows.Close(); err != nil {
		return fmt.Errorf("close Codex session messages cursor: %w", err)
	}
	return nil
}

func getTtyToTmuxPaneMapping(ctx context.Context) (map[string]tmuxPane, error) {
	cmd := subprocess.Command{
		Name: "tmux",
		Args: []string{
			"list-panes",
			"-a",
			"-F",
			`{"id":"#{pane_id}","tty":"#{pane_tty}","title":"#{pane_title}"}`,
		},
	}

	tmuxPanesResult, err := subprocess.RunSuccess(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to get tmux panes-ttys: %v", err)
	}

	lines := strings.Split(tmuxPanesResult.Stdout, "\n")
	if len(lines) > 0 {
		lines = lines[:len(lines)-1]
	}

	tmuxPanes, err := util.BulkUnmarshal[tmuxPane](lines)
	if err != nil {
		return nil, err
	}

	result := make(map[string]tmuxPane)
	for _, pane := range tmuxPanes {
		result[pane.Tty] = pane
	}

	return result, nil
}

func getCodexFilesInfo(ctx context.Context) ([]codexFilesInfo, error) {
	cmd := subprocess.Command{
		Name: "lsof",
		Args: []string{"-nP", "-c", "codex", "-F", "0pcfn"},
	}

	lsofResult, err := subprocess.Run(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to query lsof: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	// lsof returns 1 with empty output when no processes match the selection.
	if lsofResult.ExitCode == 1 && lsofResult.Stdout == "" && lsofResult.Stderr == "" {
		return nil, nil
	}
	if lsofResult.ExitCode != 0 {
		return nil, fmt.Errorf("lsof exited with code %d: %s", lsofResult.ExitCode, lsofResult.Stderr)
	}

	lockFileRegex := regexp.MustCompile(`/thread-writer-locks/([[:xdigit:]-]+)\.lock$`)

	info := codexFilesInfo{}
	var sessionIds []string
	var result []codexFilesInfo
	appendFilesInfo := func() {
		// One Codex process can own several threads, including sub-agent threads.
		// Each open writer lock proves activity, but multiple IDs cannot be mapped
		// unambiguously to a tmux pane.
		for _, sessionId := range sessionIds {
			sessionInfo := info
			sessionInfo.sessionId = sessionId
			if validCodexFilesInfo(sessionInfo) {
				result = append(result, sessionInfo)
			}
		}
	}

	for field := range bytes.SplitSeq([]byte(lsofResult.Stdout), []byte{0}) {
		field = bytes.TrimLeft(field, "\n")
		if len(field) == 0 {
			continue
		}

		tag := string(field[0])
		value := string(field[1:])

		switch tag {
		case "p":
			appendFilesInfo()
			info = codexFilesInfo{}
			sessionIds = nil

		case "c":
			info.command = value

		case "n":
			if strings.HasPrefix(value, "/dev/tty") || strings.HasPrefix(value, "/dev/pts/") {
				if info.tty != "" && info.tty != value {
					info.ambiguousTty = true
				}
				info.tty = value
			} else if id := lockFileRegex.FindStringSubmatch(value); id != nil {
				if !slices.Contains(sessionIds, id[1]) {
					sessionIds = append(sessionIds, id[1])
				}
				info.ambiguousId = len(sessionIds) > 1
			}
		}
	}

	appendFilesInfo()

	return result, nil
}

func isGuardianSessionSource(source string) bool {
	var parsed struct {
		Subagent struct {
			Other string `json:"other"`
		} `json:"subagent"`
	}
	return json.Unmarshal([]byte(source), &parsed) == nil && parsed.Subagent.Other == "guardian"
}

func validCodexFilesInfo(info codexFilesInfo) bool {
	if info.command != "codex" || info.sessionId == "" {
		return false
	}

	return true
}

func codexFilesInfoForTmux(filesInfo []codexFilesInfo) []codexFilesInfo {
	var result []codexFilesInfo
	for _, info := range filesInfo {
		if info.tty != "" && !info.ambiguousTty && !info.ambiguousId {
			result = append(result, info)
		}
	}
	return result
}

func getCodexSessionState(ctx context.Context, db *sql.DB, sessionId string) (model.State, error) {
	const query = `
		select status from thread_turns where thread_id = ?
		order by rollout_ordinal desc, started_at desc, turn_id desc limit 1
	`
	var status sql.NullString
	err := db.QueryRowContext(ctx, query, sessionId).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Unknown, nil
	}
	if err != nil {
		return "", fmt.Errorf("get Codex session %q state: %w", sessionId, err)
	}

	if status.String == "inProgress" {
		return model.Working, nil
	}
	return model.Idle, nil
}
