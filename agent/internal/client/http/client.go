package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	config Config
	client *http.Client
}

func New(config Config) Client {
	return Client{config: config, client: &http.Client{Timeout: config.Timeout}}
}

func (c Client) Request(
	ctx context.Context,
	domain, action string,
	req, dst any,
) error {
	request, err := json.Marshal(req)
	if err != nil {
		return err
	}

	body, err := json.Marshal(struct {
		Action  string          `json:"action"`
		Request json.RawMessage `json:"request"`
	}{
		Action:  action,
		Request: json.RawMessage(request),
	})
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		strings.TrimRight(c.config.URL, "/")+"/"+domain, bytes.NewBuffer(body))

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
