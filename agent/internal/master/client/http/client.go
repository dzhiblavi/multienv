package http

import (
	"context"

	transport "multienv/agent/internal/client/http"
	"multienv/common/master/operations"
)

type httpClient struct {
	client transport.Client
}

func New(config transport.Config) httpClient {
	return httpClient{client: transport.New(config)}
}

func (c httpClient) Open(
	ctx context.Context,
	req operations.OpenRequest,
) (*operations.OpenResponse, error) {
	var resp operations.OpenResponse
	if err := c.client.Post(ctx, "master", "open", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c httpClient) ShowNotification(
	ctx context.Context,
	req operations.ShowNotificationRequest,
) (*operations.ShowNotificationResponse, error) {
	var resp operations.ShowNotificationResponse

	if err := c.client.Post(ctx, "master", "show-notification", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
