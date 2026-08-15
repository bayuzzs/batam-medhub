package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"batam-medhub/internal/auth"
	"batam-medhub/internal/model"

	"gorm.io/gorm"
)

// ErrInvalidProviderSecret is returned when the supplied provider secret does not match any active credential.
var ErrInvalidProviderSecret = errors.New("provider secret is invalid or inactive")

// ProviderAuthService authenticates provider actors using stored secret hashes.
type ProviderAuthService struct {
	db *gorm.DB
}

// NewProviderAuthService constructs a ProviderAuthService with the given database handle.
func NewProviderAuthService(db *gorm.DB) *ProviderAuthService {
	return &ProviderAuthService{db: db}
}

// AuthenticateProvider validates the raw provider secret against stored credentials.
func (s *ProviderAuthService) AuthenticateProvider(ctx context.Context, rawSecret string) (*model.Provider, error) {
	secret := strings.TrimSpace(rawSecret)
	if secret == "" {
		return nil, ErrInvalidProviderSecret
	}

	secretHash := auth.HashToken(secret)
	now := time.Now().UTC()

	var cred model.ProviderCredential
	err := s.db.WithContext(ctx).
		Where("secret_hash = ? AND status = 'ACTIVE' AND (expires_at IS NULL OR expires_at > ?)", secretHash, now).
		First(&cred).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidProviderSecret
		}
		return nil, fmt.Errorf("lookup provider credential: %w", err)
	}

	var provider model.Provider
	if err := s.db.WithContext(ctx).
		Where("id = ? AND status = 'ACTIVE'", cred.ProviderID).
		First(&provider).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidProviderSecret
		}
		return nil, fmt.Errorf("lookup provider: %w", err)
	}

	return &provider, nil
}
