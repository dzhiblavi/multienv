package operations

import "multienv/common/doctor/model"

type GlobalRequest struct{}

type GlobalResponse struct {
	Report model.Report `json:"report"`
}
