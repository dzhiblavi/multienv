package service

import (
	"context"
	"fmt"

	agentclient "multienv/agent/client/agent"
	"multienv/agent/domains/codex/local"
)

type UnarchiveSessionParams struct {
	SessionId string `json:"session_id"`
	Hostname  string `json:"hostname"`
}

func (a Agent) UnarchiveSession(ctx context.Context, params UnarchiveSessionParams) error {
	if params.Hostname == a.hostname {
		return a.unarchiveLocalSession(ctx, params.SessionId)
	} else {
		return a.unarchiveRemoteSession(ctx, params)
	}
}

func (a Agent) unarchiveLocalSession(ctx context.Context, id string) error {
	err := local.UnarchiveSession(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to unarchive codex session: %v", err)
	}

	session, err := local.GetSession(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get session: %v", err)
	}

	err = a.publishSession(ctx, session)
	return err
}

func (a Agent) unarchiveRemoteSession(ctx context.Context, params UnarchiveSessionParams) error {
	_, err := agentclient.Invoke(ctx, a.hostClient, params.Hostname, "codex", "unarchive-session", params)
	if err != nil {
		return fmt.Errorf("invoke agent failed: %v", err)
	}

	return nil
}
