package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"batam-medhub/internal/auth"
	"batam-medhub/internal/config"
	"batam-medhub/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Common domain errors returned by authentication operations.
var (
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrEmailConflict       = errors.New("email already registered")
	ErrInvalidRefreshToken = errors.New("refresh token is invalid or expired")
	ErrPatientNotFound     = errors.New("patient not found")
	ErrValidationError     = errors.New("validation error")
)

// AuthService handles patient authentication, session management, and credential rotation.
type AuthService struct {
	db  *gorm.DB
	cfg config.Config
}

// NewAuthService constructs an AuthService with the provided database and config.
func NewAuthService(db *gorm.DB, cfg config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

// RegisterRequest holds the payload for registering a new patient.
type RegisterRequest struct {
	FullName          string `json:"full_name"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	PreferredCurrency string `json:"preferred_currency"`
}

// LoginRequest holds the payload for authenticating an existing patient.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RefreshRequest holds the payload for rotating or revoking a refresh token.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// PatientProfileResponse represents the public profile view of a patient.
type PatientProfileResponse struct {
	ID                string    `json:"id"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PreferredCurrency string    `json:"preferred_currency"`
	Synthetic         bool      `json:"synthetic"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// AuthSessionResponse represents the token pair and profile returned upon successful auth.
type AuthSessionResponse struct {
	TokenType        string                 `json:"token_type"`
	AccessToken      string                 `json:"access_token"`
	RefreshToken     string                 `json:"refresh_token"`
	ExpiresInSeconds int                    `json:"expires_in_seconds"`
	RefreshExpiresAt time.Time              `json:"refresh_expires_at"`
	Profile          PatientProfileResponse `json:"profile"`
}

// Register validates input, checks email uniqueness, creates a new patient, and issues tokens.
func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*AuthSessionResponse, error) {
	fullName := strings.TrimSpace(req.FullName)
	if len(fullName) < 2 || len(fullName) > 120 {
		return nil, fmt.Errorf("%w: full_name must be between 2 and 120 characters", ErrValidationError)
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if len(email) < 3 || len(email) > 254 {
		return nil, fmt.Errorf("%w: email must be between 3 and 254 characters", ErrValidationError)
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, fmt.Errorf("%w: invalid email address format", ErrValidationError)
	}

	if err := auth.ValidatePassword(req.Password); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrValidationError, err.Error())
	}

	currency := strings.ToUpper(strings.TrimSpace(req.PreferredCurrency))
	if currency == "" {
		currency = "SGD"
	}
	if currency != "SGD" && currency != "IDR" {
		return nil, fmt.Errorf("%w: preferred_currency must be SGD or IDR", ErrValidationError)
	}

	var existing model.Patient
	err := s.db.WithContext(ctx).Where("email_normalized = ?", email).First(&existing).Error
	if err == nil {
		return nil, ErrEmailConflict
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("check email uniqueness: %w", err)
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	patientID := auth.NewUUID()
	sessionID := auth.NewUUID()
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}
	refreshTokenHash := auth.HashToken(refreshToken)

	now := time.Now().UTC()
	refreshExpiresAt := now.Add(s.cfg.RefreshTokenTTL)

	patient := model.Patient{
		ID:                patientID,
		EmailNormalized:   email,
		PasswordHash:      passwordHash,
		FullName:          fullName,
		PreferredCurrency: currency,
		Status:            "ACTIVE",
		Synthetic:         true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	session := model.AuthSession{
		ID:               sessionID,
		PatientID:        patientID,
		RefreshTokenHash: refreshTokenHash,
		ExpiresAt:        refreshExpiresAt,
		CreatedAt:        now,
		LastUsedAt:       &now,
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&patient).Error; err != nil {
			return err
		}
		if err := tx.Create(&session).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("persist patient and session: %w", err)
	}

	accessToken, err := auth.IssueAccessToken(
		s.cfg.JWTSigningSecret,
		s.cfg.JWTIssuer,
		s.cfg.JWTAudience,
		patientID,
		sessionID,
		currency,
		s.cfg.JWTAccessTTL,
	)
	if err != nil {
		return nil, fmt.Errorf("issue access token: %w", err)
	}

	return &AuthSessionResponse{
		TokenType:        "Bearer",
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresInSeconds: int(s.cfg.JWTAccessTTL.Seconds()),
		RefreshExpiresAt: refreshExpiresAt,
		Profile: PatientProfileResponse{
			ID:                patient.ID,
			FullName:          patient.FullName,
			Email:             patient.EmailNormalized,
			PreferredCurrency: patient.PreferredCurrency,
			Synthetic:         patient.Synthetic,
			CreatedAt:         patient.CreatedAt,
			UpdatedAt:         patient.UpdatedAt,
		},
	}, nil
}

// Login authenticates credentials using constant-time verification and creates a new session.
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*AuthSessionResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || req.Password == "" {
		auth.DummyCheckPassword(req.Password)
		return nil, ErrInvalidCredentials
	}

	var patient model.Patient
	err := s.db.WithContext(ctx).Where("email_normalized = ? AND status = 'ACTIVE'", email).First(&patient).Error
	if err != nil {
		auth.DummyCheckPassword(req.Password)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("query patient: %w", err)
	}

	if !auth.CheckPassword(patient.PasswordHash, req.Password) {
		return nil, ErrInvalidCredentials
	}

	sessionID := auth.NewUUID()
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}
	refreshTokenHash := auth.HashToken(refreshToken)

	now := time.Now().UTC()
	refreshExpiresAt := now.Add(s.cfg.RefreshTokenTTL)

	session := model.AuthSession{
		ID:               sessionID,
		PatientID:        patient.ID,
		RefreshTokenHash: refreshTokenHash,
		ExpiresAt:        refreshExpiresAt,
		CreatedAt:        now,
		LastUsedAt:       &now,
	}

	if err := s.db.WithContext(ctx).Create(&session).Error; err != nil {
		return nil, fmt.Errorf("create auth session: %w", err)
	}

	accessToken, err := auth.IssueAccessToken(
		s.cfg.JWTSigningSecret,
		s.cfg.JWTIssuer,
		s.cfg.JWTAudience,
		patient.ID,
		sessionID,
		patient.PreferredCurrency,
		s.cfg.JWTAccessTTL,
	)
	if err != nil {
		return nil, fmt.Errorf("issue access token: %w", err)
	}

	return &AuthSessionResponse{
		TokenType:        "Bearer",
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresInSeconds: int(s.cfg.JWTAccessTTL.Seconds()),
		RefreshExpiresAt: refreshExpiresAt,
		Profile: PatientProfileResponse{
			ID:                patient.ID,
			FullName:          patient.FullName,
			Email:             patient.EmailNormalized,
			PreferredCurrency: patient.PreferredCurrency,
			Synthetic:         patient.Synthetic,
			CreatedAt:         patient.CreatedAt,
			UpdatedAt:         patient.UpdatedAt,
		},
	}, nil
}

// Refresh rotates an active refresh token atomically and returns a replacement token pair.
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*AuthSessionResponse, error) {
	token := strings.TrimSpace(refreshToken)
	if len(token) < 43 || len(token) > 256 {
		return nil, ErrInvalidRefreshToken
	}

	tokenHash := auth.HashToken(token)
	now := time.Now().UTC()

	var newSessionID string
	var newRefreshToken string
	var newRefreshTokenHash string
	var patient model.Patient
	refreshExpiresAt := now.Add(s.cfg.RefreshTokenTTL)

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var session model.AuthSession
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("refresh_token_hash = ?", tokenHash).
			First(&session).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrInvalidRefreshToken
			}
			return err
		}

		if session.RevokedAt != nil || session.ExpiresAt.Before(now) || session.ReplacedBySessionID != nil {
			return ErrInvalidRefreshToken
		}

		if err := tx.Where("id = ? AND status = 'ACTIVE'", session.PatientID).First(&patient).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrInvalidRefreshToken
			}
			return err
		}

		newSessionID = auth.NewUUID()
		var genErr error
		newRefreshToken, genErr = auth.GenerateRefreshToken()
		if genErr != nil {
			return fmt.Errorf("generate new refresh token: %w", genErr)
		}
		newRefreshTokenHash = auth.HashToken(newRefreshToken)

		// Create replacement session first so that newSessionID exists in auth_sessions
		newSession := model.AuthSession{
			ID:               newSessionID,
			PatientID:        patient.ID,
			RefreshTokenHash: newRefreshTokenHash,
			ExpiresAt:        refreshExpiresAt,
			CreatedAt:        now,
			LastUsedAt:       &now,
		}
		if err := tx.Create(&newSession).Error; err != nil {
			return fmt.Errorf("create replacement session: %w", err)
		}

		// Revoke current session and link replacement
		session.RevokedAt = &now
		session.ReplacedBySessionID = &newSessionID
		session.LastUsedAt = &now
		if err := tx.Save(&session).Error; err != nil {
			return fmt.Errorf("update old session: %w", err)
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, ErrInvalidRefreshToken) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, fmt.Errorf("refresh session: %w", err)
	}

	accessToken, err := auth.IssueAccessToken(
		s.cfg.JWTSigningSecret,
		s.cfg.JWTIssuer,
		s.cfg.JWTAudience,
		patient.ID,
		newSessionID,
		patient.PreferredCurrency,
		s.cfg.JWTAccessTTL,
	)
	if err != nil {
		return nil, fmt.Errorf("issue access token: %w", err)
	}

	return &AuthSessionResponse{
		TokenType:        "Bearer",
		AccessToken:      accessToken,
		RefreshToken:     newRefreshToken,
		ExpiresInSeconds: int(s.cfg.JWTAccessTTL.Seconds()),
		RefreshExpiresAt: refreshExpiresAt,
		Profile: PatientProfileResponse{
			ID:                patient.ID,
			FullName:          patient.FullName,
			Email:             patient.EmailNormalized,
			PreferredCurrency: patient.PreferredCurrency,
			Synthetic:         patient.Synthetic,
			CreatedAt:         patient.CreatedAt,
			UpdatedAt:         patient.UpdatedAt,
		},
	}, nil
}

// Logout revokes the session associated with the provided refresh token idempotently.
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	token := strings.TrimSpace(refreshToken)
	if len(token) < 43 || len(token) > 256 {
		return nil // Idempotent: return success without leaking token validity
	}

	tokenHash := auth.HashToken(token)
	now := time.Now().UTC()

	_ = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var session model.AuthSession
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("refresh_token_hash = ?", tokenHash).
			First(&session).Error; err != nil {
			return nil // not found, do nothing
		}

		if session.RevokedAt == nil {
			session.RevokedAt = &now
			session.LastUsedAt = &now
			_ = tx.Save(&session)
		}
		return nil
	})

	return nil
}
