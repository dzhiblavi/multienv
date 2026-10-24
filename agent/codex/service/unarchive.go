package service

import (
	"context"
	"fmt"

	"multienv/agent/internal/codex/local"
	agentoperations "multienv/common/agent/operations"
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
	err := local.UnarchiveSession(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to unarchive codex session: %v", err)
	}

	session, err := local.GetSession(ctx, id)
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
	op := agentoperations.ForwardRequest{
		Hostname: params.Hostname,
		Domain:   "codex",
		Command:  "unarchive-session",
		Input:    params,
	}

	_, err := a.agentClient.Forward(ctx, op)
	if err != nil {
		return fmt.Errorf("forward agent call failed: %v", err)
	}

	return nil
}
