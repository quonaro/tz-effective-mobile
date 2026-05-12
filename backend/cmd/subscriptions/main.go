package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"subscriptions/internal/config"
	"subscriptions/internal/handler"
	"subscriptions/internal/middleware"
	"subscriptions/internal/repository"
	"subscriptions/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func main() {
	// Load config first to get log level
	cfg, err := config.Load("")
	if err != nil {
		// Use default logger for this error
		slog.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(cfg.LogLevel),
	}))

	if err := cfg.Validate(); err != nil {
		logger.Error("invalid config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()


	repo := repository.NewSubscriptionRepository(db, logger)
	subService := service.NewSubscriptionService(repo, logger)
	router := handler.NewRouter(subService)

	// Add logging middleware manually since humachi router doesn't auto-wrap
	// Actually humachi already uses chi router, so we can wrap it
	routerWithLogger := middleware.Logger(logger)(router)

	srv := &http.Server{
		Addr:         cfg.HTTPAddress,
		Handler:      routerWithLogger,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("starting server", slog.String("address", cfg.HTTPAddress))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutting down gracefully")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("server shutdown error", slog.String("error", err.Error()))
		}
	case err := <-errCh:
		logger.Error("server error", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("server stopped")
}
