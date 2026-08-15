package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"batam-medhub/internal/service"

	"github.com/gin-gonic/gin"
)

func handleProviderDisruption(provAuth *service.ProviderAuthService, disSvc *service.DisruptionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extract Provider Secret from X-Provider-Key or X-Integration-Key header
		providerKey := strings.TrimSpace(c.GetHeader("X-Provider-Key"))
		if providerKey == "" {
			providerKey = strings.TrimSpace(c.GetHeader("X-Integration-Key"))
		}
		if providerKey == "" {
			abort(c, &apiError{
				status:  http.StatusUnauthorized,
				code:    "UNAUTHORIZED",
				message: "Missing provider secret header X-Provider-Key or X-Integration-Key",
			})
			return
		}

		provider, err := provAuth.AuthenticateProvider(c.Request.Context(), providerKey)
		if err != nil {
			abort(c, &apiError{
				status:  http.StatusUnauthorized,
				code:    "UNAUTHORIZED",
				message: "Invalid provider secret",
			})
			return
		}

		// Read raw request body for fingerprinting and deserialization
		rawBody, err := io.ReadAll(io.LimitReader(c.Request.Body, 1024*1024))
		if err != nil {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Failed to read request body",
			})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))

		var req service.ProviderEventRequest
		if err := json.Unmarshal(rawBody, &req); err != nil {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Malformed provider disruption event payload: " + err.Error(),
			})
			return
		}

		receipt, isReplayed, err := disSvc.IngestProviderEvent(c.Request.Context(), provider, rawBody, &req)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrForbiddenEventType):
				abort(c, &apiError{
					status:  http.StatusForbidden,
					code:    "FORBIDDEN",
					message: "Provider event type is incompatible with authenticated provider type",
				})
			case errors.Is(err, service.ErrForbiddenTarget):
				abort(c, &apiError{
					status:  http.StatusForbidden,
					code:    "FORBIDDEN",
					message: "Target itinerary item or reservation does not belong to provider or journey",
				})
			case errors.Is(err, service.ErrEventDuplicateConflict):
				abort(c, &apiError{
					status:  http.StatusConflict,
					code:    "CONFLICT",
					message: "Duplicate external_event_id with conflicting event payload",
				})
			case errors.Is(err, service.ErrNotFound):
				abort(c, &apiError{
					status:  http.StatusNotFound,
					code:    "NOT_FOUND",
					message: err.Error(),
				})
			case errors.Is(err, service.ErrValidationError):
				abort(c, &apiError{
					status:  http.StatusUnprocessableEntity,
					code:    "UNPROCESSABLE_ENTITY",
					message: err.Error(),
				})
			default:
				abort(c, &apiError{
					status:  http.StatusInternalServerError,
					code:    "INTERNAL_ERROR",
					message: "Failed to process provider disruption event: " + err.Error(),
				})
			}
			return
		}

		if isReplayed {
			c.Header("Idempotency-Replayed", "true")
		}
		c.JSON(http.StatusAccepted, receipt)
	}
}

func handleGetDisruption(disSvc *service.DisruptionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		patientID := c.GetString(contextPatientIDKey)
		if patientID == "" {
			abort(c, &apiError{
				status:  http.StatusUnauthorized,
				code:    "UNAUTHORIZED",
				message: "Authentication is required to view disruption details.",
			})
			return
		}
		disruptionID := c.Param("disruption_id")

		detail, err := disSvc.GetDisruption(c.Request.Context(), patientID, disruptionID)
		if err != nil {
			if errors.Is(err, service.ErrDisruptionNotFound) {
				abort(c, &apiError{
					status:  http.StatusNotFound,
					code:    "NOT_FOUND",
					message: "Disruption not found or not owned by patient",
				})
				return
			}
			abort(c, &apiError{
				status:  http.StatusInternalServerError,
				code:    "INTERNAL_ERROR",
				message: "Failed to retrieve disruption: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, detail)
	}
}

func handleApproveRecoveryOption(disSvc *service.DisruptionService, idemSvc *service.IdempotencyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		patientID := c.GetString(contextPatientIDKey)
		if patientID == "" {
			abort(c, &apiError{
				status:  http.StatusUnauthorized,
				code:    "UNAUTHORIZED",
				message: "Authentication is required to approve recovery options.",
			})
			return
		}

		recoveryOptionID := c.Param("recovery_option_id")
		if recoveryOptionID == "" {
			recoveryOptionID = c.Param("option_id")
		}

		bodyBytes, err := io.ReadAll(io.LimitReader(c.Request.Body, 1024*1024))
		if err != nil {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Failed to read request body",
			})
			return
		}

		var req approvalRequestBody
		decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&req); err != nil || req.Approved == nil || !*req.Approved {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Explicit approval with approved=true is required.",
			})
			return
		}

		// Idempotency Check
		idemKey := strings.TrimSpace(c.GetHeader("Idempotency-Key"))
		operation := "POST /v1/recovery-options/" + recoveryOptionID + "/approve"
		var fingerprint string
		if idemKey != "" && idemSvc != nil {
			if !service.ValidateIdempotencyKey(idemKey) {
				abort(c, &apiError{
					status:  http.StatusBadRequest,
					code:    "BAD_REQUEST",
					message: "Idempotency-Key header must match pattern ^[A-Za-z0-9][A-Za-z0-9._:-]{7,127}$",
				})
				return
			}
			fingerprint = service.ComputeFingerprint(bodyBytes)
			cached, err := idemSvc.Check(c.Request.Context(), patientID, operation, idemKey, fingerprint)
			if err != nil {
				if errors.Is(err, service.ErrIdempotencyConflict) {
					abort(c, &apiError{
						status:  http.StatusConflict,
						code:    "CONFLICT",
						message: "Idempotency key reused with different request payload.",
					})
					return
				}
				abort(c, &apiError{
					status:    http.StatusInternalServerError,
					code:      "INTERNAL_ERROR",
					message:   "Idempotency verification failed.",
					retryable: true,
				})
				return
			}
			if cached != nil && cached.Replayed {
				c.Header("Idempotency-Replayed", "true")
				c.Data(cached.StatusCode, "application/json; charset=utf-8", cached.ResponseBody)
				return
			}
		}

		journeyDetail, err := disSvc.ApproveRecoveryOption(c.Request.Context(), patientID, recoveryOptionID, *req.Approved)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrValidationError):
				abort(c, &apiError{
					status:  http.StatusBadRequest,
					code:    "BAD_REQUEST",
					message: err.Error(),
				})
			case errors.Is(err, service.ErrRecoveryOptionNotFound), errors.Is(err, service.ErrDisruptionNotFound):
				abort(c, &apiError{
					status:  http.StatusNotFound,
					code:    "NOT_FOUND",
					message: err.Error(),
				})
			case errors.Is(err, service.ErrRecoveryNotApplicable), errors.Is(err, service.ErrRecoveryOptionExpired):
				abort(c, &apiError{
					status:  http.StatusConflict,
					code:    "CONFLICT",
					message: err.Error(),
				})
			default:
				abort(c, &apiError{
					status:  http.StatusConflict,
					code:    "RECOVERY_FAILED",
					message: "Failed to apply recovery option: " + err.Error(),
				})
			}
			return
		}

		if idemKey != "" && idemSvc != nil {
			_ = idemSvc.Record(c.Request.Context(), patientID, operation, idemKey, fingerprint, http.StatusOK, journeyDetail, 24*time.Hour)
		}

		c.JSON(http.StatusOK, journeyDetail)
	}
}
