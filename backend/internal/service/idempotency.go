package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"batam-medhub/internal/auth"
	"batam-medhub/internal/model"

	"gorm.io/gorm"
)

var (
	ErrIdempotencyConflict = errors.New("idempotency key reused with different request payload")
	idempotencyKeyRegex    = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{7,127}$`)
)

// IdempotencyResult represents the outcome of an idempotency lookup.
type IdempotencyResult struct {
	Replayed     bool
	StatusCode   int
	ResponseBody []byte
}

// IdempotencyService handles caching and atomic replay of mutating operations.
type IdempotencyService struct {
	db *gorm.DB
}

// NewIdempotencyService constructs an IdempotencyService with the provided database handle.
func NewIdempotencyService(db *gorm.DB) *IdempotencyService {
	return &IdempotencyService{db: db}
}

// ValidateKey validates the syntax and length of an Idempotency-Key.
func ValidateIdempotencyKey(key string) bool {
	if len(key) < 8 || len(key) > 128 {
		return false
	}
	return idempotencyKeyRegex.MatchString(key)
}

// ComputeFingerprint hashes the raw request body bytes with SHA-256.
func ComputeFingerprint(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}

// Check checks whether an identical request has already been processed for this scope.
func (s *IdempotencyService) Check(ctx context.Context, authScope, operation, key, fingerprint string) (*IdempotencyResult, error) {
	var record model.IdempotencyRecord
	err := s.db.WithContext(ctx).
		Where("auth_scope = ? AND operation = ? AND idempotency_key = ?", authScope, operation, key).
		First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No existing record, proceed with operation
		}
		return nil, fmt.Errorf("lookup idempotency record: %w", err)
	}

	if record.ExpiresAt.Before(time.Now().UTC()) {
		return nil, nil // Expired record, treat as new
	}

	if record.RequestFingerprint != fingerprint {
		return nil, ErrIdempotencyConflict
	}

	return &IdempotencyResult{
		Replayed:     true,
		StatusCode:   record.ResponseStatus,
		ResponseBody: record.ResponseBody,
	}, nil
}

// Record stores the response of a successful or failed mutation for idempotent replay.
func (s *IdempotencyService) Record(ctx context.Context, authScope, operation, key, fingerprint string, statusCode int, responseBody any, ttl time.Duration) error {
	now := time.Now().UTC()
	expiresAt := now.Add(ttl)

	bodyBytes, err := json.Marshal(responseBody)
	if err != nil {
		return fmt.Errorf("marshal idempotency response: %w", err)
	}

	record := model.IdempotencyRecord{
		ID:                 auth.NewUUID(),
		AuthScope:          authScope,
		Operation:          operation,
		IdempotencyKey:     key,
		RequestFingerprint: fingerprint,
		ResponseStatus:     statusCode,
		ResponseBody:       json.RawMessage(bodyBytes),
		ExpiresAt:          expiresAt,
		CreatedAt:          now,
	}

	if err := s.db.WithContext(ctx).Create(&record).Error; err != nil {
		return fmt.Errorf("save idempotency record: %w", err)
	}

	return nil
}
