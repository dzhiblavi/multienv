package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"multienv/agent/internal/codex"
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
	session, err := codex.GetSession(ctx, id)
	if err != nil {
		return model.Session{}, fmt.Errorf("get Codex session: %w", err)
	}

	if session.State == model.Archived {
		if err := codex.UnarchiveSession(ctx, id); err != nil {
			return model.Session{}, fmt.Errorf("unarchive Codex session: %w", err)
		}

		// Check the actual state after unarchiving before allowing a new process.
		session, err = codex.GetSession(ctx, id)
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
	response, err := a.client.ForwardAgent(ctx, operations.ForwardAgentRequest{
		Hostname: params.Hostname,
		Command:  "prepare-resume",
		Input:    params,
	})
	if err != nil {
		return model.Session{}, fmt.Errorf("forward prepare-resume: %w", err)
	}

	var session model.Session
	if err := json.Unmarshal([]byte(response.Output), &session); err != nil {
		return model.Session{}, fmt.Errorf("invalid prepare-resume response: %w", err)
	}

	return session, nil
}
