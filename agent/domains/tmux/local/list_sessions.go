package local

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"multienv/agent/domains/tmux/model"
	"multienv/common/subprocess"
)

// Keep the path last so tabs in directory names are preserved by SplitN.
const sessionFormat = "#{session_id}\t#{session_name}\t#{session_windows}\t" +
	"#{session_attached}\t#{?session_last_attached,#{session_last_attached},0}\t" +
	"#{pane_current_command}\t#{pane_current_path}"

func ListSessions(ctx context.Context) (model.SessionsById, error) {
	cmd := subprocess.Command{
		Name: "tmux",
		Args: []string{
			"-u", // Preserve tabs and Unicode in the output even without a UTF-8 locale.
			"list-sessions",
			"-F",
			sessionFormat,
		},
	}

	syncedAt := time.Now()
	result, err := subprocess.RunSuccess(ctx, cmd)
	if err != nil {
		return nil, err
	}

	sessions := make(model.SessionsById)
	for line := range strings.SplitSeq(strings.TrimSuffix(result.Stdout, "\n"), "\n") {
		if line == "" {
			continue
		}
		session, err := parseSession(line, syncedAt)
		if err != nil {
			return nil, err
		}
		sessions[session.Id] = session
	}

	windows, err := listWindows(ctx, "")
	if err != nil {
		return nil, err
	}

	for id, session := range sessions {
		if list, ok := windows[id]; ok {
			session.Windows = list
			sessions[id] = session
		}
	}

	return sessions, nil
}

func parseSession(line string, syncedAt time.Time) (model.Session, error) {
	fields := strings.SplitN(line, "\t", 7)
	if len(fields) != 7 {
		return model.Session{}, fmt.Errorf("invalid tmux session: expected 7 fields, got %d", len(fields))
	}
	windowCount, err := strconv.Atoi(fields[2])
	if err != nil {
		return model.Session{}, fmt.Errorf("parse tmux window count: %w", err)
	}
	attached, err := strconv.Atoi(fields[3])
	if err != nil {
		return model.Session{}, fmt.Errorf("parse tmux attached clients: %w", err)
	}
	lastAttached, err := strconv.ParseInt(fields[4], 10, 64)
	if err != nil {
		return model.Session{}, fmt.Errorf("parse tmux last attachment: %w", err)
	}
	return model.Session{
		Id:             fields[0],
		Name:           fields[1],
		WindowCount:    windowCount,
		Windows:        []model.Window{},
		Attached:       attached,
		LastAttached:   lastAttached,
		CurrentCommand: fields[5],
		CurrentPath:    fields[6],
		SyncedAt:       syncedAt,
	}, nil
}
