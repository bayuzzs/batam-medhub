package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"batam-medhub/internal/adapter"
	"batam-medhub/internal/auth"
	"batam-medhub/internal/model"

	"gorm.io/gorm"
)

var (
	ErrPlanOptionNotFound   = errors.New("plan option not found")
	ErrInvalidPlanOption    = errors.New("plan option is not in a confirmable state")
	ErrBookingHoldFailed    = errors.New("provider hold creation failed")
	ErrBookingConfirmFailed = errors.New("provider hold confirmation failed")
	ErrApprovalRequired     = errors.New("explicit approval is required")
)

// BookingSagaService coordinates transactional multi-provider holds, confirmations, and compensation releases.
type BookingSagaService struct {
	db    *gorm.DB
	hosp  *adapter.HospitalAdapter
	ferry *adapter.FerryAdapter
	hotel *adapter.HotelAdapter
	trans *adapter.TransportAdapter
	money *MoneyService
}

// NewBookingSagaService constructs a BookingSagaService.
func NewBookingSagaService(
	db *gorm.DB,
	hosp *adapter.HospitalAdapter,
	ferry *adapter.FerryAdapter,
	hotel *adapter.HotelAdapter,
	trans *adapter.TransportAdapter,
	money *MoneyService,
) *BookingSagaService {
	return &BookingSagaService{
		db:    db,
		hosp:  hosp,
		ferry: ferry,
		hotel: hotel,
		trans: trans,
		money: money,
	}
}

type holdRecord struct {
	planItem    model.PlanItem
	provider    model.Provider
	reservation model.Reservation
	hold        *adapter.Hold
}

type confirmRecord struct {
	holdRec     holdRecord
	reservation *adapter.Reservation
}

