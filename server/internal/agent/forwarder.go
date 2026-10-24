package agent

import (
	"context"

	"multienv/common/agent/operations"
)

type Forwarder interface {
	Forward(ctx context.Context, req operations.ForwardRequest) (*operations.ForwardResponse, error)
}
