package main

import (
	"log"
	"net/http"
	"project/handlers"
	"project/internal/repository"
	"project/internal/service"
	"project/pkg/logger"
	"project/router"
)
func main() {
	loggy, err := logger.New(true)
	if err != nil {
		log.Fatal("Не удалось создать логгер:", err)
	}

	auditRepo:=repository.NewAuditRepo("data/audit.json")
	auditService:=service.NewAuditService(auditRepo)

	userRepo := repository.New("data/users.json")
	userService := service.New(userRepo, auditService)
	userHandler := handlers.New(loggy, userService)


	mux:=router.New(userHandler)


	loggy.Info("Сервер успешно запущен на порту :8080. Ожидание запросов...")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Критическая ошибка при работе сервера:", err)
	}
}


