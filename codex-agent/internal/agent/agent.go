package agent

import (
	"fmt"
	"os"

	"multienv/codex-agent/internal/client"
)

type Agent struct {
	client   client.Client
	hostname string
	config   Config
}

func New(config Config, client client.Client) (*Agent, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %v", err)
	}

	return &Agent{
		client:   client,
		hostname: hostname,
	}, nil
}
