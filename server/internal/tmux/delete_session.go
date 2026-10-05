package tmux

import (
	"context"
	"encoding/json"
	"fmt"

	"multienv/common/tmux/operations"
)

func (s Dispatcher) DeleteSession(
	ctx context.Context,
	req operations.DeleteSessionRequest,
) (*operations.DeleteSessionResponse, error) {
	hostname, _ := json.Marshal(req.Hostname)
	sessionID, _ := json.Marshal(req.SessionId)

	hostPath := fmt.Sprintf("$[%s]", hostname)
	sessionPath := fmt.Sprintf("%s[%s]", hostPath, sessionID)

	err := s.storage.Delete(ctx, sessionPath)
	if err != nil {
		return nil, err
	}

	return &operations.DeleteSessionResponse{}, nil
}
