package service

import (
	"context"
	"fmt"

	"multienv/agent/internal/tmux/local"
	agentoperations "multienv/common/agent/operations"
)

type RenameWindowParams struct {
	Hostname  string `json:"hostname"`
	SessionId string `json:"session_id"`
	WindowId  string `json:"window_id"`
	NewName   string `json:"new_name"`
}

func (a Agent) RenameWindow(ctx context.Context, params RenameWindowParams) error {
	if params.Hostname == "" || params.SessionId == "" || params.WindowId == "" || params.NewName == "" {
		return fmt.Errorf("hostname, session_id, window_id and new_name are required")
	}

	if params.Hostname == a.hostname {
		return local.RenameWindow(ctx, params.SessionId, params.WindowId, params.NewName)
	}

	_, err := a.agentClient.Forward(ctx, agentoperations.ForwardRequest{
		Hostname: params.Hostname,
		Domain:   "tmux",
		Command:  "rename-window",
		Input:    params,
	})
	return err
}
