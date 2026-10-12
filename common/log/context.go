package log

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"maps"
	"sync"
)

type contextKeyType string
type customLogDataKeyType string

var contextKey contextKeyType
var customLogDataKey customLogDataKeyType

type customLogData struct {
	mu     sync.RWMutex
	values map[string]any
}

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

func InitContext(ctx context.Context) context.Context {
	ctx = initContextId(ctx)
	return context.WithValue(ctx, customLogDataKey, &customLogData{
		values: make(map[string]any),
	})
}

func initContextId(ctx context.Context) context.Context {
	id := generateShortID()
	return context.WithValue(ctx, contextKey, id)
}

func generateShortID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func AddCustomLogData(ctx context.Context, key string, value any) {
	data := getCustomLogData(ctx)
	if data == nil {
		return
	}

	data.mu.Lock()
	defer data.mu.Unlock()
	data.values[key] = value
}

func WithCustomLogData(ctx context.Context, key string, value any) context.Context {
	data := &customLogData{values: make(map[string]any)}
	if parent := getCustomLogData(ctx); parent != nil {
		parent.mu.RLock()
		maps.Copy(data.values, parent.values)
		parent.mu.RUnlock()
	}

	data.values[key] = value
	return context.WithValue(ctx, customLogDataKey, data)
}

func getCustomLogData(ctx context.Context) *customLogData {
	if ctx == nil {
		return nil
	}
	data, _ := ctx.Value(customLogDataKey).(*customLogData)
	return data
}

func appendCustomLogData(ctx context.Context, args []any) []any {
	data := getCustomLogData(ctx)
	if data == nil {
		return args
	}

	data.mu.RLock()
	defer data.mu.RUnlock()
	if len(data.values) == 0 {
		return args
	}

	fields := make([]any, len(args), len(args)+2*len(data.values))
	copy(fields, args)
	for key, value := range data.values {
		fields = append(fields, key, value)
	}
	return fields
}
