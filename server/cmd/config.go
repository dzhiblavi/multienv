package main

import (
	"multienv/server/internal/codex"
	"multienv/server/internal/httpserver"
	"multienv/server/internal/storage/redis"
	"multienv/server/internal/tmux"
)

type Config struct {
	Server  httpserver.Config `yaml:"server"`
	Tmux    tmux.Config       `yaml:"tmux"`
	Codex   codex.Config      `yaml:"codex"`
	Storage redis.Config      `yaml:"storage"`
}

func (c *Config) Default() {
	c.Server.Default()
	c.Tmux.Default()
	c.Codex.Default()
	c.Storage.Default()
}
