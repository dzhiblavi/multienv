package service

import (
	"context"
	"fmt"
	"strings"

	agentclient "multienv/agent/client/agent"
	"multienv/agent/internal/codex/local"
	"multienv/common/codex/operations"
)

type RenameSessionParams struct {
	Hostname  string `json:"hostname"`
	SessionId string `json:"session_id"`
	NewName   string `json:"new_name"`
}

func (a Agent) RenameSession(ctx context.Context, params RenameSessionParams) error {
	if params.Hostname == "" || params.SessionId == "" || strings.TrimSpace(params.NewName) == "" {
		return fmt.Errorf("hostname, session_id and new_name are required")
	}

	if params.Hostname == a.hostname {
		return a.renameLocalSession(ctx, params.SessionId, params.NewName)
	}

	_, err := agentclient.Invoke(ctx, a.hostClient, params.Hostname, "codex", "rename-session", params)
	return err
}

func (a Agent) renameLocalSession(ctx context.Context, id, newName string) error {
	if err := local.RenameSession(ctx, id, newName); err != nil {
		return err
	}

	session, err := local.GetSession(ctx, id)
	if err != nil {
		return fmt.Errorf("get renamed Codex session: %w", err)
	}

	_, err = a.client.SetSession(ctx, operations.SetSessionRequest{
		Hostname: a.hostname,
		Session:  session,
	})
	if err != nil {
		return fmt.Errorf("session renamed, but publishing its snapshot failed: %w", err)
	}

	return nil
}
