package agent

import (
	"context"
	"fmt"

	"multienv/agent/internal/tmux"
	"multienv/common/tmux/operations"
)

func (a Agent) SyncLocal(ctx context.Context) error {
	localSessions, err := tmux.ListSessions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get local sessions: %v", err)
	}

	op := operations.SetSessionsRequest{
		Hostname: a.hostname,
		Sessions: localSessions,
	}

	_, err = a.client.SetSessions(ctx, op)
	return err
}
