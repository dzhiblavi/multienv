package redis

import (
	"context"
	"fmt"
)

func (s *redisStorage) Check(ctx context.Context) error {
	if err := s.client.Ping(ctx).Err(); err != nil {
		return err
	}

	names := []string{"JSON.GET", "JSON.SET", "JSON.DEL", "JSON.TYPE", "EVAL", "EVALSHA"}
	args := []any{"COMMAND", "INFO"}
	for _, name := range names {
		args = append(args, name)
	}

	commands, err := s.client.Do(ctx, args...).Slice()
	if err != nil {
		return err
	}

	if len(commands) != len(names) {
		return fmt.Errorf("invalid Redis command information")
	}

	for i, name := range names {
		if commands[i] == nil {
			return fmt.Errorf("Redis command %s is unavailable", name)
		}
	}

	return nil
}
