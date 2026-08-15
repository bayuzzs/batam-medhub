package config

import (
	"errors"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const defaultHTTPAddr = ":8080"

// Config holds application configuration values loaded from the environment.
type Config struct {
	HTTPAddr         string
	DatabaseURL      string
	JWTSigningSecret string
	JWTIssuer        string
	JWTAudience      string
	JWTAccessTTL     time.Duration
	RefreshTokenTTL  time.Duration
}

// Load reads and validates environment configuration for the backend application.
func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:         strings.TrimSpace(os.Getenv("HTTP_ADDR")),
		DatabaseURL:      strings.TrimSpace(os.Getenv("DATABASE_URL")),
		JWTSigningSecret: strings.TrimSpace(os.Getenv("JWT_SIGNING_SECRET")),
		JWTIssuer:        strings.TrimSpace(os.Getenv("JWT_ISSUER")),
		JWTAudience:      strings.TrimSpace(os.Getenv("JWT_AUDIENCE")),
		JWTAccessTTL:     15 * time.Minute,
		RefreshTokenTTL:  30 * 24 * time.Hour,
	}
	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = defaultHTTPAddr
	}
	if cfg.JWTIssuer == "" {
		cfg.JWTIssuer = "batam-medhub"
	}
	if cfg.JWTAudience == "" {
		cfg.JWTAudience = "batam-medhub-mobile"
	}
	if ttlStr := strings.TrimSpace(os.Getenv("JWT_ACCESS_TTL_SECONDS")); ttlStr != "" {
		if sec, err := strconv.Atoi(ttlStr); err == nil && sec >= 60 && sec <= 3600 {
			cfg.JWTAccessTTL = time.Duration(sec) * time.Second
		}
	}
	if ttlStr := strings.TrimSpace(os.Getenv("REFRESH_TOKEN_TTL_SECONDS")); ttlStr != "" {
		if sec, err := strconv.Atoi(ttlStr); err == nil && sec > 0 {
			cfg.RefreshTokenTTL = time.Duration(sec) * time.Second
		}
	}

	if err := validateHTTPAddr(cfg.HTTPAddr); err != nil {
		return Config{}, err
	}
	if err := validateDatabaseURL(cfg.DatabaseURL); err != nil {
		return Config{}, err
	}
	if err := validateJWTSigningSecret(cfg.JWTSigningSecret); err != nil {
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

func validateJWTSigningSecret(secret string) error {
	if secret == "" {
		return errors.New("JWT_SIGNING_SECRET is required")
	}
	if len([]byte(secret)) < 32 {
		return errors.New("JWT_SIGNING_SECRET must be at least 32 bytes")
	}
	return nil
}
