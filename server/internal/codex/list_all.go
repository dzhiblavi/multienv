package codex

import (
	"context"
	"fmt"

	"multienv/common/codex/model"
	"multienv/common/codex/operations"
	"multienv/server/internal/storage"
)

func (s Dispatcher) ListAll(
	ctx context.Context,
	_ operations.ListAllRequest,
) (*operations.ListAllResponse, error) {
	sessions, err := storage.Get[model.SessionsByHost](ctx, s.storage, ".")
	if err != nil {
		return nil, fmt.Errorf("storage request failed: %v", err)
	}

	if sessions == nil {
		sessions = make(model.SessionsByHost)
	}

	return &operations.ListAllResponse{Sessions: sessions}, nil
}
