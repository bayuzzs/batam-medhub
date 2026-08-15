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
	HTTPAddr                string
	DatabaseURL             string
	JWTSigningSecret        string
	JWTIssuer               string
	JWTAudience             string
	JWTAccessTTL            time.Duration
	RefreshTokenTTL         time.Duration
	HospitalBaseURL         string
	HospitalIntegrationKey  string
	FerryBaseURL            string
	FerryIntegrationKey     string
	HotelBaseURL            string
	HotelIntegrationKey     string
	TransportBaseURL        string
	TransportIntegrationKey string
	ProviderTimeout         time.Duration
}

// Load reads and validates environment configuration for the backend application.
func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:                strings.TrimSpace(os.Getenv("HTTP_ADDR")),
		DatabaseURL:             strings.TrimSpace(os.Getenv("DATABASE_URL")),
		JWTSigningSecret:        strings.TrimSpace(os.Getenv("JWT_SIGNING_SECRET")),
		JWTIssuer:               strings.TrimSpace(os.Getenv("JWT_ISSUER")),
		JWTAudience:             strings.TrimSpace(os.Getenv("JWT_AUDIENCE")),
		JWTAccessTTL:            15 * time.Minute,
		RefreshTokenTTL:         30 * 24 * time.Hour,
		HospitalBaseURL:         strings.TrimSpace(os.Getenv("HOSPITAL_BASE_URL")),
		HospitalIntegrationKey:  strings.TrimSpace(os.Getenv("HOSPITAL_INTEGRATION_KEY")),
		FerryBaseURL:            strings.TrimSpace(os.Getenv("FERRY_BASE_URL")),
		FerryIntegrationKey:     strings.TrimSpace(os.Getenv("FERRY_INTEGRATION_KEY")),
		HotelBaseURL:            strings.TrimSpace(os.Getenv("HOTEL_BASE_URL")),
		HotelIntegrationKey:     strings.TrimSpace(os.Getenv("HOTEL_INTEGRATION_KEY")),
		TransportBaseURL:        strings.TrimSpace(os.Getenv("TRANSPORT_BASE_URL")),
		TransportIntegrationKey: strings.TrimSpace(os.Getenv("TRANSPORT_INTEGRATION_KEY")),
		ProviderTimeout:         5 * time.Second,
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
	if cfg.HospitalBaseURL == "" {
		cfg.HospitalBaseURL = "http://localhost:8081"
	}
	if cfg.HospitalIntegrationKey == "" {
		cfg.HospitalIntegrationKey = "hospital_dev_secret"
	}
	if cfg.FerryBaseURL == "" {
		cfg.FerryBaseURL = "http://localhost:8082"
	}
	if cfg.FerryIntegrationKey == "" {
		cfg.FerryIntegrationKey = "ferry_dev_secret"
	}
	if cfg.HotelBaseURL == "" {
		cfg.HotelBaseURL = "http://localhost:8083"
	}
	if cfg.HotelIntegrationKey == "" {
		cfg.HotelIntegrationKey = "hotel_dev_secret"
	}
	if cfg.TransportBaseURL == "" {
		cfg.TransportBaseURL = "http://localhost:8084"
	}
	if cfg.TransportIntegrationKey == "" {
		cfg.TransportIntegrationKey = "transport_dev_secret"
	}
	if timeoutStr := strings.TrimSpace(os.Getenv("PROVIDER_HTTP_TIMEOUT_SECONDS")); timeoutStr != "" {
		if sec, err := strconv.Atoi(timeoutStr); err == nil && sec > 0 {
			cfg.ProviderTimeout = time.Duration(sec) * time.Second
		}
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
