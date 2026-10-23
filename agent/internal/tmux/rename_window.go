package tmux

import (
	"context"

	"multienv/common/subprocess"
)

func RenameWindow(ctx context.Context, sessionId, windowId, newName string) error {
	cmd := subprocess.Command{
		Name: "tmux",
		Args: []string{
			"rename-window",
			"-t", sessionId + ":" + windowId,
			"--", newName,
		},
	}

	_, err := subprocess.RunSuccess(ctx, cmd)
	return err
}
