package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

var (
	// ErrInvalidDemoScenario indicates the requested scenario is not supported.
	ErrInvalidDemoScenario = errors.New("unsupported demo scenario: only DEFAULT is supported")
	// ErrDemoResetNotConfirmed indicates confirm was not explicitly set to true.
	ErrDemoResetNotConfirmed = errors.New("demo reset must be explicitly confirmed with confirm=true")
)

// DemoResetRequest defines the payload for POST /v1/demo/reset.
type DemoResetRequest struct {
	Scenario string `json:"scenario"`
	Confirm  bool   `json:"confirm"`
}

// DemoResetResponse represents the outcome of POST /v1/demo/reset.
type DemoResetResponse struct {
	Status   string    `json:"status"`
	Scenario string    `json:"scenario"`
	ResetAt  time.Time `json:"reset_at"`
}

// DemoService provides management operations for synthetic hackathon demo state.
type DemoService struct {
	db *gorm.DB
}

// NewDemoService constructs a new DemoService instance.
func NewDemoService(db *gorm.DB) *DemoService {
	return &DemoService{db: db}
}

// ResetDemoData truncates dynamic orchestration state, synthetic patients, and active sessions
// while preserving static reference catalogs and provider configuration.
func (s *DemoService) ResetDemoData(ctx context.Context, req DemoResetRequest) (*DemoResetResponse, error) {
	if req.Scenario != "DEFAULT" {
		return nil, ErrInvalidDemoScenario
	}
	if !req.Confirm {
		return nil, ErrDemoResetNotConfirmed
	}

	resetAt := time.Now().UTC()

	// Truncate all dynamic tables in PostgreSQL.
	// TRUNCATE CASCADE cleanly deletes all rows without firing row-level immutability triggers.
	truncateSQL := `
		TRUNCATE TABLE
			recovery_items,
			recovery_options,
			disruptions,
			provider_events,
			itinerary_items,
			itinerary_versions,
			reservations,
			journeys,
			plan_items,
			plan_options,
			trip_requests,
			idempotency_records,
			auth_sessions,
			patients
		CASCADE;
	`

	if err := s.db.WithContext(ctx).Exec(truncateSQL).Error; err != nil {
		return nil, fmt.Errorf("truncate dynamic demo data: %w", err)
	}

	return &DemoResetResponse{
		Status:   "RESET",
		Scenario: "DEFAULT",
		ResetAt:  resetAt,
	}, nil
}
