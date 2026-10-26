package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	redis "github.com/redis/go-redis/v9"

	"multienv/common/storage/model"
)

func (s *redisStorage) Get(ctx context.Context, at model.Location) (json.RawMessage, bool, error) {
	if err := at.Validate(); err != nil {
		return nil, false, err
	}

	return decodeValue(s.client.JSONGet(ctx, at.Key, jsonPath(at.Path)).Result())
}

func decodeValue(value string, err error) (json.RawMessage, bool, error) {
	if errors.Is(err, redis.Nil) || (err == nil && value == "") {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}

	var matches []json.RawMessage
	if err := json.Unmarshal([]byte(value), &matches); err != nil {
		return nil, false, fmt.Errorf("decode Redis JSON value: %w", err)
	}

	if len(matches) == 0 {
		return nil, false, nil
	}

	if len(matches) != 1 {
		return nil, false, fmt.Errorf("storage location matched %d values", len(matches))
	}

	return matches[0], true, nil
}
