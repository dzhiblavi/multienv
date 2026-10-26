package storage

import (
	"context"

	"multienv/common/storage/operations"
)

func (d Dispatcher) Get(ctx context.Context, req operations.GetRequest) (*operations.GetResponse, error) {
	if err := validateLocation(req.Location); err != nil {
		return nil, err
	}

	value, found, err := d.storage.Get(ctx, req.Location)
	if err != nil {
		return nil, err
	}

	return &operations.GetResponse{Value: value, Found: found}, nil
}
