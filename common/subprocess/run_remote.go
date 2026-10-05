package subprocess

import (
	"context"
	"os"
)

func RunRemote(ctx context.Context, command RemoteCommand) (RunResult, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return RunResult{}, err
	}

	if command.Hostname == hostname {
		return RunSuccess(ctx, command.Command)
	}

	cmd := Command{
		Name:  "ssh",
		Stdin: command.Stdin,
		Args: []string{
			"-T", // no terminal allocation
			"-o", "BatchMode=yes",
			"-o", "ConnectTimeout=10",
			"--",
			command.Hostname,
			command.Name,
		},
	}
	cmd.Args = append(cmd.Args, command.Args...)

	return Run(ctx, cmd)
}
