package transport

import (
	"context"
	"errors"
	"fmt"
	"time"

	"batam-medhub/providers/internal/platform"
)

var (
	ErrInvalidRequest           = errors.New("invalid request")
	ErrProviderIdentityMismatch = errors.New("provider identity mismatch")
)

type TransportSearchCriteria struct {
	PickupLocationCode  string              `json:"pickup_location_code"`
	DropoffLocationCode string              `json:"dropoff_location_code"`
	PassengerCount      int                 `json:"passenger_count"`
	PickupWindow        platform.TimeWindow `json:"pickup_window"`
	Accessibility       []string            `json:"accessibility,omitempty"`
}

type SearchRequestPayload struct {
	ProviderType string                  `json:"provider_type"`
	Criteria     TransportSearchCriteria `json:"criteria"`
}

type TransportOfferDetails struct {
	ProviderType        string   `json:"provider_type"`
	AvailabilityID      string   `json:"availability_id"`
	VehicleType         string   `json:"vehicle_type"`
	PickupLocationCode  string   `json:"pickup_location_code"`
	DropoffLocationCode string   `json:"dropoff_location_code"`
	PassengerCapacity   int      `json:"passenger_capacity"`
	Accessibility       []string `json:"accessibility"`
}

type OfferResponse struct {
	OfferID        string                `json:"offer_id"`
	ProviderID     string                `json:"provider_id"`
	ProviderType   string                `json:"provider_type"`
	Status         string                `json:"status"`
	ServiceWindow  platform.TimeWindow   `json:"service_window"`
	AvailableUnits int                   `json:"available_units"`
	UnitPrice      platform.Money        `json:"unit_price"`
	ValidUntil     string                `json:"valid_until"`
	Synthetic      bool                  `json:"synthetic"`
	Source         string                `json:"source"`
	Details        TransportOfferDetails `json:"details"`
}

type SearchTransportResponse struct {
	ProviderID   string          `json:"provider_id"`
	ProviderType string          `json:"provider_type"`
	Offers       []OfferResponse `json:"offers"`
}

type TransportBookingRequirements struct {
	PassengerCount      int                 `json:"passenger_count"`
	PickupLocationCode  string              `json:"pickup_location_code"`
	DropoffLocationCode string              `json:"dropoff_location_code"`
	PickupWindow        platform.TimeWindow `json:"pickup_window"`
	Accessibility       []string            `json:"accessibility"`
}

type CreateHoldRequestPayload struct {
	ProviderID          string                        `json:"provider_id"`
	ProviderType        string                        `json:"provider_type"`
	OfferID             string                        `json:"offer_id"`
	Units               int                           `json:"units"`
	ExpectedUnitPrice   platform.Money                `json:"expected_unit_price"`
	ClientReference     string                        `json:"client_reference"`
	BookingRequirements *TransportBookingRequirements `json:"booking_requirements"`
}

type HoldResponse struct {
	HoldID            string              `json:"hold_id"`
	ExternalReference string              `json:"external_reference"`
	ProviderID        string              `json:"provider_id"`
	ProviderType      string              `json:"provider_type"`
	OfferID           string              `json:"offer_id"`
	ClientReference   string              `json:"client_reference"`
	Status            string              `json:"status"`
	Units             int                 `json:"units"`
	UnitPrice         platform.Money      `json:"unit_price"`
	TotalPrice        platform.Money      `json:"total_price"`
	ServiceWindow     platform.TimeWindow `json:"service_window"`
	CreatedAt         string              `json:"created_at"`
	ExpiresAt         string              `json:"expires_at"`
}

