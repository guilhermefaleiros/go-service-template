package main

import (
	"context"
	"errors"
	"guilhermefaleiros/go-service-template/internal/infrastructure/api"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		os.Setenv("ENVIRONMENT", "local")
	}
	app, err := api.NewAPI(os.Getenv("ENVIRONMENT"))
	if err != nil {
		slog.Error("Failed to initialize API: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := app.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start server: %v", err)
		}
	}()

	<-quit
	slog.Info("Signal received, shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server: %v", err)
	}

	slog.Info("Server gracefully stopped")
}
