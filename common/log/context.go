package log

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
)

type contextKeyType string
type customLogDataKeyType string

var contextKey contextKeyType
var customLogDataKey customLogDataKeyType

type ContextHandler struct {
	slog.Handler
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx != nil {
		if id, ok := ctx.Value(contextKey).(string); ok {
			r.AddAttrs(slog.String("context", id))
		}
	}

	return h.Handler.Handle(ctx, r)
}

func InjectContextId(ctx context.Context) context.Context {
	id := generateShortID()
	return context.WithValue(ctx, contextKey, id)
}

func generateShortID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func InitCustomLogData(ctx context.Context, values *map[string]any) context.Context {
	return context.WithValue(ctx, customLogDataKey, values)
}

func AddCustomLogData(ctx context.Context, key string, value any) {
	val := ctx.Value(customLogDataKey)
	if val == nil {
		return
	}

	customLogData := val.(*map[string]any)
	(*customLogData)[key] = value
}
