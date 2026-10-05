package agent

import (
	"context"
	"fmt"

	"multienv/common/tmux/operations"
	"multienv/tmux-agent/internal/tmux"
)

type RenameSessionParams struct {
	SessionId string `json:"session_id"`
	Hostname  string `json:"hostname"`
	NewName   string `json:"new_name"`
}

func (a Agent) RenameSession(ctx context.Context, params RenameSessionParams) error {
	if params.Hostname == a.hostname {
		return a.renameLocalSession(ctx, params.SessionId, params.NewName)
	} else {
		return a.renameRemoteSession(ctx, params)
	}
}

func (a Agent) renameLocalSession(ctx context.Context, id, newName string) error {
	err := tmux.RenameSession(ctx, id, newName)
	if err != nil {
		return fmt.Errorf("failed to rename session: %v", err)
	}

	sessions, err := tmux.ListSessions(ctx)
	if err != nil {
		return fmt.Errorf("failed to list sessions: %v", err)
	}

	session, ok := sessions[id]
	if !ok {
		return fmt.Errorf("failed to find renamed session")
	}

	op := operations.SetSessionRequest{
		Hostname: a.hostname,
		Session:  session,
	}

	_, err = a.client.SetSession(ctx, op)
	return err
}

func (a Agent) renameRemoteSession(ctx context.Context, params RenameSessionParams) error {
	op := operations.ForwardAgentRequest{
		Hostname: params.Hostname,
		Command:  "rename-session",
		Input:    params,
	}

	_, err := a.client.ForwardAgent(ctx, op)
	if err != nil {
		return fmt.Errorf("forward agent call failed: %v", err)
	}

	return nil
}
