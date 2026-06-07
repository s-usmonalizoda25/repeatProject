package router

import (
	"net/http"
	"project/handlers"
	"project/internal/rate_limiter"
	"project/middleware"
	"project/pkg/logger"
)

func New(userHandler *handlers.UserHandler, rl *rate_limiter.RateLimiter, loggy *logger.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /users", middleware.Logging(loggy)(middleware.Auth(http.HandlerFunc(userHandler.GetAll))))
	mux.Handle("POST /user", middleware.Logging(loggy)(middleware.Auth(http.HandlerFunc(userHandler.Create))))
	mux.Handle("PUT /user", middleware.Logging(loggy)(middleware.Auth(http.HandlerFunc(userHandler.Update))))
	mux.Handle("GET /user/{id}", middleware.Logging(loggy)(middleware.Auth(http.HandlerFunc(userHandler.GetByID))))

	mux.Handle("POST /login", middleware.Logging(loggy)(http.HandlerFunc(userHandler.Login)))
	mux.Handle("GET /users/all", middleware.Logging(loggy)(middleware.Auth(http.HandlerFunc(userHandler.GetAllArchive))))
	mux.Handle("DELETE /users/{id}", middleware.Logging(loggy)(middleware.Auth(http.HandlerFunc(userHandler.SoftDelete))))

	return middleware.RateLimit(rl)(mux)
}
