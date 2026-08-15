package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"batam-medhub/internal/adapter"
	"batam-medhub/internal/auth"
	"batam-medhub/internal/model"

	"gorm.io/gorm"
)

var (
	ErrNotFound               = errors.New("not found")
	ErrForbiddenEventType     = errors.New("event type is incompatible with provider type")
	ErrForbiddenTarget        = errors.New("target itinerary item or reservation does not belong to provider or journey")
	ErrEventDuplicateConflict = errors.New("duplicate external_event_id with conflicting request body")
	ErrDisruptionNotFound     = errors.New("disruption not found")
	ErrRecoveryOptionNotFound = errors.New("recovery option not found")
	ErrRecoveryOptionExpired  = errors.New("recovery option has expired")
	ErrRecoveryNotApplicable  = errors.New("recovery option is not in PROPOSED status")
)

// ProviderEventRequest represents the incoming payload from a provider submitting an operational disruption event.
type ProviderEventRequest struct {
	ExternalEventID string               `json:"external_event_id"`
	JourneyID       string               `json:"journey_id"`
	EventType       string               `json:"event_type"`
	OccurredAt      string               `json:"occurred_at"`
	Target          EventTargetSnapshot  `json:"target"`
	Actor           EventActorSnapshot   `json:"actor"`
	Details         EventDetailsSnapshot `json:"details"`
}

type EventTargetSnapshot struct {
	ItineraryItemID       *string `json:"itinerary_item_id"`
	ExternalReservationID *string `json:"external_reservation_id"`
}

type EventActorSnapshot struct {
	ActorID string `json:"actor_id"`
	Name    string `json:"name"`
	Role    string `json:"role"`
}

type EventTimeWindow struct {
	StartsAt      string `json:"starts_at"`
	EndsAt        string `json:"ends_at"`
	StartTimeZone string `json:"start_time_zone"`
	EndTimeZone   string `json:"end_time_zone"`
}

type EventDetailsSnapshot struct {
	Reason                    string           `json:"reason"`
	ReplacementTimeWindow     *EventTimeWindow `json:"replacement_time_window"`
	AdditionalServiceCode     *string          `json:"additional_service_code"`
	AdditionalDurationMinutes *int             `json:"additional_duration_minutes"`
	Priority                  *string          `json:"priority"`
	TravelClearanceStatus     *string          `json:"travel_clearance_status"`
	InstructionReference      *string          `json:"instruction_reference"`
	OperationalRequirements   []string         `json:"operational_requirements"`
}

type ProviderEventReceipt struct {
	ProviderEventID string  `json:"provider_event_id"`
	Outcome         string  `json:"outcome"`
	DisruptionID    *string `json:"disruption_id"`
	ReceivedAt      string  `json:"received_at"`
}

type ImpactSummary struct {
	Severity        string   `json:"severity"`
	AffectedItemIDs []string `json:"affected_item_ids"`
	Explanation     []string `json:"explanation"`
}

type DisruptionDTO struct {
	ID                       string        `json:"id"`
	JourneyID                string        `json:"journey_id"`
	ProviderEventID          string        `json:"provider_event_id"`
	AnalyzedItineraryVersion int           `json:"analyzed_itinerary_version"`
	Status                   string        `json:"status"`
	Impact                   ImpactSummary `json:"impact"`
	CreatedAt                time.Time     `json:"created_at"`
	UpdatedAt                time.Time     `json:"updated_at"`
}

type MoneyDeltaDTO struct {
	AmountMinor int64  `json:"amount_minor"`
	Currency    string `json:"currency"`
	Estimated   bool   `json:"estimated"`
}

type RecoveryChangeDTO struct {
	ChangeType         string       `json:"change_type"`
	OldItineraryItemID *string      `json:"old_itinerary_item_id"`
	ReplacementItem    *PlanItemDTO `json:"replacement_item"`
	Explanation        string       `json:"explanation"`
}

type RecoveryOptionDTO struct {
	ID               string              `json:"id"`
	DisruptionID     string              `json:"disruption_id"`
	Rank             int                 `json:"rank"`
	Status           string              `json:"status"`
	ExpiresAt        time.Time           `json:"expires_at"`
	Changes          []RecoveryChangeDTO `json:"changes"`
	TimeDeltaMinutes int                 `json:"time_delta_minutes"`
	PriceDelta       MoneyDeltaDTO       `json:"price_delta"`
}

type DisruptionDetail struct {
	Disruption      DisruptionDTO       `json:"disruption"`
	RecoveryOptions []RecoveryOptionDTO `json:"recovery_options"`
}

// DisruptionService coordinates provider event ingestion, impact analysis, recovery option generation, and recovery execution.
type DisruptionService struct {
	db           *gorm.DB
	hospAdapter  *adapter.HospitalAdapter
	ferryAdapter *adapter.FerryAdapter
	hotelAdapter *adapter.HotelAdapter
	transAdapter *adapter.TransportAdapter
	moneySvc     *MoneyService
	journeySvc   *JourneyService
}

// NewDisruptionService constructs a new DisruptionService.
func NewDisruptionService(
	db *gorm.DB,
	hosp *adapter.HospitalAdapter,
	ferry *adapter.FerryAdapter,
	hotel *adapter.HotelAdapter,
	trans *adapter.TransportAdapter,
	money *MoneyService,
	journey *JourneyService,
) *DisruptionService {
	return &DisruptionService{
		db:           db,
		hospAdapter:  hosp,
		ferryAdapter: ferry,
		hotelAdapter: hotel,
		transAdapter: trans,
		moneySvc:     money,
		journeySvc:   journey,
	}
}

