package tmux

import (
	"context"
	"errors"
	"strings"

	"multienv/common/subprocess"
)

const lastWindowError = "cannot kill the last window of a session"

func KillWindow(ctx context.Context, sessionId, windowId string) error {
	target := sessionId + ":" + windowId
	quotedTarget := "'" + strings.ReplaceAll(target, "'", "'\\''") + "'"
	cmd := subprocess.Command{
		Name: "tmux",
		Args: []string{
			"if-shell", "-F", "-t", target,
			"#{>:#{session_windows},1}",
			"unlink-window -k -t " + quotedTarget,
			"display-message -p '" + lastWindowError + "'",
		},
	}

	result, err := subprocess.RunSuccess(ctx, cmd)
	if err != nil {
		return err
	}

	if strings.TrimSpace(result.Stdout) == lastWindowError {
		return errors.New(lastWindowError)
	}

	return nil
}
