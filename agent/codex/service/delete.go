package service

import (
	"context"
	"fmt"

	agentclient "multienv/agent/client/agent"
	"multienv/agent/internal/codex/local"
	tmux "multienv/agent/internal/tmux/local"
	"multienv/common/codex/operations"
	"multienv/common/log"
)

type DeleteSessionParams struct {
	SessionId string `json:"session_id"`
	Hostname  string `json:"hostname"`
}

func (a Agent) DeleteSession(ctx context.Context, params DeleteSessionParams) error {
	if params.Hostname == a.hostname {
		return a.deleteLocalSession(ctx, params.SessionId)
	} else {
		return a.deleteRemoteSession(ctx, params)
	}
}

func (a Agent) deleteLocalSession(ctx context.Context, id string) error {
	mapping, err := local.GetSessionToTmuxPaneMapping(ctx)
	if err != nil {
		return fmt.Errorf("failed to list local active sessions: %v", err)
	}

	tmuxPane, ok := mapping[id]
	if ok {
		if err = tmux.KillPane(ctx, tmuxPane); err != nil {
			log.Warn(ctx, "failed to kill tmux pane with the session",
				"session_id", id, "pane_id", tmuxPane, "error", err)
		}
	}

	err = local.DeleteSession(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete codex session: %v", err)
	}

	op := operations.DeleteSessionRequest{
		Hostname:  a.hostname,
		SessionId: id,
	}

	_, err = a.client.DeleteSession(ctx, op)
	return err
}

func (a Agent) deleteRemoteSession(ctx context.Context, params DeleteSessionParams) error {
	_, err := agentclient.Invoke(ctx, a.hostClient, params.Hostname, "codex", "delete-session", params)
	if err != nil {
		return fmt.Errorf("invoke agent failed: %v", err)
	}

	return nil
}