// ConfirmPlanOption executes the booking saga for a chosen plan option.
func (s *BookingSagaService) ConfirmPlanOption(ctx context.Context, patientID, planOptionID, reqID, idemKey string) (*JourneyDetail, error) {
	// 1. Fetch PlanOption and verify ownership via TripRequest
	var planOption model.PlanOption
	if err := s.db.WithContext(ctx).Where("id = ?", planOptionID).First(&planOption).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPlanOptionNotFound
		}
		return nil, fmt.Errorf("lookup plan option: %w", err)
	}

	var trip model.TripRequest
	if err := s.db.WithContext(ctx).Where("id = ? AND patient_id = ?", planOption.TripRequestID, patientID).First(&trip).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPlanOptionNotFound
		}
		return nil, fmt.Errorf("lookup trip request: %w", err)
	}

	if planOption.Status != "PROPOSED" {
		return nil, fmt.Errorf("%w: plan option status is %s", ErrInvalidPlanOption, planOption.Status)
	}
	if trip.Status != "PLAN_READY" && trip.Status != "PLANNING" && trip.Status != "CONFIRMING" {
		return nil, fmt.Errorf("%w: trip request status is %s", ErrInvalidTripState, trip.Status)
	}

	// Update TripRequest to CONFIRMING state
	trip.Status = "CONFIRMING"
	trip.UpdatedAt = time.Now().UTC()
	_ = s.db.WithContext(ctx).Save(&trip).Error

	// Extract intent details for requirements
	var intent StructuredIntent
	if len(trip.StructuredIntent) > 0 {
		_ = json.Unmarshal(trip.StructuredIntent, &intent)
	}
	passengerCount := 1
	if intent.PatientCount != nil && *intent.PatientCount > 0 {
		passengerCount = *intent.PatientCount
	}
	if intent.CompanionCount != nil && *intent.CompanionCount > 0 {
		passengerCount += *intent.CompanionCount
	}
	accessibility := intent.Preferences.Accessibility
	if accessibility == nil {
		accessibility = []string{}
	}

	// 2. Load all plan items
	var planItems []model.PlanItem
	if err := s.db.WithContext(ctx).
		Where("plan_option_id = ?", planOption.ID).
		Order("sequence_number ASC").
		Find(&planItems).Error; err != nil {
		return nil, fmt.Errorf("load plan items: %w", err)
	}

	if len(planItems) == 0 {
		return nil, fmt.Errorf("%w: plan option has no items", ErrInvalidPlanOption)
	}

	// Filter bookable items
	var bookableItems []model.PlanItem
	for _, it := range planItems {
		if it.ProviderID != nil && it.ExternalOfferID != nil {
			bookableItems = append(bookableItems, it)
		}
	}

	// 3. Holds Phase
	var heldRecords []holdRecord
	var holdPhaseErr error

	for _, item := range bookableItems {
		var prov model.Provider
		if err := s.db.WithContext(ctx).Where("id = ?", *item.ProviderID).First(&prov).Error; err != nil {
			holdPhaseErr = fmt.Errorf("lookup provider %s: %w", *item.ProviderID, err)
			break
		}

		resID := auth.NewUUID()
		now := time.Now().UTC()
		res := model.Reservation{
			ID:               resID,
			TripRequestID:    trip.ID,
			PlanItemID:       &item.ID,
			ProviderID:       prov.ID,
			Status:           "PENDING",
			ExternalOfferID:  *item.ExternalOfferID,
			ProviderSnapshot: []byte(`{}`),
			CleanupRequired:  false,
			CreatedAt:        now,
			UpdatedAt:        now,
		}

		if err := s.db.WithContext(ctx).Create(&res).Error; err != nil {
			holdPhaseErr = fmt.Errorf("create pending reservation: %w", err)
			break
		}

		var minor int64 = 0
		if item.SourceAmountMinor != nil {
			minor = *item.SourceAmountMinor
		}
		var curr string = "IDR"
		if item.SourceCurrency != nil {
			curr = *item.SourceCurrency
		}

		units := 1
		unitPriceMinor := minor
		switch prov.ProviderType {
		case adapter.ProviderTypeFerry:
			units = passengerCount
			if units > 0 {
				unitPriceMinor = minor / int64(units)
			}
		case adapter.ProviderTypeHospital:
			units = 1
			if intent.PatientCount != nil && *intent.PatientCount > 0 {
				units = *intent.PatientCount
			}
			if units > 0 {
				unitPriceMinor = minor / int64(units)
			}
		case adapter.ProviderTypeTransport:
			units = 1
			unitPriceMinor = minor
		case adapter.ProviderTypeHotel:
			units = 1
			unitPriceMinor = minor
		}

		createHoldReq := adapter.CreateHoldRequest{
			ProviderID:        prov.Code,
			ProviderType:      prov.ProviderType,
			OfferID:           *item.ExternalOfferID,
			Units:             units,
			ExpectedUnitPrice: adapter.Money{AmountMinor: unitPriceMinor, Currency: curr},
			ClientReference:   fmt.Sprintf("trip-%s-item-%d", trip.ID, item.SequenceNumber),
		}

		if prov.ProviderType == adapter.ProviderTypeTransport {
			createHoldReq.BookingRequirements = &adapter.TransportBookingRequirements{
				PassengerCount:      passengerCount,
				PickupLocationCode:  safeString(item.OriginCode),
				DropoffLocationCode: safeString(item.DestinationCode),
				PickupWindow: adapter.TimeWindow{
					StartsAt:      item.StartsAt.UTC().Format(time.RFC3339),
					EndsAt:        item.EndsAt.UTC().Format(time.RFC3339),
					StartTimeZone: item.StartTimeZone,
					EndTimeZone:   item.EndTimeZone,
				},
				Accessibility: accessibility,
			}
		}

		itemHoldIdemKey := fmt.Sprintf("hold-%s-%s", planOption.ID, item.ID)
		if idemKey != "" {
			itemHoldIdemKey = fmt.Sprintf("hold-%s-%s-%s", idemKey, planOption.ID, item.ID)
		}

		holdResult, err := s.createProviderHold(ctx, prov.ProviderType, reqID, itemHoldIdemKey, createHoldReq)
		if err != nil {
			slog.Warn("provider hold failed", "provider_type", prov.ProviderType, "offer_id", *item.ExternalOfferID, "error", err)
			res.Status = "FAILED"
			res.UpdatedAt = time.Now().UTC()
			_ = s.db.WithContext(ctx).Save(&res).Error
			holdPhaseErr = fmt.Errorf("%w: provider %s hold failed: %v", ErrBookingHoldFailed, prov.ProviderType, err)
			break
		}

		holdExpiry, _ := time.Parse(time.RFC3339, holdResult.ExpiresAt)
		snapBytes, _ := json.Marshal(holdResult)

		res.Status = "HELD"
		res.ExternalHoldID = &holdResult.HoldID
		res.HoldExpiresAt = &holdExpiry
		res.ProviderSnapshot = snapBytes
		res.UpdatedAt = time.Now().UTC()
		_ = s.db.WithContext(ctx).Save(&res).Error

		heldRecords = append(heldRecords, holdRecord{
			planItem:    item,
			provider:    prov,
			reservation: res,
			hold:        holdResult,
		})
	}

	// If hold phase failed, trigger compensation on held records
	if holdPhaseErr != nil {
		uncertain := s.compensateHolds(ctx, heldRecords, reqID)
		if uncertain {
			trip.Status = "MANUAL_REVIEW"
		} else {
			trip.Status = "CONFIRMATION_FAILED"
		}
		trip.UpdatedAt = time.Now().UTC()
		_ = s.db.WithContext(ctx).Save(&trip).Error
		return nil, holdPhaseErr
	}

	// 4. Confirmation Phase
	var confirmedRecords []confirmRecord
	var confirmPhaseErr error

	for _, hr := range heldRecords {
		itemConfIdemKey := fmt.Sprintf("conf-%s-%s", planOption.ID, hr.planItem.ID)
		if idemKey != "" {
			itemConfIdemKey = fmt.Sprintf("conf-%s-%s-%s", idemKey, planOption.ID, hr.planItem.ID)
		}

		confResult, err := s.confirmProviderHold(ctx, hr.provider.ProviderType, reqID, itemConfIdemKey, hr.hold.HoldID)
		if err != nil {
			slog.Warn("provider confirm failed", "provider_type", hr.provider.ProviderType, "hold_id", hr.hold.HoldID, "error", err)
			hr.reservation.Status = "FAILED"
			hr.reservation.UpdatedAt = time.Now().UTC()
			_ = s.db.WithContext(ctx).Save(&hr.reservation).Error
			confirmPhaseErr = fmt.Errorf("%w: provider %s confirm failed: %v", ErrBookingConfirmFailed, hr.provider.ProviderType, err)
			break
		}

		snapBytes, _ := json.Marshal(confResult)
		hr.reservation.Status = "CONFIRMED"
		hr.reservation.ExternalReservationID = &confResult.ReservationID
		hr.reservation.ProviderSnapshot = snapBytes
		hr.reservation.UpdatedAt = time.Now().UTC()
		_ = s.db.WithContext(ctx).Save(&hr.reservation).Error

		confirmedRecords = append(confirmedRecords, confirmRecord{
			holdRec:     hr,
			reservation: confResult,
		})
	}

	// If confirm phase failed, trigger compensation on confirmed & unconfirmed records
	if confirmPhaseErr != nil {
		uncertain := s.compensateConfirmations(ctx, confirmedRecords, heldRecords, reqID)
		if uncertain {
			trip.Status = "MANUAL_REVIEW"
		} else {
			trip.Status = "CONFIRMATION_FAILED"
		}
		trip.UpdatedAt = time.Now().UTC()
		_ = s.db.WithContext(ctx).Save(&trip).Error
		return nil, confirmPhaseErr
	}

	// 5. Full Success: Transactional Journey & Itinerary Creation
	journeyID := auth.NewUUID()
	itineraryVersionID := auth.NewUUID()
	now := time.Now().UTC()

	var createdJourney model.Journey
	var createdVersion model.ItineraryVersion

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update all reservations with journey_id
		for _, cr := range confirmedRecords {
			if err := tx.Model(&model.Reservation{}).
				Where("id = ?", cr.holdRec.reservation.ID).
				Updates(map[string]any{
					"journey_id": journeyID,
					"updated_at": now,
				}).Error; err != nil {
				return err
			}
		}

		// Create Journey
		createdJourney = model.Journey{
			ID:                   journeyID,
			TripRequestID:        trip.ID,
			PatientID:            patientID,
			Status:               "ACTIVE",
			CurrentVersionNumber: 1,
			ActivatedAt:          now,
			CreatedAt:            now,
			UpdatedAt:            now,
		}
		if err := tx.Create(&createdJourney).Error; err != nil {
			return err
		}

		// Create ItineraryVersion in DRAFT state first (items can only be inserted in DRAFT state)
		createdVersion = model.ItineraryVersion{
			ID:                 itineraryVersionID,
			JourneyID:          journeyID,
			VersionNumber:      1,
			Status:             "DRAFT",
			ChangeReason:       "Initial confirmed itinerary version",
			SourceDisruptionID: nil,
			TotalPriceSnapshot: planOption.TotalPriceSnapshot,
			ActivatedAt:        nil,
			CreatedAt:          now,
		}
		if err := tx.Create(&createdVersion).Error; err != nil {
			return err
		}

		// Map plan items to itinerary items
		for _, pi := range planItems {
			itItemID := auth.NewUUID()
			var resID *string
			var extResID *string
			var provID *string
			status := "BUFFER"

			for _, cr := range confirmedRecords {
				if cr.holdRec.planItem.ID == pi.ID {
					resID = &cr.holdRec.reservation.ID
					extResID = &cr.reservation.ReservationID
					provID = pi.ProviderID
					status = "CONFIRMED"
					break
				}
			}

			// Preserve operational notes inside snapshot
			snapData := map[string]any{
				"operational_notes": []string{},
			}
			if len(pi.OperationalNotes) > 0 {
				var notes []string
				if err := json.Unmarshal(pi.OperationalNotes, &notes); err == nil {
					snapData["operational_notes"] = notes
				}
			}
			snapBytes, _ := json.Marshal(snapData)

			var srcMinor, dispMinor *int64
			var srcCurr, dispCurr *string
			var fxRateVal, fxSrc *string
			var fxEff *time.Time

			if status == "CONFIRMED" && pi.SourceAmountMinor != nil {
				srcMinor = pi.SourceAmountMinor
				srcCurr = pi.SourceCurrency
				dispMinor = pi.DisplayAmountMinor
				dispCurr = pi.DisplayCurrency

				// Populate FX metadata from snapshot or defaults
				fxRate := "1.000000"
				fxSourceStr := "DEMO_STATIC_2026_08"
				if len(pi.OfferSnapshot) > 0 {
					var rawPrice struct {
						FXRate        string    `json:"fx_rate"`
						FXSource      string    `json:"fx_source"`
						FXEffectiveAt time.Time `json:"fx_effective_at"`
					}
					if err := json.Unmarshal(pi.OfferSnapshot, &rawPrice); err == nil && rawPrice.FXRate != "" {
						fxRate = rawPrice.FXRate
						if rawPrice.FXSource != "" {
							fxSourceStr = rawPrice.FXSource
						}
					}
				}
				fxRateVal = &fxRate
				fxSrc = &fxSourceStr
				fxEff = &now
			}

			itItem := model.ItineraryItem{
				ID:                    itItemID,
				ItineraryVersionID:    itineraryVersionID,
				ReservationID:         resID,
				ProviderID:            provID,
				ItemType:              pi.ItemType,
				SequenceNumber:        pi.SequenceNumber,
				ExternalReservationID: extResID,
				Title:                 pi.Title,
				Status:                status,
				StartsAt:              pi.StartsAt,
				EndsAt:                pi.EndsAt,
				StartTimeZone:         pi.StartTimeZone,
				EndTimeZone:           pi.EndTimeZone,
				OriginCode:            pi.OriginCode,
				DestinationCode:       pi.DestinationCode,
				SourceAmountMinor:     srcMinor,
				SourceCurrency:        srcCurr,
				DisplayAmountMinor:    dispMinor,
				DisplayCurrency:       dispCurr,
				FXRateValue:           fxRateVal,
				FXSource:              fxSrc,
				FXEffectiveAt:         fxEff,
				Snapshot:              snapBytes,
				Synthetic:             true,
				Source:                "MOCK",
				CreatedAt:             now,
			}

			if err := tx.Create(&itItem).Error; err != nil {
				return err
			}
		}

		// Activate ItineraryVersion (transition DRAFT -> ACTIVE)
		if err := tx.Model(&model.ItineraryVersion{}).
			Where("id = ?", itineraryVersionID).
			Updates(map[string]any{
				"status":       "ACTIVE",
				"activated_at": now,
			}).Error; err != nil {
			return err
		}

		// Update PlanOption
		if err := tx.Model(&model.PlanOption{}).
			Where("id = ?", planOption.ID).
			Updates(map[string]any{
				"status":     "CONFIRMED",
				"updated_at": now,
			}).Error; err != nil {
			return err
		}

		// Update TripRequest
		trip.Status = "ACTIVE"
		trip.SelectedPlanOptionID = &planOption.ID
		trip.UpdatedAt = now
		return tx.Save(&trip).Error
	})

	if err != nil {
		return nil, fmt.Errorf("transactionally persist journey: %w", err)
	}

	// 6. Return loaded JourneyDetail
	return s.loadJourneyDetail(ctx, journeyID)
}

