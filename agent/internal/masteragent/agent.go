package agent

import client "multienv/agent/internal/client/master"

type Agent struct {
	client client.Client
}

func New(client client.Client) *Agent {
	return &Agent{client: client}
}
