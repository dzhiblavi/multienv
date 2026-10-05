package operations

import (
	"multienv/common/codex/model"
)

type SetSessionsRequest struct {
	Hostname string             `json:"hostname"`
	Sessions model.SessionsById `json:"sessions"`
}

type SetSessionsResponse struct{}
