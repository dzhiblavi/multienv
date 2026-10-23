package http

import (
	"context"
	"fmt"
	"time"

	transport "multienv/agent/internal/client/http"
	"multienv/common/doctor/operations"
)

type httpClient struct {
	client      transport.Client
	pingTimeout time.Duration
}

func New(config transport.Config) httpClient {
	return httpClient{client: transport.New(config), pingTimeout: min(config.Timeout, 5*time.Second)}
}

func (c httpClient) Ping(ctx context.Context, req operations.PingRequest) (*operations.PingResponse, error) {
	var response operations.PingResponse
	if err := c.client.Get(ctx, "/doctor/ping", c.pingTimeout, &response); err != nil {
		return nil, err
	}

	if response.Service != "multienv" || response.Hostname == "" {
		return nil, fmt.Errorf("unexpected Multienv ping response")
	}

	return &response, nil
}

func (c httpClient) Global(ctx context.Context, req operations.GlobalRequest) (*operations.GlobalResponse, error) {
	var response operations.GlobalResponse
	if err := c.client.Get(ctx, "/doctor/global", 65*time.Second, &response); err != nil {
		return nil, err
	}

	if len(response.Report.Hosts) == 0 {
		return nil, fmt.Errorf("empty doctor response")
	}

	return &response, nil
}
