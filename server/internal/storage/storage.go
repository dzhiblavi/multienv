package storage

import (
	"context"
	"encoding/json"

	"multienv/common/storage/model"
)

type Storage interface {
	Get(ctx context.Context, at model.Location) (json.RawMessage, bool, error)
	Set(ctx context.Context, at model.Location, value json.RawMessage) error
	SetIfAbsent(ctx context.Context, at model.Location, value json.RawMessage) (bool, error)
	Delete(ctx context.Context, at model.Location) (bool, error)
}
