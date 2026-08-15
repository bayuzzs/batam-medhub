package httpapi

import (
	"errors"
	"net/http"

	"batam-medhub/internal/service"

	"github.com/gin-gonic/gin"
)

func handleRegister(svc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req service.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Malformed request payload.",
			})
			return
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
				status:  http.StatusInternalServerError,
				code:    "INTERNAL_ERROR",
				message: "Registration failed unexpectedly.",
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
		if err := c.ShouldBindJSON(&req); err != nil {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Malformed request payload.",
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
				status:  http.StatusInternalServerError,
				code:    "INTERNAL_ERROR",
				message: "Login failed unexpectedly.",
			})
			return
		}

		c.JSON(http.StatusOK, session)
	}
}

func handleRefresh(svc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req service.RefreshRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Malformed request payload.",
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
				status:  http.StatusInternalServerError,
				code:    "INTERNAL_ERROR",
				message: "Session refresh failed unexpectedly.",
			})
			return
		}

		c.JSON(http.StatusOK, session)
	}
}

func handleLogout(svc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req service.RefreshRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "Malformed request payload.",
			})
			return
		}

		_ = svc.Logout(c.Request.Context(), req.RefreshToken)
		c.Status(http.StatusNoContent)
	}
}
