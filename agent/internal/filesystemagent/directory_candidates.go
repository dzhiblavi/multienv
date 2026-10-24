package agent

import (
	"context"
	"fmt"

	"multienv/agent/internal/filesystem"
	"multienv/common/filesystem/operations"
)

func (a Agent) DirectoryCandidates(
	ctx context.Context,
	req operations.DirectoryCandidatesRequest,
) (*operations.DirectoryCandidatesResponse, error) {
	if req.Hostname == "" {
		return nil, fmt.Errorf("hostname is required")
	}

	if req.Hostname != a.hostname {
		return a.client.DirectoryCandidates(ctx, req)
	}

	directories, err := filesystem.DirectoryCandidates(ctx)
	if err != nil {
		return nil, err
	}

	return &operations.DirectoryCandidatesResponse{Directories: directories}, nil
}
