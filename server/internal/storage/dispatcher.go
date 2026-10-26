package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"multienv/common/log"
	"multienv/common/storage/operations"
	"multienv/server/internal/dispatcher"
)

type Dispatcher struct {
	storage Storage
}

func NewDispatcher(storage Storage) Dispatcher {
	return Dispatcher{storage: storage}
}

func (d Dispatcher) Dispatch(ctx context.Context, payload string) (any, error) {
	var operation operations.Operation
	if err := json.Unmarshal([]byte(payload), &operation); err != nil {
		return nil, fmt.Errorf("decode operation: %w: %w", dispatcher.ErrInvalidParams, err)
	}

	log.AddCustomLogData(ctx, "action", operation.Action)
	log.AddCustomLogData(ctx, "dispatcher", "storage")
	log.Info(ctx, "dispatch")

	switch operation.Action {
	case "get":
		var req operations.GetRequest
		if err := json.Unmarshal(operation.Request, &req); err != nil {
			return nil, dispatcher.ErrInvalidParams
		}

		return d.Get(ctx, req)
	case "set":
		var req operations.SetRequest
		if err := json.Unmarshal(operation.Request, &req); err != nil {
			return nil, dispatcher.ErrInvalidParams
		}

		return d.Set(ctx, req)
	case "set-if-absent":
		var req operations.SetIfAbsentRequest
		if err := json.Unmarshal(operation.Request, &req); err != nil {
			return nil, dispatcher.ErrInvalidParams
		}

		return d.SetIfAbsent(ctx, req)
	case "delete":
		var req operations.DeleteRequest
		if err := json.Unmarshal(operation.Request, &req); err != nil {
			return nil, dispatcher.ErrInvalidParams
		}

		return d.Delete(ctx, req)
	}

	return nil, fmt.Errorf("unsupported action %s: %w", operation.Action, dispatcher.ErrUnknownAction)
}