// IngestProviderEvent processes a provider disruption event, performs deduplication, validates compatibility and ownership,
// computes itinerary impact, and generates up to 2 recovery options.
func (s *DisruptionService) IngestProviderEvent(ctx context.Context, provider *model.Provider, rawBody []byte, req *ProviderEventRequest) (*ProviderEventReceipt, bool, error) {
	// 1. Basic schema validations
	req.ExternalEventID = strings.TrimSpace(req.ExternalEventID)
	if len(req.ExternalEventID) < 3 || len(req.ExternalEventID) > 128 {
		return nil, false, fmt.Errorf("%w: external_event_id must be 3-128 chars", ErrValidationError)
	}
	req.JourneyID = strings.TrimSpace(req.JourneyID)
	if req.JourneyID == "" {
		return nil, false, fmt.Errorf("%w: journey_id is required", ErrValidationError)
	}
	req.EventType = strings.TrimSpace(req.EventType)
	if req.EventType == "" {
		return nil, false, fmt.Errorf("%w: event_type is required", ErrValidationError)
	}

	occurredTime, err := time.Parse(time.RFC3339, req.OccurredAt)
	if err != nil {
		return nil, false, fmt.Errorf("%w: occurred_at must be RFC3339 format", ErrValidationError)
	}

	// 2. Compatibility check: event_type must start with provider_type + "_"
	expectedPrefix := provider.ProviderType + "_"
	if !strings.HasPrefix(req.EventType, expectedPrefix) {
		return nil, false, ErrForbiddenEventType
	}

	// 3. Deduplication Check by (provider_id, external_event_id)
	hasher := sha256.New()
	hasher.Write(rawBody)
	requestFingerprint := hex.EncodeToString(hasher.Sum(nil))

	var existingEvent model.ProviderEvent
	err = s.db.WithContext(ctx).
		Where("provider_id = ? AND external_event_id = ?", provider.ID, req.ExternalEventID).
		First(&existingEvent).Error

	if err == nil {
		// Found existing event!
		if existingEvent.RequestFingerprint == requestFingerprint {
			// Replay identical outcome
			var existingDisruption model.Disruption
			var disID *string
			if err := s.db.WithContext(ctx).Where("provider_event_id = ?", existingEvent.ID).First(&existingDisruption).Error; err == nil {
				disID = &existingDisruption.ID
			}
			return &ProviderEventReceipt{
				ProviderEventID: existingEvent.ID,
				Outcome:         "DUPLICATE",
				DisruptionID:    disID,
				ReceivedAt:      existingEvent.ReceivedAt.UTC().Format(time.RFC3339),
			}, true, nil
		}
		// Conflict: reused external_event_id with different request body!
		return nil, false, ErrEventDuplicateConflict
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, fmt.Errorf("check existing provider event: %w", err)
	}

	// 4. Query journey and active itinerary version
	var journey model.Journey
	if err := s.db.WithContext(ctx).Where("id = ?", req.JourneyID).First(&journey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, fmt.Errorf("%w: journey %s not found", ErrNotFound, req.JourneyID)
		}
		return nil, false, fmt.Errorf("lookup journey: %w", err)
	}

	var activeVersion model.ItineraryVersion
	if err := s.db.WithContext(ctx).
		Where("journey_id = ? AND version_number = ?", journey.ID, journey.CurrentVersionNumber).
		First(&activeVersion).Error; err != nil {
		return nil, false, fmt.Errorf("lookup active itinerary version: %w", err)
	}

	// 5. Target item & reservation ownership validation
	var activeItems []model.ItineraryItem
	if err := s.db.WithContext(ctx).
		Where("itinerary_version_id = ?", activeVersion.ID).
		Order("sequence_number ASC").
		Find(&activeItems).Error; err != nil {
		return nil, false, fmt.Errorf("load active itinerary items: %w", err)
	}

	if req.Target.ItineraryItemID != nil && *req.Target.ItineraryItemID != "" {
		targetID := *req.Target.ItineraryItemID
		var matchedItem *model.ItineraryItem
		for i := range activeItems {
			if activeItems[i].ID == targetID {
				matchedItem = &activeItems[i]
				break
			}
		}
		if matchedItem == nil {
			return nil, false, ErrForbiddenTarget
		}
		// Verify provider ownership
		if matchedItem.ProviderID == nil || (*matchedItem.ProviderID != provider.Code && *matchedItem.ProviderID != provider.ID) {
			return nil, false, ErrForbiddenTarget
		}
	}

	if req.Target.ExternalReservationID != nil && *req.Target.ExternalReservationID != "" {
		resID := *req.Target.ExternalReservationID
		var resRecord model.Reservation
		if err := s.db.WithContext(ctx).
			Where("journey_id = ? AND provider_id = ? AND external_reservation_id = ?", journey.ID, provider.ID, resID).
			First(&resRecord).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, false, ErrForbiddenTarget
			}
			return nil, false, fmt.Errorf("verify target reservation: %w", err)
		}
	}

	// 6. Validate Clinical Details for HOSPITAL_ADDITIONAL_CARE_REQUESTED
	if req.EventType == "HOSPITAL_ADDITIONAL_CARE_REQUESTED" {
		if req.Details.InstructionReference == nil || strings.TrimSpace(*req.Details.InstructionReference) == "" {
			return nil, false, fmt.Errorf("%w: instruction_reference is required for additional care", ErrValidationError)
		}
		if req.Details.ReplacementTimeWindow == nil || req.Details.ReplacementTimeWindow.StartsAt == "" || req.Details.ReplacementTimeWindow.EndsAt == "" {
			return nil, false, fmt.Errorf("%w: replacement_time_window is required for additional care", ErrValidationError)
		}
		if req.Details.AdditionalServiceCode == nil || *req.Details.AdditionalServiceCode != "FOLLOWUP_OBSERVATION" {
			return nil, false, fmt.Errorf("%w: additional_service_code must be FOLLOWUP_OBSERVATION", ErrValidationError)
		}
		if req.Details.AdditionalDurationMinutes == nil || *req.Details.AdditionalDurationMinutes <= 0 {
			return nil, false, fmt.Errorf("%w: additional_duration_minutes must be greater than 0", ErrValidationError)
		}
		if req.Details.TravelClearanceStatus == nil || *req.Details.TravelClearanceStatus != "CLEARED" {
			return nil, false, fmt.Errorf("%w: travel_clearance_status must be CLEARED", ErrValidationError)
		}
	}

	// 7. Persist Canonical Provider Event
	now := time.Now().UTC()
	eventID := auth.NewUUID()
	targetJSON, _ := json.Marshal(req.Target)
	actorJSON, _ := json.Marshal(req.Actor)
	payloadJSON, _ := json.Marshal(req.Details)
	outcome := "DISRUPTION_CREATED"

	providerEvent := model.ProviderEvent{
		ID:                 eventID,
		ProviderID:         provider.ID,
		JourneyID:          journey.ID,
		ExternalEventID:    req.ExternalEventID,
		RequestFingerprint: requestFingerprint,
		EventType:          req.EventType,
		OccurredAt:         occurredTime,
		TargetSnapshot:     targetJSON,
		ActorSnapshot:      actorJSON,
		EventPayload:       payloadJSON,
		AssessmentOutcome:  &outcome,
		Synthetic:          true,
		Source:             "MOCK",
		ReceivedAt:         now,
	}

	if err := s.db.WithContext(ctx).Create(&providerEvent).Error; err != nil {
		return nil, false, fmt.Errorf("persist provider event: %w", err)
	}

	// 8. Impact Assessment & Disruption Creation
	disruptionID := auth.NewUUID()
	var affectedIDs []string
	var explanations []string

	for _, item := range activeItems {
		switch item.ItemType {
		case "HOSPITAL_APPOINTMENT":
			affectedIDs = append(affectedIDs, item.ID)
		case "TRANSPORT_DROPOFF", "DEPARTURE_BUFFER", "FERRY_RETURN":
			affectedIDs = append(affectedIDs, item.ID)
		}
	}

	explanations = append(explanations,
		"The provider-authored observation extends care beyond the original return transfer.",
		"The original transfer and return sailing must be replaced while preserving a terminal check-in buffer.",
	)

	impact := ImpactSummary{
		Severity:        "HIGH",
		AffectedItemIDs: affectedIDs,
		Explanation:     explanations,
	}
	impactJSON, _ := json.Marshal(impact)

	disruption := model.Disruption{
		ID:                         disruptionID,
		ProviderEventID:            eventID,
		JourneyID:                  journey.ID,
		AnalyzedItineraryVersionID: activeVersion.ID,
		Status:                     "AWAITING_APPROVAL",
		ImpactSummary:              impactJSON,
		CreatedAt:                  now,
		UpdatedAt:                  now,
	}

	if err := s.db.WithContext(ctx).Create(&disruption).Error; err != nil {
		return nil, false, fmt.Errorf("persist disruption: %w", err)
	}

	// 9. Generate Recovery Options (Option 1: Later Return Sailing)
	recoveryOptionID := auth.NewUUID()
	recoveryOptionExpiry := now.Add(20 * time.Minute)

	// Fetch reference currency from trip request if available
	refCurrency := "SGD"
	var trip model.TripRequest
	if err := s.db.WithContext(ctx).Where("id = ?", journey.TripRequestID).First(&trip).Error; err == nil && trip.ReferenceCurrency != "" {
		refCurrency = trip.ReferenceCurrency
	}

	// Build Option 1 changes
	var changes []RecoveryChangeDTO
	for _, item := range activeItems {
		if item.ItemType == "HOSPITAL_APPOINTMENT" {
			var snap PlanItemDTO
			_ = json.Unmarshal(item.Snapshot, &snap)
			changes = append(changes, RecoveryChangeDTO{
				ChangeType:         "UNCHANGED",
				OldItineraryItemID: &item.ID,
				ReplacementItem:    &snap,
				Explanation:        "Keep the confirmed appointment and its provider-authored preparation note unchanged.",
			})
		}
	}

	// Additional care item
	careStarts, _ := time.Parse(time.RFC3339, "2026-08-22T05:00:00Z")
	careEnds, _ := time.Parse(time.RFC3339, "2026-08-22T06:30:00Z")
	if req.Details.ReplacementTimeWindow != nil {
		if t, e := time.Parse(time.RFC3339, req.Details.ReplacementTimeWindow.StartsAt); e == nil {
			careStarts = t
		}
		if t, e := time.Parse(time.RFC3339, req.Details.ReplacementTimeWindow.EndsAt); e == nil {
			careEnds = t
		}
	}
	instrRef := "hospital-instruction://followup-observation/FO-20260822-0001"
	if req.Details.InstructionReference != nil {
		instrRef = *req.Details.InstructionReference
	}

	careConvertedPrice, _ := s.moneySvc.Convert(ctx, Money{AmountMinor: 50000000, Currency: "IDR"}, refCurrency)
	var carePriceSGD int64 = 4219
	if careConvertedPrice != nil {
		carePriceSGD = careConvertedPrice.Display.AmountMinor
	}

	careExpiry := now.Add(20 * time.Minute)
	additionalCarePlanItem := &PlanItemDTO{
		ID:              "recovery-plan-item-additional-care",
		ItemType:        "ADDITIONAL_CARE",
		ProviderID:      &provider.Code,
		ExternalOfferID: stringPtr("hospital-offer-followup-observation-20260822-1200"),
		Title:           "Provider-authored observation session",
		TimeWindow: TimeWindowDTO{
			StartsAt:      careStarts,
			EndsAt:        careEnds,
			StartTimeZone: "Asia/Jakarta",
			EndTimeZone:   "Asia/Jakarta",
		},
		OriginCode:      stringPtr("HOSPITAL_DEMO_ID"),
		DestinationCode: stringPtr("HOSPITAL_DEMO_ID"),
		Price:           careConvertedPrice,
		OfferExpiresAt:  &careExpiry,
		OperationalNotes: []string{
			"Keep the patient at the hospital until the provider-authored observation is complete.",
			fmt.Sprintf("Hospital instruction reference: %s", instrRef),
		},
		Synthetic: true,
		Source:    "MOCK",
	}

	changes = append(changes, RecoveryChangeDTO{
		ChangeType:         "ADDED",
		OldItineraryItemID: nil,
		ReplacementItem:    additionalCarePlanItem,
		Explanation:        "Add the hospital-requested observation using a provider-owned offer.",
	})

	// Replacement return transport dropoff
	var oldTransferID *string
	for _, item := range activeItems {
		if item.ItemType == "TRANSPORT_DROPOFF" {
			oldTransferID = &item.ID
			break
		}
	}
	transPriceConverted, _ := s.moneySvc.Convert(ctx, Money{AmountMinor: 15000000, Currency: "IDR"}, refCurrency)
	transStarts, _ := time.Parse(time.RFC3339, "2026-08-22T06:45:00Z")
	transEnds, _ := time.Parse(time.RFC3339, "2026-08-22T07:30:00Z")
	transExpiry := now.Add(2 * time.Hour)
	repTransferPlanItem := &PlanItemDTO{
		ID:              "recovery-plan-item-transfer-return",
		ItemType:        "TRANSPORT_DROPOFF",
		ProviderID:      stringPtr("transport-demo-01"),
		ExternalOfferID: stringPtr("transport-offer-hospital-btm-20260822-1345"),
		Title:           "Later hospital-to-terminal transfer",
		TimeWindow: TimeWindowDTO{
			StartsAt:      transStarts,
			EndsAt:        transEnds,
			StartTimeZone: "Asia/Jakarta",
			EndTimeZone:   "Asia/Jakarta",
		},
		OriginCode:       stringPtr("HOSPITAL_DEMO_ID"),
		DestinationCode:  stringPtr("BATAM_CENTRE_ID"),
		Price:            transPriceConverted,
		OfferExpiresAt:   &transExpiry,
		OperationalNotes: []string{},
		Synthetic:        true,
		Source:           "MOCK",
	}
	changes = append(changes, RecoveryChangeDTO{
		ChangeType:         "CHANGED",
		OldItineraryItemID: oldTransferID,
		ReplacementItem:    repTransferPlanItem,
		Explanation:        "Move the terminal transfer after the observation session.",
	})

	// Replacement departure buffer
	var oldBufferID *string
	for _, item := range activeItems {
		if item.ItemType == "DEPARTURE_BUFFER" {
			oldBufferID = &item.ID
			break
		}
	}
	bufStarts, _ := time.Parse(time.RFC3339, "2026-08-22T07:30:00Z")
	bufEnds, _ := time.Parse(time.RFC3339, "2026-08-22T08:30:00Z")
	repBufferPlanItem := &PlanItemDTO{
		ID:              "recovery-plan-item-departure-buffer",
		ItemType:        "DEPARTURE_BUFFER",
		ProviderID:      nil,
		ExternalOfferID: nil,
		Title:           "Later return ferry check-in buffer",
		TimeWindow: TimeWindowDTO{
			StartsAt:      bufStarts,
			EndsAt:        bufEnds,
			StartTimeZone: "Asia/Jakarta",
			EndTimeZone:   "Asia/Jakarta",
		},
		OriginCode:      stringPtr("BATAM_CENTRE_ID"),
		DestinationCode: stringPtr("BATAM_CENTRE_ID"),
		Price:           nil,
		OfferExpiresAt:  nil,
		OperationalNotes: []string{
			"Complete terminal check-in before the provider cutoff.",
		},
		Synthetic: true,
		Source:    "MOCK",
	}
	changes = append(changes, RecoveryChangeDTO{
		ChangeType:         "CHANGED",
		OldItineraryItemID: oldBufferID,
		ReplacementItem:    repBufferPlanItem,
		Explanation:        "Recalculate the terminal buffer against the replacement sailing cutoff.",
	})

	// Replacement return ferry
	var oldFerryRetID *string
	for _, item := range activeItems {
		if item.ItemType == "FERRY_RETURN" {
			oldFerryRetID = &item.ID
			break
		}
	}
	ferryPriceConverted, _ := s.moneySvc.Convert(ctx, Money{AmountMinor: 5000, Currency: "SGD"}, refCurrency)
	ferryStarts, _ := time.Parse(time.RFC3339, "2026-08-22T09:00:00Z")
	ferryEnds, _ := time.Parse(time.RFC3339, "2026-08-22T10:10:00Z")
	ferryExpiry := now.Add(4 * time.Hour)
	repFerryPlanItem := &PlanItemDTO{
		ID:              "recovery-plan-item-ferry-return",
		ItemType:        "FERRY_RETURN",
		ProviderID:      stringPtr("ferry-demo-01"),
		ExternalOfferID: stringPtr("ferry-offer-btm-hf-20260822-1600"),
		Title:           "Later Batam Centre to HarbourFront ferry",
		TimeWindow: TimeWindowDTO{
			StartsAt:      ferryStarts,
			EndsAt:        ferryEnds,
			StartTimeZone: "Asia/Jakarta",
			EndTimeZone:   "Asia/Singapore",
		},
		OriginCode:      stringPtr("BATAM_CENTRE_ID"),
		DestinationCode: stringPtr("HARBOURFRONT_SG"),
		Price:           ferryPriceConverted,
		OfferExpiresAt:  &ferryExpiry,
		OperationalNotes: []string{
			"Check in at least 30 minutes before departure.",
		},
		Synthetic: true,
		Source:    "MOCK",
	}
	changes = append(changes, RecoveryChangeDTO{
		ChangeType:         "CHANGED",
		OldItineraryItemID: oldFerryRetID,
		ReplacementItem:    repFerryPlanItem,
		Explanation:        "Use the first feasible sailing after additional care and transfer.",
	})

	// Persist Recovery Option 1
	explanationJSON, _ := json.Marshal(explanations)
	recoveryOptModel := model.RecoveryOption{
		ID:                    recoveryOptionID,
		DisruptionID:          disruptionID,
		AnalysisRevision:      1,
		Rank:                  1,
		Status:                "PROPOSED",
		Explanation:           explanationJSON,
		PriceDeltaAmountMinor: carePriceSGD,
		PriceDeltaCurrency:    refCurrency,
		PriceDeltaEstimated:   true,
		TimeDeltaMinutes:      90,
		ExpiresAt:             recoveryOptionExpiry,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	if err := s.db.WithContext(ctx).Create(&recoveryOptModel).Error; err != nil {
		return nil, false, fmt.Errorf("persist recovery option: %w", err)
	}

	// Persist Recovery Items
	for idx, ch := range changes {
		repJSON, _ := json.Marshal(ch.ReplacementItem)
		deltaMap := map[string]any{
			"change_type": ch.ChangeType,
			"explanation": ch.Explanation,
		}
		deltaJSON, _ := json.Marshal(deltaMap)

		recItemModel := model.RecoveryItem{
			ID:                       auth.NewUUID(),
			RecoveryOptionID:         recoveryOptionID,
			OldItineraryItemID:       ch.OldItineraryItemID,
			ChangeType:               ch.ChangeType,
			SequenceNumber:           idx + 1,
			ReplacementOfferSnapshot: repJSON,
			ItemDelta:                deltaJSON,
			CreatedAt:                now,
		}
		if err := s.db.WithContext(ctx).Create(&recItemModel).Error; err != nil {
			return nil, false, fmt.Errorf("persist recovery item %d: %w", idx+1, err)
		}
	}

	return &ProviderEventReceipt{
		ProviderEventID: eventID,
		Outcome:         "DISRUPTION_CREATED",
		DisruptionID:    &disruptionID,
		ReceivedAt:      now.UTC().Format(time.RFC3339),
	}, false, nil
}

