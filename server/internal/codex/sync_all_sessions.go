package codex

import (
	"context"

	"multienv/common/codex/operations"
	hostoperations "multienv/common/host/operations"
	"multienv/common/log"
	"multienv/common/subprocess"
)

func (s Dispatcher) SyncAllSessions(
	ctx context.Context,
	req operations.SyncAllSessionsRequest,
) (*operations.SyncAllSessionsResponse, error) {
	if len(s.config.Hosts) == 0 {
		return &operations.SyncAllSessionsResponse{}, nil
	}

	response, err := s.invoker.Invoke(ctx, hostoperations.InvokeRequest{
		Hostnames: s.config.Hosts,
		Command:   subprocess.Command{Name: "multienv-agent", Args: []string{"codex", "sync-local"}},
	})
	if err != nil {
		return nil, err
	}

	for _, result := range response.Results {
		if err := result.Err(); err != nil {
			log.Error(ctx, "failed to sync host", "host", result.Hostname, "error", err)
		}
	}

	return &operations.SyncAllSessionsResponse{}, nil
}
