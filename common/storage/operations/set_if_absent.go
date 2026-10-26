package operations

import (
	"encoding/json"

	"multienv/common/storage/model"
)

type SetIfAbsentRequest struct {
	model.Location
	Value json.RawMessage `json:"value"`
}

type SetIfAbsentResponse struct {
	Created bool `json:"created"`
}
