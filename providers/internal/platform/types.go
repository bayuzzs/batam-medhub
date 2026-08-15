package platform

import (
	"encoding/json"
	"time"
)

// FormatUTC formats a time.Time as an RFC3339 UTC instant ending in 'Z'.
func FormatUTC(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// Money represents integer minor currency units.
type Money struct {
	AmountMinor int64  `json:"amount_minor"`
	Currency    string `json:"currency"`
}

// TimeWindow represents a scheduled interval with source IANA time zones.
type TimeWindow struct {
	StartsAt      string `json:"starts_at"`
	EndsAt        string `json:"ends_at"`
	StartTimeZone string `json:"start_time_zone"`
	EndTimeZone   string `json:"end_time_zone"`
}

// ReleaseResult represents the outcome of releasing a hold or reservation.
type ReleaseResult struct {
	ResourceType      string `json:"resource_type"`
	ResourceID        string `json:"resource_id"`
	ExternalReference string `json:"external_reference"`
	ProviderID        string `json:"provider_id"`
	ProviderType      string `json:"provider_type"`
	Status            string `json:"status"`
	ReleasedAt        string `json:"released_at"`
}

// HealthResponse represents the health status of a provider service.
type HealthResponse struct {
	Status         string `json:"status"`
	ProviderID     string `json:"provider_id"`
	ProviderType   string `json:"provider_type"`
	DatabaseStatus string `json:"database_status"`
	CheckedAt      string `json:"checked_at"`
}

// ErrorDetail represents a specific field-level validation or domain error detail.
type ErrorDetail struct {
	Field  string `json:"field,omitempty"`
	Reason string `json:"reason"`
}

// ErrorBody represents the standard error structure.
type ErrorBody struct {
	Code      string        `json:"code"`
	Message   string        `json:"message"`
	Retryable bool          `json:"retryable"`
	RequestID string        `json:"request_id"`
	Details   []ErrorDetail `json:"details"`
}

// ErrorEnvelope represents the standard error response wrapper.
type ErrorEnvelope struct {
	Error ErrorBody `json:"error"`
}

// Custom JSON marshalling for ErrorBody to ensure Details is always an array `[]`, never `null`.
func (e ErrorBody) MarshalJSON() ([]byte, error) {
	type Alias ErrorBody
	details := e.Details
	if details == nil {
		details = []ErrorDetail{}
	}
	return json.Marshal(&struct {
		Alias
		Details []ErrorDetail `json:"details"`
	}{
		Alias:   Alias(e),
		Details: details,
	})
}
