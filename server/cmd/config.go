package main

import (
	"os"

	"multienv/server/internal/httpserver"
	"multienv/server/internal/storage/redis"
)

type Config struct {
	Server  httpserver.Config `yaml:"server"`
	Hosts   []string          `yaml:"hosts"`
	Storage redis.Config      `yaml:"storage"`
}

func (c *Config) Default() {
	c.Server.Default()
	c.Storage.Default()

	hostname, _ := os.Hostname()
	c.Hosts = []string{hostname}
}
