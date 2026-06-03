package router

import (
	"net/http"
	"project/handlers"
	"project/middleware"
)

func New(userHandler *handlers.UserHandler) *http.ServeMux{
	mux:=http.NewServeMux()
	mux.Handle("GET /users", middleware.Logging(middleware.Auth(http.HandlerFunc(userHandler.GetAll))))
	mux.Handle("POST /user", middleware.Logging(middleware.Auth(http.HandlerFunc(userHandler.Create))))
	mux.Handle("PUT /user", middleware.Logging(middleware.Auth(http.HandlerFunc(userHandler.Update))))
	mux.Handle("GET /user/{id}", middleware.Logging(middleware.Auth(http.HandlerFunc(userHandler.GetByID))))

	return mux
}

