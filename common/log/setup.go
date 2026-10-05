package log

import (
	"log/slog"
	"os"
)

func Setup(file *os.File) {
	baseHandler := slog.NewJSONHandler(file, nil)
	logger := slog.New(ContextHandler{Handler: baseHandler})
	slog.SetDefault(logger)
}
