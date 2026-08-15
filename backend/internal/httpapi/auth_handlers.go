package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"batam-medhub/internal/service"

	"github.com/gin-gonic/gin"
)

func handleRegister(svc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req service.RegisterRequest
		if err := bindStrictJSON(c, &req); err != nil {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Malformed request payload.",
			})
			return
		}

		fullName := strings.TrimSpace(req.FullName)
		if len(fullName) < 2 || len(fullName) > 120 {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "full_name must be between 2 and 120 characters.",
			})
			return
		}

		if !isValidEmail(req.Email) {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "email must be a valid email address.",
			})
			return
		}

		if !isValidPassword(req.Password) {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "password must be between 8 and 72 UTF-8 bytes.",
			})
			return
		}

		if req.PreferredCurrency != "" {
			curr := strings.ToUpper(strings.TrimSpace(req.PreferredCurrency))
			if curr != "SGD" && curr != "IDR" {
				abort(c, &apiError{
					status:  http.StatusBadRequest,
					code:    "BAD_REQUEST",
					message: "preferred_currency must be SGD or IDR.",
				})
				return
			}
		}

		session, err := svc.Register(c.Request.Context(), req)
		if err != nil {
			if errors.Is(err, service.ErrEmailConflict) {
				abort(c, &apiError{
					status:  http.StatusConflict,
					code:    "CONFLICT",
					message: "A patient with this email already exists.",
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
				message:   "Registration failed unexpectedly.",
				retryable: true,
			})
			return
		}

		c.Header("Location", "/v1/profile")
		c.JSON(http.StatusCreated, session)
	}
}

func handleLogin(svc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req service.LoginRequest
		if err := bindStrictJSON(c, &req); err != nil {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Malformed request payload.",
			})
			return
		}

		if !isValidEmail(req.Email) || !isValidPassword(req.Password) {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Invalid login payload format.",
			})
			return
		}

		session, err := svc.Login(c.Request.Context(), req)
		if err != nil {
			if errors.Is(err, service.ErrInvalidCredentials) {
				abort(c, &apiError{
					status:  http.StatusUnauthorized,
					code:    "INVALID_CREDENTIALS",
					message: "Invalid email or password.",
				})
				return
			}
			abort(c, &apiError{
				status:    http.StatusInternalServerError,
				code:      "INTERNAL_ERROR",
				message:   "Login failed unexpectedly.",
				retryable: true,
			})
			return
		}

		c.JSON(http.StatusOK, session)
	}
}

func handleRefresh(svc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req service.RefreshRequest
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

		session, err := svc.Refresh(c.Request.Context(), req.RefreshToken)
		if err != nil {
			if errors.Is(err, service.ErrInvalidRefreshToken) {
				abort(c, &apiError{
					status:  http.StatusUnauthorized,
					code:    "INVALID_CREDENTIALS",
					message: "Refresh token is invalid or expired.",
				})
				return
			}
			abort(c, &apiError{
				status:    http.StatusInternalServerError,
				code:      "INTERNAL_ERROR",
				message:   "Session refresh failed unexpectedly.",
				retryable: true,
			})
			return
		}

		c.JSON(http.StatusOK, session)
	}
}

func handleLogout(svc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req service.RefreshRequest
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

		if err := svc.Logout(c.Request.Context(), req.RefreshToken); err != nil {
			abort(c, &apiError{
				status:    http.StatusInternalServerError,
				code:      "INTERNAL_ERROR",
				message:   "Logout failed unexpectedly.",
				retryable: true,
			})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
