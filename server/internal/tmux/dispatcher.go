package tmux

import (
	"context"
	"encoding/json"
	"fmt"

	"multienv/common/log"
	"multienv/common/tmux/operations"
	"multienv/server/internal/dispatcher"
	"multienv/server/internal/storage"
)

type Dispatcher struct {
	config  Config
	storage storage.Storage
}

func New(config Config, storage storage.Storage) Dispatcher {
	return Dispatcher{config, storage}
}

func (s Dispatcher) Dispatch(ctx context.Context, payload string) (any, error) {
	var env operations.Operation
	if err := json.Unmarshal([]byte(payload), &env); err != nil {
		return "", err
	}

	log.Info(ctx, "dispatch", "dispatcher", "tmux", "action", env.Action)
	log.AddCustomLogData(ctx, "action", env.Action)

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

	case "forward-agent":
		var req operations.ForwardAgentRequest
		if err := json.Unmarshal(env.Request, &req); err != nil {
			return nil, dispatcher.ErrInvalidParams
		}
		return s.ForwardAgent(ctx, req)

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
