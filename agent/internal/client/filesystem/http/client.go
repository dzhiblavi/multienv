package http

import (
	"context"

	transport "multienv/agent/internal/client/http"
	"multienv/common/filesystem/operations"
)

type httpClient struct {
	client transport.Client
}

func New(config transport.Config) httpClient {
	return httpClient{client: transport.New(config)}
}

func (c httpClient) DirectoryCandidates(
	ctx context.Context,
	req operations.DirectoryCandidatesRequest,
) (*operations.DirectoryCandidatesResponse, error) {
	var response operations.DirectoryCandidatesResponse
	if err := c.client.Post(ctx, "filesystem", "directory-candidates", req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
