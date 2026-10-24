package client

import (
	"context"

	"multienv/common/master/operations"
)

type Client interface {
	Open(
		context.Context,
		operations.OpenRequest,
	) (*operations.OpenResponse, error)

	ShowNotification(
		context.Context,
		operations.ShowNotificationRequest,
	) (*operations.ShowNotificationResponse, error)
}
