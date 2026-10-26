package operations

import (
	"encoding/json"

	"multienv/common/storage/model"
)

type GetRequest struct {
	model.Location
}

type GetResponse struct {
	Found bool            `json:"found"`
	Value json.RawMessage `json:"value"`
}