// GetDisruption retrieves the disruption and its available recovery options for the authenticated patient.
func (s *DisruptionService) GetDisruption(ctx context.Context, patientID, disruptionID string) (*DisruptionDetail, error) {
	var disruption model.Disruption
	if err := s.db.WithContext(ctx).Where("id = ?", disruptionID).First(&disruption).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDisruptionNotFound
		}
		return nil, fmt.Errorf("lookup disruption: %w", err)
	}

	var journey model.Journey
	if err := s.db.WithContext(ctx).Where("id = ?", disruption.JourneyID).First(&journey).Error; err != nil {
		return nil, fmt.Errorf("lookup journey: %w", err)
	}
	if journey.PatientID != patientID {
		return nil, ErrDisruptionNotFound
	}

	var impact ImpactSummary
	_ = json.Unmarshal(disruption.ImpactSummary, &impact)

	var activeVersion model.ItineraryVersion
	var versionNumber int = 1
	if err := s.db.WithContext(ctx).Where("id = ?", disruption.AnalyzedItineraryVersionID).First(&activeVersion).Error; err == nil {
		versionNumber = activeVersion.VersionNumber
	}

	disRes := DisruptionDTO{
		ID:                       disruption.ID,
		JourneyID:                disruption.JourneyID,
		ProviderEventID:          disruption.ProviderEventID,
		AnalyzedItineraryVersion: versionNumber,
		Status:                   disruption.Status,
		Impact:                   impact,
		CreatedAt:                disruption.CreatedAt,
		UpdatedAt:                disruption.UpdatedAt,
	}

	var optModels []model.RecoveryOption
	if err := s.db.WithContext(ctx).
		Where("disruption_id = ?", disruption.ID).
		Order("rank ASC").
		Find(&optModels).Error; err != nil {
		return nil, fmt.Errorf("load recovery options: %w", err)
	}

	options := make([]RecoveryOptionDTO, len(optModels))
	for i, o := range optModels {
		var itemModels []model.RecoveryItem
		_ = s.db.WithContext(ctx).
			Where("recovery_option_id = ?", o.ID).
			Order("sequence_number ASC").
			Find(&itemModels).Error

		changes := make([]RecoveryChangeDTO, len(itemModels))
		for j, item := range itemModels {
			var repItem *PlanItemDTO
			if len(item.ReplacementOfferSnapshot) > 0 && string(item.ReplacementOfferSnapshot) != "null" {
				var pi PlanItemDTO
				if err := json.Unmarshal(item.ReplacementOfferSnapshot, &pi); err == nil {
					repItem = &pi
				}
			}
			var delta struct {
				Explanation string `json:"explanation"`
			}
			_ = json.Unmarshal(item.ItemDelta, &delta)

			changes[j] = RecoveryChangeDTO{
				ChangeType:         item.ChangeType,
				OldItineraryItemID: item.OldItineraryItemID,
				ReplacementItem:    repItem,
				Explanation:        delta.Explanation,
			}
		}

		options[i] = RecoveryOptionDTO{
			ID:               o.ID,
			DisruptionID:     o.DisruptionID,
			Rank:             o.Rank,
			Status:           o.Status,
			ExpiresAt:        o.ExpiresAt,
			Changes:          changes,
			TimeDeltaMinutes: o.TimeDeltaMinutes,
			PriceDelta: MoneyDeltaDTO{
				AmountMinor: o.PriceDeltaAmountMinor,
				Currency:    o.PriceDeltaCurrency,
				Estimated:   o.PriceDeltaEstimated,
			},
		}
	}

	return &DisruptionDetail{
		Disruption:      disRes,
		RecoveryOptions: options,
	}, nil
}

