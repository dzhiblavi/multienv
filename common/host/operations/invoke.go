package operations

import (
	"errors"

	"multienv/common/host/model"
	"multienv/common/subprocess"
)

type InvokeRequest struct {
	subprocess.Command

	Hostnames []string `json:"hostnames"`
	Timeout   string   `json:"timeout,omitempty"`
	Compress  bool     `json:"compress,omitempty"`
}

type InvokeResponse struct {
	Results []model.Invocation `json:"results"`
}

func (r InvokeResponse) Err() error {
	var failures []error
	for _, result := range r.Results {
		failures = append(failures, result.Err())
	}

	return errors.Join(failures...)
}
