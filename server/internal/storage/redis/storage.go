package redis

import (
	"context"
	"fmt"

	redis "github.com/redis/go-redis/v9"

	"multienv/common/util"
	"multienv/server/internal/storage"
)

type redisStorage struct {
	client redis.UniversalClient
}

var _ storage.Storage = (*redisStorage)(nil)

func New(ctx context.Context, config Config) (*redisStorage, error) {
	if config.CredentialsPath == "" {
		return nil, fmt.Errorf("credentials_path is required")
	}

	if err := util.ReadYaml(ctx, config.CredentialsPath, &config.credentials); err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Address, config.Port),
		Username: config.credentials.Username,
		Password: config.credentials.Password,
		DB:       config.DB,
	})
	store := &redisStorage{client: client}
	if err := store.Check(ctx); err != nil {
		client.Close()
		return nil, err
	}

	return store, nil
}

func (s redisStorage) Close() error {
	return s.client.Close()
}
