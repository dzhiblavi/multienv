package agent

import (
	"context"
	"fmt"
	"log"

	"multienv/codex-agent/internal/codex"
	"multienv/common/codex/model"
	"multienv/common/codex/operations"
)

func (a Agent) ListAll(ctx context.Context) (model.SessionsByHost, error) {
	sessions := make(model.SessionsByHost)

	response, err := a.client.ListAll(ctx, operations.ListAllRequest{})
	if err != nil {
		log.Printf("failed to get all sessions: %v, ignoring", err)
	} else {
		sessions = response.Sessions
	}

	localSessions, err := codex.ListSessions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get local sessions: %v", err)
	}

	delete(sessions, a.hostname)
	sessions[a.hostname] = localSessions

	return sessions, nil
}
