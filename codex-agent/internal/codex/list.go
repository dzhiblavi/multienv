package codex

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"multienv/common/codex/model"
	"multienv/common/subprocess"
	"multienv/common/util"
)

type tmuxPane struct {
	Id  string `json:"id"`
	Tty string `json:"tty"`
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
	if len(filesInfo) == 0 {
		return result, nil
	}

	paneByTty, err := getTtyToTmuxPaneMapping(ctx)
	if err != nil {
		return nil, fmt.Errorf("get tmux panes: %w", err)
	}

	for _, info := range filesInfo {
		if paneId, ok := paneByTty[info.tty]; ok {
			result[info.sessionId] = paneId
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
	db, err := openCodexDatabase(ctx, "state")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("open Codex sessions database: %w", err)
	}

	if err == nil {
		defer db.Close()
		const query = `
			select id, coalesce(archived, 0), coalesce(cwd, ''),
				substr(replace(replace(replace(replace(
					coalesce(nullif(name, ''), nullif(title, ''), nullif(preview, ''), ''),
					char(9), ' '), char(10), ' '), char(13), ' '), char(27), ''), 1, 100)
			from threads
		`
		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("query Codex sessions: %w", err)
		}

		defer rows.Close()

		for rows.Next() {
			var sessionId, workDir, title string
			var archived int64
			if err := rows.Scan(&sessionId, &archived, &workDir, &title); err != nil {
				return nil, fmt.Errorf("read Codex session: %w", err)
			}
			state := model.NotLoaded
			if archived != 0 {
				state = model.Archived
			}
			result[sessionId] = model.Session{
				Hostname: hostname,
				Id:       sessionId,
				WorkDir:  workDir,
				Title:    title,
				State:    state,
				SyncedAt: syncedAt,
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
	if len(filesInfo) == 0 {
		return result, nil
	}

	paneByTty, err := getTtyToTmuxPaneMapping(ctx)
	if err != nil {
		return nil, fmt.Errorf("get tmux panes: %w", err)
	}

	threadsDB, err := openCodexDatabase(ctx, "thread_history")
	if err != nil {
		return nil, fmt.Errorf("failed to open threads DB: %v", err)
	}
	defer threadsDB.Close()

	for _, info := range filesInfo {
		if !validCodexFilesInfo(info) {
			continue
		}

		paneId, ok := paneByTty[info.tty]
		if !ok {
			continue
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

		session.TmuxPaneId = paneId

		state, err := getCodexSessionState(ctx, threadsDB, info.sessionId)
		if err != nil {
			return nil, fmt.Errorf("get active Codex session state: %w", err)
		}
		session.State = state

		result[info.sessionId] = session
	}

	return result, nil
}

func getTtyToTmuxPaneMapping(ctx context.Context) (map[string]string, error) {
	cmd := subprocess.Command{
		Name: "tmux",
		Args: []string{
			"list-panes",
			"-a",
			"-F",
			`{"id":"#{pane_id}","tty":"#{pane_tty}"}`,
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

	result := make(map[string]string)
	for _, pane := range tmuxPanes {
		result[pane.Tty] = pane.Id
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
	var result []codexFilesInfo

	for field := range bytes.SplitSeq([]byte(lsofResult.Stdout), []byte{0}) {
		field = bytes.TrimLeft(field, "\n")
		if len(field) == 0 {
			continue
		}

		tag := string(field[0])
		value := string(field[1:])

		switch tag {
		case "p":
			if validCodexFilesInfo(info) {
				result = append(result, info)
			}
			info = codexFilesInfo{}

		case "c":
			info.command = value

		case "n":
			if strings.HasPrefix(value, "/dev/tty") || strings.HasPrefix(value, "/dev/pts/") {
				if info.tty != "" && info.tty != value {
					info.ambiguousTty = true
				}
				info.tty = value
			} else if id := lockFileRegex.FindStringSubmatch(value); id != nil {
				if info.sessionId != "" && info.sessionId != id[1] {
					info.ambiguousId = true
				}
				info.sessionId = id[1]
			}
		}
	}

	if validCodexFilesInfo(info) {
		result = append(result, info)
	}

	return result, nil
}

func validCodexFilesInfo(info codexFilesInfo) bool {
	if info.command != "codex" || info.sessionId == "" || info.tty == "" {
		return false
	}

	if info.ambiguousId || info.ambiguousTty {
		return false
	}

	return true
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
