package http

import (
	"bytes"
	"compress/gzip"
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
	compress bool,
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

	if compress {
		body, err = compressBody(body)
		if err != nil {
			return fmt.Errorf("compress request: %w", err)
		}
	}

	return c.send(ctx, "/"+domain, bytes.NewReader(body), dst, compress)
}

func (c Client) send(
	ctx context.Context,
	path string,
	body io.Reader,
	dst any,
	compress bool,
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
	req.Header.Set("Accept-Encoding", "identity")
	if compress {
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Accept-Encoding", "gzip")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	var responseBody io.Reader = resp.Body
	switch encoding := strings.ToLower(strings.TrimSpace(resp.Header.Get("Content-Encoding"))); encoding {
	case "", "identity":
	case "gzip":
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return fmt.Errorf("decode gzip response: %w", err)
		}
		defer reader.Close()
		responseBody = reader
	default:
		return fmt.Errorf("unsupported response Content-Encoding %q", encoding)
	}

	data, err := io.ReadAll(responseBody)
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

func compressBody(body []byte) ([]byte, error) {
	var compressed bytes.Buffer
	writer := gzip.NewWriter(&compressed)
	if _, err := writer.Write(body); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	return compressed.Bytes(), nil
}