type ReservationResponse struct {
	ReservationID     string              `json:"reservation_id"`
	ExternalReference string              `json:"external_reference"`
	HoldID            string              `json:"hold_id"`
	ProviderID        string              `json:"provider_id"`
	ProviderType      string              `json:"provider_type"`
	OfferID           string              `json:"offer_id"`
	ClientReference   string              `json:"client_reference"`
	Status            string              `json:"status"`
	Units             int                 `json:"units"`
	TotalPrice        platform.Money      `json:"total_price"`
	ServiceWindow     platform.TimeWindow `json:"service_window"`
	ConfirmedAt       string              `json:"confirmed_at"`
	ReleasedAt        *string             `json:"released_at"`
}

type Service struct {
	providerID string
	repo       *Repository
}

func NewService(providerID string, repo *Repository) *Service {
	return &Service{
		providerID: providerID,
		repo:       repo,
	}
}

func (s *Service) SearchOffers(ctx context.Context, req SearchRequestPayload, now time.Time) (*SearchTransportResponse, []platform.ErrorDetail, error) {
	if req.ProviderType != "TRANSPORT" {
		return nil, []platform.ErrorDetail{{Field: "provider_type", Reason: "must be TRANSPORT"}}, ErrInvalidRequest
	}

	crit := req.Criteria
	var details []platform.ErrorDetail

	if !platform.ValidateLocationCode(crit.PickupLocationCode) {
		details = append(details, platform.ErrorDetail{Field: "criteria.pickup_location_code", Reason: "must match ^[A-Z][A-Z0-9_]*$ (length 2..64)"})
	}
	if !platform.ValidateLocationCode(crit.DropoffLocationCode) {
		details = append(details, platform.ErrorDetail{Field: "criteria.dropoff_location_code", Reason: "must match ^[A-Z][A-Z0-9_]*$ (length 2..64)"})
	}
	if crit.PickupLocationCode != "" && crit.DropoffLocationCode != "" && crit.PickupLocationCode == crit.DropoffLocationCode {
		details = append(details, platform.ErrorDetail{Field: "criteria.dropoff_location_code", Reason: "pickup and dropoff locations must be different"})
	}

	if crit.PassengerCount < 1 || crit.PassengerCount > 20 {
		details = append(details, platform.ErrorDetail{Field: "criteria.passenger_count", Reason: "must be between 1 and 20"})
	}

	windowDetails := platform.ValidateTimeWindow(crit.PickupWindow, "criteria.pickup_window")
	details = append(details, windowDetails...)

	if len(crit.Accessibility) > 0 {
		details = append(details, platform.ValidateAccessibility(crit.Accessibility, "criteria.accessibility")...)
	}

	if len(details) > 0 {
		return nil, details, ErrInvalidRequest
	}

	windowStart, _ := time.Parse(time.RFC3339, crit.PickupWindow.StartsAt)
	windowEnd, _ := time.Parse(time.RFC3339, crit.PickupWindow.EndsAt)

	offersWithUnits, err := s.repo.SearchAvailableOffers(
		ctx,
		crit.PickupLocationCode,
		crit.DropoffLocationCode,
		windowStart,
		windowEnd,
		crit.PassengerCount,
		crit.Accessibility,
		now,
	)
	if err != nil {
		return nil, nil, err
	}

	offers := make([]OfferResponse, 0, len(offersWithUnits))
	for _, item := range offersWithUnits {
		offer := item.Offer
		accessibility := []string(offer.Accessibility)
		if accessibility == nil {
			accessibility = []string{}
		}

		offers = append(offers, OfferResponse{
			OfferID:      offer.ExternalOfferID,
			ProviderID:   s.providerID,
			ProviderType: "TRANSPORT",
			Status:       "AVAILABLE",
			ServiceWindow: platform.TimeWindow{
				StartsAt:      platform.FormatUTC(offer.ServiceStartsAt),
				EndsAt:        platform.FormatUTC(offer.ServiceEndsAt),
				StartTimeZone: offer.StartTimeZone,
				EndTimeZone:   offer.EndTimeZone,
			},
			AvailableUnits: item.AvailableUnits,
			UnitPrice: platform.Money{
				AmountMinor: offer.PriceAmountMinor,
				Currency:    offer.PriceCurrency,
			},
			ValidUntil: platform.FormatUTC(offer.ValidUntil),
			Synthetic:  true,
			Source:     "MOCK",
			Details: TransportOfferDetails{
				ProviderType:        "TRANSPORT",
				AvailabilityID:      offer.ExternalAvailabilityID,
				VehicleType:         offer.VehicleType,
				PickupLocationCode:  offer.PickupLocationCode,
				DropoffLocationCode: offer.DropoffLocationCode,
				PassengerCapacity:   offer.PassengerCapacity,
				Accessibility:       accessibility,
			},
		})
	}

	return &SearchTransportResponse{
		ProviderID:   s.providerID,
		ProviderType: "TRANSPORT",
		Offers:       offers,
	}, nil, nil
}

