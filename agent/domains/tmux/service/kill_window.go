package service

import (
	"context"
	"fmt"

	agentclient "multienv/agent/client/agent"
	"multienv/agent/domains/tmux/local"
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
		if err := local.KillWindow(ctx, params.SessionId, params.WindowId); err != nil {
			return err
		}

		return a.SyncLocal(ctx)
	}

	_, err := agentclient.Invoke(ctx, a.hostClient, params.Hostname, "tmux", "kill-window", params)
	return err
}
