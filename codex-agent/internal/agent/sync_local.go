package agent

import (
	"context"
	"fmt"

	"multienv/common/codex/operations"
	"multienv/codex-agent/internal/codex"
)

func (a Agent) SyncLocal(ctx context.Context) error {
	localSessions, err := codex.ListSessions(ctx)
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
