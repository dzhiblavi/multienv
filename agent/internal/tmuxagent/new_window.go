package agent

import (
	"context"
	"fmt"

	client "multienv/agent/internal/client/tmux"
	"multienv/agent/internal/tmux"
	"multienv/common/tmux/operations"
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
		paneId, err := tmux.NewWindow(
			ctx, params.SessionId, params.Name, params.WorkDir, params.Command,
		)
		if err != nil {
			return nil, err
		}
		return &NewWindowResult{TmuxPaneId: paneId}, nil
	}

	result, err := client.ForwardAgent[NewWindowResult](ctx, a.client, operations.ForwardAgentRequest{
		Hostname: params.Hostname,
		Command:  "new-window",
		Input:    params,
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}