func (s *BookingSagaService) createProviderHold(
	ctx context.Context,
	provType, reqID, idemKey string,
	req adapter.CreateHoldRequest,
) (*adapter.Hold, error) {
	switch provType {
	case adapter.ProviderTypeHospital:
		return s.hosp.CreateHold(ctx, reqID, idemKey, req)
	case adapter.ProviderTypeFerry:
		return s.ferry.CreateHold(ctx, reqID, idemKey, req)
	case adapter.ProviderTypeHotel:
		return s.hotel.CreateHold(ctx, reqID, idemKey, req)
	case adapter.ProviderTypeTransport:
		return s.trans.CreateHold(ctx, reqID, idemKey, req)
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", provType)
	}
}

func (s *BookingSagaService) confirmProviderHold(
	ctx context.Context,
	provType, reqID, idemKey, holdID string,
) (*adapter.Reservation, error) {
	switch provType {
	case adapter.ProviderTypeHospital:
		return s.hosp.ConfirmHold(ctx, reqID, idemKey, holdID)
	case adapter.ProviderTypeFerry:
		return s.ferry.ConfirmHold(ctx, reqID, idemKey, holdID)
	case adapter.ProviderTypeHotel:
		return s.hotel.ConfirmHold(ctx, reqID, idemKey, holdID)
	case adapter.ProviderTypeTransport:
		return s.trans.ConfirmHold(ctx, reqID, idemKey, holdID)
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", provType)
	}
}

