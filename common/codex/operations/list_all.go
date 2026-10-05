package operations

import (
	"multienv/common/codex/model"
)

type ListAllRequest struct{}

type ListAllResponse struct {
	Sessions model.SessionsByHost `json:"sessions"`
}
