package local

import (
	"context"

	"multienv/common/subprocess"
)

func KillSession(ctx context.Context, id string) error {
	cmd := subprocess.Command{
		Name: "tmux",
		Args: []string{
			"kill-session",
			"-t",
			id,
		},
	}

	_, err := subprocess.RunSuccess(ctx, cmd)
	return err
}
