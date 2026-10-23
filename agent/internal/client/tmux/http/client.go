package http

import (
	"context"

	transport "multienv/agent/internal/client/http"
	"multienv/common/tmux/operations"
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

	if err := c.client.Post(ctx, "tmux", "set-sessions", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c httpClient) ListAll(
	ctx context.Context,
	req operations.ListAllRequest,
) (*operations.ListAllResponse, error) {
	var resp operations.ListAllResponse

	if err := c.client.Post(ctx, "tmux", "list-all", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c httpClient) SetSession(
	ctx context.Context,
	req operations.SetSessionRequest,
) (*operations.SetSessionResponse, error) {
	var resp operations.SetSessionResponse

	if err := c.client.Post(ctx, "tmux", "set-session", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c httpClient) DeleteSession(
	ctx context.Context,
	req operations.DeleteSessionRequest,
) (*operations.DeleteSessionResponse, error) {
	var resp operations.DeleteSessionResponse

	if err := c.client.Post(ctx, "tmux", "delete-session", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c httpClient) SyncAllSessions(
	ctx context.Context,
	req operations.SyncAllSessionsRequest,
) (*operations.SyncAllSessionsResponse, error) {
	var resp operations.SyncAllSessionsResponse

	if err := c.client.Post(ctx, "tmux", "sync-all", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c httpClient) ForwardAgent(
	ctx context.Context,
	req operations.ForwardAgentRequest,
) (*operations.ForwardAgentResponse, error) {
	var resp operations.ForwardAgentResponse

	if err := c.client.Post(ctx, "tmux", "forward-agent", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c httpClient) PanePreview(
	ctx context.Context,
	req operations.PanePreviewRequest,
) (*operations.PanePreviewResponse, error) {
	var response operations.PanePreviewResponse
	if err := c.client.Post(ctx, "tmux", "pane-preview", req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
