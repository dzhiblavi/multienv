package local

import (
	"context"

	"multienv/common/subprocess"
)

func KillPane(ctx context.Context, paneId string) error {
	cmd := subprocess.Command{
		Name: "tmux",
		Args: []string{
			"kill-pane",
			"-t",
			paneId,
		},
	}

	_, err := subprocess.RunSuccess(ctx, cmd)
	return err
}
