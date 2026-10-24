package client

import (
	"context"

	"multienv/common/doctor/operations"
)

type Client interface {
	Ping(context.Context, operations.PingRequest) (*operations.PingResponse, error)
	Global(context.Context, operations.GlobalRequest) (*operations.GlobalResponse, error)
}
