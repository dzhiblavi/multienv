package service

import (
	"fmt"
	"os"

	agentclient "multienv/agent/client/agent"
	"multienv/agent/codex/client"
)

type Agent struct {
	client      client.Client
	agentClient agentclient.Client
	hostname    string
}

func New(client client.Client, agentClient agentclient.Client) (*Agent, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %v", err)
	}

	return &Agent{
		client:      client,
		agentClient: agentClient,
		hostname:    hostname,
	}, nil
}
