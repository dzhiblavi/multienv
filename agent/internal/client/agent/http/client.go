package http

import (
	"context"

	transport "multienv/agent/internal/client/http"
	"multienv/common/agent/operations"
)

type httpClient struct {
	client transport.Client
}

func New(config transport.Config) httpClient {
	return httpClient{client: transport.New(config)}
}

func (c httpClient) Forward(
	ctx context.Context,
	req operations.ForwardRequest,
) (*operations.ForwardResponse, error) {
	var response operations.ForwardResponse
	if err := c.client.Post(ctx, "agent", "forward", req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
