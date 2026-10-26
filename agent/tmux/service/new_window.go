package service

import (
	"context"
	"fmt"

	agentclient "multienv/agent/client/agent"
	"multienv/agent/internal/tmux/local"
)

type NewWindowParams struct {
	Hostname  string   `json:"hostname"`
	SessionId string   `json:"session_id"`
	Name      string   `json:"name,omitempty"`
	WorkDir   string   `json:"work_dir,omitempty"`
	Command   []string `json:"command,omitempty"`
}

type NewWindowResult struct {
	TmuxPaneId string `json:"tmux_pane_id"`
}

func (a Agent) NewWindow(ctx context.Context, params NewWindowParams) (*NewWindowResult, error) {
	if params.Hostname == "" || params.SessionId == "" {
		return nil, fmt.Errorf("hostname and session_id are required")
	}

	if params.Hostname == a.hostname {
		paneId, err := local.NewWindow(
			ctx, params.SessionId, params.Name, params.WorkDir, params.Command,
		)
		if err != nil {
			return nil, err
		}
		return &NewWindowResult{TmuxPaneId: paneId}, nil
	}

	result, err := agentclient.InvokeDecode[NewWindowResult](
		ctx, a.hostClient, params.Hostname, "tmux", "new-window", params,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
