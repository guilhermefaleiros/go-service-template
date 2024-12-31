package main

import (
	"fmt"
	"guilhermefaleiros/go-service-template/internal/infrastructure/api"
	"log/slog"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	app, err := api.NewAPI("development")
	if err != nil {
		slog.Error(fmt.Sprintf("failed to create api: %v", err))
		panic(err)
		return
	}
	err = app.Start()
	if err != nil {
		slog.Error(fmt.Sprintf("failed to start api: %v", err))
		panic(err)
	}
}
