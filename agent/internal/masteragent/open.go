package agent

import (
	"context"
	"fmt"

	"multienv/common/master/operations"
)

func (a Agent) Open(ctx context.Context, target string) error {
	_, err := a.client.Open(ctx, operations.OpenRequest{Target: target})
	if err != nil {
		return fmt.Errorf("open target: %w", err)
	}
	return nil
}
