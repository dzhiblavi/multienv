package service

import (
	"context"
	"fmt"

	"multienv/agent/domains/codex/local"
	storageclient "multienv/agent/domains/storage/client"
	storagemodel "multienv/common/storage/model"
)

func (a Agent) SyncLocal(ctx context.Context) error {
	localSessions, err := local.ListSessions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get local sessions: %v", err)
	}

	at := storagemodel.Location{Key: sessionsKey, Path: []string{a.hostname}}
	return storageclient.SetEncode(ctx, a.storage, at, localSessions)
}
