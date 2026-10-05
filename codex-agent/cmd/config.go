package main

import (
	"multienv/codex-agent/internal/agent"
	"multienv/codex-agent/internal/client/http"
)

type Config struct {
	Agent  agent.Config `yaml:"agent"`
	Client http.Config  `yaml:"client"`
}

func (c *Config) Default() {
	c.Agent.Default()
	c.Client.Default()
}
