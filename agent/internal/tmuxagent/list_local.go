package agent

import (
	"context"

	"multienv/agent/internal/tmux"
	"multienv/common/tmux/model"
)

func (a Agent) ListLocal(ctx context.Context) (model.SessionsById, error) {
	return tmux.ListSessions(ctx)
}
