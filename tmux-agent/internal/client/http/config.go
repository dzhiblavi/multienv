package http

import (
	"time"
)

type Config struct {
	Timeout time.Duration `yaml:"timeout"`
	URL     string        `yaml:"url"`
}

func (c *Config) Default() {
	c.Timeout = 2 * time.Second
	c.URL = "http://127.0.0.1:12123"
}
