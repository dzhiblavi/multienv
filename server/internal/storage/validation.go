package storage

import (
	"encoding/json"
	"fmt"

	"multienv/common/storage/model"
	"multienv/server/internal/dispatcher"
)

func validateLocation(at model.Location) error {
	if err := at.Validate(); err != nil {
		return fmt.Errorf("%w: %w", dispatcher.ErrInvalidParams, err)
	}

	return nil
}

func validateWrite(at model.Location, value json.RawMessage) error {
	if err := validateLocation(at); err != nil {
		return err
	}

	if !json.Valid(value) {
		return fmt.Errorf("a JSON value is required: %w", dispatcher.ErrInvalidParams)
	}

	return nil
}
