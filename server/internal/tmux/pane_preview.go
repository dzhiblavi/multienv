package tmux

import (
	"context"
	"encoding/json"
	"fmt"

	"multienv/common/tmux/operations"
	"multienv/server/internal/dispatcher"
)

func (s Dispatcher) PanePreview(
	ctx context.Context,
	req operations.PanePreviewRequest,
) (*operations.PanePreviewResponse, error) {
	if req.Hostname == "" || req.Target == "" {
		return nil, fmt.Errorf("hostname and target are required: %w", dispatcher.ErrInvalidParams)
	}

	forwarded, err := s.ForwardAgent(ctx, operations.ForwardAgentRequest{
		Hostname: req.Hostname,
		Command:  "pane-preview",
		Input:    req,
	})
	if err != nil {
		return nil, fmt.Errorf("request pane preview: %w: %w", dispatcher.ErrInternal, err)
	}

	var response operations.PanePreviewResponse
	if err := json.Unmarshal([]byte(forwarded.Output), &response); err != nil {
		return nil, fmt.Errorf("decode pane preview: %w: %w", dispatcher.ErrInternal, err)
	}

	return &response, nil
}
