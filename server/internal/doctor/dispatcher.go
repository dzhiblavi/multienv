package doctor

import (
	"context"
	"encoding/json"
	"fmt"

	"multienv/common/doctor/operations"
	"multienv/common/log"
	"multienv/server/internal/dispatcher"
)

type Dispatcher struct {
	service *Service
}

func NewDispatcher(service *Service) Dispatcher {
	return Dispatcher{service: service}
}

func (d Dispatcher) Dispatch(ctx context.Context, payload string) (any, error) {
	var operation operations.Operation
	if err := json.Unmarshal([]byte(payload), &operation); err != nil {
		return nil, fmt.Errorf("decode operation: %w: %w", dispatcher.ErrInvalidParams, err)
	}

	log.AddCustomLogData(ctx, "action", operation.Action)
	log.AddCustomLogData(ctx, "dispatcher", "doctor")
	log.Info(ctx, "dispatch")

	switch operation.Action {
	case "ping":
		var req operations.PingRequest
		if err := json.Unmarshal(operation.Request, &req); err != nil {
			return nil, dispatcher.ErrInvalidParams
		}

		return d.Ping(ctx, req)
	case "global":
		var req operations.GlobalRequest
		if err := json.Unmarshal(operation.Request, &req); err != nil {
			return nil, dispatcher.ErrInvalidParams
		}

		return d.Global(ctx, req)
	}

	return nil, fmt.Errorf("unsupported action %s: %w", operation.Action, dispatcher.ErrUnknownAction)
}
