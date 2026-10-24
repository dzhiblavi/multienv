package codex

import (
	"context"
	"sync"

	agentoperations "multienv/common/agent/operations"
	"multienv/common/codex/operations"
	"multienv/common/log"
)

func (s Dispatcher) SyncAllSessions(
	ctx context.Context,
	req operations.SyncAllSessionsRequest,
) (*operations.SyncAllSessionsResponse, error) {
	var wg sync.WaitGroup

	for _, host := range s.config.Hosts {
		wg.Add(1)

		go func(host string) {
			ctx := log.WithCustomLogData(ctx, "host", host)
			log.Info(ctx, "syncing host")

			req := agentoperations.ForwardRequest{
				Hostname: host,
				Domain:   "codex",
				Command:  "sync-local",
			}

			_, err := s.forwarder.Forward(ctx, req)
			if err != nil {
				log.Error(ctx, "failed to sync host", "error", err)
			}

			wg.Done()
		}(host)
	}

	wg.Wait()

	return &operations.SyncAllSessionsResponse{}, nil
}
