package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	redis "github.com/redis/go-redis/v9"

	"multienv/common/storage/model"
)

// JSON.SET in Redis 8.2 does not unescape quoted names when JSON.SET creates or replaces a field.
func (s *redisStorage) writeEscapedPath(
	ctx context.Context,
	at model.Location,
	value json.RawMessage,
	ifAbsent bool,
) (bool, error) {
	for {
		if err := ctx.Err(); err != nil {
			return false, err
		}

		applied := false
		err := s.client.Watch(ctx, func(tx *redis.Tx) error {
			current, found, err := decodeValue(tx.JSONGet(ctx, at.Key, "$").Result())
			if err != nil {
				return err
			}

			updated, changed, err := replaceObjectValue(current, found, at.Path, value, ifAbsent)
			if err != nil || !changed {
				return err
			}

			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.JSONSet(ctx, at.Key, "$", string(updated))
				return nil
			})
			applied = err == nil
			return err
		}, at.Key)
		if !errors.Is(err, redis.TxFailedErr) {
			return applied, err
		}
	}
}

func replaceObjectValue(
	current json.RawMessage,
	found bool,
	path []string,
	value json.RawMessage,
	ifAbsent bool,
) (json.RawMessage, bool, error) {
	if len(path) == 0 {
		if ifAbsent && found {
			return current, false, nil
		}

		return value, true, nil
	}

	object := map[string]json.RawMessage{}
	if found {
		if err := json.Unmarshal(current, &object); err != nil || object == nil {
			return nil, false, fmt.Errorf("storage parent is not an object")
		}
	}

	child, exists := object[path[0]]
	updated, changed, err := replaceObjectValue(child, exists, path[1:], value, ifAbsent)
	if err != nil || !changed {
		return current, false, err
	}

	object[path[0]] = updated
	encoded, err := json.Marshal(object)
	return encoded, true, err
}
