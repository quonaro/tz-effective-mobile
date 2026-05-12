package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"subscriptions/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var direction string
	flag.StringVar(&direction, "direction", "up", "Migration direction: up or down")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg, err := config.Load("")
	if err != nil {
		logger.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}
	if err := cfg.Validate(); err != nil {
		logger.Error("invalid config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	ctx := context.Background()
	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	// Create migrations table if not exists
	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		logger.Error("failed to create migrations table", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Load migration files
	migrationsDir := "migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		// Try to find migrations in parent directories
		for _, dir := range []string{"../migrations", "../../migrations", "./migrations"} {
			if _, err := os.Stat(dir); err == nil {
				migrationsDir = dir
				break
			}
		}
	}

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		logger.Error("failed to read migrations directory", slog.String("error", err.Error()))
		os.Exit(1)
	}

	appliedCount := 0
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		version := strings.TrimSuffix(file.Name(), ".sql")

		// Check if already applied
		var exists bool
		err := db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", version).Scan(&exists)
		if err != nil {
			logger.Error("failed to check migration status", slog.String("version", version), slog.String("error", err.Error()))
			continue
		}

		if direction == "up" && exists {
			logger.Info("migration already applied", slog.String("version", version))
			continue
		}

		if direction == "down" && !exists {
			logger.Info("migration not applied, skipping rollback", slog.String("version", version))
			continue
		}

		// Read and execute migration
		content, err := os.ReadFile(filepath.Join(migrationsDir, file.Name()))
		if err != nil {
			logger.Error("failed to read migration file", slog.String("version", version), slog.String("error", err.Error()))
			continue
		}

		_, err = db.Exec(ctx, string(content))
		if err != nil {
			logger.Error("failed to execute migration", slog.String("version", version), slog.String("error", err.Error()))
			os.Exit(1)
		}

		if direction == "up" {
			_, err = db.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", version)
			if err != nil {
				logger.Error("failed to record migration", slog.String("version", version), slog.String("error", err.Error()))
				os.Exit(1)
			}
			logger.Info("migration applied", slog.String("version", version))
		} else {
			_, err = db.Exec(ctx, "DELETE FROM schema_migrations WHERE version = $1", version)
			if err != nil {
				logger.Error("failed to remove migration record", slog.String("version", version), slog.String("error", err.Error()))
				os.Exit(1)
			}
			logger.Info("migration rolled back", slog.String("version", version))
		}
		appliedCount++
	}

	fmt.Printf("\n%d migrations %s successfully\n", appliedCount, direction)
}
