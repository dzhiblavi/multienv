package filesystem

import (
	"context"
	"encoding/json"
	"fmt"

	"multienv/common/filesystem/operations"
	"multienv/common/subprocess"
	"multienv/server/internal/dispatcher"
)

func (s Dispatcher) DirectoryCandidates(
	ctx context.Context,
	req operations.DirectoryCandidatesRequest,
) (*operations.DirectoryCandidatesResponse, error) {
	if req.Hostname == "" {
		return nil, fmt.Errorf("hostname is required: %w", dispatcher.ErrInvalidParams)
	}

	input, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w: %w", dispatcher.ErrInternal, err)
	}

	result, err := subprocess.RunRemote(ctx, subprocess.RemoteCommand{
		Hostname: req.Hostname,
		Command: subprocess.Command{
			Name:  "~/.bin/multienv-agent",
			Args:  []string{"filesystem", "directory-candidates"},
			Stdin: string(input),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("request directory candidates: %w: %w", dispatcher.ErrInternal, err)
	}
	if result.ExitCode != 0 {
		return nil, fmt.Errorf("agent exited with code %d: %s: %w",
			result.ExitCode, result.Stderr, dispatcher.ErrInternal)
	}

	var response operations.DirectoryCandidatesResponse
	if err := json.Unmarshal([]byte(result.Stdout), &response); err != nil {
		return nil, fmt.Errorf("decode directory candidates: %w: %w", dispatcher.ErrInternal, err)
	}

	return &response, nil
}
