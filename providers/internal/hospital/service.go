package hospital

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"batam-medhub/providers/internal/platform"
)

var (
	ErrProviderIdentityMismatch = errors.New("provider identity mismatch")
	ErrInvalidRequest           = errors.New("invalid request")
	ianaTimeZoneRegex           = regexp.MustCompile(`^[A-Za-z_]+(/[A-Za-z0-9._+-]+)+$`)
	serviceCodeRegex            = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
)

type HospitalOfferDetails struct {
	ProviderType      string   `json:"provider_type"`
	ServiceCode       string   `json:"service_code"`
	AppointmentSlotID string   `json:"appointment_slot_id"`
	FacilityName      string   `json:"facility_name"`
	DurationMinutes   int      `json:"duration_minutes"`
	Languages         []string `json:"languages"`
	Accessibility     []string `json:"accessibility"`
}

type OfferResponse struct {
	OfferID        string               `json:"offer_id"`
	ProviderID     string               `json:"provider_id"`
	ProviderType   string               `json:"provider_type"`
	Status         string               `json:"status"`
	ServiceWindow  platform.TimeWindow  `json:"service_window"`
	AvailableUnits int                  `json:"available_units"`
	UnitPrice      platform.Money       `json:"unit_price"`
	ValidUntil     string               `json:"valid_until"`
	Synthetic      bool                 `json:"synthetic"`
	Source         string               `json:"source"`
	Details        HospitalOfferDetails `json:"details"`
}

type SearchHospitalResponse struct {
	ProviderID   string          `json:"provider_id"`
	ProviderType string          `json:"provider_type"`
	Offers       []OfferResponse `json:"offers"`
}

type HospitalSearchCriteria struct {
	ServiceCode       string              `json:"service_code"`
	PatientCount      int                 `json:"patient_count"`
	AppointmentWindow platform.TimeWindow `json:"appointment_window"`
	Accessibility     []string            `json:"accessibility"`
}

type SearchRequestPayload struct {
	ProviderType string                 `json:"provider_type"`
	Criteria     HospitalSearchCriteria `json:"criteria"`
}

