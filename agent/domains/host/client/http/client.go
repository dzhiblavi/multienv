package http

import (
	"context"

	transport "multienv/agent/client/http"
	"multienv/common/host/operations"
)

type httpClient struct {
	client transport.Client
}

func New(config transport.Config) httpClient {
	return httpClient{client: transport.New(config)}
}

func (c httpClient) Invoke(ctx context.Context, req operations.InvokeRequest) (*operations.InvokeResponse, error) {
	var response operations.InvokeResponse
	if err := c.client.Post(ctx, "host", "invoke", req, &response, req.Compress); err != nil {
		return nil, err
	}

	return &response, nil
}
