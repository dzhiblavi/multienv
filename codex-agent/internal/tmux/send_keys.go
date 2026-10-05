package tmux

import (
	"context"

	"multienv/common/subprocess"
)

func SendKeys(ctx context.Context, target string, keys []string) error {
	args := []string{"send-keys", "-t", target}
	args = append(args, keys...)

	cmd := subprocess.Command{
		Name: "tmux",
		Args: args,
	}

	_, err := subprocess.RunSuccess(ctx, cmd)
	return err
}
