package service

import (
	"context"
	"fmt"

	"multienv/agent/internal/tmux/local"
	agentoperations "multienv/common/agent/operations"
	"multienv/common/tmux/operations"
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

	op := operations.DeleteSessionRequest{
		Hostname:  a.hostname,
		SessionId: id,
	}

	_, err = a.client.DeleteSession(ctx, op)
	return err
}

func (a Agent) killRemoteSession(ctx context.Context, params KillSessionParams) error {
	op := agentoperations.ForwardRequest{
		Hostname: params.Hostname,
		Domain:   "tmux",
		Command:  "kill-session",
		Input:    params,
	}

	_, err := a.agentClient.Forward(ctx, op)
	if err != nil {
		return fmt.Errorf("forward agent call failed: %v", err)
	}

	return nil
}