func (s *BookingSagaService) compensateHolds(ctx context.Context, heldRecords []holdRecord, reqID string) bool {
	uncertain := false
	for _, hr := range heldRecords {
		idemKey := fmt.Sprintf("rel-hold-%s", hr.hold.HoldID)
		var relErr error
		switch hr.provider.ProviderType {
		case adapter.ProviderTypeHospital:
			_, relErr = s.hosp.ReleaseHold(ctx, reqID, idemKey, hr.hold.HoldID)
		case adapter.ProviderTypeFerry:
			_, relErr = s.ferry.ReleaseHold(ctx, reqID, idemKey, hr.hold.HoldID)
		case adapter.ProviderTypeHotel:
			_, relErr = s.hotel.ReleaseHold(ctx, reqID, idemKey, hr.hold.HoldID)
		case adapter.ProviderTypeTransport:
			_, relErr = s.trans.ReleaseHold(ctx, reqID, idemKey, hr.hold.HoldID)
		}

		if relErr != nil {
			slog.Error("compensation release hold failed", "hold_id", hr.hold.HoldID, "error", relErr)
			hr.reservation.CleanupRequired = true
			uncertain = true
		} else {
			hr.reservation.Status = "RELEASED"
		}
		hr.reservation.UpdatedAt = time.Now().UTC()
		_ = s.db.WithContext(ctx).Save(&hr.reservation).Error
	}
	return uncertain
}

