package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"batam-medhub/internal/config"
	"batam-medhub/internal/service"

	"github.com/gin-gonic/gin"
)

// handleDemoReset handles POST /v1/demo/reset.
func handleDemoReset(cfg config.Config, demoSvc *service.DemoService, idemSvc *service.IdempotencyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Authenticate Demo Key / Secret
		demoKey := strings.TrimSpace(c.GetHeader("X-Demo-Key"))
		if demoKey == "" {
			demoKey = strings.TrimSpace(c.GetHeader("X-Demo-Secret"))
		}

		expectedSecret := cfg.DemoSecret
		if expectedSecret == "" {
			expectedSecret = "demo_dev_secret"
		}

		if demoKey == "" || demoKey != expectedSecret {
			abort(c, &apiError{
				status:  http.StatusUnauthorized,
				code:    "UNAUTHORIZED",
				message: "Missing or invalid demo secret key",
			})
			return
		}

		// Read raw request body
		rawBody, err := io.ReadAll(io.LimitReader(c.Request.Body, 64*1024))
		if err != nil {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Failed to read request body",
			})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))

		// 2. Parse and Validate Request Payload
		var req service.DemoResetRequest
		decoder := json.NewDecoder(bytes.NewReader(rawBody))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&req); err != nil {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Malformed demo reset request body: " + err.Error(),
			})
			return
		}

		if req.Scenario != "DEFAULT" || !req.Confirm {
			details := make([]errorDetail, 0)
			if req.Scenario != "DEFAULT" {
				details = append(details, errorDetail{
					Field:  "scenario",
					Reason: "Only scenario 'DEFAULT' is supported",
				})
			}
			if !req.Confirm {
				details = append(details, errorDetail{
					Field:  "confirm",
					Reason: "Must be explicitly set to true",
				})
			}
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Invalid demo reset parameters",
				details: details,
			})
			return
		}

		// 3. Handle Idempotency if key is present
		idemKey := strings.TrimSpace(c.GetHeader("Idempotency-Key"))
		operation := "POST /v1/demo/reset"
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
			fingerprint = service.ComputeFingerprint(rawBody)
			cached, err := idemSvc.Check(c.Request.Context(), "demo-admin", operation, idemKey, fingerprint)
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

		// 4. Perform Data Reset
		res, err := demoSvc.ResetDemoData(c.Request.Context(), req)
		if err != nil {
			if errors.Is(err, service.ErrInvalidDemoScenario) || errors.Is(err, service.ErrDemoResetNotConfirmed) {
				abort(c, &apiError{
					status:  http.StatusBadRequest,
					code:    "BAD_REQUEST",
					message: err.Error(),
				})
				return
			}
			abort(c, &apiError{
				status:    http.StatusInternalServerError,
				code:      "INTERNAL_ERROR",
				message:   "Failed to execute demo reset: " + err.Error(),
				retryable: true,
			})
			return
		}

		respBytes, _ := json.Marshal(res)

		// 5. Store Idempotency Record
		if idemKey != "" && idemSvc != nil {
			_ = idemSvc.Record(c.Request.Context(), "demo-admin", operation, idemKey, fingerprint, http.StatusOK, res, 24*time.Hour)
		}

		c.Data(http.StatusOK, "application/json; charset=utf-8", respBytes)
	}
}
