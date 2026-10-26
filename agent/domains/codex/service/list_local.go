package service

import (
	"context"

	"multienv/agent/domains/codex/local"
	"multienv/agent/domains/codex/model"
)

func (a Agent) ListLocal(ctx context.Context) (model.SessionsById, error) {
	return local.ListSessions(ctx)
}
