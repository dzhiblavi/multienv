package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"multienv/common/log"
)

func main() {
	log.Setup(os.Stderr)
	cmd := newRootCmd()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	ctx = log.InjectContextId(ctx)

	if err := cmd.ExecuteContext(ctx); err != nil {
		log.Error(ctx, "tmux-agent failed", "error", err)
		os.Exit(1)
	}
}
