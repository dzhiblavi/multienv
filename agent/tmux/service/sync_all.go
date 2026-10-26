package service

import (
	"context"

	"multienv/common/host/model"
	"multienv/common/host/operations"
	"multienv/common/log"
	"multienv/common/subprocess"
)

func (a Agent) SyncAll(ctx context.Context) error {
	response, err := a.hostClient.Invoke(ctx, operations.InvokeRequest{
		Hostnames: []string{model.All},
		Command:   subprocess.Command{Name: "multienv-agent", Args: []string{"tmux", "sync-local"}},
	})
	if err != nil {
		return err
	}

	for _, result := range response.Results {
		if err := result.Err(); err != nil {
			log.Error(ctx, "failed to sync host", "host", result.Hostname, "error", err)
		}
	}

	return nil
}
