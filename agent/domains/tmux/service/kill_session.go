package service

import (
	"context"
	"fmt"
	storagemodel "multienv/common/storage/model"
	storageops "multienv/common/storage/operations"

	agentclient "multienv/agent/client/agent"
	"multienv/agent/domains/tmux/local"
)

type KillSessionParams struct {
	SessionId string `json:"session_id"`
	Hostname  string `json:"hostname"`
}

func (a Agent) KillSession(ctx context.Context, params KillSessionParams) error {
	if params.Hostname == a.hostname {
		return a.killLocalSession(ctx, params.SessionId)
	} else {
		return a.killRemoteSession(ctx, params)
	}
}

func (a Agent) killLocalSession(ctx context.Context, id string) error {
	err := local.KillSession(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to kill session: %v", err)
	}

	at := storagemodel.Location{Key: sessionsKey, Path: []string{a.hostname, id}}
	_, err = a.storage.Delete(ctx, storageops.DeleteRequest{Location: at})
	return err
}

func (a Agent) killRemoteSession(ctx context.Context, params KillSessionParams) error {
	_, err := agentclient.Invoke(ctx, a.hostClient, params.Hostname, "tmux", "kill-session", params)
	if err != nil {
		return fmt.Errorf("invoke agent failed: %v", err)
	}

	return nil
}
