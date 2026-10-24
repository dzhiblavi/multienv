package service

import (
	"context"
	"fmt"

	agentclient "multienv/agent/client/agent"
	"multienv/agent/internal/tmux/local"
	agentoperations "multienv/common/agent/operations"
	"multienv/common/tmux/model"
	"multienv/common/tmux/operations"
)

type NewSessionParams struct {
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
}

func (a Agent) NewSession(ctx context.Context, params NewSessionParams) (*model.Session, error) {
	if params.Hostname == a.hostname {
		return a.newLocalSession(ctx, params.Name)
	} else {
		return a.newRemoteSession(ctx, params)
	}
}

func (a Agent) newLocalSession(ctx context.Context, name string) (*model.Session, error) {
	session, err := local.NewSession(ctx, name)

	if err != nil {
		return nil, fmt.Errorf("failed to create a new session: %v", err)
	}

	op := operations.SetSessionRequest{
		Hostname: a.hostname,
		Session:  session,
	}

	_, err = a.client.SetSession(ctx, op)
	return &session, err
}

func (a Agent) newRemoteSession(ctx context.Context, params NewSessionParams) (*model.Session, error) {
	session, err := agentclient.Forward[model.Session](ctx, a.agentClient, agentoperations.ForwardRequest{
		Hostname: params.Hostname,
		Domain:   "tmux",
		Command:  "new-session",
		Input:    params,
	})
	if err != nil {
		return nil, err
	}

	return &session, nil
}
