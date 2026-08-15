package config

import (
	"errors"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const defaultHTTPAddr = ":8080"

type Config struct {
	HTTPAddr    string
	DatabaseURL string
}

func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:    strings.TrimSpace(os.Getenv("HTTP_ADDR")),
		DatabaseURL: strings.TrimSpace(os.Getenv("DATABASE_URL")),
	}
	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = defaultHTTPAddr
	}

	if err := validateHTTPAddr(cfg.HTTPAddr); err != nil {
		return Config{}, err
	}
	if err := validateDatabaseURL(cfg.DatabaseURL); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func validateHTTPAddr(address string) error {
	_, port, err := net.SplitHostPort(address)
	if err != nil {
		return errors.New("HTTP_ADDR must be a host:port address")
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return errors.New("HTTP_ADDR port must be between 1 and 65535")
	}
	return nil
}

func validateDatabaseURL(databaseURL string) error {
	if databaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}

	parsed, err := url.ParseRequestURI(databaseURL)
	if err != nil || (parsed.Scheme != "postgres" && parsed.Scheme != "postgresql") || strings.Trim(parsed.Path, "/") == "" {
		return errors.New("DATABASE_URL must be a PostgreSQL URL with a database name")
	}
	return nil
}