func (s *Service) CreateHold(
	ctx context.Context,
	req CreateHoldRequestPayload,
	idempotencyKey string,
	requestFingerprint string,
	now time.Time,
) (*HoldResponse, bool, *CapacityConflictDetails, *OfferExpiredDetails, []platform.ErrorDetail, error) {
	if req.ProviderID != s.providerID || req.ProviderType != "TRANSPORT" {
		return nil, false, nil, nil, nil, ErrProviderIdentityMismatch
	}

	var details []platform.ErrorDetail

	if req.BookingRequirements == nil {
		details = append(details, platform.ErrorDetail{
			Field:  "booking_requirements",
			Reason: "is required for TRANSPORT provider type",
		})
	} else {
		br := req.BookingRequirements
		if br.PassengerCount < 1 || br.PassengerCount > 20 {
			details = append(details, platform.ErrorDetail{
				Field:  "booking_requirements.passenger_count",
				Reason: "must be between 1 and 20",
			})
		}
		if !platform.ValidateLocationCode(br.PickupLocationCode) {
			details = append(details, platform.ErrorDetail{
				Field:  "booking_requirements.pickup_location_code",
				Reason: "must match ^[A-Z][A-Z0-9_]*$ (length 2..64)",
			})
		}
		if !platform.ValidateLocationCode(br.DropoffLocationCode) {
			details = append(details, platform.ErrorDetail{
				Field:  "booking_requirements.dropoff_location_code",
				Reason: "must match ^[A-Z][A-Z0-9_]*$ (length 2..64)",
			})
		}
		if br.PickupLocationCode != "" && br.DropoffLocationCode != "" && br.PickupLocationCode == br.DropoffLocationCode {
			details = append(details, platform.ErrorDetail{
				Field:  "booking_requirements.dropoff_location_code",
				Reason: "pickup and dropoff locations must be different",
			})
		}
		windowDetails := platform.ValidateTimeWindow(br.PickupWindow, "booking_requirements.pickup_window")
		details = append(details, windowDetails...)

		if len(br.Accessibility) > 0 {
			details = append(details, platform.ValidateAccessibility(br.Accessibility, "booking_requirements.accessibility")...)
		}
	}

	if req.Units < 1 || req.Units > 20 {
		details = append(details, platform.ErrorDetail{
			Field:  "units",
			Reason: "must be between 1 and 20",
		})
	}
	if !platform.ValidateResourceId(req.OfferID) {
		details = append(details, platform.ErrorDetail{
			Field:  "offer_id",
			Reason: "is invalid or missing",
		})
	}
	if !platform.ValidateClientReference(req.ClientReference) {
		details = append(details, platform.ErrorDetail{
			Field:  "client_reference",
			Reason: "is invalid or missing",
		})
	}

	moneyDetails := platform.ValidateMoney(req.ExpectedUnitPrice, "expected_unit_price")
	details = append(details, moneyDetails...)

	if len(details) > 0 {
		return nil, false, nil, nil, details, ErrInvalidRequest
	}

	br := req.BookingRequirements
	reqPickupStart, _ := time.Parse(time.RFC3339, br.PickupWindow.StartsAt)
	reqPickupEnd, _ := time.Parse(time.RFC3339, br.PickupWindow.EndsAt)

	result, err := s.repo.CreateHoldTx(ctx, CreateHoldParams{
		ClientScope:             s.providerID,
		Operation:               "POST /v1/holds",
		IdempotencyKey:          idempotencyKey,
		RequestFingerprint:      requestFingerprint,
		OfferID:                 req.OfferID,
		Units:                   req.Units,
		ExpectedUnitPrice:       req.ExpectedUnitPrice.AmountMinor,
		ExpectedCurrency:        req.ExpectedUnitPrice.Currency,
		ClientReference:         req.ClientReference,
		PassengerCount:          br.PassengerCount,
		PickupLocationCode:      br.PickupLocationCode,
		DropoffLocationCode:     br.DropoffLocationCode,
		RequestedPickupStartsAt: reqPickupStart,
		RequestedPickupEndsAt:   reqPickupEnd,
		RequestedStartTimeZone:  br.PickupWindow.StartTimeZone,
		RequestedEndTimeZone:    br.PickupWindow.EndTimeZone,
		Accessibility:           br.Accessibility,
		HoldDuration:            10 * time.Minute,
		Now:                     now,
	})
	if err != nil {
		if errors.Is(err, ErrCapacityConflict) && result != nil {
			return nil, false, result.ConflictDetails, nil, nil, err
		}
		if errors.Is(err, ErrOfferExpired) && result != nil {
			return nil, false, nil, result.OfferExpiredDetails, nil, err
		}
		return nil, false, nil, nil, nil, err
	}

	hold := result.Hold
	resp := &HoldResponse{
		HoldID:            hold.ExternalHoldID,
		ExternalReference: hold.ExternalReference,
		ProviderID:        s.providerID,
		ProviderType:      "TRANSPORT",
		OfferID:           hold.OfferID,
		ClientReference:   hold.ClientReference,
		Status:            "HELD",
		Units:             hold.Units,
		UnitPrice: platform.Money{
			AmountMinor: hold.UnitPriceAmountMinor,
			Currency:    hold.PriceCurrency,
		},
		TotalPrice: platform.Money{
			AmountMinor: hold.TotalPriceAmountMinor,
			Currency:    hold.PriceCurrency,
		},
		ServiceWindow: platform.TimeWindow{
			StartsAt:      platform.FormatUTC(hold.ServiceStartsAt),
			EndsAt:        platform.FormatUTC(hold.ServiceEndsAt),
			StartTimeZone: hold.StartTimeZone,
			EndTimeZone:   hold.EndTimeZone,
		},
		CreatedAt: platform.FormatUTC(hold.CreatedAt),
		ExpiresAt: platform.FormatUTC(hold.ExpiresAt),
	}

	return resp, result.Replayed, nil, nil, nil, nil
}

