package local

import (
	"context"

	"multienv/common/subprocess"
)

func RenameSession(ctx context.Context, id, newName string) error {
	cmd := subprocess.Command{
		Name: "tmux",
		Args: []string{
			"rename-session",
			"-t", id,
			newName,
		},
	}

	_, err := subprocess.RunSuccess(ctx, cmd)
	return err
}
