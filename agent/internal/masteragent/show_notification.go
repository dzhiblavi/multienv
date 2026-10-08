package agent

import (
	"context"
	"fmt"

	"multienv/common/master/model"
	"multienv/common/master/operations"
)

func (a Agent) ShowNotification(ctx context.Context, notification model.Notification) error {
	_, err := a.client.ShowNotification(ctx, operations.ShowNotificationRequest{
		Notification: notification,
	})
	if err != nil {
		return fmt.Errorf("show notification: %w", err)
	}
	return nil
}
