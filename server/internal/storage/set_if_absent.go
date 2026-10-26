package storage

import (
	"context"

	"multienv/common/storage/operations"
)

func (d Dispatcher) SetIfAbsent(
	ctx context.Context,
	req operations.SetIfAbsentRequest,
) (*operations.SetIfAbsentResponse, error) {
	if err := validateWrite(req.Location, req.Value); err != nil {
		return nil, err
	}

	created, err := d.storage.SetIfAbsent(ctx, req.Location, req.Value)
	if err != nil {
		return nil, err
	}

	return &operations.SetIfAbsentResponse{Created: created}, nil
}
