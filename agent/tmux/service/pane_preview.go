package service

import (
	"context"
	"fmt"

	agentclient "multienv/agent/client/agent"
	"multienv/agent/internal/tmux/local"
	"multienv/common/tmux/operations"
)

func (a Agent) PanePreview(
	ctx context.Context,
	req operations.PanePreviewRequest,
) (*operations.PanePreviewResponse, error) {
	if req.Hostname == "" || req.Target == "" {
		return nil, fmt.Errorf("hostname and target are required")
	}

	if req.Hostname != a.hostname {
		response, err := agentclient.InvokeDecode[operations.PanePreviewResponse](
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

	return &operations.PanePreviewResponse{Content: content}, nil
}
