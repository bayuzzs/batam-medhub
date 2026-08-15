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

	"batam-medhub/providers/internal/ferry"
	"batam-medhub/providers/internal/hospital"
	"batam-medhub/providers/internal/hotel"
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

	integrationKey, err := resolveIntegrationKey(identity.Type)
	if err != nil {
		return err
	}

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
	router.Use(platform.SafeRecoveryMiddleware(), platform.RequestIDMiddleware())

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
	} else if identity.Type == "FERRY" {
		repo := ferry.NewRepository(db)
		svc := ferry.NewService(identity.ID, repo)
		handler := ferry.NewHandler(svc)
		handler.RegisterRoutes(v1)
	} else if identity.Type == "HOTEL" {
		repo := hotel.NewRepository(db)
		svc := hotel.NewService(identity.ID, repo)
		handler := hotel.NewHandler(svc)
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

func resolveIntegrationKey(providerType string) (string, error) {
	if key := os.Getenv("PROVIDER_INTEGRATION_KEY"); key != "" {
		return key, nil
	}
	var envVar string
	switch strings.ToUpper(providerType) {
	case "HOSPITAL":
		envVar = "HOSPITAL_INTEGRATION_KEY"
	case "FERRY":
		envVar = "FERRY_INTEGRATION_KEY"
	case "HOTEL":
		envVar = "HOTEL_INTEGRATION_KEY"
	case "TRANSPORT":
		envVar = "TRANSPORT_INTEGRATION_KEY"
	default:
		return "", fmt.Errorf("unknown provider type: %s", providerType)
	}
	if key := os.Getenv(envVar); key != "" {
		return key, nil
	}
	return "", fmt.Errorf("integration key is required (set PROVIDER_INTEGRATION_KEY or %s)", envVar)
}
