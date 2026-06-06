package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"project/handlers"
	"project/internal/consumer"
	"project/internal/rate_limiter"
	"project/internal/repository"
	"project/internal/service"
	"project/internal/service/eventBus"
	"project/middleware"
	"project/pkg/logger"
	"sync"
	"syscall"
	"time"
	"project/router"
)

func main() {
	loggy, err := logger.New(true)
	if err != nil {
		log.Fatal("Failed to create logger", err)
	}

	bus := eventBus.NewBus(100)

	ctx, cancel := context.WithCancel(context.Background())


	var wg sync.WaitGroup
	consumer.StartAuditConsumer(ctx, &wg, bus, loggy)


	rl:=ratelimiter.New()
	go rl.WorkerClear()

	const fileName = "data/users.json"
	userRepo := repository.New(fileName)
	userService := service.New(userRepo, bus)
	userHandler := handlers.New(loggy, userService)


	mux := http.NewServeMux()
	mux.Handle("GET /users", middleware.Logging(middleware.Auth(http.HandlerFunc(userHandler.GetAll))))
	mux.Handle("POST /user", middleware.Logging(middleware.Auth(http.HandlerFunc(userHandler.Create))))
	mux.Handle("PUT /user", middleware.Logging(middleware.Auth(http.HandlerFunc(userHandler.Update))))
	mux.Handle("GET /user/{id}", middleware.Logging(middleware.Auth(http.HandlerFunc(userHandler.GetByID))))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	siteHandler := router.New(userHandler, rl)

	server=&http.Server{
		Addr:    ":8080",
		Handler: siteHandler,
	}
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		loggy.Info("Получен сигнал остановки сервера. Завершаем работу...")
		cancel()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		
		if err := server.Shutdown(shutdownCtx); err != nil {
			loggy.Error("Ошибка при остановке сервера: " + err.Error())
		}
	}()

	loggy.Info("Сервер успешно запущен на порту :8080...")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal("Failed to start server: ", err)
	}
	wg.Wait()
	loggy.Info("Все фоновые задачи завершены. Программа успешно закрыта.")
}