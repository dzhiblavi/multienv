package host

import (
	"context"
	"encoding/json"
	"fmt"

	"multienv/common/host/operations"
	"multienv/common/log"
	"multienv/server/internal/dispatcher"
)

type Dispatcher struct {
	hosts  []string
	master string
}

func (d Dispatcher) Dispatch(ctx context.Context, payload string) (any, error) {
	var operation operations.Operation
	if err := json.Unmarshal([]byte(payload), &operation); err != nil {
		return nil, fmt.Errorf("decode operation: %w: %w", dispatcher.ErrInvalidParams, err)
	}

	log.AddCustomLogData(ctx, "action", operation.Action)
	log.AddCustomLogData(ctx, "dispatcher", "host")
	log.Info(ctx, "dispatch")

	switch operation.Action {
	case "invoke":
		var req operations.InvokeRequest
		if err := json.Unmarshal(operation.Request, &req); err != nil {
			return nil, fmt.Errorf("decode invoke request: %w: %w", dispatcher.ErrInvalidParams, err)
		}

		return d.Invoke(ctx, req)
	}

	return nil, fmt.Errorf("unsupported action %s: %w", operation.Action, dispatcher.ErrUnknownAction)
}
