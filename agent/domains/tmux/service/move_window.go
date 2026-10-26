package service

import (
	"context"
	"fmt"

	agentclient "multienv/agent/client/agent"
	"multienv/agent/domains/tmux/local"
)

type MoveWindowParams struct {
	Hostname  string `json:"hostname"`
	SessionId string `json:"session_id"`
	WindowId  string `json:"window_id"`
	NewIndex  *int   `json:"new_index"`
}

func (a Agent) MoveWindow(ctx context.Context, params MoveWindowParams) error {
	if params.Hostname == "" || params.SessionId == "" || params.WindowId == "" {
		return fmt.Errorf("hostname, session_id and window_id are required")
	}

	if params.NewIndex == nil || *params.NewIndex < 0 {
		return fmt.Errorf("new_index is required and must be non-negative")
	}

	if params.Hostname == a.hostname {
		return local.MoveWindow(ctx, params.SessionId, params.WindowId, *params.NewIndex)
	}

	_, err := agentclient.Invoke(ctx, a.hostClient, params.Hostname, "tmux", "move-window", params)
	return err
}
