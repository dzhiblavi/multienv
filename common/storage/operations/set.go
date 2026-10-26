package operations

import (
	"encoding/json"

	"multienv/common/storage/model"
)

type SetRequest struct {
	model.Location
	Value json.RawMessage `json:"value"`
}

type SetResponse struct{}
