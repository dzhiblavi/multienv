package agent

import (
	"context"
	"fmt"

	"multienv/agent/internal/codex"
	"multienv/common/codex/operations"
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
	err := codex.UnarchiveSession(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to unarchive codex session: %v", err)
	}

	session, err := codex.GetSession(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get session: %v", err)
	}

	op := operations.SetSessionRequest{
		Hostname: a.hostname,
		Session:  session,
	}

	_, err = a.client.SetSession(ctx, op)
	return err
}

func (a Agent) unarchiveRemoteSession(ctx context.Context, params UnarchiveSessionParams) error {
	op := operations.ForwardAgentRequest{
		Hostname: params.Hostname,
		Command:  "unarchive-session",
		Input:    params,
	}

	_, err := a.client.ForwardAgent(ctx, op)
	if err != nil {
		return fmt.Errorf("forward agent call failed: %v", err)
	}

	return nil
}
