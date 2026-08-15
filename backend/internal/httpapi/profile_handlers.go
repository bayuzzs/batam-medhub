package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"batam-medhub/internal/service"

	"github.com/gin-gonic/gin"
)

func handleGetProfile(svc *service.ProfileService) gin.HandlerFunc {
	return func(c *gin.Context) {
		patientID := c.GetString(contextPatientIDKey)
		profile, err := svc.GetProfile(c.Request.Context(), patientID)
		if err != nil {
			if errors.Is(err, service.ErrPatientNotFound) {
				abort(c, &apiError{
					status:  http.StatusNotFound,
					code:    "NOT_FOUND",
					message: "Patient profile not found.",
				})
				return
			}
			abort(c, &apiError{
				status:    http.StatusInternalServerError,
				code:      "INTERNAL_ERROR",
				message:   "Failed to retrieve profile.",
				retryable: true,
			})
			return
		}

		c.JSON(http.StatusOK, profile)
	}
}

func handleUpdateProfile(svc *service.ProfileService) gin.HandlerFunc {
	return func(c *gin.Context) {
		patientID := c.GetString(contextPatientIDKey)
		sessionID := c.GetString(contextSessionIDKey)

		var req service.UpdateProfileRequest
		if err := bindStrictJSON(c, &req); err != nil {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Malformed request payload.",
			})
			return
		}

		if !isValidRefreshToken(req.RefreshToken) {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "refresh_token must be between 43 and 256 characters.",
			})
			return
		}

		if req.FullName == nil && req.PreferredCurrency == nil {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "At least one field (full_name or preferred_currency) must be provided.",
			})
			return
		}

		if req.FullName != nil {
			fullName := strings.TrimSpace(*req.FullName)
			if len(fullName) < 2 || len(fullName) > 120 {
				abort(c, &apiError{
					status:  http.StatusBadRequest,
					code:    "BAD_REQUEST",
					message: "full_name must be between 2 and 120 characters.",
				})
				return
			}
		}

		if req.PreferredCurrency != nil {
			curr := strings.ToUpper(strings.TrimSpace(*req.PreferredCurrency))
			if curr != "SGD" && curr != "IDR" {
				abort(c, &apiError{
					status:  http.StatusBadRequest,
					code:    "BAD_REQUEST",
					message: "preferred_currency must be SGD or IDR.",
				})
				return
			}
		}

		session, err := svc.UpdateProfile(c.Request.Context(), patientID, sessionID, req)
		if err != nil {
			if errors.Is(err, service.ErrInvalidRefreshToken) {
				abort(c, &apiError{
					status:  http.StatusUnauthorized,
					code:    "INVALID_CREDENTIALS",
					message: "Refresh token is invalid or does not match the active session.",
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
				message:   "Failed to update profile.",
				retryable: true,
			})
			return
		}

		c.JSON(http.StatusOK, session)
	}
}
