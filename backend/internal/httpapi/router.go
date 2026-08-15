package httpapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	requestIDHeader = "X-Request-ID"
	requestIDKey    = "request_id"
)

var requestIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{7,127}$`)

type healthResponse struct {
	Status         string    `json:"status"`
	DatabaseStatus string    `json:"database_status"`
	CheckedAt      time.Time `json:"checked_at"`
}

type errorDetail struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

type errorBody struct {
	Code      string        `json:"code"`
	Message   string        `json:"message"`
	Retryable bool          `json:"retryable"`
	RequestID string        `json:"request_id"`
	Details   []errorDetail `json:"details"`
}

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type apiError struct {
	status    int
	code      string
	message   string
	retryable bool
	details   []errorDetail
}

func (e *apiError) Error() string { return e.code }

func New(db *gorm.DB, logger *slog.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(requestID(), requestLogger(logger), structuredErrors(), recoverPanics(logger))
	router.GET("/healthz", liveness)
	router.GET("/readyz", readiness(db))
	return router
}

func liveness(c *gin.Context) {
	c.JSON(http.StatusOK, healthResponse{
		Status:         "UP",
		DatabaseStatus: "NOT_CHECKED",
		CheckedAt:      time.Now().UTC(),
	})
}

func readiness(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		sqlDB, err := db.DB()
		if err == nil {
			err = sqlDB.PingContext(ctx)
		}
		if err != nil {
			abort(c, &apiError{
				status:    http.StatusServiceUnavailable,
				code:      "SERVICE_UNAVAILABLE",
				message:   "Core database is unavailable.",
				retryable: true,
			})
			return
		}

		c.JSON(http.StatusOK, healthResponse{
			Status:         "UP",
			DatabaseStatus: "UP",
			CheckedAt:      time.Now().UTC(),
		})
	}
}

func requestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(requestIDHeader)
		if !requestIDPattern.MatchString(id) {
			id = newRequestID()
		}
		c.Set(requestIDKey, id)
		c.Header(requestIDHeader, id)
		c.Next()
	}
}

func newRequestID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err == nil {
		return "req-" + hex.EncodeToString(bytes)
	}
	return "req-" + strconv.FormatInt(time.Now().UnixNano(), 10)
}

func requestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()
		c.Next()
		logger.InfoContext(c.Request.Context(), "request completed",
			"request_id", getRequestID(c),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration_ms", time.Since(startedAt).Milliseconds(),
		)
	}
}

func structuredErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Written() || len(c.Errors) == 0 {
			return
		}

		responseError := &apiError{
			status:  http.StatusInternalServerError,
			code:    "INTERNAL_ERROR",
			message: "The backend failed unexpectedly.",
		}
		var contractedError *apiError
		if errors.As(c.Errors.Last().Err, &contractedError) {
			responseError = contractedError
		}
		writeError(c, responseError)
	}
}

func recoverPanics(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				logger.ErrorContext(c.Request.Context(), "panic recovered",
					"request_id", getRequestID(c),
					"panic_type", fmt.Sprintf("%T", recovered),
				)
				abort(c, &apiError{
					status:  http.StatusInternalServerError,
					code:    "INTERNAL_ERROR",
					message: "The backend failed unexpectedly.",
				})
			}
		}()
		c.Next()
	}
}

func abort(c *gin.Context, err *apiError) {
	_ = c.Error(err)
	c.Abort()
}

func writeError(c *gin.Context, err *apiError) {
	details := err.details
	if details == nil {
		details = []errorDetail{}
	}
	c.JSON(err.status, errorEnvelope{Error: errorBody{
		Code:      err.code,
		Message:   err.message,
		Retryable: err.retryable,
		RequestID: getRequestID(c),
		Details:   details,
	}})
}

func getRequestID(c *gin.Context) string {
	requestID, _ := c.Get(requestIDKey)
	return requestID.(string)
}
