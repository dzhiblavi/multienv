package service

import (
	"context"
	"fmt"

	agentclient "multienv/agent/client/agent"
	"multienv/agent/domains/filesystem/local"
)

type DirectoryCandidatesParams struct {
	Hostname string `json:"hostname"`
}

type DirectoryCandidatesResult struct {
	Directories []string `json:"directories"`
}

func (a Agent) DirectoryCandidates(
	ctx context.Context,
	req DirectoryCandidatesParams,
) (*DirectoryCandidatesResult, error) {
	if req.Hostname == "" {
		return nil, fmt.Errorf("hostname is required")
	}

	if req.Hostname != a.hostname {
		response, err := agentclient.InvokeDecode[DirectoryCandidatesResult](
			ctx, a.hostClient, req.Hostname, "filesystem", "directory-candidates", req,
		)
		if err != nil {
			return nil, err
		}

		return &response, nil
	}

	directories, err := local.DirectoryCandidates(ctx)
	if err != nil {
		return nil, err
	}

	return &DirectoryCandidatesResult{Directories: directories}, nil
}
