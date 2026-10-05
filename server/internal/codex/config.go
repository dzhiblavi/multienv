package codex

import (
	"os"
)

type Config struct {
	Hosts []string `yaml:"hosts"`
}

func (c *Config) Default() {
	hostname, _ := os.Hostname()
	c.Hosts = []string{hostname}
}
