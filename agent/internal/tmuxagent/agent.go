package agent

import (
	"fmt"
	"os"

	client "multienv/agent/internal/client/tmux"
)

type Agent struct {
	client   client.Client
	hostname string
}

func New(client client.Client) (*Agent, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %v", err)
	}

	return &Agent{
		client:   client,
		hostname: hostname,
	}, nil
}
