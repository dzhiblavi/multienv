package agent

import (
	"context"
	"fmt"

	"multienv/agent/internal/tmux"
	"multienv/common/tmux/operations"
)

type KillWindowParams struct {
	Hostname  string `json:"hostname"`
	SessionId string `json:"session_id"`
	WindowId  string `json:"window_id"`
}

func (a Agent) KillWindow(ctx context.Context, params KillWindowParams) error {
	if params.Hostname == "" || params.SessionId == "" || params.WindowId == "" {
		return fmt.Errorf("hostname, session_id and window_id are required")
	}

	if params.Hostname == a.hostname {
		if err := tmux.KillWindow(ctx, params.SessionId, params.WindowId); err != nil {
			return err
		}

		return a.SyncLocal(ctx)
	}

	_, err := a.client.ForwardAgent(ctx, operations.ForwardAgentRequest{
		Hostname: params.Hostname,
		Command:  "kill-window",
		Input:    params,
	})
	return err
}
