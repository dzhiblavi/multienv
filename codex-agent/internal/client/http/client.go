package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"multienv/common/codex/operations"
)

type httpClient struct {
	config Config
	client *http.Client
}

func New(config Config) httpClient {
	return httpClient{
		config: config,
		client: &http.Client{Timeout: config.Timeout},
	}
}

func (c httpClient) SetSessions(
	ctx context.Context,
	req operations.SetSessionsRequest,
) (*operations.SetSessionsResponse, error) {
	var resp operations.SetSessionsResponse

	if err := httpRequest(c, ctx, "set-sessions", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c httpClient) ListAll(
	ctx context.Context,
	req operations.ListAllRequest,
) (*operations.ListAllResponse, error) {
	var resp operations.ListAllResponse

	if err := httpRequest(c, ctx, "list-all", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c httpClient) SetSession(
	ctx context.Context,
	req operations.SetSessionRequest,
) (*operations.SetSessionResponse, error) {
	var resp operations.SetSessionResponse

	if err := httpRequest(c, ctx, "set-session", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c httpClient) DeleteSession(
	ctx context.Context,
	req operations.DeleteSessionRequest,
) (*operations.DeleteSessionResponse, error) {
	var resp operations.DeleteSessionResponse

	if err := httpRequest(c, ctx, "delete-session", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c httpClient) SyncAllSessions(
	ctx context.Context,
	req operations.SyncAllSessionsRequest,
) (*operations.SyncAllSessionsResponse, error) {
	var resp operations.SyncAllSessionsResponse

	if err := httpRequest(c, ctx, "sync-all", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c httpClient) ForwardAgent(
	ctx context.Context,
	req operations.ForwardAgentRequest,
) (*operations.ForwardAgentResponse, error) {
	var resp operations.ForwardAgentResponse

	if err := httpRequest(c, ctx, "forward-agent", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func httpRequest[Req any, Resp any](
	c httpClient,
	ctx context.Context,
	action string,
	req Req,
	dst Resp,
) error {
	request, err := json.Marshal(req)
	if err != nil {
		return err
	}

	body, err := json.Marshal(operations.Operation{
		Action:  action,
		Request: json.RawMessage(request),
	})
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		strings.TrimRight(c.config.URL, "/")+"/codex", bytes.NewBuffer(body))

	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %s: %s", resp.Status, data)
	}

	return json.Unmarshal(data, dst)
}
