package client

import (
	"context"
	"encoding/json"
	"fmt"

	"multienv/agent/domains/host/client"
	"multienv/common/host/model"
	"multienv/common/host/operations"
	"multienv/common/subprocess"
)

func Invoke(
	ctx context.Context,
	client client.Client,
	hostname, domain, command string,
	input any,
) (*operations.InvokeResponse, error) {
	stdin, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req := operations.InvokeRequest{
		Hostnames: []string{hostname},
		Command: subprocess.Command{
			Name: "multienv-agent", Args: []string{domain, command}, Stdin: string(stdin),
		},
	}
	response, err := client.Invoke(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("invoke agent: %w", err)
	}

	if len(response.Results) != 1 {
		return nil, fmt.Errorf("invalid agent invocation response for %s", hostname)
	}

	if hostname != model.Master && response.Results[0].Hostname != hostname {
		return nil, fmt.Errorf("invalid agent invocation response for %s", hostname)
	}

	return response, response.Err()
}

func InvokeDecode[T any](
	ctx context.Context,
	client client.Client,
	hostname, domain, command string,
	input any,
) (T, error) {
	var result T
	response, err := Invoke(ctx, client, hostname, domain, command, input)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal([]byte(response.Results[0].Stdout), &result); err != nil {
		return result, fmt.Errorf("decode %s response: %w", command, err)
	}

	return result, nil
}
