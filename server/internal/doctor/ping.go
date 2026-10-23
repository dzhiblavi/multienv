package doctor

import (
	"context"

	"multienv/common/doctor/operations"
)

func (d Dispatcher) Ping(ctx context.Context, req operations.PingRequest) (*operations.PingResponse, error) {
	return &operations.PingResponse{Service: "multienv", Hostname: d.service.hostname}, nil
}
