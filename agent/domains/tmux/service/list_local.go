package service

import (
	"context"

	"multienv/agent/domains/tmux/local"
	"multienv/agent/domains/tmux/model"
)

func (a Agent) ListLocal(ctx context.Context) (model.SessionsById, error) {
	return local.ListSessions(ctx)
}
