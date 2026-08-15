package httpapi

import (
	"errors"
	"net/http"
	"strconv"

	"batam-medhub/internal/service"

	"github.com/gin-gonic/gin"
)

func handleGetJourneyItinerary(journeySvc *service.JourneyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		patientID := c.GetString(contextPatientIDKey)
		if patientID == "" {
			abort(c, &apiError{
				status:  http.StatusUnauthorized,
				code:    "UNAUTHORIZED",
				message: "Authentication is required to view journeys.",
			})
			return
		}

		journeyID := c.Param("journey_id")
		if journeyID == "" {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "journey_id is required.",
			})
			return
		}

		detail, err := journeySvc.GetJourneyItinerary(c.Request.Context(), patientID, journeyID)
		if err != nil {
			if errors.Is(err, service.ErrJourneyNotFound) {
				abort(c, &apiError{
					status:  http.StatusNotFound,
					code:    "NOT_FOUND",
					message: "Journey not found.",
				})
				return
			}

			abort(c, &apiError{
				status:  http.StatusInternalServerError,
				code:    "INTERNAL_ERROR",
				message: "Failed to retrieve journey.",
			})
			return
		}

		c.JSON(http.StatusOK, detail)
	}
}

func handleGetJourneyItineraryVersion(journeySvc *service.JourneyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		patientID := c.GetString(contextPatientIDKey)
		if patientID == "" {
			abort(c, &apiError{
				status:  http.StatusUnauthorized,
				code:    "UNAUTHORIZED",
				message: "Authentication is required to view itinerary versions.",
			})
			return
		}

		journeyID := c.Param("journey_id")
		versionStr := c.Param("version")

		versionNumber, err := strconv.Atoi(versionStr)
		if err != nil || versionNumber < 1 {
			abort(c, &apiError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "version must be a positive integer.",
			})
			return
		}

		versionDTO, err := journeySvc.GetJourneyItineraryVersion(c.Request.Context(), patientID, journeyID, versionNumber)
		if err != nil {
			if errors.Is(err, service.ErrJourneyNotFound) || errors.Is(err, service.ErrItineraryVersionNotFound) {
				abort(c, &apiError{
					status:  http.StatusNotFound,
					code:    "NOT_FOUND",
					message: "Itinerary version not found.",
				})
				return
			}

			abort(c, &apiError{
				status:  http.StatusInternalServerError,
				code:    "INTERNAL_ERROR",
				message: "Failed to retrieve itinerary version.",
			})
			return
		}

		c.JSON(http.StatusOK, versionDTO)
	}
}

func handleListJourneys(journeySvc *service.JourneyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		patientID := c.GetString(contextPatientIDKey)
		if patientID == "" {
			abort(c, &apiError{
				status:  http.StatusUnauthorized,
				code:    "UNAUTHORIZED",
				message: "Authentication is required to list journeys.",
			})
			return
		}

		list, err := journeySvc.ListJourneys(c.Request.Context(), patientID)
		if err != nil {
			abort(c, &apiError{
				status:  http.StatusInternalServerError,
				code:    "INTERNAL_ERROR",
				message: "Failed to list journeys.",
			})
			return
		}

		c.JSON(http.StatusOK, list)
	}
}
