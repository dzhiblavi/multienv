package service

import (
	"fmt"
	"os"

	agentclient "multienv/agent/internal/client/agent"
)

type Agent struct {
	agentClient agentclient.Client
	hostname    string
}

func New(agentClient agentclient.Client) (*Agent, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	return &Agent{agentClient: agentClient, hostname: hostname}, nil
}
