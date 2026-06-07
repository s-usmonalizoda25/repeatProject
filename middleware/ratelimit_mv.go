package middleware

import (
	"net/http"
	"project/internal/rate_limiter"
	"strings"
)

func RateLimit(rl *rate_limiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if idx := strings.LastIndex(ip, ":"); idx != -1 {
				ip = ip[:idx]
			}
			if ip == "" {
				ip = "unknown"
			}

			allow, err := rl.Allow(ip)
			if err != nil || !allow {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
