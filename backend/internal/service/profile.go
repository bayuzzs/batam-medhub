package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"batam-medhub/internal/auth"
	"batam-medhub/internal/config"
	"batam-medhub/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ProfileService handles patient profile queries and updates.
type ProfileService struct {
	db  *gorm.DB
	cfg config.Config
}

// NewProfileService constructs a ProfileService with the provided database and config.
func NewProfileService(db *gorm.DB, cfg config.Config) *ProfileService {
	return &ProfileService{db: db, cfg: cfg}
}

// UpdateProfileRequest holds the payload for updating patient name and preferred currency.
type UpdateProfileRequest struct {
	RefreshToken      string  `json:"refresh_token"`
	FullName          *string `json:"full_name,omitempty"`
	PreferredCurrency *string `json:"preferred_currency,omitempty"`
}

// GetProfile retrieves the patient profile for the given patient ID.
func (s *ProfileService) GetProfile(ctx context.Context, patientID string) (*PatientProfileResponse, error) {
	var patient model.Patient
	err := s.db.WithContext(ctx).Where("id = ? AND status = 'ACTIVE'", patientID).First(&patient).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPatientNotFound
		}
		return nil, fmt.Errorf("get profile: %w", err)
	}

	return &PatientProfileResponse{
		ID:                patient.ID,
		FullName:          patient.FullName,
		Email:             patient.EmailNormalized,
		PreferredCurrency: patient.PreferredCurrency,
		Synthetic:         patient.Synthetic,
		CreatedAt:         patient.CreatedAt,
		UpdatedAt:         patient.UpdatedAt,
	}, nil
}

// UpdateProfile validates credentials, updates the patient profile, rotates the active session, and returns new tokens.
func (s *ProfileService) UpdateProfile(ctx context.Context, patientID, currentSessionID string, req UpdateProfileRequest) (*AuthSessionResponse, error) {
	token := strings.TrimSpace(req.RefreshToken)
	if len(token) < 43 || len(token) > 256 {
		return nil, fmt.Errorf("%w: refresh_token is required and must be valid", ErrValidationError)
	}

	if req.FullName == nil && req.PreferredCurrency == nil {
		return nil, fmt.Errorf("%w: either full_name or preferred_currency must be provided", ErrValidationError)
	}

	var newFullName *string
	if req.FullName != nil {
		trimmed := strings.TrimSpace(*req.FullName)
		if len(trimmed) < 2 || len(trimmed) > 120 {
			return nil, fmt.Errorf("%w: full_name must be between 2 and 120 characters", ErrValidationError)
		}
		newFullName = &trimmed
	}

	var newCurrency *string
	if req.PreferredCurrency != nil {
		curr := strings.ToUpper(strings.TrimSpace(*req.PreferredCurrency))
		if curr != "SGD" && curr != "IDR" {
			return nil, fmt.Errorf("%w: preferred_currency must be SGD or IDR", ErrValidationError)
		}
		newCurrency = &curr
	}

	tokenHash := auth.HashToken(token)
	now := time.Now().UTC()
	refreshExpiresAt := now.Add(s.cfg.RefreshTokenTTL)

	var newSessionID string
	var newRefreshToken string
	var newRefreshTokenHash string
	var patient model.Patient

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var session model.AuthSession
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", currentSessionID).
			First(&session).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrInvalidRefreshToken
			}
			return err
		}

		if session.PatientID != patientID || session.RefreshTokenHash != tokenHash ||
			session.RevokedAt != nil || session.ExpiresAt.Before(now) || session.ReplacedBySessionID != nil {
			return ErrInvalidRefreshToken
		}

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND status = 'ACTIVE'", patientID).
			First(&patient).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrPatientNotFound
			}
			return err
		}

		if newFullName != nil {
			patient.FullName = *newFullName
		}
		if newCurrency != nil {
			patient.PreferredCurrency = *newCurrency
		}
		patient.UpdatedAt = now
		if err := tx.Save(&patient).Error; err != nil {
			return fmt.Errorf("save patient: %w", err)
		}

		newSessionID = auth.NewUUID()
		var genErr error
		newRefreshToken, genErr = auth.GenerateRefreshToken()
		if genErr != nil {
			return fmt.Errorf("generate refresh token: %w", genErr)
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
			return fmt.Errorf("revoke current session: %w", err)
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, ErrInvalidRefreshToken) || errors.Is(err, ErrPatientNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("update profile: %w", err)
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
