package agent

import (
	"context"
	"fmt"

	"multienv/codex-agent/internal/codex"
	"multienv/codex-agent/internal/tmux"
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
	mapping, err := codex.GetSessionToTmuxPaneMapping(ctx)
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

	err = codex.DeleteSession(ctx, id)
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
	op := operations.ForwardAgentRequest{
		Hostname: params.Hostname,
		Command:  "delete-session",
		Input:    params,
	}

	_, err := a.client.ForwardAgent(ctx, op)
	if err != nil {
		return fmt.Errorf("forward agent call failed: %v", err)
	}

	return nil
}
