package middleware

import (
	"net/http"
	"time"

	"project/pkg/logger"

	"go.uber.org/zap"
)

type responseWriterSpy struct {
	http.ResponseWriter
	statusCode int
}

func (spy *responseWriterSpy) WriteHeader(code int) {
	spy.statusCode = code
	spy.ResponseWriter.WriteHeader(code)
}

func (spy *responseWriterSpy) Write(b []byte) (int, error) {
	return spy.ResponseWriter.Write(b)
}

func Logging(loggy *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			spy := &responseWriterSpy{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}
			next.ServeHTTP(spy, r)
			loggy.Info("HTTP Request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", spy.statusCode),
				zap.Duration("latency", time.Since(start)),
			)
		})
	}
}
