package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"multienv/agent/internal/tmux"
	"multienv/common/tmux/operations"
)

type NewWindowParams struct {
	Hostname  string   `json:"hostname"`
	SessionId string   `json:"session_id"`
	Name      string   `json:"name"`
	WorkDir   string   `json:"work_dir,omitempty"`
	Command   []string `json:"command,omitempty"`
}

type NewWindowResult struct {
	TmuxPaneId string `json:"tmux_pane_id"`
}

func (a Agent) NewWindow(ctx context.Context, params NewWindowParams) (*NewWindowResult, error) {
	if params.Hostname == "" || params.SessionId == "" || params.Name == "" {
		return nil, fmt.Errorf("hostname, session_id and name are required")
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

	op := operations.ForwardAgentRequest{
		Hostname: params.Hostname,
		Command:  "new-window",
		Input:    params,
	}
	response, err := a.client.ForwardAgent(ctx, op)
	if err != nil {
		return nil, fmt.Errorf("forward agent call failed: %w", err)
	}

	var result NewWindowResult
	if err := json.Unmarshal([]byte(response.Output), &result); err != nil {
		return nil, fmt.Errorf("agent's response is invalid: %w", err)
	}
	return &result, nil
}
