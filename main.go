package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"project/handlers"
	"project/internal/config"
	"project/internal/consumer"
	"project/internal/rate_limiter"
	"project/internal/repository"
	"project/internal/service"
	"project/internal/service/eventBus"
	"project/pkg/db"
	"project/pkg/logger"
	"project/router"

	"go.uber.org/zap"
)

func main() {
	loggy, err := logger.New(true)
	if err != nil {
		return
	}
	defer loggy.Sync()

	bus := eventBus.NewBus(100)
	defer bus.Close()

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// const userFileName = "data/users.json"
	const auditFileName = "data/audit.json"

	auditRepo := repository.NewAuditRepo(auditFileName)
	auditService := service.NewAuditService(auditRepo)

	consumer.StartAuditConsumer(ctx, &wg, bus, loggy, auditService)

	rl := rate_limiter.New()
	go rl.WorkerClear(ctx, &wg)

	cfg, err := config.New("config/config.env")
	if err != nil {
		log.Fatal("config.New", err)
	}
	database, err := db.New(db.Options{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	})
	if err != nil {
		loggy.Fatal("db.New", zap.Error(err))
		return
	}
	defer database.Close()

	userRepo := repository.New(database)
	userService := service.New(userRepo)
	userHandler := handlers.New(loggy, userService)

	siteHandler := router.New(userHandler, rl, loggy)

	server := &http.Server{
		Addr:    ":8080",
		Handler: siteHandler,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		loggy.Info("Получен сигнал остановки. Завершаем работу сервера...")
		cancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			loggy.Error("Ошибка при остановке сервера: " + err.Error())
		}
	}()

	loggy.Info("Сервер успешно запущен на порту :8080...")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		loggy.Fatal("Ошибка запуска сервера: " + err.Error())
	}

	wg.Wait()
	loggy.Info("Все фоновые задачи успешно завершены. Программа закрыта.")
}
