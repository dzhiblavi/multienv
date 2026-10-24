package client

import (
	"context"

	"multienv/common/filesystem/operations"
)

type Client interface {
	DirectoryCandidates(
		ctx context.Context,
		req operations.DirectoryCandidatesRequest,
	) (*operations.DirectoryCandidatesResponse, error)
}
