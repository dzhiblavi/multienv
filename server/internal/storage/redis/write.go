package redis

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	redis "github.com/redis/go-redis/v9"

	"multienv/common/storage/model"
)

//go:embed set.lua
var setSource string

var setScript = redis.NewScript(setSource)

func (s *redisStorage) Set(ctx context.Context, at model.Location, value json.RawMessage) error {
	_, err := s.write(ctx, at, value, "")
	return err
}

func (s *redisStorage) SetIfAbsent(ctx context.Context, at model.Location, value json.RawMessage) (bool, error) {
	return s.write(ctx, at, value, "NX")
}

func (s *redisStorage) write(
	ctx context.Context,
	at model.Location,
	value json.RawMessage,
	mode string,
) (bool, error) {
	if err := at.Validate(); err != nil {
		return false, err
	}

	if !json.Valid(value) {
		return false, fmt.Errorf("invalid JSON value")
	}

	for _, part := range at.Path {
		if strings.ContainsAny(part, "\\\"") {
			return s.writeEscapedPath(ctx, at, value, mode == "NX")
		}
	}

	args := []any{mode, string(value)}
	for _, part := range at.Path {
		args = append(args, quoteKey(part), quotePathKey(part))
	}

	applied, err := setScript.Run(ctx, s.client, []string{at.Key}, args...).Int()
	return applied == 1, err
}
