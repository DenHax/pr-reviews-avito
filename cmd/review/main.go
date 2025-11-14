package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/DenHax/pr-reviews-avito/docs"
	"github.com/DenHax/pr-reviews-avito/internal/http/handler"
	"github.com/DenHax/pr-reviews-avito/internal/http/server"
	"github.com/DenHax/pr-reviews-avito/internal/logger/slogger"
	"github.com/DenHax/pr-reviews-avito/internal/repo"
	"github.com/DenHax/pr-reviews-avito/internal/service"
	"github.com/DenHax/pr-reviews-avito/internal/storage/postgres"
	storage "github.com/DenHax/pr-reviews-avito/internal/storage/postgres"
)

// @title API Documentation for the PR Reviews Service for Avito Tech
// @version 1.0
// @description This is the API documentation for the PR Reviews Service
// @host localhost:8080
// @BasePath /
func main() {
	slogger.InitLogging()

	storageConfig, err := postgres.SetupConfig()
	if err != nil {
		slog.Error("fail to get connection url", slog.String("error", err.Error()))
		os.Exit(1)
	}
	storage, err := storage.New(*storageConfig)
	if err != nil {
		slog.Error("failed to init storage", slog.String("error", err.Error()))
		os.Exit(2)
	}
	slog.Debug("db connection", slog.String("connection url", storageConfig.URL))

	repos := repo.NewRepository(storage)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	serverConfig, err := server.SetupConfig()
	if err != nil {
		slog.Error("Failed to setup config", slog.String("error", err.Error()))
		os.Exit(3)
	}
	slog.Info("starting server", slog.String("address", serverConfig.Address))
	srv := server.New(*serverConfig, handlers.Init())

	go func() {
		if err := srv.Run(); err != nil {
			slog.Error("failed to stop server", slog.String("error", err.Error()))
		}
	}()

	slog.Info("server started")

	<-done
	slog.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("failed to stop server", slog.String("error", err.Error()))
		return
	}
}