func (s *Service) ConfirmHold(
	ctx context.Context,
	externalHoldID string,
	idempotencyKey string,
	requestFingerprint string,
	now time.Time,
) (*ReservationResponse, bool, *HoldExpiredDetails, error) {
	result, err := s.repo.ConfirmHoldTx(ctx, ConfirmHoldParams{
		ClientScope:        s.providerID,
		Operation:          fmt.Sprintf("POST /v1/holds/%s/confirm", externalHoldID),
		IdempotencyKey:     idempotencyKey,
		RequestFingerprint: requestFingerprint,
		ExternalHoldID:     externalHoldID,
		Now:                now,
	})
	if err != nil {
		if errors.Is(err, ErrHoldExpired) && result != nil {
			return nil, false, result.HoldExpiredDetails, err
		}
		return nil, false, nil, err
	}

	res := result.Reservation
	resp := s.mapReservationResponse(res, externalHoldID)
	return resp, result.Replayed, nil, nil
}

func (s *Service) ReleaseHold(
	ctx context.Context,
	externalHoldID string,
	idempotencyKey string,
	requestFingerprint string,
	now time.Time,
) (*platform.ReleaseResult, bool, error) {
	result, err := s.repo.ReleaseHoldTx(ctx, ReleaseHoldParams{
		ClientScope:        s.providerID,
		Operation:          fmt.Sprintf("POST /v1/holds/%s/release", externalHoldID),
		IdempotencyKey:     idempotencyKey,
		RequestFingerprint: requestFingerprint,
		ExternalHoldID:     externalHoldID,
		Now:                now,
	})
	if err != nil {
		return nil, false, err
	}

	return &platform.ReleaseResult{
		ResourceType:      "HOLD",
		ResourceID:        result.Hold.ExternalHoldID,
		ExternalReference: result.Hold.ExternalReference,
		ProviderID:        s.providerID,
		ProviderType:      "TRANSPORT",
		Status:            "RELEASED",
		ReleasedAt:        platform.FormatUTC(result.ReleasedAt),
	}, result.Replayed, nil
}

