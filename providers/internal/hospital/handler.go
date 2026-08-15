package hospital

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"batam-medhub/providers/internal/platform"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/offers/search", h.SearchOffers)
	rg.POST("/holds", h.CreateHold)
	rg.POST("/holds/:hold_id/confirm", h.ConfirmHold)
	rg.POST("/holds/:hold_id/release", h.ReleaseHold)
	rg.GET("/reservations/:reservation_id", h.GetReservation)
	rg.POST("/reservations/:reservation_id/release", h.ReleaseReservation)
}

func (h *Handler) SearchOffers(c *gin.Context) {
	bodyBytes, err := platform.ReadAndRestoreBody(c)
	if err != nil {
		platform.RespondBadRequest(c, "The request failed validation.", []platform.ErrorDetail{
			{Reason: "failed to read request body"},
		})
		return
	}

	payload, details, err := platform.DecodeStrictJSON[SearchRequestPayload](bodyBytes, platform.DefaultMaxBodyBytes)
	if err != nil {
		platform.RespondBadRequest(c, "The request failed validation.", details)
		return
	}

	resp, details, err := h.service.SearchOffers(c.Request.Context(), payload, time.Now().UTC())
	if err != nil {
		if errors.Is(err, ErrInvalidRequest) {
			platform.RespondBadRequest(c, "The request failed validation.", details)
			return
		}
		platform.RespondInternalError(c, "The provider could not complete the request.")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) CreateHold(c *gin.Context) {
	idempKey := c.GetHeader("Idempotency-Key")
	if !platform.ValidateIdempotencyKey(idempKey) {
		platform.RespondBadRequest(c, "The request failed validation.", []platform.ErrorDetail{
			{Field: "Idempotency-Key", Reason: "is invalid or missing"},
		})
		return
	}

	bodyBytes, err := platform.ReadAndRestoreBody(c)
	if err != nil {
		platform.RespondBadRequest(c, "The request failed validation.", []platform.ErrorDetail{
			{Reason: "failed to read request body"},
		})
		return
	}

	payload, details, err := platform.DecodeStrictJSON[CreateHoldRequestPayload](bodyBytes, platform.DefaultMaxBodyBytes)
	if err != nil {
		platform.RespondBadRequest(c, "The request failed validation.", details)
		return
	}

	fingerprint := platform.ComputeFingerprint("POST", c.FullPath(), bodyBytes)
	resp, replayed, capacityDetails, offerExpiredDetails, validationDetails, err := h.service.CreateHold(
		c.Request.Context(),
		payload,
		idempKey,
		fingerprint,
		time.Now().UTC(),
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrProviderIdentityMismatch):
			platform.RespondForbidden(c, "The asserted provider identity does not match this provider service.")
		case errors.Is(err, ErrInvalidRequest):
			platform.RespondBadRequest(c, "The request failed validation.", validationDetails)
		case errors.Is(err, ErrIdempotencyConflict):
			platform.RespondIdempotencyConflict(c, "The idempotency key was already used with a different request.")
		case errors.Is(err, ErrNotFound):
			platform.RespondNotFound(c, "The requested provider resource was not found.")
		case errors.Is(err, ErrCapacityConflict):
			var details []platform.ErrorDetail
			if capacityDetails != nil {
				details = []platform.ErrorDetail{
					{
						Field:  "units",
						Reason: fmt.Sprintf("requested=%d available=%d", capacityDetails.Requested, capacityDetails.Available),
					},
				}
			}
			platform.RespondCapacityConflict(c, "The requested units are no longer available.", details)
		case errors.Is(err, ErrOfferChanged):
			platform.RespondOfferChanged(c, "The offer price has changed.", nil)
		case errors.Is(err, ErrOfferExpired):
			var details []platform.ErrorDetail
			if offerExpiredDetails != nil {
				details = []platform.ErrorDetail{
					{
						Field:  "offer_id",
						Reason: fmt.Sprintf("valid_until=%s", platform.FormatUTC(offerExpiredDetails.ValidUntil)),
					},
				}
			}
			platform.RespondOfferExpired(c, "The offer expired before it could be held.", details)
		default:
			platform.RespondInternalError(c, "The provider could not complete the request.")
		}
		return
	}

	c.Header("Idempotency-Replayed", strconv.FormatBool(replayed))
	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) ConfirmHold(c *gin.Context) {
	holdID := c.Param("hold_id")
	if !platform.ValidateResourceId(holdID) {
		platform.RespondBadRequest(c, "The request failed validation.", []platform.ErrorDetail{
			{Field: "hold_id", Reason: "is invalid"},
		})
		return
	}

	idempKey := c.GetHeader("Idempotency-Key")
	if !platform.ValidateIdempotencyKey(idempKey) {
		platform.RespondBadRequest(c, "The request failed validation.", []platform.ErrorDetail{
			{Field: "Idempotency-Key", Reason: "is invalid or missing"},
		})
		return
	}

	bodyBytes, err := platform.ReadAndRestoreBody(c)
	if err != nil {
		platform.RespondBadRequest(c, "The request failed validation.", []platform.ErrorDetail{
			{Reason: "failed to read request body"},
		})
		return
	}

	fingerprint := platform.ComputeFingerprint("POST", c.FullPath(), bodyBytes)

	resp, replayed, expiredDetails, err := h.service.ConfirmHold(
		c.Request.Context(),
		holdID,
		idempKey,
		fingerprint,
		time.Now().UTC(),
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrIdempotencyConflict):
			platform.RespondIdempotencyConflict(c, "The idempotency key was already used with a different request.")
		case errors.Is(err, ErrNotFound):
			platform.RespondNotFound(c, "The requested provider resource was not found.")
		case errors.Is(err, ErrHoldExpired):
			var details []platform.ErrorDetail
			if expiredDetails != nil {
				details = []platform.ErrorDetail{
					{
						Field:  "hold_id",
						Reason: fmt.Sprintf("expired_at=%s", platform.FormatUTC(expiredDetails.ExpiredAt)),
					},
				}
			}
			platform.RespondHoldExpired(c, "The hold expired before confirmation.", details)
		case errors.Is(err, ErrInvalidState):
			platform.RespondInvalidState(c, "The resource cannot perform this transition from its current status.")
		default:
			platform.RespondInternalError(c, "The provider could not complete the request.")
		}
		return
	}

	c.Header("Idempotency-Replayed", strconv.FormatBool(replayed))
	c.Header("Location", "/v1/reservations/"+resp.ReservationID)
	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) ReleaseHold(c *gin.Context) {
	holdID := c.Param("hold_id")
	if !platform.ValidateResourceId(holdID) {
		platform.RespondBadRequest(c, "The request failed validation.", []platform.ErrorDetail{
			{Field: "hold_id", Reason: "is invalid"},
		})
		return
	}

	idempKey := c.GetHeader("Idempotency-Key")
	if !platform.ValidateIdempotencyKey(idempKey) {
		platform.RespondBadRequest(c, "The request failed validation.", []platform.ErrorDetail{
			{Field: "Idempotency-Key", Reason: "is invalid or missing"},
		})
		return
	}

	bodyBytes, err := platform.ReadAndRestoreBody(c)
	if err != nil {
		platform.RespondBadRequest(c, "The request failed validation.", []platform.ErrorDetail{
			{Reason: "failed to read request body"},
		})
		return
	}

	fingerprint := platform.ComputeFingerprint("POST", c.FullPath(), bodyBytes)

	resp, replayed, err := h.service.ReleaseHold(
		c.Request.Context(),
		holdID,
		idempKey,
		fingerprint,
		time.Now().UTC(),
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrIdempotencyConflict):
			platform.RespondIdempotencyConflict(c, "The idempotency key was already used with a different request.")
		case errors.Is(err, ErrNotFound):
			platform.RespondNotFound(c, "The requested provider resource was not found.")
		case errors.Is(err, ErrInvalidState):
			platform.RespondInvalidState(c, "The resource cannot perform this transition from its current status.")
		default:
			platform.RespondInternalError(c, "The provider could not complete the request.")
		}
		return
	}

	c.Header("Idempotency-Replayed", strconv.FormatBool(replayed))
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetReservation(c *gin.Context) {
	reservationID := c.Param("reservation_id")
	if !platform.ValidateResourceId(reservationID) {
		platform.RespondBadRequest(c, "The request failed validation.", []platform.ErrorDetail{
			{Field: "reservation_id", Reason: "is invalid"},
		})
		return
	}

	resp, err := h.service.GetReservation(c.Request.Context(), reservationID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			platform.RespondNotFound(c, "The requested provider resource was not found.")
			return
		}
		platform.RespondInternalError(c, "The provider could not complete the request.")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ReleaseReservation(c *gin.Context) {
	reservationID := c.Param("reservation_id")
	if !platform.ValidateResourceId(reservationID) {
		platform.RespondBadRequest(c, "The request failed validation.", []platform.ErrorDetail{
			{Field: "reservation_id", Reason: "is invalid"},
		})
		return
	}

	idempKey := c.GetHeader("Idempotency-Key")
	if !platform.ValidateIdempotencyKey(idempKey) {
		platform.RespondBadRequest(c, "The request failed validation.", []platform.ErrorDetail{
			{Field: "Idempotency-Key", Reason: "is invalid or missing"},
		})
		return
	}

	bodyBytes, err := platform.ReadAndRestoreBody(c)
	if err != nil {
		platform.RespondBadRequest(c, "The request failed validation.", []platform.ErrorDetail{
			{Reason: "failed to read request body"},
		})
		return
	}

	fingerprint := platform.ComputeFingerprint("POST", c.FullPath(), bodyBytes)

	resp, replayed, err := h.service.ReleaseReservation(
		c.Request.Context(),
		reservationID,
		idempKey,
		fingerprint,
		time.Now().UTC(),
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrIdempotencyConflict):
			platform.RespondIdempotencyConflict(c, "The idempotency key was already used with a different request.")
		case errors.Is(err, ErrNotFound):
			platform.RespondNotFound(c, "The requested provider resource was not found.")
		case errors.Is(err, ErrInvalidState):
			platform.RespondInvalidState(c, "The resource cannot perform this transition from its current status.")
		default:
			platform.RespondInternalError(c, "The provider could not complete the request.")
		}
		return
	}

	c.Header("Idempotency-Replayed", strconv.FormatBool(replayed))
	c.JSON(http.StatusOK, resp)
}
