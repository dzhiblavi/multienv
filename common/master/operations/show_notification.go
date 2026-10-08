package operations

import (
	"multienv/common/master/model"
)

type ShowNotificationRequest struct {
	Notification model.Notification `json:"notification"`
}

type ShowNotificationResponse struct{}
