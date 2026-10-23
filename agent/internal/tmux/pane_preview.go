package tmux

import (
	"context"
	"fmt"

	"multienv/common/subprocess"
)

func PanePreview(ctx context.Context, target string) (string, error) {
	result, err := subprocess.RunSuccess(ctx, subprocess.Command{
		Name: "tmux",
		Args: []string{"-u", "capture-pane", "-p", "-e", "-t", target},
	})
	if err != nil {
		return "", fmt.Errorf("capture pane preview: %w", err)
	}

	return result.Stdout, nil
}
