package codex

import (
	"context"

	"multienv/common/subprocess"
)

func ArchiveSession(ctx context.Context, id string) error {
	cmd := subprocess.Command{
		Name: "codex",
		Args: []string{
			"--no-daemon",
			"archive",
			id,
		},
	}

	_, err := subprocess.RunSuccess(ctx, cmd)
	return err
}
