package adapter

import "fmt"

// ProviderError represents a structured error returned by an external provider service.
type ProviderError struct {
	StatusCode int
	Code       string
	Message    string
	Retryable  bool
	RequestID  string
	Details    []ErrorDetail
}

// Error implements the standard error interface.
func (e *ProviderError) Error() string {
	return fmt.Sprintf("provider error (status %d, code %s): %s", e.StatusCode, e.Code, e.Message)
}

// IsNotFound checks if the error indicates a 404 NOT_FOUND.
func (e *ProviderError) IsNotFound() bool {
	return e.StatusCode == 404 || e.Code == "NOT_FOUND"
}

// IsConflict checks if the error indicates a 409 conflict.
func (e *ProviderError) IsConflict() bool {
	return e.StatusCode == 409
}

// IsExpired checks if the error indicates an expired offer or hold (410 GONE or expired code).
func (e *ProviderError) IsExpired() bool {
	return e.StatusCode == 410 || e.Code == "OFFER_EXPIRED" || e.Code == "HOLD_EXPIRED"
}

// IsUnauthorized checks if the error indicates integration authentication failure (401).
func (e *ProviderError) IsUnauthorized() bool {
	return e.StatusCode == 401 || e.Code == "AUTHENTICATION_FAILED"
}

// IsUnavailable checks if the provider or its datastore is unavailable (503).
func (e *ProviderError) IsUnavailable() bool {
	return e.StatusCode == 503 || e.Code == "SERVICE_UNAVAILABLE"
}
