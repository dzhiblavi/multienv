package codex

import (
	"context"
	"encoding/json"
	"fmt"

	"multienv/common/codex/operations"
	"multienv/common/log"
	"multienv/server/internal/agent"
	"multienv/server/internal/dispatcher"
	"multienv/server/internal/storage"
)

type Dispatcher struct {
	config    Config
	storage   storage.Storage
	forwarder agent.Forwarder
}

func New(config Config, storage storage.Storage, forwarder agent.Forwarder) Dispatcher {
	return Dispatcher{config: config, storage: storage, forwarder: forwarder}
}

func (s Dispatcher) Dispatch(ctx context.Context, payload string) (any, error) {
	var env operations.Operation
	if err := json.Unmarshal([]byte(payload), &env); err != nil {
		return "", err
	}

	log.AddCustomLogData(ctx, "action", env.Action)
	log.AddCustomLogData(ctx, "dispatcher", "codex")
	log.Info(ctx, "dispatch")

	switch env.Action {
	case "list-all":
		var req operations.ListAllRequest
		if err := json.Unmarshal(env.Request, &req); err != nil {
			return nil, dispatcher.ErrInvalidParams
		}
		return s.ListAll(ctx, req)

	case "set-sessions":
		var req operations.SetSessionsRequest
		if err := json.Unmarshal(env.Request, &req); err != nil {
			return nil, dispatcher.ErrInvalidParams
		}
		return s.SetSessions(ctx, req)

	case "set-session":
		var req operations.SetSessionRequest
		if err := json.Unmarshal(env.Request, &req); err != nil {
			return nil, dispatcher.ErrInvalidParams
		}
		return s.SetSession(ctx, req)

	case "delete-session":
		var req operations.DeleteSessionRequest
		if err := json.Unmarshal(env.Request, &req); err != nil {
			return nil, dispatcher.ErrInvalidParams
		}
		return s.DeleteSession(ctx, req)

	case "sync-all":
		var req operations.SyncAllSessionsRequest
		if err := json.Unmarshal(env.Request, &req); err != nil {
			return nil, dispatcher.ErrInvalidParams
		}
		return s.SyncAllSessions(ctx, req)
	}

	return nil, fmt.Errorf(
		"unsupported action %s: %w",
		env.Action,
		dispatcher.ErrUnknownAction,
	)
}
