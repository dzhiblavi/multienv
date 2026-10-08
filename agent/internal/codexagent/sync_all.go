package agent

import (
	"context"

	"multienv/common/codex/operations"
)

func (a Agent) SyncAll(ctx context.Context) error {
	_, err := a.client.SyncAllSessions(ctx, operations.SyncAllSessionsRequest{})
	return err
}
