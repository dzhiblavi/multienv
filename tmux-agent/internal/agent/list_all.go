package agent

import (
	"context"
	"fmt"
	"log"

	"multienv/common/tmux/model"
	"multienv/common/tmux/operations"
	"multienv/tmux-agent/internal/tmux"
)

func (a Agent) ListAll(ctx context.Context) (model.SessionsByHost, error) {
	sessions := make(model.SessionsByHost)

	response, err := a.client.ListAll(ctx, operations.ListAllRequest{})
	if err != nil {
		log.Printf("failed to get all sessions: %v, ignoring", err)
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
