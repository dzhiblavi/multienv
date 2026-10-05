package storage

import (
	"context"
)

// JSON-like storage
type Storage interface {
	Delete(ctx context.Context, path string) error
	Get(ctx context.Context, path string, dst any) error
	Set(ctx context.Context, path string, value any) error
	SetIfAbsent(ctx context.Context, path string, value any) error
	Init(ctx context.Context) error

	Close()
}

func Get[T any](ctx context.Context, storage Storage, path string) (T, error) {
	var result T
	err := storage.Get(ctx, path, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
