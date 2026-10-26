package subprocess

import (
	"context"
	"os"
	"strings"

	"multienv/common/log"
)

func RunRemote(ctx context.Context, command RemoteCommand) (RunResult, error) {
	ctx = log.WithCustomLogData(ctx, "host", command.Hostname)

	hostname, err := os.Hostname()
	if err != nil {
		return RunResult{}, err
	}

	if command.Hostname == hostname {
		return Run(ctx, command.Command)
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
			remoteCommandLine(command.Command),
		},
	}

	return Run(ctx, cmd)
}

func remoteCommandLine(command Command) string {
	name := quoteRemoteArgument(command.Name)
	if after, ok := strings.CutPrefix(command.Name, "~/"); ok {
		name = "\"$HOME\"/" + quoteRemoteArgument(after)
	}

	parts := []string{name}
	for _, arg := range command.Args {
		parts = append(parts, quoteRemoteArgument(arg))
	}

	return strings.Join(parts, " ")
}

func quoteRemoteArgument(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}
