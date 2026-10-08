package tmux

import (
	"context"
	"strings"
	"time"

	"multienv/common/subprocess"
	"multienv/common/tmux/model"
)

func NewSession(ctx context.Context, name string) (model.Session, error) {
	cmd := subprocess.Command{
		Name: "tmux",
		Args: []string{
			"new-session", "-d",
			"-s", name,
			"-P", "-F", sessionFormat,
		},
	}

	result, err := subprocess.RunSuccess(ctx, cmd)
	if err != nil {
		return model.Session{}, err
	}

	return parseSession(strings.TrimSuffix(result.Stdout, "\n"), time.Now())
}
