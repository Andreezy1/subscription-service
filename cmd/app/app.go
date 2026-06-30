package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"subscription-service/config"
	"subscription-service/internal/handler"
	"subscription-service/internal/repository"
	"subscription-service/internal/service"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Server *http.Server
	DB     *pgxpool.Pool
}

func newApp(cfg *config.Config, logger *slog.Logger) (*App, error) {
	dbCtx, dbCandel := context.WithTimeout(context.Background(), 5*time.Second)
	defer dbCandel()

	db, err := repository.NewPostgresConnect(dbCtx, cfg)
	if err != nil {
		return nil, err
	}

	repo := repository.NewPostgresRepository(db)

	service := service.NewSubscriptionService(repo)

	h := handler.NewHandler(
		service,
		logger,
	)
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	router.Use(handler.NewLoggerMiddleware(logger))
	h.RegisterRoutes(router)

	server := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	return &App{
		Server: server,
		DB:     db,
	}, nil
}

func run(app *App, logger *slog.Logger) {
	go func() {
		logger.Info(
			"server started",
			slog.String("addr", app.Server.Addr),
		)
		if err := app.Server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			logger.Error(
				"http server",
				slog.Any("error", err),
			)
		}
	}()
	waitShutdown(app, logger)
}

func waitShutdown(app *App, logger *slog.Logger) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Server.Shutdown(shutdownCtx); err != nil {
		logger.Error(
			"shutdown server",
			slog.Any("error", err),
		)
	}
	app.DB.Close()
	logger.Info("database connection closed")
	logger.Info("server stopped")
}
