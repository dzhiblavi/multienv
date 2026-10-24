package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	config Config
	client *http.Client
}

func New(config Config) Client {
	return Client{config: config, client: &http.Client{Timeout: config.Timeout}}
}

func (c Client) Post(
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

	return c.send(ctx, "/"+domain, bytes.NewReader(body), dst)
}

func (c Client) send(
	ctx context.Context,
	path string,
	body io.Reader,
	dst any,
) error {
	endpoint, err := url.Parse(c.config.URL)
	if err != nil || endpoint.Host == "" || (endpoint.Scheme != "http" && endpoint.Scheme != "https") {
		return fmt.Errorf("invalid server URL %q", c.config.URL)
	}

	endpoint.Path = strings.TrimRight(endpoint.Path, "/") + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), body)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
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

	if err := json.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}
