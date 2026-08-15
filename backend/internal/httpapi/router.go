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

	"batam-medhub/internal/adapter"
	"batam-medhub/internal/ai"
	"batam-medhub/internal/config"
	"batam-medhub/internal/service"

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

// New constructs and configures the Gin HTTP engine with middleware and routes.
func New(db *gorm.DB, cfg config.Config, logger *slog.Logger) *gin.Engine {
	if logger == nil {
		logger = slog.Default()
	}
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.HandleMethodNotAllowed = true
	_ = router.SetTrustedProxies(nil)

	router.Use(requestID(), requestLogger(logger), structuredErrors(), recoverPanics(logger))

	router.NoRoute(func(c *gin.Context) {
		abort(c, &apiError{
			status:  http.StatusNotFound,
			code:    "NOT_FOUND",
			message: "The requested resource was not found.",
		})
	})

	router.NoMethod(func(c *gin.Context) {
		abort(c, &apiError{
			status:  http.StatusMethodNotAllowed,
			code:    "METHOD_NOT_ALLOWED",
			message: "Method not allowed for this route.",
		})
	})

	router.GET("/healthz", liveness)
	router.GET("/readyz", readiness(db))

	authSvc := service.NewAuthService(db, cfg)
	profileSvc := service.NewProfileService(db, cfg)
	catalogSvc := service.NewCatalogService(db)
	moneySvc := service.NewMoneyService(db)
	idemSvc := service.NewIdempotencyService(db)

	hospAdapter := adapter.NewHospitalAdapter(cfg.HospitalBaseURL, cfg.HospitalIntegrationKey, cfg.ProviderTimeout)
	ferryAdapter := adapter.NewFerryAdapter(cfg.FerryBaseURL, cfg.FerryIntegrationKey, cfg.ProviderTimeout)
	hotelAdapter := adapter.NewHotelAdapter(cfg.HotelBaseURL, cfg.HotelIntegrationKey, cfg.ProviderTimeout)
	transAdapter := adapter.NewTransportAdapter(cfg.TransportBaseURL, cfg.TransportIntegrationKey, cfg.ProviderTimeout)
	aggregator := adapter.NewAggregator(hospAdapter, ferryAdapter, hotelAdapter, transAdapter)

	aiClient := ai.NewClient(cfg.CloudflareAIBaseURL, cfg.CloudflareAccountID, cfg.CloudflareAPIToken, cfg.CloudflareAIModel, cfg.CloudflareAITimeout)
	aiExtractor := ai.NewExtractor(aiClient, catalogSvc, logger)

	tripSvc := service.NewTripService(db, catalogSvc, moneySvc, aggregator, aiExtractor)
	bookingSvc := service.NewBookingSagaService(db, hospAdapter, ferryAdapter, hotelAdapter, transAdapter, moneySvc)
	journeySvc := service.NewJourneyService(db)
	provAuthSvc := service.NewProviderAuthService(db)
	disruptionSvc := service.NewDisruptionService(db, hospAdapter, ferryAdapter, hotelAdapter, transAdapter, moneySvc, journeySvc)

	registerLimiter := newIPRateLimiter(30, time.Minute)
	loginLimiter := newIPRateLimiter(30, time.Minute)
	refreshLimiter := newIPRateLimiter(60, time.Minute)

	v1 := router.Group("/v1")
	{
		authGroup := v1.Group("/auth")
		authGroup.Use(noStoreHeader())
		{
			authGroup.POST("/register", rateLimitMiddleware(registerLimiter), handleRegister(authSvc))
			authGroup.POST("/login", rateLimitMiddleware(loginLimiter), handleLogin(authSvc))
			authGroup.POST("/refresh", rateLimitMiddleware(refreshLimiter), handleRefresh(authSvc))
			authGroup.POST("/logout", handleLogout(authSvc))
		}

		profileGroup := v1.Group("/profile")
		profileGroup.Use(patientBearerAuth(db, cfg))
		{
			profileGroup.GET("", handleGetProfile(profileSvc))
			profileGroup.PATCH("", noStoreHeader(), handleUpdateProfile(profileSvc))
		}

		catalogGroup := v1.Group("/medical-services")
		catalogGroup.Use(patientBearerAuth(db, cfg))
		{
			catalogGroup.GET("", handleListMedicalServices(catalogSvc))
		}

		tripGroup := v1.Group("/trip-requests")
		tripGroup.Use(patientBearerAuth(db, cfg))
		{
			tripGroup.POST("", handleCreateTripRequest(tripSvc, idemSvc))
			tripGroup.GET("/:trip_request_id", handleGetTripRequest(tripSvc))
			tripGroup.PATCH("/:trip_request_id/intent", handleAmendTripRequestIntent(tripSvc, idemSvc))
			tripGroup.POST("/:trip_request_id/plans", handleGenerateTripPlans(tripSvc, idemSvc))
			tripGroup.POST("/:trip_request_id/select-plan", handleSelectPlanForTrip(bookingSvc, tripSvc, idemSvc))
		}

		planGroup := v1.Group("/plan-options")
		planGroup.Use(patientBearerAuth(db, cfg))
		{
			planGroup.POST("/:plan_option_id/confirm", handleConfirmPlanOption(bookingSvc, idemSvc))
		}

		journeyGroup := v1.Group("/journeys")
		journeyGroup.Use(patientBearerAuth(db, cfg))
		{
			journeyGroup.GET("", handleListJourneys(journeySvc))
			journeyGroup.GET("/:journey_id", handleGetJourneyItinerary(journeySvc))
			journeyGroup.GET("/:journey_id/itinerary", handleGetJourneyItinerary(journeySvc))
			journeyGroup.GET("/:journey_id/itineraries/:version", handleGetJourneyItineraryVersion(journeySvc))
			journeyGroup.POST("/:journey_id/recovery-options/:option_id/select", handleApproveRecoveryOption(disruptionSvc, idemSvc))
		}

		// Provider Disruption Ingestion (OpenAPI: POST /v1/provider/disruptions and alias POST /v1/events/disruptions)
		v1.POST("/provider/disruptions", handleProviderDisruption(provAuthSvc, disruptionSvc))
		v1.POST("/events/disruptions", handleProviderDisruption(provAuthSvc, disruptionSvc))

		// Disruption Detail (Auth: PatientBearer)
		disruptionGroup := v1.Group("/disruptions")
		disruptionGroup.Use(patientBearerAuth(db, cfg))
		{
			disruptionGroup.GET("/:disruption_id", handleGetDisruption(disruptionSvc))
		}

		// Recovery Options Approval (OpenAPI: POST /v1/recovery-options/:recovery_option_id/approve)
		recoveryGroup := v1.Group("/recovery-options")
		recoveryGroup.Use(patientBearerAuth(db, cfg))
		{
			recoveryGroup.POST("/:recovery_option_id/approve", handleApproveRecoveryOption(disruptionSvc, idemSvc))
		}

		// Demo Reset Endpoint (Auth: DemoSecret)
		v1.POST("/demo/reset", handleDemoReset(cfg, service.NewDemoService(db), idemSvc))
	}

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
