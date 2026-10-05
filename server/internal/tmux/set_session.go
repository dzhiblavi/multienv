package tmux

import (
	"context"
	"encoding/json"
	"fmt"

	"multienv/common/tmux/operations"
)

func (s Dispatcher) SetSession(
	ctx context.Context,
	req operations.SetSessionRequest,
) (*operations.SetSessionResponse, error) {
	hostname, _ := json.Marshal(req.Hostname)
	sessionID, _ := json.Marshal(req.Session.Id)

	hostPath := fmt.Sprintf("$[%s]", hostname)
	sessionPath := fmt.Sprintf("%s[%s]", hostPath, sessionID)

	err := s.storage.SetIfAbsent(ctx, hostPath, map[string]any{})
	if err != nil {
		return nil, fmt.Errorf("initialize host: %w", err)
	}

	err = s.storage.Set(ctx, sessionPath, req.Session)
	if err != nil {
		return nil, fmt.Errorf("storage request failed: %w", err)
	}

	return &operations.SetSessionResponse{}, nil
}
