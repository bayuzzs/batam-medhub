package platform

import (
	"log"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// SafeRecoveryMiddleware recovers from panics and emits the contracted INTERNAL_ERROR response.
// It never logs request headers or authentication secrets.
func SafeRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				reqID := GetRequestID(c)
				log.Printf("[PANIC RECOVERED] req_id=%s method=%s path=%s err=%v\nstack:\n%s",
					reqID, c.Request.Method, c.Request.URL.Path, r, string(debug.Stack()))
				if !c.Writer.Written() {
					RespondInternalError(c, "The provider could not complete the request.")
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}
