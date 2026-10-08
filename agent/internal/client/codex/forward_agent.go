package client

import (
	"context"
	"encoding/json"
	"fmt"

	"multienv/common/codex/operations"
)

func ForwardAgent[T any](ctx context.Context, client Client, req operations.ForwardAgentRequest) (T, error) {
	var zero T
	response, err := client.ForwardAgent(ctx, req)
	if err != nil {
		return zero, fmt.Errorf("forward %s: %w", req.Command, err)
	}

	var result T
	if err := json.Unmarshal([]byte(response.Output), &result); err != nil {
		return zero, fmt.Errorf("invalid %s response: %w", req.Command, err)
	}
	return result, nil
}
