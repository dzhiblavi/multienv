package local

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"multienv/common/subprocess"
	"multienv/common/tmux/model"
)

// Include the session ID because a window may be linked to multiple sessions.
const windowFormat = "#{session_id}\t#{window_id}\t#{window_index}\t#{window_name}\t" +
	"#{window_active}\t#{pane_current_command}\t#{pane_current_path}"

func listWindows(ctx context.Context, sessionId string) (map[string][]model.Window, error) {
	args := []string{"-u", "list-windows", "-F", windowFormat}
	if sessionId == "" {
		args = append(args, "-a")
	} else {
		args = append(args, "-t", sessionId)
	}
	result, err := subprocess.RunSuccess(ctx, subprocess.Command{Name: "tmux", Args: args})
	if err != nil {
		return nil, fmt.Errorf("list tmux windows: %w", err)
	}

	windows := make(map[string][]model.Window)
	for line := range strings.SplitSeq(strings.TrimSuffix(result.Stdout, "\n"), "\n") {
		if line == "" {
			continue
		}

		fields := strings.SplitN(line, "\t", 7)
		if len(fields) != 7 {
			return nil, fmt.Errorf("invalid tmux window: expected 7 fields, got %d", len(fields))
		}

		index, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, fmt.Errorf("parse tmux window index: %w", err)
		}

		active, err := strconv.ParseBool(fields[4])
		if err != nil {
			return nil, fmt.Errorf("parse tmux active window: %w", err)
		}

		windows[fields[0]] = append(windows[fields[0]], model.Window{
			Id:             fields[1],
			Index:          index,
			Name:           fields[3],
			Active:         active,
			CurrentCommand: fields[5],
			CurrentPath:    fields[6],
		})
	}

	return windows, nil
}
