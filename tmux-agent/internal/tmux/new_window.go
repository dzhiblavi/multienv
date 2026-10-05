package tmux

import (
	"context"
	"fmt"
	"strings"

	"multienv/common/subprocess"
)

func NewWindow(ctx context.Context, sessionId, name string) (string, error) {
	cmd := subprocess.Command{
		Name: "tmux",
		Args: []string{
			"new-window", "-d", "-P", "-F", "#{pane_id}",
			"-t", sessionId + ":", "-n", name,
		},
	}

	result, err := subprocess.RunSuccess(ctx, cmd)
	if err != nil {
		return "", fmt.Errorf("create tmux window: %w", err)
	}
	return strings.TrimSpace(result.Stdout), nil
}
