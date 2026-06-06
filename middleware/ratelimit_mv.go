package middleware

import (
	"net/http"
	"strconv"
	"project/internal/rate_limiter"
)

func RateLimit(rl *ratelimiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idStr := r.PathValue("id")
			id:=1
			if idStr != "" {
				parsedID, err := strconv.Atoi(idStr)
				if err== nil {
					id=parsedID
				}
			}
			allow, err := rl.Allow(id)
			if err != nil || !allow {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}