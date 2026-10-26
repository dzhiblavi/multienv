package service

import (
	"fmt"
	"os"

	"multienv/agent/codex/client"
	hostclient "multienv/agent/host/client"
)

type Agent struct {
	client     client.Client
	hostClient hostclient.Client
	hostname   string
}

func New(client client.Client, hostClient hostclient.Client) (*Agent, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %v", err)
	}

	return &Agent{
		client:     client,
		hostClient: hostClient,
		hostname:   hostname,
	}, nil
}
