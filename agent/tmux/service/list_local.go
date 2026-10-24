package service

import (
	"context"

	"multienv/agent/internal/tmux/local"
	"multienv/common/tmux/model"
)

func (a Agent) ListLocal(ctx context.Context) (model.SessionsById, error) {
	return local.ListSessions(ctx)
}
