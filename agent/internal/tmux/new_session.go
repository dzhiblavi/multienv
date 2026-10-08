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
			"-P", "-F", "#{session_id}",
		},
	}

	result, err := subprocess.RunSuccess(ctx, cmd)
	if err != nil {
		return model.Session{}, err
	}

	return model.Session{
		Id:          strings.TrimSpace(result.Stdout),
		Name:        name,
		WindowCount: 1,
		SyncedAt:    time.Now(),
	}, nil
}
