package agent

import (
	"context"
	"fmt"

	"multienv/agent/internal/tmux"
	"multienv/common/tmux/operations"
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
		return tmux.RenameWindow(ctx, params.SessionId, params.WindowId, params.NewName)
	}

	_, err := a.client.ForwardAgent(ctx, operations.ForwardAgentRequest{
		Hostname: params.Hostname,
		Command:  "rename-window",
		Input:    params,
	})
	return err
}
