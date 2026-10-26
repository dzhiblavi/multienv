package service

import (
	"context"

	storageclient "multienv/agent/storage/client"
	"multienv/agent/tmux/model"
	storagemodel "multienv/common/storage/model"
)

const sessionsKey = "@tmux"

func (a Agent) publishSession(ctx context.Context, session model.Session) error {
	at := storagemodel.Location{Key: sessionsKey, Path: []string{a.hostname, session.Id}}
	return storageclient.SetEncode(ctx, a.storage, at, session)
}
