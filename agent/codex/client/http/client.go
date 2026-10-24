package http

import (
	"context"

	transport "multienv/agent/client/http"
	"multienv/common/codex/operations"
)

type httpClient struct {
	client transport.Client
}

func New(config transport.Config) httpClient {
	return httpClient{client: transport.New(config)}
}

func (c httpClient) SetSessions(
	ctx context.Context,
	req operations.SetSessionsRequest,
) (*operations.SetSessionsResponse, error) {
	var resp operations.SetSessionsResponse

	if err := c.client.Post(ctx, "codex", "set-sessions", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c httpClient) ListAll(
	ctx context.Context,
	req operations.ListAllRequest,
) (*operations.ListAllResponse, error) {
	var resp operations.ListAllResponse

	if err := c.client.Post(ctx, "codex", "list-all", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c httpClient) SetSession(
	ctx context.Context,
	req operations.SetSessionRequest,
) (*operations.SetSessionResponse, error) {
	var resp operations.SetSessionResponse

	if err := c.client.Post(ctx, "codex", "set-session", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c httpClient) DeleteSession(
	ctx context.Context,
	req operations.DeleteSessionRequest,
) (*operations.DeleteSessionResponse, error) {
	var resp operations.DeleteSessionResponse

	if err := c.client.Post(ctx, "codex", "delete-session", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c httpClient) SyncAllSessions(
	ctx context.Context,
	req operations.SyncAllSessionsRequest,
) (*operations.SyncAllSessionsResponse, error) {
	var resp operations.SyncAllSessionsResponse

	if err := c.client.Post(ctx, "codex", "sync-all", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
