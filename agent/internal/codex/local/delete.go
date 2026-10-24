package local

import (
	"context"

	"multienv/common/subprocess"
)

func DeleteSession(ctx context.Context, id string) error {
	cmd := subprocess.Command{
		Name: "codex",
		Args: []string{
			"--no-daemon",
			"delete",
			"--force",
			id,
		},
	}

	_, err := subprocess.RunSuccess(ctx, cmd)
	return err
}
