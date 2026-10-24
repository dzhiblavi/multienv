package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"multienv/common/agent/operations"
	"multienv/common/subprocess"
	"multienv/server/internal/dispatcher"
)

var commandName = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

func (s Dispatcher) Forward(
	ctx context.Context,
	req operations.ForwardRequest,
) (*operations.ForwardResponse, error) {
	if req.Hostname == "" || !commandName.MatchString(req.Domain) || !commandName.MatchString(req.Command) {
		return nil, fmt.Errorf("hostname and valid domain and command are required: %w", dispatcher.ErrInvalidParams)
	}

	input, err := json.Marshal(req.Input)
	if err != nil {
		return nil, err
	}

	cmd := subprocess.RemoteCommand{
		Command: subprocess.Command{
			Name:  "~/.bin/multienv-agent",
			Args:  []string{req.Domain, req.Command},
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

	return &operations.ForwardResponse{Output: result.Stdout}, nil
}
