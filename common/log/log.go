package log

import (
	"context"
	"log/slog"
)

func Info(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, msg, appendCustomLogData(ctx, args)...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, msg, appendCustomLogData(ctx, args)...)
}

func Error(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, appendCustomLogData(ctx, args)...)
}
