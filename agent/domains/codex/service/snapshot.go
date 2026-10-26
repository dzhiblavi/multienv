package service

import (
	"context"

	"multienv/agent/domains/codex/model"
	storageclient "multienv/agent/domains/storage/client"
	storagemodel "multienv/common/storage/model"
)

const sessionsKey = "@codex"

func (a Agent) publishSession(ctx context.Context, session model.Session) error {
	at := storagemodel.Location{Key: sessionsKey, Path: []string{a.hostname, session.Id}}
	return storageclient.SetEncode(ctx, a.storage, at, session)
}
