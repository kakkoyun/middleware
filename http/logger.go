package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// LoggerMiddleware TODO
type LoggerMiddleware struct {
	logger log.Logger
}

// NewLogger provides new logging middleware.
func NewLogger(logger log.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{
		logger: log.With(logger, "protocol", "http", "http.component", "server"),
	}
}

// NewHandler TODO
func (m *LoggerMiddleware) NewHandler(handlerName string, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		keyvals := []interface{}{
			"request", r.Header.Get("X-Request-ID"),
			"proto", r.Proto,
			"method", r.Method,
			"status", ww.Status(),
			"content", r.Header.Get("Content-Type"),
			"path", r.URL.Path,
			"duration", time.Since(start),
			"bytes", ww.BytesWritten(),
		}

		if ww.Status()/100 == 5 {
			level.Warn(m.logger).Log(keyvals...)
			return
		}
		level.Debug(m.logger).Log(keyvals...)
	}
}
