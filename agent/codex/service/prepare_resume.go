package service

import (
	"context"
	"fmt"
	"os"

	agentclient "multienv/agent/client/agent"
	"multienv/agent/internal/codex/local"
	agentoperations "multienv/common/agent/operations"
	"multienv/common/codex/model"
	"multienv/common/codex/operations"
)

type PrepareResumeParams struct {
	Hostname  string `json:"hostname"`
	SessionId string `json:"session_id"`
}

func (a Agent) PrepareResume(
	ctx context.Context,
	params PrepareResumeParams,
) (model.Session, error) {
	if params.Hostname == "" || params.SessionId == "" {
		return model.Session{}, fmt.Errorf("hostname and session_id are required")
	}

	if params.Hostname == a.hostname {
		return a.prepareLocalResume(ctx, params.SessionId)
	}

	return a.prepareRemoteResume(ctx, params)
}

func (a Agent) prepareLocalResume(ctx context.Context, id string) (model.Session, error) {
	session, err := local.GetSession(ctx, id)
	if err != nil {
		return model.Session{}, fmt.Errorf("get Codex session: %w", err)
	}

	if session.State == model.Archived {
		if err := local.UnarchiveSession(ctx, id); err != nil {
			return model.Session{}, fmt.Errorf("unarchive Codex session: %w", err)
		}

		// Check the actual state after unarchiving before allowing a new process.
		session, err = local.GetSession(ctx, id)
		if err != nil {
			return model.Session{}, fmt.Errorf("get unarchived Codex session: %w", err)
		}
	}

	if session.State != model.NotLoaded || session.TmuxPaneId != "" {
		return model.Session{}, fmt.Errorf("session %s is not ready to resume: %s", id, session.State)
	}

	info, err := os.Stat(session.WorkDir)
	if err != nil {
		return model.Session{}, fmt.Errorf("check session directory: %w", err)
	}
	if !info.IsDir() {
		return model.Session{}, fmt.Errorf("session directory is not a directory: %s", session.WorkDir)
	}

	_, err = a.client.SetSession(ctx, operations.SetSessionRequest{
		Hostname: a.hostname,
		Session:  session,
	})
	if err != nil {
		return model.Session{}, fmt.Errorf("publish prepared Codex session: %w", err)
	}

	return session, nil
}

func (a Agent) prepareRemoteResume(
	ctx context.Context,
	params PrepareResumeParams,
) (model.Session, error) {
	return agentclient.Forward[model.Session](ctx, a.agentClient, agentoperations.ForwardRequest{
		Hostname: params.Hostname,
		Domain:   "codex",
		Command:  "prepare-resume",
		Input:    params,
	})
}
