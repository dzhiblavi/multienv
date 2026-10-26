package host

import (
	"context"

	"multienv/common/host/operations"
)

type Invoker interface {
	Invoke(context.Context, operations.InvokeRequest) (*operations.InvokeResponse, error)
}
