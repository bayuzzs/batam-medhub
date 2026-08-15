package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"batam-medhub/internal/config"
	"batam-medhub/internal/database"
	"batam-medhub/internal/httpapi"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		logger.Error("database initialization failed")
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("database pool initialization failed")
		os.Exit(1)
	}
	defer sqlDB.Close()

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           httpapi.New(db, logger),
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Info("API listening", "address", cfg.HTTPAddr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("API stopped", "error", err)
		os.Exit(1)
	}
}
