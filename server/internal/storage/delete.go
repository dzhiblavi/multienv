package storage

import (
	"context"

	"multienv/common/storage/operations"
)

func (d Dispatcher) Delete(ctx context.Context, req operations.DeleteRequest) (*operations.DeleteResponse, error) {
	if err := validateLocation(req.Location); err != nil {
		return nil, err
	}

	deleted, err := d.storage.Delete(ctx, req.Location)
	if err != nil {
		return nil, err
	}

	return &operations.DeleteResponse{Deleted: deleted}, nil
}
