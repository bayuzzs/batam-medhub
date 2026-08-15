package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// allowedCORSHeaders are the headers the patient mobile web app may send.
// The app authenticates with a Bearer token and uses the Idempotency-Key
// header for booking/recovery operations.
var allowedCORSHeaders = []string{
	"Authorization",
	"Content-Type",
	"Idempotency-Key",
	"X-Request-ID",
	"X-Demo-Key",
}

var allowedCORSMethods = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodOptions,
}

// cors allows browser clients (e.g. the Flutter web build served from a local
// dev server) to call the core API. It reflects the request origin and answers
// OPTIONS preflight requests so authenticated cross-origin calls succeed.
// Integration/demo only: it is deliberately permissive for local development.
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}

		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", joinComma(allowedCORSHeaders))
		c.Header("Access-Control-Allow-Methods", joinComma(allowedCORSMethods))
		c.Header("Access-Control-Max-Age", "600")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func joinComma(items []string) string {
	out := ""
	for i, item := range items {
		if i > 0 {
			out += ", "
		}
		out += item
	}
	return out
}
