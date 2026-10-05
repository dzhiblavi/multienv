package agent

import (
	"context"

	"multienv/codex-agent/internal/codex"
	"multienv/common/codex/model"
)

func (a Agent) ListLocal(ctx context.Context) (model.SessionsById, error) {
	return codex.ListSessions(ctx)
}
