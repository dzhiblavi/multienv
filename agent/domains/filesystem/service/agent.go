package service

import (
	"fmt"
	"os"

	hostclient "multienv/agent/domains/host/client"
)

type Agent struct {
	hostClient hostclient.Client
	hostname   string
}

func New(hostClient hostclient.Client) (*Agent, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	return &Agent{hostClient: hostClient, hostname: hostname}, nil
}
