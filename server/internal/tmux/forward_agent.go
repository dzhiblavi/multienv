package tmux

import (
	"context"
	"encoding/json"
	"fmt"

	"multienv/common/subprocess"
	"multienv/common/tmux/operations"
)

func (s Dispatcher) ForwardAgent(
	ctx context.Context,
	req operations.ForwardAgentRequest,
) (*operations.ForwardAgentResponse, error) {
	input, err := json.Marshal(req.Input)
	if err != nil {
		return nil, err
	}

	cmd := subprocess.RemoteCommand{
		Command: subprocess.Command{
			Name:  "~/.bin/tmux-agent",
			Args:  []string{req.Command},
			Stdin: string(input),
		},
		Hostname: req.Hostname,
	}

	result, err := subprocess.RunRemote(ctx, cmd)
	if err != nil {
		return nil, err
	}

	if result.ExitCode != 0 {
		return nil, fmt.Errorf(
			"non-zero exit code: %d, stderr: %s, stdout: %s",
			result.ExitCode,
			result.Stderr,
			result.Stdout,
		)
	}

	return &operations.ForwardAgentResponse{Output: result.Stdout}, nil
}
