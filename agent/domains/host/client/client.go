package client

import (
	"context"

	"multienv/common/host/operations"
)

type Client interface {
	Invoke(context.Context, operations.InvokeRequest) (*operations.InvokeResponse, error)
}
