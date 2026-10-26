package storage

import (
	"context"

	"multienv/common/storage/operations"
)

func (d Dispatcher) Set(ctx context.Context, req operations.SetRequest) (*operations.SetResponse, error) {
	if err := validateWrite(req.Location, req.Value); err != nil {
		return nil, err
	}

	if err := d.storage.Set(ctx, req.Location, req.Value); err != nil {
		return nil, err
	}

	return &operations.SetResponse{}, nil
}