func (s *Service) GetReservation(ctx context.Context, externalReservationID string) (*ReservationResponse, error) {
	res, err := s.repo.GetReservation(ctx, externalReservationID)
	if err != nil {
		return nil, err
	}

	holdExternalID := ""
	if res.Hold != nil {
		holdExternalID = res.Hold.ExternalHoldID
	}

	return s.mapReservationResponse(res, holdExternalID), nil
}

func (s *Service) ReleaseReservation(
	ctx context.Context,
	externalReservationID string,
	idempotencyKey string,
	requestFingerprint string,
	now time.Time,
) (*platform.ReleaseResult, bool, error) {
	result, err := s.repo.ReleaseReservationTx(ctx, ReleaseReservationParams{
		ClientScope:           s.providerID,
		Operation:             fmt.Sprintf("POST /v1/reservations/%s/release", externalReservationID),
		IdempotencyKey:        idempotencyKey,
		RequestFingerprint:    requestFingerprint,
		ExternalReservationID: externalReservationID,
		Now:                   now,
	})
	if err != nil {
		return nil, false, err
	}

	return &platform.ReleaseResult{
		ResourceType:      "RESERVATION",
		ResourceID:        result.Reservation.ExternalReservationID,
		ExternalReference: result.Reservation.ExternalReference,
		ProviderID:        s.providerID,
		ProviderType:      "TRANSPORT",
		Status:            "RELEASED",
		ReleasedAt:        platform.FormatUTC(result.ReleasedAt),
	}, result.Replayed, nil
}

func (s *Service) mapReservationResponse(res *Reservation, holdExternalID string) *ReservationResponse {
	var releasedAtStr *string
	if res.ReleasedAt != nil {
		str := platform.FormatUTC(*res.ReleasedAt)
		releasedAtStr = &str
	}

	return &ReservationResponse{
		ReservationID:     res.ExternalReservationID,
		ExternalReference: res.ExternalReference,
		HoldID:            holdExternalID,
		ProviderID:        s.providerID,
		ProviderType:      "TRANSPORT",
		OfferID:           res.OfferID,
		ClientReference:   res.ClientReference,
		Status:            res.Status,
		Units:             res.Units,
		TotalPrice: platform.Money{
			AmountMinor: res.TotalPriceAmountMinor,
			Currency:    res.PriceCurrency,
		},
		ServiceWindow: platform.TimeWindow{
			StartsAt:      platform.FormatUTC(res.ServiceStartsAt),
			EndsAt:        platform.FormatUTC(res.ServiceEndsAt),
			StartTimeZone: res.StartTimeZone,
			EndTimeZone:   res.EndTimeZone,
		},
		ConfirmedAt: platform.FormatUTC(res.ConfirmedAt),
		ReleasedAt:  releasedAtStr,
	}
}
