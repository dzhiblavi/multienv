package service

import (
	"fmt"
	"os"

	hostclient "multienv/agent/domains/host/client"
	storageclient "multienv/agent/domains/storage/client"
)

type Agent struct {
	storage    storageclient.Client
	hostClient hostclient.Client
	hostname   string
}

func New(storage storageclient.Client, hostClient hostclient.Client) (*Agent, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %v", err)
	}

	return &Agent{
		storage:    storage,
		hostClient: hostClient,
		hostname:   hostname,
	}, nil
}
