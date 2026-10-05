package agent

import (
	"context"
	"fmt"
	"log"

	"multienv/codex-agent/internal/codex"
	"multienv/codex-agent/internal/tmux"
	"multienv/common/codex/operations"
)

type ArchiveSessionParams struct {
	SessionId string `json:"session_id"`
	Hostname  string `json:"hostname"`
}

func (a Agent) ArchiveSession(ctx context.Context, params ArchiveSessionParams) error {
	if params.Hostname == a.hostname {
		return a.archiveLocalSession(ctx, params.SessionId)
	} else {
		return a.archiveRemoteSession(ctx, params)
	}
}

func (a Agent) archiveLocalSession(ctx context.Context, id string) error {
	mapping, err := codex.GetSessionToTmuxPaneMapping(ctx)
	if err != nil {
		return fmt.Errorf("failed to list local active sessions: %v", err)
	}

	tmuxPane, ok := mapping[id]
	if ok {
		if err = tmux.KillPane(ctx, tmuxPane); err != nil {
			log.Printf("failed to kill tmux pane with the session: %v", err)
		}
	}

	err = codex.ArchiveSession(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to archive codex session: %v", err)
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

func (a Agent) archiveRemoteSession(ctx context.Context, params ArchiveSessionParams) error {
	op := operations.ForwardAgentRequest{
		Hostname: params.Hostname,
		Command:  "archive-session",
		Input:    params,
	}

	_, err := a.client.ForwardAgent(ctx, op)
	if err != nil {
		return fmt.Errorf("forward agent call failed: %v", err)
	}

	return nil
}
