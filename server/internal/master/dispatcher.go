package master

import (
	"context"
	"encoding/json"
	"fmt"

	"multienv/common/log"
	"multienv/common/master/operations"
	"multienv/server/internal/dispatcher"
)

type Dispatcher struct{}

func New() Dispatcher {
	return Dispatcher{}
}

func (s Dispatcher) Dispatch(ctx context.Context, payload string) (any, error) {
	var env operations.Operation
	if err := json.Unmarshal([]byte(payload), &env); err != nil {
		return nil, fmt.Errorf("decode operation: %w: %w", dispatcher.ErrInvalidParams, err)
	}

	log.Info(ctx, "dispatch", "dispatcher", "master", "action", env.Action)
	log.AddCustomLogData(ctx, "action", env.Action)

	switch env.Action {
	case "open":
		var req operations.OpenRequest
		if err := json.Unmarshal(env.Request, &req); err != nil {
			return nil, dispatcher.ErrInvalidParams
		}
		return s.Open(ctx, req)

	case "show-notification":
		var req operations.ShowNotificationRequest
		if err := json.Unmarshal(env.Request, &req); err != nil {
			return nil, dispatcher.ErrInvalidParams
		}
		return s.ShowNotification(ctx, req)
	}

	return nil, fmt.Errorf("unsupported action %s: %w", env.Action, dispatcher.ErrUnknownAction)
}
