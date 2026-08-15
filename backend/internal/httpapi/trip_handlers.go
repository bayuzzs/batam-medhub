package httpapi

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"batam-medhub/internal/service"

	"github.com/gin-gonic/gin"
)

const idempotencyTTL = 24 * time.Hour

func handleCreateTripRequest(tripSvc *service.TripService, idemSvc *service.IdempotencyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		patientID := c.GetString(contextPatientIDKey)

		idemKey := c.GetHeader("Idempotency-Key")
		if !service.ValidateIdempotencyKey(idemKey) {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Idempotency-Key header is required and must match pattern ^[A-Za-z0-9][A-Za-z0-9._:-]{7,127}$",
			})
			return
		}

		// Read raw request body for idempotency fingerprint calculation
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Failed to read request body.",
			})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		fingerprint := service.ComputeFingerprint(bodyBytes)
		operation := "POST /v1/trip-requests"

		// Check idempotency cache
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

		var req service.CreateTripRequestInput
		if err := bindStrictJSON(c, &req); err != nil {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Malformed request payload.",
			})
			return
		}

		prompt := strings.TrimSpace(req.Prompt)
		if len(prompt) < 3 || len(prompt) > 2000 {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "prompt must be between 3 and 2000 characters.",
			})
			return
		}

		if req.Locale != "en" && req.Locale != "id" {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "locale must be 'en' or 'id'.",
			})
			return
		}

		// Retrieve patient reference currency
		patientObj, _ := c.Get(contextPatientKey)
		refCurr := "SGD"
		if p, ok := patientObj.(*service.PatientProfileResponse); ok && p != nil {
			refCurr = p.PreferredCurrency
		}

		detail, err := tripSvc.CreateTripRequest(c.Request.Context(), patientID, prompt, req.Locale, refCurr)
		if err != nil {
			if errors.Is(err, service.ErrValidationError) {
				abort(c, &apiError{
					status:  http.StatusUnprocessableEntity,
					code:    "UNPROCESSABLE_ENTITY",
					message: err.Error(),
				})
				return
			}
			abort(c, &apiError{
				status:    http.StatusInternalServerError,
				code:      "INTERNAL_ERROR",
				message:   "Failed to create trip request.",
				retryable: true,
			})
			return
		}

		c.Header("Location", "/v1/trip-requests/"+detail.TripRequest.ID)
		c.Header("Idempotency-Replayed", "false")

		_ = idemSvc.Record(c.Request.Context(), patientID, operation, idemKey, fingerprint, http.StatusCreated, detail, idempotencyTTL)

		c.JSON(http.StatusCreated, detail)
	}
}

func handleGetTripRequest(tripSvc *service.TripService) gin.HandlerFunc {
	return func(c *gin.Context) {
		patientID := c.GetString(contextPatientIDKey)
		tripID := c.Param("trip_request_id")

		detail, err := tripSvc.GetTripRequest(c.Request.Context(), patientID, tripID)
		if err != nil {
			if errors.Is(err, service.ErrTripRequestNotFound) {
				abort(c, &apiError{
					status:  http.StatusNotFound,
					code:    "NOT_FOUND",
					message: "Trip request not found.",
				})
				return
			}
			abort(c, &apiError{
				status:    http.StatusInternalServerError,
				code:      "INTERNAL_ERROR",
				message:   "Failed to retrieve trip request.",
				retryable: true,
			})
			return
		}

		c.JSON(http.StatusOK, detail)
	}
}

func handleAmendTripRequestIntent(tripSvc *service.TripService, idemSvc *service.IdempotencyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		patientID := c.GetString(contextPatientIDKey)
		tripID := c.Param("trip_request_id")

		idemKey := c.GetHeader("Idempotency-Key")
		if !service.ValidateIdempotencyKey(idemKey) {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Idempotency-Key header is required and must match pattern ^[A-Za-z0-9][A-Za-z0-9._:-]{7,127}$",
			})
			return
		}

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Failed to read request body.",
			})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		fingerprint := service.ComputeFingerprint(bodyBytes)
		operation := "PATCH /v1/trip-requests/" + tripID + "/intent"

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

		var req service.AmendIntentRequest
		if err := bindStrictJSON(c, &req); err != nil {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Malformed request payload.",
			})
			return
		}

		if (req.Answer == nil || *req.Answer == "") && req.Corrections == nil {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "At least one of answer or corrections must be provided.",
			})
			return
		}

		detail, err := tripSvc.AmendIntent(c.Request.Context(), patientID, tripID, req)
		if err != nil {
			if errors.Is(err, service.ErrTripRequestNotFound) {
				abort(c, &apiError{
					status:  http.StatusNotFound,
					code:    "NOT_FOUND",
					message: "Trip request not found.",
				})
				return
			}
			if errors.Is(err, service.ErrValidationError) {
				abort(c, &apiError{
					status:  http.StatusUnprocessableEntity,
					code:    "UNPROCESSABLE_ENTITY",
					message: err.Error(),
				})
				return
			}
			abort(c, &apiError{
				status:    http.StatusInternalServerError,
				code:      "INTERNAL_ERROR",
				message:   "Failed to amend trip intent.",
				retryable: true,
			})
			return
		}

		c.Header("Idempotency-Replayed", "false")
		_ = idemSvc.Record(c.Request.Context(), patientID, operation, idemKey, fingerprint, http.StatusOK, detail, idempotencyTTL)

		c.JSON(http.StatusOK, detail)
	}
}

func handleGenerateTripPlans(tripSvc *service.TripService, idemSvc *service.IdempotencyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		patientID := c.GetString(contextPatientIDKey)
		tripID := c.Param("trip_request_id")

		idemKey := c.GetHeader("Idempotency-Key")
		if !service.ValidateIdempotencyKey(idemKey) {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Idempotency-Key header is required and must match pattern ^[A-Za-z0-9][A-Za-z0-9._:-]{7,127}$",
			})
			return
		}

		fingerprint := service.ComputeFingerprint([]byte{})
		operation := "POST /v1/trip-requests/" + tripID + "/plans"

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

		result, err := tripSvc.GeneratePlans(c.Request.Context(), patientID, tripID)
		if err != nil {
			if errors.Is(err, service.ErrTripRequestNotFound) {
				abort(c, &apiError{
					status:  http.StatusNotFound,
					code:    "NOT_FOUND",
					message: "Trip request not found.",
				})
				return
			}
			if errors.Is(err, service.ErrInvalidTripState) {
				abort(c, &apiError{
					status:  http.StatusUnprocessableEntity,
					code:    "UNPROCESSABLE_ENTITY",
					message: err.Error(),
				})
				return
			}
			abort(c, &apiError{
				status:    http.StatusInternalServerError,
				code:      "INTERNAL_ERROR",
				message:   "Failed to generate plans.",
				retryable: true,
			})
			return
		}

		c.Header("Idempotency-Replayed", "false")
		_ = idemSvc.Record(c.Request.Context(), patientID, operation, idemKey, fingerprint, http.StatusOK, result, idempotencyTTL)

		c.JSON(http.StatusOK, result)
	}
}
