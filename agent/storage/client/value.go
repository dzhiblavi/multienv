package client

import (
	"context"
	"encoding/json"
	"fmt"

	"multienv/common/storage/model"
	"multienv/common/storage/operations"
)

func GetDecode[T any](ctx context.Context, client Client, at model.Location) (T, bool, error) {
	var result T
	response, err := client.Get(ctx, operations.GetRequest{Location: at})
	if err != nil {
		return result, false, err
	}

	if !response.Found {
		return result, false, nil
	}

	if err := json.Unmarshal(response.Value, &result); err != nil {
		return result, true, fmt.Errorf("decode storage value: %w", err)
	}

	return result, true, nil
}

func SetEncode(ctx context.Context, client Client, at model.Location, value any) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = client.Set(ctx, operations.SetRequest{Location: at, Value: raw})
	return err
}
