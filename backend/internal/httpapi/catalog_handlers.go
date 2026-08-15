package httpapi

import (
	"net/http"

	"batam-medhub/internal/service"

	"github.com/gin-gonic/gin"
)

func handleListMedicalServices(svc *service.CatalogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		catalog, err := svc.ListMedicalServices(c.Request.Context())
		if err != nil {
			abort(c, &apiError{
				status:  http.StatusInternalServerError,
				code:    "INTERNAL_ERROR",
				message: "Failed to retrieve medical services catalog.",
			})
			return
		}

		c.JSON(http.StatusOK, catalog)
	}
}
