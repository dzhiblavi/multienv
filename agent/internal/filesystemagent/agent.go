package agent

import (
	"fmt"
	"os"

	client "multienv/agent/internal/client/filesystem"
)

type Agent struct {
	client   client.Client
	hostname string
}

func New(client client.Client) (*Agent, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	return &Agent{client: client, hostname: hostname}, nil
}
