package subprocess

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func Run(ctx context.Context, command Command) (result RunResult, runErr error) {
	start := time.Now()
	logCommand(ctx, command)
	defer func() {
		logCommandResult(ctx, command, result, runErr, time.Since(start))
	}()

	name, err := expandCommandName(command.Name)
	if err != nil {
		return RunResult{}, fmt.Errorf("expand command path: %w", err)
	}

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	stdinBuf := bytes.NewBufferString(command.Stdin)

	cmd := exec.CommandContext(ctx, name, command.Args...)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	cmd.Stdin = stdinBuf

	err = cmd.Run()

	stdoutStr := stdoutBuf.String()
	stderrStr := stderrBuf.String()
	exitCode := 0

	if err != nil {
		var exitErr *exec.ExitError

		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		} else {
			return RunResult{}, fmt.Errorf("failed to start command: %v", err)
		}
	}

	return RunResult{
		Stdout:   stdoutStr,
		Stderr:   stderrStr,
		ExitCode: exitCode,
	}, nil
}

func RunSuccess(ctx context.Context, command Command) (RunResult, error) {
	result, err := Run(ctx, command)
	if err != nil {
		return RunResult{}, err
	}

	if result.ExitCode != 0 {
		return result, fmt.Errorf(
			"command %s finished with non-zero exit code %d: stderr: %s stdout: %s",
			command.Name,
			result.ExitCode,
			result.Stderr,
			result.Stdout,
		)
	}

	return result, nil
}

func expandCommandName(name string) (string, error) {
	var expandErr error
	name = os.Expand(name, func(variable string) string {
		value, exists := os.LookupEnv(variable)
		if !exists && expandErr == nil {
			expandErr = fmt.Errorf("environment variable %q is not set", variable)
		}
		return value
	})
	if expandErr != nil {
		return "", expandErr
	}

	if name == "~" || strings.HasPrefix(name, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("get home directory: %w", err)
		}
		name = filepath.Join(home, strings.TrimPrefix(name, "~"))
	}

	return name, nil
}
