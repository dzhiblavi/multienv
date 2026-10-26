package model

import "fmt"

type Location struct {
	Key  string   `json:"key"`
	Path []string `json:"path,omitempty"`
}

func (l Location) Validate() error {
	if l.Key == "" {
		return fmt.Errorf("storage key is required")
	}

	return nil
}
