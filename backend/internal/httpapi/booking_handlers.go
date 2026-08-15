package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"batam-medhub/internal/service"

	"github.com/gin-gonic/gin"
)

type approvalRequestBody struct {
	Approved *bool `json:"approved"`
}

type selectPlanRequestBody struct {
	PlanOptionID *string `json:"plan_option_id"`
	Approved     *bool   `json:"approved"`
}

func handleConfirmPlanOption(bookingSvc *service.BookingSagaService, idemSvc *service.IdempotencyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		patientID := c.GetString(contextPatientIDKey)
		if patientID == "" {
			abort(c, &apiError{
				status:  http.StatusUnauthorized,
				code:    "UNAUTHORIZED",
				message: "Authentication is required to confirm plan options.",
			})
			return
		}

		planOptionID := c.Param("plan_option_id")
		if planOptionID == "" {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "plan_option_id is required.",
			})
			return
		}

		bodyBytes, err := io.ReadAll(io.LimitReader(c.Request.Body, 1024*1024))
		if err != nil {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Failed to read request body.",
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

		// Check Idempotency
		idemKey := c.GetHeader("Idempotency-Key")
		operation := "POST /v1/plan-options/" + planOptionID + "/confirm"
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

		reqID := getRequestID(c)
		detail, err := bookingSvc.ConfirmPlanOption(c.Request.Context(), patientID, planOptionID, reqID, idemKey)
		if err != nil {
			if errors.Is(err, service.ErrPlanOptionNotFound) {
				abort(c, &apiError{
					status:  http.StatusNotFound,
					code:    "NOT_FOUND",
					message: "Plan option not found.",
				})
				return
			}
			if errors.Is(err, service.ErrInvalidPlanOption) || errors.Is(err, service.ErrInvalidTripState) {
				abort(c, &apiError{
					status:  http.StatusUnprocessableEntity,
					code:    "UNPROCESSABLE_ENTITY",
					message: err.Error(),
				})
				return
			}
			if errors.Is(err, service.ErrBookingHoldFailed) || errors.Is(err, service.ErrBookingConfirmFailed) {
				abort(c, &apiError{
					status:  http.StatusConflict,
					code:    "CONFLICT",
					message: err.Error(),
				})
				return
			}

			slog.ErrorContext(c.Request.Context(), "confirm plan option internal error", "error", err)
			abort(c, &apiError{
				status:  http.StatusInternalServerError,
				code:    "INTERNAL_ERROR",
				message: "An unexpected error occurred while confirming the booking.",
			})
			return
		}

		loc := "/v1/journeys/" + detail.Journey.ID + "/itinerary"
		c.Header("Location", loc)

		if idemKey != "" && idemSvc != nil {
			_ = idemSvc.Record(c.Request.Context(), patientID, operation, idemKey, fingerprint, http.StatusCreated, detail, idempotencyTTL)
		}

		c.JSON(http.StatusCreated, detail)
	}
}

func handleSelectPlanForTrip(bookingSvc *service.BookingSagaService, tripSvc *service.TripService, idemSvc *service.IdempotencyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		patientID := c.GetString(contextPatientIDKey)
		if patientID == "" {
			abort(c, &apiError{
				status:  http.StatusUnauthorized,
				code:    "UNAUTHORIZED",
				message: "Authentication is required.",
			})
			return
		}

		tripID := c.Param("trip_request_id")
		if tripID == "" {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "trip_request_id is required.",
			})
			return
		}

		bodyBytes, err := io.ReadAll(io.LimitReader(c.Request.Body, 1024*1024))
		if err != nil {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Failed to read request body.",
			})
			return
		}

		var req selectPlanRequestBody
		if len(bodyBytes) > 0 {
			decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
			decoder.DisallowUnknownFields()
			_ = decoder.Decode(&req)
		}

		if req.Approved == nil || !*req.Approved {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Explicit approval with approved=true is required.",
			})
			return
		}

		targetOptionID := ""
		if req.PlanOptionID != nil && *req.PlanOptionID != "" {
			targetOptionID = *req.PlanOptionID
		} else {
			// Find active plan option (Rank 1)
			tripDetail, err := tripSvc.GetTripRequest(c.Request.Context(), patientID, tripID)
			if err != nil {
				abort(c, &apiError{
					status:  http.StatusNotFound,
					code:    "NOT_FOUND",
					message: "Trip request not found.",
				})
				return
			}
			if len(tripDetail.PlanOptions) == 0 {
				abort(c, &apiError{
					status:  http.StatusUnprocessableEntity,
					code:    "UNPROCESSABLE_ENTITY",
					message: "No plan options available to select.",
				})
				return
			}
			targetOptionID = tripDetail.PlanOptions[0].ID
		}

		// Check Idempotency
		idemKey := c.GetHeader("Idempotency-Key")
		operation := "POST /v1/trip-requests/" + tripID + "/select-plan"
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

		reqID := getRequestID(c)
		detail, err := bookingSvc.ConfirmPlanOption(c.Request.Context(), patientID, targetOptionID, reqID, idemKey)
		if err != nil {
			if errors.Is(err, service.ErrPlanOptionNotFound) {
				abort(c, &apiError{
					status:  http.StatusNotFound,
					code:    "NOT_FOUND",
					message: "Plan option not found.",
				})
				return
			}
			if errors.Is(err, service.ErrInvalidPlanOption) || errors.Is(err, service.ErrInvalidTripState) {
				abort(c, &apiError{
					status:  http.StatusUnprocessableEntity,
					code:    "UNPROCESSABLE_ENTITY",
					message: err.Error(),
				})
				return
			}
			if errors.Is(err, service.ErrBookingHoldFailed) || errors.Is(err, service.ErrBookingConfirmFailed) {
				abort(c, &apiError{
					status:  http.StatusConflict,
					code:    "CONFLICT",
					message: err.Error(),
				})
				return
			}

			abort(c, &apiError{
				status:  http.StatusInternalServerError,
				code:    "INTERNAL_ERROR",
				message: "An unexpected error occurred while confirming the booking.",
			})
			return
		}

		loc := "/v1/journeys/" + detail.Journey.ID + "/itinerary"
		c.Header("Location", loc)

		if idemKey != "" && idemSvc != nil {
			_ = idemSvc.Record(c.Request.Context(), patientID, operation, idemKey, fingerprint, http.StatusCreated, detail, idempotencyTTL)
		}

		c.JSON(http.StatusCreated, detail)
	}
}
