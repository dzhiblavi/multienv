package codex

import (
	"context"
	"log"
	"sync"

	"multienv/common/codex/operations"
)

func (s Dispatcher) SyncAllSessions(
	ctx context.Context,
	req operations.SyncAllSessionsRequest,
) (*operations.SyncAllSessionsResponse, error) {
	var wg sync.WaitGroup

	for _, host := range s.config.Hosts {
		wg.Add(1)

		go func(host string) {
			log.Printf("syncing host %s", host)

			req := operations.ForwardAgentRequest{
				Hostname: host,
				Command:  "sync-local",
			}

			_, err := s.ForwardAgent(ctx, req)
			if err != nil {
				log.Printf("failed to sync host %s: %v", host, err)
			}

			wg.Done()
		}(host)
	}

	wg.Wait()

	return &operations.SyncAllSessionsResponse{}, nil
}
