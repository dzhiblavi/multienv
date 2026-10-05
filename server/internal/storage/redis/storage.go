package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	redis "github.com/redis/go-redis/v9"

	"multienv/common/util"
	"multienv/server/internal/storage"
)

type redisStorage struct {
	client redis.UniversalClient
	key    string
}

var _ storage.Storage = (*redisStorage)(nil)

func New(ctx context.Context, config Config, key string) (*redisStorage, error) {
	if config.CredentialsPath == "" {
		return nil, fmt.Errorf("CredentialsPath is required")
	}

	err := util.ReadYaml(ctx, config.CredentialsPath, &config.credentials)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Address, config.Port),
		Username: config.credentials.Username,
		Password: config.credentials.Password,
		DB:       config.DB,
	})

	return &redisStorage{client: client, key: key}, nil
}

func (s redisStorage) Close() {
	s.client.Close()
}

func (s *redisStorage) Init(ctx context.Context) error {
	return s.SetIfAbsent(ctx, "$", map[string]any{})
}

// Get decodes the JSON.GET response into dst, which must be a non-nil pointer.
// Paths use RedisJSON syntax: "." reads the document, ".field" reads one value,
// and "$" or "$.field" returns an array of matching values.
func (s *redisStorage) Get(ctx context.Context, path string, dst any) error {
	repr, err := s.client.JSONGet(ctx, s.key, path).Result()
	if err != nil {
		return fmt.Errorf("redis JSONGet failed: %v", err)
	}

	if repr == "" {
		// Not found
		return nil
	}

	err = json.Unmarshal([]byte(repr), dst)
	if err != nil {
		return fmt.Errorf("redis response unmarshal failed: %v", err)
	}

	return nil
}

func (s *redisStorage) Set(ctx context.Context, path string, value any) error {
	repr, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.client.JSONSet(ctx, s.key, path, repr).Err()
}

func (s *redisStorage) SetIfAbsent(ctx context.Context, path string, value any) error {
	repr, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = s.client.JSONSetMode(ctx, s.key, path, repr, "NX").Err()
	if errors.Is(err, redis.Nil) {
		return nil
	}

	return err
}

func (s *redisStorage) Delete(ctx context.Context, path string) error {
	res := s.client.JSONDel(ctx, s.key, path)

	if res.Err() != nil {
		return res.Err()
	}

	return nil
}
