package operations

import (
	"multienv/common/tmux/model"
)

type SetSessionsRequest struct {
	Hostname string             `json:"hostname"`
	Sessions model.SessionsById `json:"sessions"`
}

type SetSessionsResponse struct{}