type CreateHoldRequestPayload struct {
	ProviderID          string           `json:"provider_id"`
	ProviderType        string           `json:"provider_type"`
	OfferID             string           `json:"offer_id"`
	Units               int              `json:"units"`
	ExpectedUnitPrice   platform.Money   `json:"expected_unit_price"`
	ClientReference     string           `json:"client_reference"`
	BookingRequirements *json.RawMessage `json:"booking_requirements,omitempty"`
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

func (s *Service) SearchOffers(ctx context.Context, req SearchRequestPayload, now time.Time) (*SearchHospitalResponse, []platform.ErrorDetail, error) {
	if req.ProviderType != "HOSPITAL" {
		return nil, []platform.ErrorDetail{{Field: "provider_type", Reason: "must be HOSPITAL"}}, ErrInvalidRequest
	}

	crit := req.Criteria
	var details []platform.ErrorDetail

	if !platform.ValidateServiceCode(crit.ServiceCode) {
		details = append(details, platform.ErrorDetail{Field: "criteria.service_code", Reason: "is invalid or missing"})
	}
	if crit.PatientCount < 1 || crit.PatientCount > 20 {
		details = append(details, platform.ErrorDetail{Field: "criteria.patient_count", Reason: "must be between 1 and 20"})
	}
	if !platform.ValidateIanaTimezone(crit.AppointmentWindow.StartTimeZone) {
		details = append(details, platform.ErrorDetail{Field: "criteria.appointment_window.start_time_zone", Reason: "must be a valid IANA time zone"})
	}
	if !platform.ValidateIanaTimezone(crit.AppointmentWindow.EndTimeZone) {
		details = append(details, platform.ErrorDetail{Field: "criteria.appointment_window.end_time_zone", Reason: "must be a valid IANA time zone"})
	}

	if !platform.ValidateRFC3339UTC(crit.AppointmentWindow.StartsAt) {
		details = append(details, platform.ErrorDetail{Field: "criteria.appointment_window.starts_at", Reason: "must be an RFC 3339 UTC timestamp ending in 'Z'"})
	}
	if !platform.ValidateRFC3339UTC(crit.AppointmentWindow.EndsAt) {
		details = append(details, platform.ErrorDetail{Field: "criteria.appointment_window.ends_at", Reason: "must be an RFC 3339 UTC timestamp ending in 'Z'"})
	}

	if len(crit.Accessibility) > 0 {
		details = append(details, platform.ValidateAccessibility(crit.Accessibility, "criteria.accessibility")...)
	}

	if len(details) > 0 {
		return nil, details, ErrInvalidRequest
	}

	windowStart, _ := time.Parse(time.RFC3339, crit.AppointmentWindow.StartsAt)
	windowEnd, _ := time.Parse(time.RFC3339, crit.AppointmentWindow.EndsAt)
	if !windowStart.Before(windowEnd) {
		return nil, []platform.ErrorDetail{{Field: "criteria.appointment_window", Reason: "starts_at must precede ends_at"}}, ErrInvalidRequest
	}

	service, err := s.repo.FindMedicalServiceByCode(ctx, crit.ServiceCode)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return &SearchHospitalResponse{
				ProviderID:   s.providerID,
				ProviderType: "HOSPITAL",
				Offers:       []OfferResponse{},
			}, nil, nil
		}
		return nil, nil, err
	}

	slotsWithUnits, err := s.repo.SearchAvailableSlots(
		ctx,
		service.ID,
		windowStart,
		windowEnd,
		now,
		crit.Accessibility,
		crit.PatientCount,
	)
	if err != nil {
		return nil, nil, err
	}

	offers := make([]OfferResponse, 0, len(slotsWithUnits))
	for _, item := range slotsWithUnits {
		slot := item.Slot
		languages := []string(slot.Languages)
		if languages == nil {
			languages = []string{}
		}
		accessibility := []string(slot.Accessibility)
		if accessibility == nil {
			accessibility = []string{}
		}

		offers = append(offers, OfferResponse{
			OfferID:      slot.OfferID,
			ProviderID:   s.providerID,
			ProviderType: "HOSPITAL",
			Status:       "AVAILABLE",
			ServiceWindow: platform.TimeWindow{
				StartsAt:      platform.FormatUTC(slot.StartsAt),
				EndsAt:        platform.FormatUTC(slot.EndsAt),
				StartTimeZone: slot.StartTimeZone,
				EndTimeZone:   slot.EndTimeZone,
			},
			AvailableUnits: item.AvailableUnits,
			UnitPrice: platform.Money{
				AmountMinor: slot.PriceAmountMinor,
				Currency:    slot.PriceCurrency,
			},
			ValidUntil: platform.FormatUTC(slot.ValidUntil),
			Synthetic:  true,
			Source:     "MOCK",
			Details: HospitalOfferDetails{
				ProviderType:      "HOSPITAL",
				ServiceCode:       service.Code,
				AppointmentSlotID: slot.ExternalSlotID,
				FacilityName:      slot.FacilityName,
				DurationMinutes:   slot.DurationMinutes,
				Languages:         languages,
				Accessibility:     accessibility,
			},
		})
	}

	return &SearchHospitalResponse{
		ProviderID:   s.providerID,
		ProviderType: "HOSPITAL",
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
	if req.ProviderID != s.providerID || req.ProviderType != "HOSPITAL" {
		return nil, false, nil, nil, nil, ErrProviderIdentityMismatch
	}

	var details []platform.ErrorDetail

	if req.BookingRequirements != nil && string(*req.BookingRequirements) != "null" && string(*req.BookingRequirements) != "" {
		details = append(details, platform.ErrorDetail{
			Field:  "booking_requirements",
			Reason: "is rejected for HOSPITAL provider type",
		})
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

	result, err := s.repo.CreateHoldTx(ctx, CreateHoldParams{
		ClientScope:        s.providerID,
		Operation:          "POST /v1/holds",
		IdempotencyKey:     idempotencyKey,
		RequestFingerprint: requestFingerprint,
		OfferID:            req.OfferID,
		Units:              req.Units,
		ExpectedUnitPrice:  req.ExpectedUnitPrice.AmountMinor,
		ExpectedCurrency:   req.ExpectedUnitPrice.Currency,
		ClientReference:    req.ClientReference,
		HoldDuration:       10 * time.Minute,
		Now:                now,
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
		ProviderType:      "HOSPITAL",
		OfferID:           hold.OfferID,
		ClientReference:   hold.ClientReference,
		Status:            "HELD",
		Units:             hold.PatientCount,
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
	var releasedAtStr *string
	if res.ReleasedAt != nil {
		str := platform.FormatUTC(*res.ReleasedAt)
		releasedAtStr = &str
	}

	resp := &ReservationResponse{
		ReservationID:     res.ExternalReservationID,
		ExternalReference: res.ExternalReference,
		HoldID:            externalHoldID,
		ProviderID:        s.providerID,
		ProviderType:      "HOSPITAL",
		OfferID:           res.OfferID,
		ClientReference:   res.ClientReference,
		Status:            res.Status,
		Units:             res.PatientCount,
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
		ProviderType:      "HOSPITAL",
		Status:            "RELEASED",
		ReleasedAt:        platform.FormatUTC(result.ReleasedAt),
	}, result.Replayed, nil
}

func (s *Service) GetReservation(ctx context.Context, externalReservationID string) (*ReservationResponse, error) {
	res, err := s.repo.GetReservation(ctx, externalReservationID)
	if err != nil {
		return nil, err
	}

	holdID := res.HoldID
	if res.Hold != nil && res.Hold.ExternalHoldID != "" {
		holdID = res.Hold.ExternalHoldID
	}

	var releasedAtStr *string
	if res.ReleasedAt != nil {
		str := platform.FormatUTC(*res.ReleasedAt)
		releasedAtStr = &str
	}

	return &ReservationResponse{
		ReservationID:     res.ExternalReservationID,
		ExternalReference: res.ExternalReference,
		HoldID:            holdID,
		ProviderID:        s.providerID,
		ProviderType:      "HOSPITAL",
		OfferID:           res.OfferID,
		ClientReference:   res.ClientReference,
		Status:            res.Status,
		Units:             res.PatientCount,
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
	}, nil
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
		ProviderType:      "HOSPITAL",
		Status:            "RELEASED",
		ReleasedAt:        platform.FormatUTC(result.ReleasedAt),
	}, result.Replayed, nil
}
