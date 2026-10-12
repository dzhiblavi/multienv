package subprocess

import (
	"context"
	"time"

	"multienv/common/log"
)

func logCommandResult(
	ctx context.Context,
	result RunResult,
	err error,
	elapsed time.Duration,
) {
	args := []any{"exec_sec", elapsed.Seconds()}

	if err != nil {
		args = append(args, "error", err)
		log.Error(ctx, "subprocess failed to start", args...)
		return
	}

	args = append(args, "exit_code", result.ExitCode)
	log.Info(ctx, "subprocess done", args...)
}
