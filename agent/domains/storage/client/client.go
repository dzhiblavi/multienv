package client

import (
	"context"

	"multienv/common/storage/operations"
)

type Client interface {
	Get(context.Context, operations.GetRequest) (*operations.GetResponse, error)
	Set(context.Context, operations.SetRequest) (*operations.SetResponse, error)
	SetIfAbsent(context.Context, operations.SetIfAbsentRequest) (*operations.SetIfAbsentResponse, error)
	Delete(context.Context, operations.DeleteRequest) (*operations.DeleteResponse, error)
}
