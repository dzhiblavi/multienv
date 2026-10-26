package service

import (
	"context"
	"fmt"

	agentclient "multienv/agent/client/agent"
	"multienv/agent/internal/codex/local"
	tmux "multienv/agent/internal/tmux/local"
	"multienv/common/log"
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

	err = local.ArchiveSession(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to archive codex session: %v", err)
	}

	session, err := local.GetSession(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get session: %v", err)
	}

	err = a.publishSession(ctx, session)
	return err
}

func (a Agent) archiveRemoteSession(ctx context.Context, params ArchiveSessionParams) error {
	_, err := agentclient.Invoke(ctx, a.hostClient, params.Hostname, "codex", "archive-session", params)
	if err != nil {
		return fmt.Errorf("invoke agent failed: %v", err)
	}

	return nil
}
