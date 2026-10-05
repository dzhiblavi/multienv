package operations

import (
	"multienv/common/tmux/model"
)

type ListAllRequest struct{}

type ListAllResponse struct {
	Sessions model.SessionsByHost `json:"sessions"`
}
