package httpserver

import (
	"net/http"
	"time"

	"multienv/common/log"
)

const (
	Access = "access"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(writer http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{writer, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

func AccessLogMiddleware() Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			respWriter := newLoggingResponseWriter(w)

			customLogData := make(map[string]any)
			ctx := r.Context()
			ctx = log.InitCustomLogData(ctx, &customLogData)
			ctx = log.InjectContextId(ctx)

			r = r.WithContext(ctx)
			handler(respWriter, r)

			args := []any{
				"path", r.URL.Path,
				"code", respWriter.statusCode,
				"exec_sec", time.Since(start).Seconds(),
			}
			for key, value := range customLogData {
				args = append(args, key, value)
			}

			log.Info(ctx, "request processed", args...)
		}
	}
}
