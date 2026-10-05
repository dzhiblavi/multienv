package tmux

import (
	"context"
	"encoding/json"
	"fmt"

	"multienv/common/tmux/operations"
)

func (s Dispatcher) SetSessions(
	ctx context.Context,
	req operations.SetSessionsRequest,
) (*operations.SetSessionsResponse, error) {
	hostname, _ := json.Marshal(req.Hostname)
	path := fmt.Sprintf("$[%s]", hostname)
	err := s.storage.Set(ctx, path, req.Sessions)

	if err != nil {
		return nil, fmt.Errorf("storage request failed: %w", err)
	}

	return &operations.SetSessionsResponse{}, nil
}
