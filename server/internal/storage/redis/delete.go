package redis

import (
	"context"

	"multienv/common/storage/model"
)

func (s *redisStorage) Delete(ctx context.Context, at model.Location) (bool, error) {
	if err := at.Validate(); err != nil {
		return false, err
	}

	deleted, err := s.client.JSONDel(ctx, at.Key, jsonPath(at.Path)).Result()
	return deleted > 0, err
}
