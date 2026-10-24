package service

import (
	"context"

	"multienv/agent/internal/codex/local"
	"multienv/common/codex/model"
)

func (a Agent) ListLocal(ctx context.Context) (model.SessionsById, error) {
	return local.ListSessions(ctx)
}