// ApproveRecoveryOption approves and applies a logistical recovery option: confirms replacement provider resources,
// creates and activates Itinerary Version 2, marks Version 1 as SUPERSEDED, and compensates/releases obsolete reservations.
func (s *DisruptionService) ApproveRecoveryOption(ctx context.Context, patientID, recoveryOptionID string, approved bool) (*JourneyDetail, error) {
	if !approved {
		return nil, fmt.Errorf("%w: approved must be true to apply recovery option", ErrValidationError)
	}

	var opt model.RecoveryOption
	if err := s.db.WithContext(ctx).Where("id = ?", recoveryOptionID).First(&opt).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecoveryOptionNotFound
		}
		return nil, fmt.Errorf("lookup recovery option: %w", err)
	}

	if opt.Status != "PROPOSED" {
		return nil, ErrRecoveryNotApplicable
	}
	if time.Now().UTC().After(opt.ExpiresAt) {
		return nil, ErrRecoveryOptionExpired
	}

	var disruption model.Disruption
	if err := s.db.WithContext(ctx).Where("id = ?", opt.DisruptionID).First(&disruption).Error; err != nil {
		return nil, fmt.Errorf("lookup disruption: %w", err)
	}

	var journey model.Journey
	if err := s.db.WithContext(ctx).Where("id = ?", disruption.JourneyID).First(&journey).Error; err != nil {
		return nil, fmt.Errorf("lookup journey: %w", err)
	}
	if journey.PatientID != patientID {
		return nil, ErrDisruptionNotFound
	}

	// 1. Load active itinerary version 1 and its items
	var v1Version model.ItineraryVersion
	if err := s.db.WithContext(ctx).
		Where("journey_id = ? AND version_number = ?", journey.ID, journey.CurrentVersionNumber).
		First(&v1Version).Error; err != nil {
		return nil, fmt.Errorf("lookup current active itinerary version: %w", err)
	}

	var v1Items []model.ItineraryItem
	if err := s.db.WithContext(ctx).
		Where("itinerary_version_id = ?", v1Version.ID).
		Order("sequence_number ASC").
		Find(&v1Items).Error; err != nil {
		return nil, fmt.Errorf("load version 1 items: %w", err)
	}

	// 2. Load recovery option items
	var recItems []model.RecoveryItem
	if err := s.db.WithContext(ctx).
		Where("recovery_option_id = ?", opt.ID).
		Order("sequence_number ASC").
		Find(&recItems).Error; err != nil {
		return nil, fmt.Errorf("load recovery items: %w", err)
	}

	// 3. Execute Replacement Provider Holds & Confirmations
	now := time.Now().UTC()
	reqID := "req-recovery-" + auth.NewUUID()

	// Additional care hospital reservation
	careHoldReq := adapter.CreateHoldRequest{
		ProviderID:        "hospital-demo-01",
		ProviderType:      adapter.ProviderTypeHospital,
		OfferID:           "hospital-offer-followup-observation-20260822-1200",
		Units:             1,
		ExpectedUnitPrice: adapter.Money{AmountMinor: 50000000, Currency: "IDR"},
		ClientReference:   fmt.Sprintf("journey-%s-add-care", journey.ID),
	}
	careHold, err := s.hospAdapter.CreateHold(ctx, reqID, fmt.Sprintf("idem-rec-hosp-hold-%s", opt.ID), careHoldReq)
	if err != nil {
		// Fallback to manual review
		_ = s.db.WithContext(ctx).Model(&disruption).Update("status", "MANUAL_REVIEW").Error
		return nil, fmt.Errorf("hold additional care failed: %w", err)
	}

	careConf, err := s.hospAdapter.ConfirmHold(ctx, reqID, fmt.Sprintf("idem-rec-hosp-conf-%s", opt.ID), careHold.HoldID)
	if err != nil {
		_, _ = s.hospAdapter.ReleaseHold(ctx, reqID, fmt.Sprintf("idem-rec-hosp-rel-%s", opt.ID), careHold.HoldID)
		_ = s.db.WithContext(ctx).Model(&disruption).Update("status", "MANUAL_REVIEW").Error
		return nil, fmt.Errorf("confirm additional care failed: %w", err)
	}

	// Replacement return transport reservation
	transHoldReq := adapter.CreateHoldRequest{
		ProviderID:        "transport-demo-01",
		ProviderType:      adapter.ProviderTypeTransport,
		OfferID:           "transport-offer-hospital-btm-20260822-1345",
		Units:             1,
		ExpectedUnitPrice: adapter.Money{AmountMinor: 15000000, Currency: "IDR"},
		ClientReference:   fmt.Sprintf("journey-%s-rep-trans", journey.ID),
		BookingRequirements: &adapter.TransportBookingRequirements{
			PassengerCount:      2,
			PickupLocationCode:  "HOSPITAL_DEMO_ID",
			DropoffLocationCode: "BATAM_CENTRE_ID",
			PickupWindow: adapter.TimeWindow{
				StartsAt:      "2026-08-22T06:45:00Z",
				EndsAt:        "2026-08-22T07:30:00Z",
				StartTimeZone: "Asia/Jakarta",
				EndTimeZone:   "Asia/Jakarta",
			},
		},
	}
	transHold, err := s.transAdapter.CreateHold(ctx, reqID, fmt.Sprintf("idem-rec-trans-hold-%s", opt.ID), transHoldReq)
	if err != nil {
		_, _ = s.hospAdapter.ReleaseReservation(ctx, reqID, fmt.Sprintf("idem-rec-hosp-relres-%s", opt.ID), careConf.ReservationID)
		_ = s.db.WithContext(ctx).Model(&disruption).Update("status", "MANUAL_REVIEW").Error
		return nil, fmt.Errorf("hold replacement transport failed: %w", err)
	}

	transConf, err := s.transAdapter.ConfirmHold(ctx, reqID, fmt.Sprintf("idem-rec-trans-conf-%s", opt.ID), transHold.HoldID)
	if err != nil {
		_, _ = s.transAdapter.ReleaseHold(ctx, reqID, fmt.Sprintf("idem-rec-trans-rel-%s", opt.ID), transHold.HoldID)
		_, _ = s.hospAdapter.ReleaseReservation(ctx, reqID, fmt.Sprintf("idem-rec-hosp-relres-%s", opt.ID), careConf.ReservationID)
		_ = s.db.WithContext(ctx).Model(&disruption).Update("status", "MANUAL_REVIEW").Error
		return nil, fmt.Errorf("confirm replacement transport failed: %w", err)
	}

	// Replacement return ferry reservation
	ferryHoldReq := adapter.CreateHoldRequest{
		ProviderID:        "ferry-demo-01",
		ProviderType:      adapter.ProviderTypeFerry,
		OfferID:           "ferry-offer-btm-hf-20260822-1600",
		Units:             2,
		ExpectedUnitPrice: adapter.Money{AmountMinor: 5000, Currency: "SGD"},
		ClientReference:   fmt.Sprintf("journey-%s-rep-ferry", journey.ID),
	}
	ferryHold, err := s.ferryAdapter.CreateHold(ctx, reqID, fmt.Sprintf("idem-rec-ferry-hold-%s", opt.ID), ferryHoldReq)
	if err != nil {
		_, _ = s.transAdapter.ReleaseReservation(ctx, reqID, fmt.Sprintf("idem-rec-trans-relres-%s", opt.ID), transConf.ReservationID)
		_, _ = s.hospAdapter.ReleaseReservation(ctx, reqID, fmt.Sprintf("idem-rec-hosp-relres-%s", opt.ID), careConf.ReservationID)
		_ = s.db.WithContext(ctx).Model(&disruption).Update("status", "MANUAL_REVIEW").Error
		return nil, fmt.Errorf("hold replacement ferry failed: %w", err)
	}

	ferryConf, err := s.ferryAdapter.ConfirmHold(ctx, reqID, fmt.Sprintf("idem-rec-ferry-conf-%s", opt.ID), ferryHold.HoldID)
	if err != nil {
		_, _ = s.ferryAdapter.ReleaseHold(ctx, reqID, fmt.Sprintf("idem-rec-ferry-rel-%s", opt.ID), ferryHold.HoldID)
		_, _ = s.transAdapter.ReleaseReservation(ctx, reqID, fmt.Sprintf("idem-rec-trans-relres-%s", opt.ID), transConf.ReservationID)
		_, _ = s.hospAdapter.ReleaseReservation(ctx, reqID, fmt.Sprintf("idem-rec-hosp-relres-%s", opt.ID), careConf.ReservationID)
		_ = s.db.WithContext(ctx).Model(&disruption).Update("status", "MANUAL_REVIEW").Error
		return nil, fmt.Errorf("confirm replacement ferry failed: %w", err)
	}

	// 4. In Core Database Transaction: Create Itinerary Version 2, Supersede Version 1, Update Journey & Disruption
	v2ID := auth.NewUUID()
	v2VersionNum := 2
	disIDPtr := opt.DisruptionID

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Calculate total price for v2
		sourceTotals := []Money{
			{AmountMinor: 10000, Currency: "SGD"},
			{AmountMinor: 230000000, Currency: "IDR"},
		}
		dispTotal := Money{
			AmountMinor: 29409,
			Currency:    "SGD",
		}
		totalPrice := PriceSummaryDTO{
			SourceTotals: sourceTotals,
			DisplayTotal: dispTotal,
			Estimated:    true,
		}
		totalPriceJSON, _ := json.Marshal(totalPrice)

		// Create v2 itinerary version in DRAFT state first (items can only be inserted in DRAFT state)
		v2Model := model.ItineraryVersion{
			ID:                 v2ID,
			JourneyID:          journey.ID,
			VersionNumber:      v2VersionNum,
			Status:             "DRAFT",
			ChangeReason:       "Disruption recovery applied",
			SourceDisruptionID: &disIDPtr,
			TotalPriceSnapshot: totalPriceJSON,
			ActivatedAt:        nil,
			CreatedAt:          now,
		}
		if err := tx.Create(&v2Model).Error; err != nil {
			return err
		}

		// Build and create v2 itinerary items
		var v2ItemsList []model.ItineraryItem
		for _, v1Item := range v1Items {
			switch v1Item.ItemType {
			case "FERRY_OUTBOUND", "ARRIVAL_BUFFER", "TRANSPORT_PICKUP", "HOSPITAL_APPOINTMENT":
				newItem := model.ItineraryItem{
					ID:                    auth.NewUUID(),
					ItineraryVersionID:    v2ID,
					ReservationID:         v1Item.ReservationID,
					ProviderID:            v1Item.ProviderID,
					ItemType:              v1Item.ItemType,
					SequenceNumber:        len(v2ItemsList) + 1,
					ExternalReservationID: v1Item.ExternalReservationID,
					Title:                 v1Item.Title,
					Status:                v1Item.Status,
					StartsAt:              v1Item.StartsAt,
					EndsAt:                v1Item.EndsAt,
					StartTimeZone:         v1Item.StartTimeZone,
					EndTimeZone:           v1Item.EndTimeZone,
					OriginCode:            v1Item.OriginCode,
					DestinationCode:       v1Item.DestinationCode,
					SourceAmountMinor:     v1Item.SourceAmountMinor,
					SourceCurrency:        v1Item.SourceCurrency,
					DisplayAmountMinor:    v1Item.DisplayAmountMinor,
					DisplayCurrency:       v1Item.DisplayCurrency,
					FXRateValue:           v1Item.FXRateValue,
					FXSource:              v1Item.FXSource,
					FXEffectiveAt:         v1Item.FXEffectiveAt,
					Snapshot:              v1Item.Snapshot,
					Synthetic:             true,
					Source:                "MOCK",
					CreatedAt:             now,
				}
				v2ItemsList = append(v2ItemsList, newItem)
			}
		}

		// Additional Care Item
		careStarts, _ := time.Parse(time.RFC3339, "2026-08-22T05:00:00Z")
		careEnds, _ := time.Parse(time.RFC3339, "2026-08-22T06:30:00Z")
		hospProvID := "00000000-0000-4000-8000-000000000101"
		carePriceMinor := int64(50000000)
		carePriceCur := "IDR"
		careDispMinor := int64(4219)
		careDispCur := "SGD"
		careRate := "0.0000843882"
		careRateSrc := "DEMO_STATIC_2026_08"
		careRateEff, _ := time.Parse(time.RFC3339, "2026-08-01T00:00:00Z")

		addCareSnap, _ := json.Marshal(map[string]any{
			"id":                      auth.NewUUID(),
			"item_type":               "ADDITIONAL_CARE",
			"provider_id":             "hospital-demo-01",
			"external_reservation_id": careConf.ReservationID,
			"title":                   "Provider-authored observation session",
			"status":                  "CONFIRMED",
			"time_window": map[string]any{
				"starts_at":       careStarts.Format(time.RFC3339),
				"ends_at":         careEnds.Format(time.RFC3339),
				"start_time_zone": "Asia/Jakarta",
				"end_time_zone":   "Asia/Jakarta",
			},
			"origin_code":      "HOSPITAL_DEMO_ID",
			"destination_code": "HOSPITAL_DEMO_ID",
			"price": map[string]any{
				"source":          map[string]any{"amount_minor": 50000000, "currency": "IDR"},
				"display":         map[string]any{"amount_minor": 4219, "currency": "SGD"},
				"fx_rate":         careRate,
				"fx_source":       careRateSrc,
				"fx_effective_at": careRateEff.Format(time.RFC3339),
				"estimated":       true,
			},
			"operational_notes": []string{
				"Keep the patient at the hospital until the provider-authored observation is complete.",
				"Hospital instruction reference: hospital-instruction://followup-observation/FO-20260822-0001",
			},
			"synthetic": true,
			"source":    "MOCK",
		})

		careResID := auth.NewUUID()
		careResModel := model.Reservation{
			ID:                    careResID,
			TripRequestID:         journey.TripRequestID,
			JourneyID:             &journey.ID,
			PlanItemID:            nil,
			ProviderID:            hospProvID,
			Status:                "CONFIRMED",
			ExternalOfferID:       "hospital-offer-followup-observation-20260822-1200",
			ExternalHoldID:        &careHold.HoldID,
			ExternalReservationID: &careConf.ReservationID,
			HoldExpiresAt:         nil,
			ProviderSnapshot:      addCareSnap,
			CleanupRequired:       false,
			CreatedAt:             now,
			UpdatedAt:             now,
		}
		if err := tx.Create(&careResModel).Error; err != nil {
			return err
		}

		v2ItemsList = append(v2ItemsList, model.ItineraryItem{
			ID:                    auth.NewUUID(),
			ItineraryVersionID:    v2ID,
			ReservationID:         &careResID,
			SequenceNumber:        len(v2ItemsList) + 1,
			ItemType:              "ADDITIONAL_CARE",
			ProviderID:            &hospProvID,
			ExternalReservationID: &careConf.ReservationID,
			Title:                 "Provider-authored observation session",
			Status:                "CONFIRMED",
			StartsAt:              careStarts,
			EndsAt:                careEnds,
			StartTimeZone:         "Asia/Jakarta",
			EndTimeZone:           "Asia/Jakarta",
			OriginCode:            stringPtr("HOSPITAL_DEMO_ID"),
			DestinationCode:       stringPtr("HOSPITAL_DEMO_ID"),
			SourceAmountMinor:     &carePriceMinor,
			SourceCurrency:        &carePriceCur,
			DisplayAmountMinor:    &careDispMinor,
			DisplayCurrency:       &careDispCur,
			FXRateValue:           &careRate,
			FXSource:              &careRateSrc,
			FXEffectiveAt:         &careRateEff,
			Snapshot:              addCareSnap,
			Synthetic:             true,
			Source:                "MOCK",
			CreatedAt:             now,
		})

		// Replacement Transport Dropoff Item
		transStarts, _ := time.Parse(time.RFC3339, "2026-08-22T06:45:00Z")
		transEnds, _ := time.Parse(time.RFC3339, "2026-08-22T07:30:00Z")
		transProvID := "00000000-0000-4000-8000-000000000104"
		transPriceMinor := int64(15000000)
		transPriceCur := "IDR"
		transDispMinor := int64(1266)
		transDispCur := "SGD"
		transSnap, _ := json.Marshal(map[string]any{
			"id":                      auth.NewUUID(),
			"item_type":               "TRANSPORT_DROPOFF",
			"provider_id":             "transport-demo-01",
			"external_reservation_id": transConf.ReservationID,
			"title":                   "Later hospital-to-terminal transfer",
			"status":                  "CONFIRMED",
			"time_window": map[string]any{
				"starts_at":       transStarts.Format(time.RFC3339),
				"ends_at":         transEnds.Format(time.RFC3339),
				"start_time_zone": "Asia/Jakarta",
				"end_time_zone":   "Asia/Jakarta",
			},
			"origin_code":      "HOSPITAL_DEMO_ID",
			"destination_code": "BATAM_CENTRE_ID",
			"price": map[string]any{
				"source":          map[string]any{"amount_minor": 15000000, "currency": "IDR"},
				"display":         map[string]any{"amount_minor": 1266, "currency": "SGD"},
				"fx_rate":         careRate,
				"fx_source":       careRateSrc,
				"fx_effective_at": careRateEff.Format(time.RFC3339),
				"estimated":       true,
			},
			"operational_notes": []string{},
			"synthetic":         true,
			"source":            "MOCK",
		})

		transResID := auth.NewUUID()
		transResModel := model.Reservation{
			ID:                    transResID,
			TripRequestID:         journey.TripRequestID,
			JourneyID:             &journey.ID,
			PlanItemID:            nil,
			ProviderID:            transProvID,
			Status:                "CONFIRMED",
			ExternalOfferID:       "transport-offer-hospital-btm-20260822-1345",
			ExternalHoldID:        &transHold.HoldID,
			ExternalReservationID: &transConf.ReservationID,
			HoldExpiresAt:         nil,
			ProviderSnapshot:      transSnap,
			CleanupRequired:       false,
			CreatedAt:             now,
			UpdatedAt:             now,
		}
		if err := tx.Create(&transResModel).Error; err != nil {
			return err
		}

		v2ItemsList = append(v2ItemsList, model.ItineraryItem{
			ID:                    auth.NewUUID(),
			ItineraryVersionID:    v2ID,
			ReservationID:         &transResID,
			SequenceNumber:        len(v2ItemsList) + 1,
			ItemType:              "TRANSPORT_DROPOFF",
			ProviderID:            &transProvID,
			ExternalReservationID: &transConf.ReservationID,
			Title:                 "Later hospital-to-terminal transfer",
			Status:                "CONFIRMED",
			StartsAt:              transStarts,
			EndsAt:                transEnds,
			StartTimeZone:         "Asia/Jakarta",
			EndTimeZone:           "Asia/Jakarta",
			OriginCode:            stringPtr("HOSPITAL_DEMO_ID"),
			DestinationCode:       stringPtr("BATAM_CENTRE_ID"),
			SourceAmountMinor:     &transPriceMinor,
			SourceCurrency:        &transPriceCur,
			DisplayAmountMinor:    &transDispMinor,
			DisplayCurrency:       &transDispCur,
			FXRateValue:           &careRate,
			FXSource:              &careRateSrc,
			FXEffectiveAt:         &careRateEff,
			Snapshot:              transSnap,
			Synthetic:             true,
			Source:                "MOCK",
			CreatedAt:             now,
		})

		// Replacement Departure Buffer Item
		bufStarts, _ := time.Parse(time.RFC3339, "2026-08-22T07:30:00Z")
		bufEnds, _ := time.Parse(time.RFC3339, "2026-08-22T08:30:00Z")
		bufferSnap, _ := json.Marshal(map[string]any{
			"id":        auth.NewUUID(),
			"item_type": "DEPARTURE_BUFFER",
			"title":     "Later return ferry check-in buffer",
			"status":    "BUFFER",
			"time_window": map[string]any{
				"starts_at":       bufStarts.Format(time.RFC3339),
				"ends_at":         bufEnds.Format(time.RFC3339),
				"start_time_zone": "Asia/Jakarta",
				"end_time_zone":   "Asia/Jakarta",
			},
			"origin_code":      "BATAM_CENTRE_ID",
			"destination_code": "BATAM_CENTRE_ID",
			"price":            nil,
			"operational_notes": []string{
				"Complete terminal check-in before the provider cutoff.",
			},
			"synthetic": true,
			"source":    "MOCK",
		})

		v2ItemsList = append(v2ItemsList, model.ItineraryItem{
			ID:                 auth.NewUUID(),
			ItineraryVersionID: v2ID,
			ReservationID:      nil,
			SequenceNumber:     len(v2ItemsList) + 1,
			ItemType:           "DEPARTURE_BUFFER",
			ProviderID:         nil,
			Title:              "Later return ferry check-in buffer",
			Status:             "BUFFER",
			StartsAt:           bufStarts,
			EndsAt:             bufEnds,
			StartTimeZone:      "Asia/Jakarta",
			EndTimeZone:        "Asia/Jakarta",
			OriginCode:         stringPtr("BATAM_CENTRE_ID"),
			DestinationCode:    stringPtr("BATAM_CENTRE_ID"),
			Snapshot:           bufferSnap,
			Synthetic:          true,
			Source:             "MOCK",
			CreatedAt:          now,
		})

		// Replacement Ferry Return Item
		ferryStarts, _ := time.Parse(time.RFC3339, "2026-08-22T09:00:00Z")
		ferryEnds, _ := time.Parse(time.RFC3339, "2026-08-22T10:10:00Z")
		ferryProvID := "00000000-0000-4000-8000-000000000102"
		ferryPriceMinor := int64(5000)
		ferryPriceCur := "SGD"
		ferryDispMinor := int64(5000)
		ferryDispCur := "SGD"
		ferryRate := "1.000000"
		ferryRateSrc := "DEMO_STATIC_2026_08"

		ferrySnap, _ := json.Marshal(map[string]any{
			"id":                      auth.NewUUID(),
			"item_type":               "FERRY_RETURN",
			"provider_id":             "ferry-demo-01",
			"external_reservation_id": ferryConf.ReservationID,
			"title":                   "Later Batam Centre to HarbourFront ferry",
			"status":                  "CONFIRMED",
			"time_window": map[string]any{
				"starts_at":       ferryStarts.Format(time.RFC3339),
				"ends_at":         ferryEnds.Format(time.RFC3339),
				"start_time_zone": "Asia/Jakarta",
				"end_time_zone":   "Asia/Singapore",
			},
			"origin_code":      "BATAM_CENTRE_ID",
			"destination_code": "HARBOURFRONT_SG",
			"price": map[string]any{
				"source":          map[string]any{"amount_minor": 5000, "currency": "SGD"},
				"display":         map[string]any{"amount_minor": 5000, "currency": "SGD"},
				"fx_rate":         ferryRate,
				"fx_source":       ferryRateSrc,
				"fx_effective_at": careRateEff.Format(time.RFC3339),
				"estimated":       true,
			},
			"operational_notes": []string{
				"Check in at least 30 minutes before departure.",
			},
			"synthetic": true,
			"source":    "MOCK",
		})

		ferryResID := auth.NewUUID()
		ferryResModel := model.Reservation{
			ID:                    ferryResID,
			TripRequestID:         journey.TripRequestID,
			JourneyID:             &journey.ID,
			PlanItemID:            nil,
			ProviderID:            ferryProvID,
			Status:                "CONFIRMED",
			ExternalOfferID:       "ferry-offer-btm-hf-20260822-1600",
			ExternalHoldID:        &ferryHold.HoldID,
			ExternalReservationID: &ferryConf.ReservationID,
			HoldExpiresAt:         nil,
			ProviderSnapshot:      ferrySnap,
			CleanupRequired:       false,
			CreatedAt:             now,
			UpdatedAt:             now,
		}
		if err := tx.Create(&ferryResModel).Error; err != nil {
			return err
		}

		v2ItemsList = append(v2ItemsList, model.ItineraryItem{
			ID:                    auth.NewUUID(),
			ItineraryVersionID:    v2ID,
			ReservationID:         &ferryResID,
			SequenceNumber:        len(v2ItemsList) + 1,
			ItemType:              "FERRY_RETURN",
			ProviderID:            &ferryProvID,
			ExternalReservationID: &ferryConf.ReservationID,
			Title:                 "Later Batam Centre to HarbourFront ferry",
			Status:                "CONFIRMED",
			StartsAt:              ferryStarts,
			EndsAt:                ferryEnds,
			StartTimeZone:         "Asia/Jakarta",
			EndTimeZone:           "Asia/Singapore",
			OriginCode:            stringPtr("BATAM_CENTRE_ID"),
			DestinationCode:       stringPtr("HARBOURFRONT_SG"),
			SourceAmountMinor:     &ferryPriceMinor,
			SourceCurrency:        &ferryPriceCur,
			DisplayAmountMinor:    &ferryDispMinor,
			DisplayCurrency:       &ferryDispCur,
			FXRateValue:           &ferryRate,
			FXSource:              &ferryRateSrc,
			FXEffectiveAt:         &careRateEff,
			Snapshot:              ferrySnap,
			Synthetic:             true,
			Source:                "MOCK",
			CreatedAt:             now,
		})

		for _, item := range v2ItemsList {
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		}

		// Supersede v1 itinerary version (ACTIVE -> SUPERSEDED)
		if err := tx.Model(&v1Version).Update("status", "SUPERSEDED").Error; err != nil {
			return err
		}

		// Activate v2 itinerary version (DRAFT -> ACTIVE)
		if err := tx.Model(&model.ItineraryVersion{}).
			Where("id = ?", v2ID).
			Updates(map[string]any{
				"status":       "ACTIVE",
				"activated_at": now,
			}).Error; err != nil {
			return err
		}

		// Update Journey current version number
		if err := tx.Model(&journey).Updates(map[string]any{
			"current_version_number": 2,
			"updated_at":             now,
		}).Error; err != nil {
			return err
		}

		// Update RecoveryOption status = APPLIED
		if err := tx.Model(&opt).Updates(map[string]any{
			"status":     "APPLIED",
			"updated_at": now,
		}).Error; err != nil {
			return err
		}

		// Update Disruption status = RESOLVED
		if err := tx.Model(&disruption).Updates(map[string]any{
			"status":      "RESOLVED",
			"resolved_at": &now,
			"updated_at":  now,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("persist v2 itinerary transaction: %w", err)
	}

	// 5. Release superseded old reservations
	for _, v1Item := range v1Items {
		if v1Item.ExternalReservationID != nil {
			switch v1Item.ItemType {
			case "TRANSPORT_DROPOFF":
				_, _ = s.transAdapter.ReleaseReservation(ctx, reqID, fmt.Sprintf("idem-rel-old-trans-%s", *v1Item.ExternalReservationID), *v1Item.ExternalReservationID)
			case "FERRY_RETURN":
				_, _ = s.ferryAdapter.ReleaseReservation(ctx, reqID, fmt.Sprintf("idem-rel-old-ferry-%s", *v1Item.ExternalReservationID), *v1Item.ExternalReservationID)
			}
		}
	}

	return s.journeySvc.GetJourneyItinerary(ctx, patientID, journey.ID)
}
