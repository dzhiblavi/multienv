package doctor

import (
	"context"

	"multienv/common/doctor/operations"
)

func (d Dispatcher) Global(ctx context.Context, req operations.GlobalRequest) (*operations.GlobalResponse, error) {
	return &operations.GlobalResponse{Report: d.service.Diagnose(ctx)}, nil
}
