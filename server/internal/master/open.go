package master

import (
	"context"
	"fmt"

	"multienv/common/master/operations"
	"multienv/common/subprocess"
	"multienv/server/internal/dispatcher"
)

func (s Dispatcher) Open(ctx context.Context, req operations.OpenRequest) (*operations.OpenResponse, error) {
	if req.Target == "" {
		return nil, fmt.Errorf("open target is required: %w", dispatcher.ErrInvalidParams)
	}

	_, err := subprocess.RunSuccess(ctx, subprocess.Command{
		Name: "/usr/bin/open",
		Args: []string{"--", req.Target},
	})
	if err != nil {
		return nil, fmt.Errorf("open target: %w", err)
	}
	return &operations.OpenResponse{}, nil
}
