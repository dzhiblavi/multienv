package tmux

import (
	"context"
	"fmt"

	"multienv/common/subprocess"
)

func MoveWindow(ctx context.Context, sessionId, windowId string, newIndex int) error {
	cmd := subprocess.Command{
		Name: "tmux",
		Args: []string{
			"move-window", "-d",
			"-s", sessionId + ":" + windowId,
			"-t", fmt.Sprintf("%s:%d", sessionId, newIndex),
		},
	}

	_, err := subprocess.RunSuccess(ctx, cmd)
	return err
}
