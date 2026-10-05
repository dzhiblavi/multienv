package agent

import (
	"context"
	"fmt"

	"multienv/common/tmux/operations"
	"multienv/tmux-agent/internal/tmux"
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
