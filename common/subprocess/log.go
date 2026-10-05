package subprocess

import (
	"context"
	"time"

	"multienv/common/log"
)

func logCommand(ctx context.Context, command Command) {
	log.Info(ctx, "starting subprocess", "command", command.Name, "args", command.Args)
}

func logCommandResult(
	ctx context.Context,
	command Command,
	result RunResult,
	err error,
	elapsed time.Duration,
) {
	args := []any{
		"command", command.Name,
		"args", command.Args,
		"exec_sec", elapsed.Seconds(),
	}

	if err != nil {
		args = append(args, "error", err)
		log.Error(ctx, "subprocess failed to start", args...)
		return
	}

	args = append(args, "exit_code", result.ExitCode)
	log.Info(ctx, "subprocess done", args...)
}
