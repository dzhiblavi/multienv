package service

import (
	"context"

	"multienv/agent/codex/model"
	"multienv/agent/internal/codex/local"
)

func (a Agent) ListLocal(ctx context.Context) (model.SessionsById, error) {
	return local.ListSessions(ctx)
}
