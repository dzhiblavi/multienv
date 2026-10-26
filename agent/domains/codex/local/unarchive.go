package local

import (
	"context"

	"multienv/common/subprocess"
)

func UnarchiveSession(ctx context.Context, id string) error {
	cmd := subprocess.Command{
		Name: "codex",
		Args: []string{
			"--no-daemon",
			"unarchive",
			id,
		},
	}

	_, err := subprocess.RunSuccess(ctx, cmd)
	return err
}
