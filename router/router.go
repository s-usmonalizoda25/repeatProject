package router

import (
	"net/http"
	"project/handlers"
	"project/internal/rate_limiter"
	"project/middleware"
)

func New(userHandler *handlers.UserHandler, rl *ratelimiter.RateLimiter) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /users", middleware.Logging(middleware.Auth(http.HandlerFunc(userHandler.GetAll))))
	mux.Handle("POST /user", middleware.Logging(middleware.Auth(http.HandlerFunc(userHandler.Create))))
	mux.Handle("PUT /user", middleware.Logging(middleware.Auth(http.HandlerFunc(userHandler.Update))))
	mux.Handle("GET /user/{id}", middleware.Logging(middleware.Auth(http.HandlerFunc(userHandler.GetByID))))
	return middleware.RateLimit(rl)(mux)
}