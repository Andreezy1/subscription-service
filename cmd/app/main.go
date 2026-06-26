// Package main Subscription Service API
//
// @title Subscription Service API
// @version 1.0
// @description REST API для управления подписками.
//
// @host localhost:8080
// @BasePath /

package main

import (
	"log/slog"
	"os"
	"subscription-service/config"
	_ "subscription-service/docs"

	"github.com/joho/godotenv"
)

func main() {
	logger := newLogger()
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		logger.Error("load config", slog.Any("error", err))
		os.Exit(1)
	}
	app, err := newApp(cfg, logger)
	if err != nil {
		logger.Error("create application", slog.Any("error", err))
		os.Exit(1)
	}
	run(app, logger)
}
