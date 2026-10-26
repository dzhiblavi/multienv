package operations

import "multienv/common/storage/model"

type DeleteRequest struct {
	model.Location
}

type DeleteResponse struct {
	Deleted bool `json:"deleted"`
}
