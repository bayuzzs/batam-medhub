package platform

import (
	"crypto/subtle"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifies the X-Integration-Key header against the configured secret.
func AuthMiddleware(expectedKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-Integration-Key")
		if key == "" || subtle.ConstantTimeCompare([]byte(key), []byte(expectedKey)) != 1 {
			RespondUnauthorized(c, "Provider integration authentication failed.")
			c.Abort()
			return
		}
		c.Next()
	}
}