func (s *BookingSagaService) compensateConfirmations(
	ctx context.Context,
	confirmed []confirmRecord,
	allHeld []holdRecord,
	reqID string,
) bool {
	uncertain := false

	// Compensate confirmed reservations
	for _, cr := range confirmed {
		idemKey := fmt.Sprintf("rel-res-%s", cr.reservation.ReservationID)
		var relErr error
		switch cr.holdRec.provider.ProviderType {
		case adapter.ProviderTypeHospital:
			_, relErr = s.hosp.ReleaseReservation(ctx, reqID, idemKey, cr.reservation.ReservationID)
		case adapter.ProviderTypeFerry:
			_, relErr = s.ferry.ReleaseReservation(ctx, reqID, idemKey, cr.reservation.ReservationID)
		case adapter.ProviderTypeHotel:
			_, relErr = s.hotel.ReleaseReservation(ctx, reqID, idemKey, cr.reservation.ReservationID)
		case adapter.ProviderTypeTransport:
			_, relErr = s.trans.ReleaseReservation(ctx, reqID, idemKey, cr.reservation.ReservationID)
		}

		if relErr != nil {
			slog.Error("compensation release reservation failed", "res_id", cr.reservation.ReservationID, "error", relErr)
			cr.holdRec.reservation.CleanupRequired = true
			uncertain = true
		} else {
			cr.holdRec.reservation.Status = "RELEASED"
		}
		cr.holdRec.reservation.UpdatedAt = time.Now().UTC()
		_ = s.db.WithContext(ctx).Save(&cr.holdRec.reservation).Error
	}

	// Compensate unconfirmed holds
	for _, hr := range allHeld {
		isConfirmed := false
		for _, cr := range confirmed {
			if cr.holdRec.reservation.ID == hr.reservation.ID {
				isConfirmed = true
				break
			}
		}
		if isConfirmed {
			continue
		}

		idemKey := fmt.Sprintf("rel-hold-%s", hr.hold.HoldID)
		var relErr error
		switch hr.provider.ProviderType {
		case adapter.ProviderTypeHospital:
			_, relErr = s.hosp.ReleaseHold(ctx, reqID, idemKey, hr.hold.HoldID)
		case adapter.ProviderTypeFerry:
			_, relErr = s.ferry.ReleaseHold(ctx, reqID, idemKey, hr.hold.HoldID)
		case adapter.ProviderTypeHotel:
			_, relErr = s.hotel.ReleaseHold(ctx, reqID, idemKey, hr.hold.HoldID)
		case adapter.ProviderTypeTransport:
			_, relErr = s.trans.ReleaseHold(ctx, reqID, idemKey, hr.hold.HoldID)
		}

		if relErr != nil {
			slog.Error("compensation release hold failed", "hold_id", hr.hold.HoldID, "error", relErr)
			hr.reservation.CleanupRequired = true
			uncertain = true
		} else {
			hr.reservation.Status = "RELEASED"
		}
		hr.reservation.UpdatedAt = time.Now().UTC()
		_ = s.db.WithContext(ctx).Save(&hr.reservation).Error
	}

	return uncertain
}

func (s *BookingSagaService) loadJourneyDetail(ctx context.Context, journeyID string) (*JourneyDetail, error) {
	journeySvc := NewJourneyService(s.db)
	var j model.Journey
	if err := s.db.WithContext(ctx).Where("id = ?", journeyID).First(&j).Error; err != nil {
		return nil, err
	}
	return journeySvc.GetJourneyItinerary(ctx, j.PatientID, journeyID)
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
