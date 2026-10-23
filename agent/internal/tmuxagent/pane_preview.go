package agent

import (
	"context"
	"fmt"

	"multienv/agent/internal/tmux"
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
		return a.client.PanePreview(ctx, req)
	}

	content, err := tmux.PanePreview(ctx, req.Target)
	if err != nil {
		return nil, err
	}

	return &operations.PanePreviewResponse{Content: content}, nil
}
