package agent

import (
	"context"
	"fmt"

	"multienv/common/log"
	"multienv/common/tmux/model"
	"multienv/common/tmux/operations"
	"multienv/tmux-agent/internal/tmux"
)

func (a Agent) ListAll(ctx context.Context) (model.SessionsByHost, error) {
	sessions := make(model.SessionsByHost)

	response, err := a.client.ListAll(ctx, operations.ListAllRequest{})
	if err != nil {
		log.Warn(ctx, "failed to get all sessions; using local sessions", "error", err)
	} else {
		sessions = response.Sessions
	}

	localSessions, err := tmux.ListSessions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get local sessions: %v", err)
	}

	delete(sessions, a.hostname)
	sessions[a.hostname] = localSessions

	return sessions, nil
}
