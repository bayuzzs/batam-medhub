package platform

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

var requestIDRegex = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$`)

const RequestIDContextKey = "request_id"

// RequestIDMiddleware validates the X-Request-ID header against the contracted pattern.
// If absent or invalid, it generates a valid replacement.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if !requestIDRegex.MatchString(requestID) {
			requestID = NewRequestID()
		}
		c.Set(RequestIDContextKey, requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// GetRequestID retrieves the validated request ID from the Gin context.
func GetRequestID(c *gin.Context) string {
	if val, ok := c.Get(RequestIDContextKey); ok {
		if id, ok := val.(string); ok && id != "" {
			return id
		}
	}
	return NewRequestID()
}

// NewRequestID generates a compliant random request ID matching ^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$.
func NewRequestID() string {
	var value [12]byte
	if _, err := rand.Read(value[:]); err != nil {
		return fmt.Sprintf("req-%d", time.Now().UTC().UnixNano())
	}
	return "req-" + hex.EncodeToString(value[:])
}
