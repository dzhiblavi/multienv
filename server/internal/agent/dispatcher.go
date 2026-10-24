package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"multienv/common/agent/operations"
	"multienv/common/log"
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

	log.AddCustomLogData(ctx, "action", env.Action)
	log.AddCustomLogData(ctx, "dispatcher", "agent")
	log.Info(ctx, "dispatch")

	switch env.Action {
	case "forward":
		var req operations.ForwardRequest
		if err := json.Unmarshal(env.Request, &req); err != nil {
			return nil, dispatcher.ErrInvalidParams
		}

		return s.Forward(ctx, req)
	}

	return nil, fmt.Errorf("unsupported action %s: %w", env.Action, dispatcher.ErrUnknownAction)
}
