package http

import (
	"context"

	transport "multienv/agent/client/http"
	"multienv/common/storage/operations"
)

type httpClient struct {
	client transport.Client
}

func New(config transport.Config) httpClient {
	return httpClient{client: transport.New(config)}
}

func (c httpClient) Get(
	ctx context.Context,
	req operations.GetRequest,
) (*operations.GetResponse, error) {
	var response operations.GetResponse
	if err := c.client.Post(ctx, "storage", "get", req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c httpClient) Set(
	ctx context.Context,
	req operations.SetRequest,
) (*operations.SetResponse, error) {
	var response operations.SetResponse
	if err := c.client.Post(ctx, "storage", "set", req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c httpClient) SetIfAbsent(
	ctx context.Context,
	req operations.SetIfAbsentRequest,
) (*operations.SetIfAbsentResponse, error) {
	var response operations.SetIfAbsentResponse
	if err := c.client.Post(ctx, "storage", "set-if-absent", req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c httpClient) Delete(
	ctx context.Context,
	req operations.DeleteRequest,
) (*operations.DeleteResponse, error) {
	var response operations.DeleteResponse
	if err := c.client.Post(ctx, "storage", "delete", req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
