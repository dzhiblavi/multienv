package util

import (
	"encoding/json"
)

func BulkUnmarshal[T any](values []string) ([]T, error) {
	var result []T

	for _, repr := range values {
		var item T

		err := json.Unmarshal([]byte(repr), &item)
		if err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, nil
}
