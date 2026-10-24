package service

import "multienv/agent/master/client"

type Agent struct {
	client client.Client
}

func New(client client.Client) *Agent {
	return &Agent{client: client}
}
