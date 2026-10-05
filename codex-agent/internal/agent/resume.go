package agent

import (
	"context"
	"fmt"
	"log"

	"multienv/codex-agent/internal/codex"
	"multienv/codex-agent/internal/tmux"
	"multienv/common/codex/model"
	"multienv/common/codex/operations"
)

type ResumeSessionParams struct {
	Hostname   string `json:"hostname"`
	SessionId  string `json:"session_id"`
	TmuxPaneId string `json:"tmux_pane_id"`
}

func (a Agent) ResumeSession(ctx context.Context, params ResumeSessionParams) error {
	if params.Hostname == a.hostname {
		return a.resumeLocalSession(ctx, params)
	} else {
		return a.resumeRemoteSession(ctx, params)
	}
}

func (a Agent) resumeLocalSession(ctx context.Context, params ResumeSessionParams) error {
	session, err := codex.GetSession(ctx, params.SessionId)
	if err != nil {
		return err
	}

	if session.State == model.Archived {
		err = codex.UnarchiveSession(ctx, params.SessionId)
		if err != nil {
			return fmt.Errorf("failed to unarchive session: %v", err)
		}

		session.State = model.NotLoaded
	}

	if session.TmuxPaneId != "" || session.State != model.NotLoaded {
		log.Printf("session is already live in another pane")
		return nil
	}

	err = tmux.SendKeys(
		ctx,
		params.TmuxPaneId,
		[]string{
			fmt.Sprintf(
				"cd %s && codex --no-daemon resume %s",
				session.WorkDir,
				params.SessionId,
			), "Enter",
		},
	)
	if err != nil {
		return fmt.Errorf("resuming session in tmux pane: %v", err)
	}

	session.TmuxPaneId = params.TmuxPaneId

	op := operations.SetSessionRequest{
		Hostname: session.Hostname,
		Session:  session,
	}

	_, err = a.client.SetSession(ctx, op)
	return err
}

func (a Agent) resumeRemoteSession(ctx context.Context, params ResumeSessionParams) error {
	op := operations.ForwardAgentRequest{
		Hostname: params.Hostname,
		Command:  "resume-session",
		Input:    params,
	}

	_, err := a.client.ForwardAgent(ctx, op)
	if err != nil {
		return fmt.Errorf("forward agent call failed: %v", err)
	}

	return nil
}
