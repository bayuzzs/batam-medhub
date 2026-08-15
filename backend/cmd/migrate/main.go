package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"batam-medhub/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if err := run(); err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}
	slog.Info("migrations applied")
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	runner, err := migrate.New("file://migrations", cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("initialize migrations: %w", err)
	}

	upErr := runner.Up()
	sourceErr, databaseErr := runner.Close()
	if upErr != nil && !errors.Is(upErr, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", upErr)
	}
	if sourceErr != nil {
		return fmt.Errorf("close migration source: %w", sourceErr)
	}
	if databaseErr != nil {
		return fmt.Errorf("close migration database: %w", databaseErr)
	}
	return nil
}
