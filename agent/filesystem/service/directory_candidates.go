package service

import (
	"context"
	"fmt"

	agentclient "multienv/agent/client/agent"
	"multienv/agent/internal/filesystem/local"
	agentoperations "multienv/common/agent/operations"
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
		response, err := agentclient.Forward[operations.DirectoryCandidatesResponse](
			ctx, a.agentClient, agentoperations.ForwardRequest{
				Hostname: req.Hostname,
				Domain:   "filesystem",
				Command:  "directory-candidates",
				Input:    req,
			})
		if err != nil {
			return nil, err
		}

		return &response, nil
	}

	directories, err := local.DirectoryCandidates(ctx)
	if err != nil {
		return nil, err
	}

	return &operations.DirectoryCandidatesResponse{Directories: directories}, nil
}
