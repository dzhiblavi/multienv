package service

import (
	"context"
	"fmt"

	"multienv/agent/domains/codex/local"
	"multienv/agent/domains/codex/model"
	storageclient "multienv/agent/domains/storage/client"
	"multienv/common/log"
	storagemodel "multienv/common/storage/model"
)

func (a Agent) ListAll(ctx context.Context) (model.SessionsByHost, error) {
	sessions := make(model.SessionsByHost)

	snapshot, _, err := storageclient.GetDecode[model.SessionsByHost](
		ctx, a.storage, storagemodel.Location{Key: sessionsKey},
	)
	if err != nil {
		log.Warn(ctx, "failed to get all sessions; using local sessions", "error", err)
	} else if snapshot != nil {
		sessions = snapshot
	}

	localSessions, err := local.ListSessions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get local sessions: %v", err)
	}

	delete(sessions, a.hostname)
	sessions[a.hostname] = localSessions

	return sessions, nil
}
