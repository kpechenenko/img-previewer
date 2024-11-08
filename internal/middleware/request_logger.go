package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// RequestLogger логирует обработанные http запросы.
type RequestLogger struct {
	handler http.Handler
}

func (l *RequestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	lw := newLoggingResponseWriter(w)
	l.handler.ServeHTTP(lw, r)

	slog.Info(
		fmt.Sprintf("ip: %s, method: %s, path: %s, proto: %s, statusCode: %d, latency: %s, userAgent: %s",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			r.Proto,
			lw.statusCode,
			time.Since(start),
			r.UserAgent(),
		),
	)
}

func NewRequestLogger(handler http.Handler) *RequestLogger {
	return &RequestLogger{handler: handler}
}
