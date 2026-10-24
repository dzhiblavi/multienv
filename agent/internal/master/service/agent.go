package service

import "multienv/agent/internal/master/client"

type Agent struct {
	client client.Client
}

func New(client client.Client) *Agent {
	return &Agent{client: client}
}
