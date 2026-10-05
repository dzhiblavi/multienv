package operations

import (
	"multienv/common/codex/model"
)

type SetSessionRequest struct {
	Hostname string        `json:"hostname"`
	Session  model.Session `json:"session"`
}

type SetSessionResponse struct{}
