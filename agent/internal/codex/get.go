package codex

import (
	"context"
	"fmt"

	"multienv/common/codex/model"
)

func GetSession(ctx context.Context, id string) (model.Session, error) {
	allSessions, err := getAllCodexSessions(ctx)
	if err != nil {
		return model.Session{}, err
	}

	session, ok := allSessions[id]
	if !ok {
		return model.Session{}, fmt.Errorf("session %s not found", id)
	}

	return session, nil
}
