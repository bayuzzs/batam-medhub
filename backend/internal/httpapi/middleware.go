package httpapi

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"batam-medhub/internal/auth"
	"batam-medhub/internal/config"
	"batam-medhub/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	contextPatientIDKey = "patient_id"
	contextSessionIDKey = "session_id"
	contextPatientKey   = "patient"
)

// patientBearerAuth verifies the HS256 JWT access token and active session in the database.
func patientBearerAuth(db *gorm.DB, cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			abort(c, &apiError{
				status:  http.StatusUnauthorized,
				code:    "UNAUTHORIZED",
				message: "Authentication required.",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			abort(c, &apiError{
				status:  http.StatusUnauthorized,
				code:    "UNAUTHORIZED",
				message: "Authorization header must use Bearer scheme.",
			})
			return
		}

		rawToken := strings.TrimSpace(parts[1])
		claims, err := auth.ValidateAccessToken(rawToken, cfg.JWTSigningSecret, cfg.JWTIssuer, cfg.JWTAudience)
		if err != nil {
			abort(c, &apiError{
				status:  http.StatusUnauthorized,
				code:    "UNAUTHORIZED",
				message: "Access token is invalid or expired.",
			})
			return
		}

		now := time.Now().UTC()
		var session model.AuthSession
		if err := db.WithContext(c.Request.Context()).
			Where("id = ? AND revoked_at IS NULL AND expires_at > ?", claims.SessionID, now).
			First(&session).Error; err != nil {
			abort(c, &apiError{
				status:  http.StatusUnauthorized,
				code:    "UNAUTHORIZED",
				message: "Session has been revoked or expired.",
			})
			return
		}

		var patient model.Patient
		if err := db.WithContext(c.Request.Context()).
			Where("id = ? AND status = 'ACTIVE'", claims.Subject).
			First(&patient).Error; err != nil {
			abort(c, &apiError{
				status:  http.StatusUnauthorized,
				code:    "UNAUTHORIZED",
				message: "Patient account is not active.",
			})
			return
		}

		// Reject bearer access when JWT preferred_currency claim differs from persisted profile
		if patient.PreferredCurrency != claims.PreferredCurrency {
			abort(c, &apiError{
				status:  http.StatusUnauthorized,
				code:    "STALE_CREDENTIALS",
				message: "Preferred currency claim does not match current profile. Please refresh your session.",
			})
			return
		}

		c.Set(contextPatientIDKey, patient.ID)
		c.Set(contextSessionIDKey, session.ID)
		c.Set(contextPatientKey, &patient)
		c.Next()
	}
}

// ipRateLimiter implements a bounded in-memory sliding window rate limiter per client IP.
type ipRateLimiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	clients map[string][]time.Time
}

func newIPRateLimiter(limit int, window time.Duration) *ipRateLimiter {
	return &ipRateLimiter{
		limit:   limit,
		window:  window,
		clients: make(map[string][]time.Time),
	}
}

func (rl *ipRateLimiter) allow(ip string) (bool, time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	timestamps := rl.clients[ip]
	valid := make([]time.Time, 0, len(timestamps))
	for _, t := range timestamps {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.limit {
		rl.clients[ip] = valid
		retryAfter := valid[0].Add(rl.window).Sub(now)
		if retryAfter < time.Second {
			retryAfter = time.Second
		}
		return false, retryAfter
	}

	valid = append(valid, now)
	rl.clients[ip] = valid
	return true, 0
}

func rateLimitMiddleware(limiter *ipRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		allowed, retryAfter := limiter.allow(ip)
		if !allowed {
			seconds := int(retryAfter.Seconds())
			if seconds < 1 {
				seconds = 1
			}
			c.Header("Retry-After", strconv.Itoa(seconds))
			abort(c, &apiError{
				status:    http.StatusTooManyRequests,
				code:      "TOO_MANY_REQUESTS",
				message:   "Rate limit exceeded. Please try again later.",
				retryable: true,
			})
			return
		}
		c.Next()
	}
}

func noStoreHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Next()
	}
}
