package platform

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RespondError writes a structured ErrorEnvelope response with the given status code.
func RespondError(c *gin.Context, statusCode int, code string, message string, retryable bool, details []ErrorDetail) {
	if details == nil {
		details = []ErrorDetail{}
	}
	c.JSON(statusCode, ErrorEnvelope{
		Error: ErrorBody{
			Code:      code,
			Message:   message,
			Retryable: retryable,
			RequestID: GetRequestID(c),
			Details:   details,
		},
	})
}

// Common error helpers matching specs/provider-openapi.yaml.

func RespondBadRequest(c *gin.Context, message string, details []ErrorDetail) {
	if message == "" {
		message = "The request failed validation."
	}
	RespondError(c, http.StatusBadRequest, "INVALID_REQUEST", message, false, details)
}

func RespondUnauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Provider integration authentication failed."
	}
	RespondError(c, http.StatusUnauthorized, "AUTHENTICATION_FAILED", message, false, nil)
}

func RespondForbidden(c *gin.Context, message string) {
	if message == "" {
		message = "The asserted provider identity does not match this provider service."
	}
	RespondError(c, http.StatusForbidden, "PROVIDER_IDENTITY_MISMATCH", message, false, nil)
}

func RespondNotFound(c *gin.Context, message string) {
	if message == "" {
		message = "The requested provider resource was not found."
	}
	RespondError(c, http.StatusNotFound, "NOT_FOUND", message, false, nil)
}

func RespondCapacityConflict(c *gin.Context, message string, details []ErrorDetail) {
	if message == "" {
		message = "The requested units are no longer available."
	}
	RespondError(c, http.StatusConflict, "CAPACITY_CONFLICT", message, true, details)
}

func RespondOfferChanged(c *gin.Context, message string, details []ErrorDetail) {
	if message == "" {
		message = "The offer price has changed."
	}
	RespondError(c, http.StatusConflict, "OFFER_CHANGED", message, true, details)
}

func RespondInvalidState(c *gin.Context, message string) {
	if message == "" {
		message = "The resource cannot perform this transition from its current status."
	}
	RespondError(c, http.StatusConflict, "INVALID_STATE", message, false, nil)
}

func RespondIdempotencyConflict(c *gin.Context, message string) {
	if message == "" {
		message = "The idempotency key was already used with a different request."
	}
	RespondError(c, http.StatusConflict, "IDEMPOTENCY_CONFLICT", message, false, nil)
}

func RespondOfferExpired(c *gin.Context, message string, details []ErrorDetail) {
	if message == "" {
		message = "The offer expired before it could be held."
	}
	RespondError(c, http.StatusGone, "OFFER_EXPIRED", message, true, details)
}

func RespondHoldExpired(c *gin.Context, message string, details []ErrorDetail) {
	if message == "" {
		message = "The hold expired before confirmation."
	}
	RespondError(c, http.StatusGone, "HOLD_EXPIRED", message, true, details)
}

func RespondInternalError(c *gin.Context, message string) {
	if message == "" {
		message = "The provider could not complete the request."
	}
	RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", message, true, nil)
}

func RespondServiceUnavailable(c *gin.Context, message string) {
	if message == "" {
		message = "The provider is temporarily unavailable."
	}
	c.Header("Retry-After", "5")
	RespondError(c, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", message, true, nil)
}
