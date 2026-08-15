package providerapp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"batam-medhub/providers/internal/hospital"
	"batam-medhub/providers/internal/platform"
)

type Identity struct {
	ID   string
	Type string
}

func Run(identity Identity) error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	integrationKey := resolveIntegrationKey(identity.Type)

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get database handle: %w", err)
	}
	defer sqlDB.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), platform.RequestIDMiddleware())

	// Unauthenticated health endpoint
	router.GET("/healthz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(ctx); err != nil {
			platform.RespondServiceUnavailable(c, "The provider is temporarily unavailable.")
			return
		}

		c.JSON(http.StatusOK, platform.HealthResponse{
			Status:         "UP",
			ProviderID:     identity.ID,
			ProviderType:   identity.Type,
			DatabaseStatus: "UP",
			CheckedAt:      platform.FormatUTC(time.Now()),
		})
	})

	// Authenticated v1 API group
	v1 := router.Group("/v1")
	v1.Use(platform.AuthMiddleware(integrationKey))

	if identity.Type == "HOSPITAL" {
		repo := hospital.NewRepository(db)
		svc := hospital.NewService(identity.ID, repo)
		handler := hospital.NewHandler(svc)
		handler.RegisterRoutes(v1)
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("%s provider (%s) listening on :%s", identity.Type, identity.ID, port)
		serverErrors <- server.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(stop)

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-stop:
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown server: %w", err)
		}
	}

	return nil
}

func resolveIntegrationKey(providerType string) string {
	if key := os.Getenv("PROVIDER_INTEGRATION_KEY"); key != "" {
		return key
	}
	switch strings.ToUpper(providerType) {
	case "HOSPITAL":
		if key := os.Getenv("HOSPITAL_INTEGRATION_KEY"); key != "" {
			return key
		}
		return "hospital_dev_secret"
	case "FERRY":
		if key := os.Getenv("FERRY_INTEGRATION_KEY"); key != "" {
			return key
		}
		return "ferry_dev_secret"
	case "HOTEL":
		if key := os.Getenv("HOTEL_INTEGRATION_KEY"); key != "" {
			return key
		}
		return "hotel_dev_secret"
	case "TRANSPORT":
		if key := os.Getenv("TRANSPORT_INTEGRATION_KEY"); key != "" {
			return key
		}
		return "transport_dev_secret"
	default:
		return "dev_secret"
	}
}
