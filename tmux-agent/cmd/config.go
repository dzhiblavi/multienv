package main

import (
	"multienv/tmux-agent/internal/agent"
	"multienv/tmux-agent/internal/client/http"
)

type Config struct {
	Agent  agent.Config `yaml:"agent"`
	Client http.Config  `yaml:"client"`
}

func (c *Config) Default() {
	c.Agent.Default()
	c.Client.Default()
}
