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

			ctx := log.InitContext(r.Context())

			r = r.WithContext(ctx)
			handler(respWriter, r)

			log.Info(ctx, "request processed",
				"path", r.URL.Path,
				"code", respWriter.statusCode,
				"exec_sec", time.Since(start).Seconds(),
			)
		}
	}
}
