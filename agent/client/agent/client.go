package client

import (
	"context"

	"multienv/common/agent/operations"
)

type Client interface {
	Forward(ctx context.Context, req operations.ForwardRequest) (*operations.ForwardResponse, error)
}
