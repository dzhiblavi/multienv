package config

import "multienv/agent/client/http"

type Config struct {
	Client http.Config `yaml:"client"`
}

func (c *Config) Default() {
	c.Client.Default()
}
