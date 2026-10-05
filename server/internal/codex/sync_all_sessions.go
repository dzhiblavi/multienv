package codex

import (
	"context"
	"sync"

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
			log.Info(ctx, "syncing host", "dispatcher", "codex", "host", host)

			req := operations.ForwardAgentRequest{
				Hostname: host,
				Command:  "sync-local",
			}

			_, err := s.ForwardAgent(ctx, req)
			if err != nil {
				log.Error(ctx, "failed to sync host", "dispatcher", "codex",
					"host", host, "error", err)
			}

			wg.Done()
		}(host)
	}

	wg.Wait()

	return &operations.SyncAllSessionsResponse{}, nil
}
