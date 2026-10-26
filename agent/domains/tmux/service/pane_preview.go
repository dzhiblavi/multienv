package service

import (
	"context"
	"fmt"

	agentclient "multienv/agent/client/agent"
	"multienv/agent/domains/tmux/local"
)

type PanePreviewParams struct {
	Hostname string `json:"hostname"`
	Target   string `json:"target"`
}

type PanePreviewResult struct {
	Content string `json:"content"`
}

func (a Agent) PanePreview(
	ctx context.Context,
	req PanePreviewParams,
) (*PanePreviewResult, error) {
	if req.Hostname == "" || req.Target == "" {
		return nil, fmt.Errorf("hostname and target are required")
	}

	if req.Hostname != a.hostname {
		response, err := agentclient.InvokeDecode[PanePreviewResult](
			ctx, a.hostClient, req.Hostname, "tmux", "pane-preview", req,
		)
		if err != nil {
			return nil, err
		}

		return &response, nil
	}

	content, err := local.PanePreview(ctx, req.Target)
	if err != nil {
		return nil, err
	}

	return &PanePreviewResult{Content: content}, nil
}
