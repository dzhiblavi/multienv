package client

import (
	"context"
	"multienv/common/codex/operations"
)

type Client interface {
	ListAll(
		ctx context.Context,
		op operations.ListAllRequest,
	) (*operations.ListAllResponse, error)

	SetSessions(
		ctx context.Context,
		op operations.SetSessionsRequest,
	) (*operations.SetSessionsResponse, error)

	SetSession(
		ctx context.Context,
		op operations.SetSessionRequest,
	) (*operations.SetSessionResponse, error)

	DeleteSession(
		ctx context.Context,
		op operations.DeleteSessionRequest,
	) (*operations.DeleteSessionResponse, error)

	SyncAllSessions(
		ctx context.Context,
		op operations.SyncAllSessionsRequest,
	) (*operations.SyncAllSessionsResponse, error)
}
