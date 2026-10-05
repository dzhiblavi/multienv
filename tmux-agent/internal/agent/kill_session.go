package agent

import (
	"context"
	"fmt"

	"multienv/common/tmux/operations"
	"multienv/tmux-agent/internal/tmux"
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
	err := tmux.KillSession(ctx, id)
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
	op := operations.ForwardAgentRequest{
		Hostname: params.Hostname,
		Command:  "kill-session",
		Input:    params,
	}

	_, err := a.client.ForwardAgent(ctx, op)
	if err != nil {
		return fmt.Errorf("forward agent call failed: %v", err)
	}

	return nil
}
