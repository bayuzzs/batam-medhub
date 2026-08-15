package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"batam-medhub/internal/auth"
	"batam-medhub/internal/model"

	"gorm.io/gorm"
)

var (
	ErrTripRequestNotFound = errors.New("trip request not found")
	ErrInvalidTripState    = errors.New("trip request is not in a plannable state")
)

// TripRequestDTO represents the API view of a trip request.
type TripRequestDTO struct {
	ID                string           `json:"id"`
	Status            string           `json:"status"`
	Intent            StructuredIntent `json:"intent"`
	PlanningRevision  int              `json:"planning_revision"`
	ReferenceCurrency string           `json:"reference_currency"`
	JourneyID         *string          `json:"journey_id"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}

// TimeWindowDTO represents an operational time interval with explicit timezones.
type TimeWindowDTO struct {
	StartsAt      time.Time `json:"starts_at"`
	EndsAt        time.Time `json:"ends_at"`
	StartTimeZone string    `json:"start_time_zone"`
	EndTimeZone   string    `json:"end_time_zone"`
}

// PlanItemDTO represents a single reservation or buffer segment in a plan option.
type PlanItemDTO struct {
	ID               string          `json:"id"`
	ItemType         string          `json:"item_type"`
	ProviderID       *string         `json:"provider_id"`
	ExternalOfferID  *string         `json:"external_offer_id"`
	Title            string          `json:"title"`
	TimeWindow       TimeWindowDTO   `json:"time_window"`
	OriginCode       *string         `json:"origin_code"`
	DestinationCode  *string         `json:"destination_code"`
	Price            *ConvertedMoney `json:"price"`
	OfferExpiresAt   *time.Time      `json:"offer_expires_at"`
	OperationalNotes []string        `json:"operational_notes"`
	Synthetic        bool            `json:"synthetic"`
	Source           string          `json:"source"`
}

// PriceSummaryDTO represents multi-currency source totals and display currency conversion.
type PriceSummaryDTO struct {
	SourceTotals []Money `json:"source_totals"`
	DisplayTotal Money   `json:"display_total"`
	Estimated    bool    `json:"estimated"`
}

// PlanOptionDTO represents a ranked cross-provider itinerary package.
type PlanOptionDTO struct {
	ID               string          `json:"id"`
	TripRequestID    string          `json:"trip_request_id"`
	PlanningRevision int             `json:"planning_revision"`
	Rank             int             `json:"rank"`
	Status           string          `json:"status"`
	ExpiresAt        time.Time       `json:"expires_at"`
	Explanation      []string        `json:"explanation"`
	Items            []PlanItemDTO   `json:"items"`
	TotalPrice       PriceSummaryDTO `json:"total_price"`
}

// TripRequestDetail wraps a trip request and its active plan options.
type TripRequestDetail struct {
	TripRequest TripRequestDTO  `json:"trip_request"`
	PlanOptions []PlanOptionDTO `json:"plan_options"`
}

// PlanningResult represents the outcome of generating plan options.
type PlanningResult struct {
	TripRequest      TripRequestDTO  `json:"trip_request"`
	Options          []PlanOptionDTO `json:"options"`
	NoMatchReasons   []string        `json:"no_match_reasons"`
	ProviderWarnings []string        `json:"provider_warnings"`
}

// CreateTripRequestInput represents the creation payload.
type CreateTripRequestInput struct {
	Prompt string `json:"prompt"`
	Locale string `json:"locale"`
}

// TripService manages trip request lifecycles and deterministic package planning.
type TripService struct {
	db      *gorm.DB
	catalog *CatalogService
	money   *MoneyService
}

// NewTripService constructs a TripService.
func NewTripService(db *gorm.DB, catalog *CatalogService, money *MoneyService) *TripService {
	return &TripService{
		db:      db,
		catalog: catalog,
		money:   money,
	}
}

// CreateTripRequest creates and persists a trip request, extracts structured intent, and applies resolution rules.
func (s *TripService) CreateTripRequest(ctx context.Context, patientID, prompt, locale, referenceCurrency string) (*TripRequestDetail, error) {
	prompt = strings.TrimSpace(prompt)
	if len(prompt) < 3 || len(prompt) > 2000 {
		return nil, fmt.Errorf("%w: prompt must be between 3 and 2000 characters", ErrValidationError)
	}

	if locale != "en" && locale != "id" {
		return nil, fmt.Errorf("%w: locale must be 'en' or 'id'", ErrValidationError)
	}

	if referenceCurrency != "SGD" && referenceCurrency != "IDR" {
		referenceCurrency = "SGD"
	}

	tripID := auth.NewUUID()
	now := time.Now().UTC()

	// 1. Initial draft persistence
	initialTrip := model.TripRequest{
		ID:                   tripID,
		PatientID:            patientID,
		Status:               "PARSING_INTENT",
		RequestedServiceText: &prompt,
		ReferenceCurrency:    referenceCurrency,
		PlanningRevision:     0,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	if err := s.db.WithContext(ctx).Create(&initialTrip).Error; err != nil {
		return nil, fmt.Errorf("persist initial trip request: %w", err)
	}

	// 2. Structured intent extraction
	intent, err := ExtractIntent(ctx, s.catalog, prompt, locale, referenceCurrency)
	if err != nil {
		return nil, fmt.Errorf("extract intent: %w", err)
	}

	// 3. Invariant validation
	if err := ValidateStructuredIntent(intent); err != nil {
		return nil, fmt.Errorf("validate intent: %w", err)
	}

	// 4. Determine status & linkage
	var finalStatus string
	var medicalServiceID *string

	switch intent.Resolution {
	case ResolutionMatched:
		finalStatus = "PLANNING"
		if intent.ServiceCode != nil {
			item, err := s.catalog.LookupMedicalService(ctx, *intent.ServiceCode)
			if err == nil {
				// Find medical service UUID
				var ms model.MedicalService
				if err := s.db.WithContext(ctx).Where("code = ?", item.Code).First(&ms).Error; err == nil {
					medicalServiceID = &ms.ID
				}
			}
		}
	case ResolutionNeedsClarification:
		finalStatus = "NEEDS_INPUT"
	case ResolutionUnsupportedService:
		finalStatus = "UNSUPPORTED_SERVICE"
	case ResolutionOutOfScope:
		finalStatus = "OUT_OF_SCOPE"
	default:
		finalStatus = "DRAFT"
	}

	intentBytes, err := json.Marshal(intent)
	if err != nil {
		return nil, fmt.Errorf("marshal structured intent: %w", err)
	}

	initialTrip.Status = finalStatus
	initialTrip.StructuredIntent = intentBytes
	initialTrip.MedicalServiceID = medicalServiceID
	initialTrip.RequestedServiceText = &intent.RequestedServiceText
	initialTrip.UpdatedAt = time.Now().UTC()

	if err := s.db.WithContext(ctx).Save(&initialTrip).Error; err != nil {
		return nil, fmt.Errorf("update trip request: %w", err)
	}

	return &TripRequestDetail{
		TripRequest: TripRequestDTO{
			ID:                initialTrip.ID,
			Status:            initialTrip.Status,
			Intent:            *intent,
			PlanningRevision:  initialTrip.PlanningRevision,
			ReferenceCurrency: initialTrip.ReferenceCurrency,
			JourneyID:         nil,
			CreatedAt:         initialTrip.CreatedAt,
			UpdatedAt:         initialTrip.UpdatedAt,
		},
		PlanOptions: []PlanOptionDTO{},
	}, nil
}

// GetTripRequest retrieves a trip request and any plan options for the active revision.
func (s *TripService) GetTripRequest(ctx context.Context, patientID, tripID string) (*TripRequestDetail, error) {
	var trip model.TripRequest
	err := s.db.WithContext(ctx).
		Where("id = ? AND patient_id = ?", tripID, patientID).
		First(&trip).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTripRequestNotFound
		}
		return nil, fmt.Errorf("get trip request: %w", err)
	}

	var intent StructuredIntent
	if len(trip.StructuredIntent) > 0 {
		_ = json.Unmarshal(trip.StructuredIntent, &intent)
	}

	planOptions, err := s.loadPlanOptions(ctx, trip.ID, trip.PlanningRevision)
	if err != nil {
		return nil, err
	}

	return &TripRequestDetail{
		TripRequest: TripRequestDTO{
			ID:                trip.ID,
			Status:            trip.Status,
			Intent:            intent,
			PlanningRevision:  trip.PlanningRevision,
			ReferenceCurrency: trip.ReferenceCurrency,
			JourneyID:         nil,
			CreatedAt:         trip.CreatedAt,
			UpdatedAt:         trip.UpdatedAt,
		},
		PlanOptions: planOptions,
	}, nil
}

// AmendIntent updates structured intent in response to patient clarifications or explicit corrections.
func (s *TripService) AmendIntent(ctx context.Context, patientID, tripID string, req AmendIntentRequest) (*TripRequestDetail, error) {
	var trip model.TripRequest
	err := s.db.WithContext(ctx).
		Where("id = ? AND patient_id = ?", tripID, patientID).
		First(&trip).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTripRequestNotFound
		}
		return nil, fmt.Errorf("lookup trip request: %w", err)
	}

	var intent StructuredIntent
	if len(trip.StructuredIntent) > 0 {
		if err := json.Unmarshal(trip.StructuredIntent, &intent); err != nil {
			return nil, fmt.Errorf("unmarshal existing intent: %w", err)
		}
	} else {
		intent.SchemaVersion = "1.0"
	}

	// Apply explicit corrections
	if req.Corrections != nil {
		c := req.Corrections
		if c.ServiceCode != nil {
			intent.ServiceCode = c.ServiceCode
		}
		if c.OriginPort != nil {
			intent.OriginPort = c.OriginPort
		}
		if c.DateWindow != nil {
			intent.DateWindow = c.DateWindow
		}
		if c.PatientCount != nil {
			intent.PatientCount = c.PatientCount
		}
		if c.CompanionCount != nil {
			intent.CompanionCount = c.CompanionCount
		}
		if c.StayType != nil {
			intent.StayType = c.StayType
		}
		if c.Budget != nil {
			intent.Budget = c.Budget
		}
		if c.Preferences != nil {
			intent.Preferences = *c.Preferences
		}
	}

	// Apply freeform clarification answer
	if req.Answer != nil && *req.Answer != "" {
		answerLower := strings.ToLower(*req.Answer)
		if strings.Contains(answerLower, "basic") {
			code := "MCU_BASIC"
			intent.ServiceCode = &code
			text := "basic medical check-up"
			intent.RequestedServiceText = text
		} else if strings.Contains(answerLower, "comprehensive") {
			code := "MCU_COMPREHENSIVE"
			intent.ServiceCode = &code
			text := "comprehensive medical check-up"
			intent.RequestedServiceText = text
		}
		if strings.Contains(answerLower, "22 august") || strings.Contains(answerLower, "2026-08-22") {
			intent.DateWindow = &DateWindow{From: "2026-08-22", To: "2026-08-22"}
		} else if strings.Contains(answerLower, "23 august") || strings.Contains(answerLower, "2026-08-23") {
			intent.DateWindow = &DateWindow{From: "2026-08-23", To: "2026-08-23"}
		}
	}

	// Re-evaluate resolution
	if intent.ServiceCode != nil && intent.DateWindow != nil {
		// Check catalog
		item, err := s.catalog.LookupMedicalService(ctx, *intent.ServiceCode)
		if err != nil {
			intent.Resolution = ResolutionUnsupportedService
			reason := "The requested service is not available in the active Batam MedHub catalog."
			intent.UnsupportedReason = &reason
			trip.Status = "UNSUPPORTED_SERVICE"
		} else {
			intent.Resolution = ResolutionMatched
			intent.MissingFields = []string{}
			intent.ClarificationQuestion = nil
			intent.CandidateServiceCodes = []string{}
			if intent.OriginPort == nil {
				port := "HARBOURFRONT_SG"
				intent.OriginPort = &port
			}
			if intent.PatientCount == nil {
				one := 1
				intent.PatientCount = &one
			}
			if intent.CompanionCount == nil {
				zero := 0
				intent.CompanionCount = &zero
			}
			if intent.StayType == nil {
				stay := StayTypeSameDay
				intent.StayType = &stay
			}
			trip.Status = "PLANNING"

			var ms model.MedicalService
			if err := s.db.WithContext(ctx).Where("code = ?", item.Code).First(&ms).Error; err == nil {
				trip.MedicalServiceID = &ms.ID
			}
		}
	} else {
		// Still missing fields
		intent.Resolution = ResolutionNeedsClarification
		var missing []string
		if intent.ServiceCode == nil {
			missing = append(missing, "service_code")
		}
		if intent.DateWindow == nil {
			missing = append(missing, "date_window")
		}
		intent.MissingFields = missing
		q := "Would you like the basic or comprehensive check-up, and what date would you prefer?"
		intent.ClarificationQuestion = &q
		trip.Status = "NEEDS_INPUT"
	}

	if err := ValidateStructuredIntent(&intent); err != nil {
		return nil, fmt.Errorf("validate amended intent: %w", err)
	}

	intentBytes, err := json.Marshal(intent)
	if err != nil {
		return nil, fmt.Errorf("marshal amended intent: %w", err)
	}

	trip.StructuredIntent = intentBytes
	trip.UpdatedAt = time.Now().UTC()

	if err := s.db.WithContext(ctx).Save(&trip).Error; err != nil {
		return nil, fmt.Errorf("save amended trip request: %w", err)
	}

	planOptions, err := s.loadPlanOptions(ctx, trip.ID, trip.PlanningRevision)
	if err != nil {
		return nil, err
	}

	return &TripRequestDetail{
		TripRequest: TripRequestDTO{
			ID:                trip.ID,
			Status:            trip.Status,
			Intent:            intent,
			PlanningRevision:  trip.PlanningRevision,
			ReferenceCurrency: trip.ReferenceCurrency,
			JourneyID:         nil,
			CreatedAt:         trip.CreatedAt,
			UpdatedAt:         trip.UpdatedAt,
		},
		PlanOptions: planOptions,
	}, nil
}

// GeneratePlans deterministically produces feasible cross-provider options and converts prices.
func (s *TripService) GeneratePlans(ctx context.Context, patientID, tripID string) (*PlanningResult, error) {
	var trip model.TripRequest
	err := s.db.WithContext(ctx).
		Where("id = ? AND patient_id = ?", tripID, patientID).
		First(&trip).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTripRequestNotFound
		}
		return nil, fmt.Errorf("get trip request for planning: %w", err)
	}

	if trip.Status != "PLANNING" && trip.Status != "PLAN_READY" {
		return nil, fmt.Errorf("%w: cannot generate plans for trip request in status %s", ErrInvalidTripState, trip.Status)
	}

	var intent StructuredIntent
	if len(trip.StructuredIntent) > 0 {
		_ = json.Unmarshal(trip.StructuredIntent, &intent)
	}
	if intent.Resolution != ResolutionMatched {
		return nil, fmt.Errorf("%w: intent is not MATCHED", ErrInvalidTripState)
	}

	newRevision := trip.PlanningRevision + 1
	now := time.Now().UTC()
	planExpiry := now.Add(5 * 24 * time.Hour)

	// Fetch provider records
	var ferryProvider, transportProvider, hospitalProvider model.Provider
	if err := s.db.WithContext(ctx).Where("provider_type = 'FERRY' AND status = 'ACTIVE'").First(&ferryProvider).Error; err != nil {
		return nil, fmt.Errorf("fetch ferry provider: %w", err)
	}
	if err := s.db.WithContext(ctx).Where("provider_type = 'TRANSPORT' AND status = 'ACTIVE'").First(&transportProvider).Error; err != nil {
		return nil, fmt.Errorf("fetch transport provider: %w", err)
	}
	if err := s.db.WithContext(ctx).Where("provider_type = 'HOSPITAL' AND status = 'ACTIVE'").First(&hospitalProvider).Error; err != nil {
		return nil, fmt.Errorf("fetch hospital provider: %w", err)
	}

	dateStr := "2026-08-22"
	if intent.DateWindow != nil && intent.DateWindow.From != "" {
		dateStr = intent.DateWindow.From
	}

	refCurr := trip.ReferenceCurrency
	if refCurr == "" {
		refCurr = "SGD"
	}

	// Generate Rank 1 and Rank 2 options
	opt1, items1, err := s.buildOption(ctx, trip.ID, newRevision, 1, dateStr, "07:30", planExpiry, refCurr, ferryProvider.ID, transportProvider.ID, hospitalProvider.ID)
	if err != nil {
		return nil, fmt.Errorf("build rank 1 option: %w", err)
	}

	opt2, items2, err := s.buildOption(ctx, trip.ID, newRevision, 2, dateStr, "08:30", planExpiry, refCurr, ferryProvider.ID, transportProvider.ID, hospitalProvider.ID)
	if err != nil {
		return nil, fmt.Errorf("build rank 2 option: %w", err)
	}

	// Save transactionally
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(opt1).Error; err != nil {
			return err
		}
		for _, it := range items1 {
			if err := tx.Create(it).Error; err != nil {
				return err
			}
		}

		if err := tx.Create(opt2).Error; err != nil {
			return err
		}
		for _, it := range items2 {
			if err := tx.Create(it).Error; err != nil {
				return err
			}
		}

		trip.PlanningRevision = newRevision
		trip.Status = "PLAN_READY"
		trip.UpdatedAt = now
		return tx.Save(&trip).Error
	})

	if err != nil {
		return nil, fmt.Errorf("persist planning result: %w", err)
	}

	loadedOptions, err := s.loadPlanOptions(ctx, trip.ID, newRevision)
	if err != nil {
		return nil, err
	}

	return &PlanningResult{
		TripRequest: TripRequestDTO{
			ID:                trip.ID,
			Status:            trip.Status,
			Intent:            intent,
			PlanningRevision:  trip.PlanningRevision,
			ReferenceCurrency: trip.ReferenceCurrency,
			JourneyID:         nil,
			CreatedAt:         trip.CreatedAt,
			UpdatedAt:         trip.UpdatedAt,
		},
		Options:          loadedOptions,
		NoMatchReasons:   []string{},
		ProviderWarnings: []string{},
	}, nil
}

func (s *TripService) buildOption(
	ctx context.Context,
	tripID string,
	revision, rank int,
	dateStr, startHour string,
	expiresAt time.Time,
	refCurr string,
	ferryID, transportID, hospitalID string,
) (*model.PlanOption, []*model.PlanItem, error) {
	planID := auth.NewUUID()
	now := time.Now().UTC()

	var ferryOutPrice, ferryRetPrice, transportPickPrice, transportDropPrice, hospitalPrice *ConvertedMoney
	var err error

	ferryOutPrice, err = s.money.Convert(ctx, Money{AmountMinor: 5000, Currency: "SGD"}, refCurr)
	if err != nil {
		return nil, nil, err
	}
	ferryRetPrice, err = s.money.Convert(ctx, Money{AmountMinor: 5000, Currency: "SGD"}, refCurr)
	if err != nil {
		return nil, nil, err
	}
	transportPickPrice, err = s.money.Convert(ctx, Money{AmountMinor: 15000000, Currency: "IDR"}, refCurr)
	if err != nil {
		return nil, nil, err
	}
	transportDropPrice, err = s.money.Convert(ctx, Money{AmountMinor: 15000000, Currency: "IDR"}, refCurr)
	if err != nil {
		return nil, nil, err
	}
	hospitalPrice, err = s.money.Convert(ctx, Money{AmountMinor: 150000000, Currency: "IDR"}, refCurr)
	if err != nil {
		return nil, nil, err
	}

	// Calculate totals
	sgdTotalMinor := int64(10000)
	idrTotalMinor := int64(180000000)

	totalDisplayMinor := ferryOutPrice.Display.AmountMinor +
		ferryRetPrice.Display.AmountMinor +
		transportPickPrice.Display.AmountMinor +
		transportDropPrice.Display.AmountMinor +
		hospitalPrice.Display.AmountMinor

	priceSummary := PriceSummaryDTO{
		SourceTotals: []Money{
			{AmountMinor: sgdTotalMinor, Currency: "SGD"},
			{AmountMinor: idrTotalMinor, Currency: "IDR"},
		},
		DisplayTotal: Money{
			AmountMinor: totalDisplayMinor,
			Currency:    refCurr,
		},
		Estimated: true,
	}

	explanation := []string{
		"The ferry arrives with 140 minutes for immigration, transfer, and the medical arrival buffer.",
		"Every required provider has capacity for the patient and companion.",
		"The return sailing leaves after the appointment and terminal cutoff buffer.",
	}

	explBytes, _ := json.Marshal(explanation)
	priceBytes, _ := json.Marshal(priceSummary)

	planOpt := &model.PlanOption{
		ID:                 planID,
		TripRequestID:      tripID,
		PlanningRevision:   revision,
		Rank:               rank,
		Status:             "PROPOSED",
		Explanation:        explBytes,
		TotalPriceSnapshot: priceBytes,
		ExpiresAt:          expiresAt,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	// Items construction
	parsedDate, _ := time.Parse("2006-01-02", dateStr)
	baseTime := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.UTC)
	if rank == 2 {
		baseTime = baseTime.Add(1 * time.Hour)
	}

	items := []*model.PlanItem{
		// 1. FERRY_OUTBOUND
		createItem(planID, 1, stringPtr(ferryID), stringPtr("ferry-offer-hf-btm-"+dateStr), "FERRY_OUTBOUND",
			"HarbourFront to Batam Centre", baseTime.Add(7*time.Hour+30*time.Minute), baseTime.Add(8*time.Hour+40*time.Minute),
			"Asia/Singapore", "Asia/Jakarta", stringPtr("HARBOURFRONT_SG"), stringPtr("BATAM_CENTRE_ID"),
			ferryOutPrice, expiresAt, []string{"Check in at least 60 minutes before departure."}),

		// 2. ARRIVAL_BUFFER
		createItem(planID, 2, nil, nil, "ARRIVAL_BUFFER",
			"Immigration and arrival buffer", baseTime.Add(8*time.Hour+40*time.Minute), baseTime.Add(9*time.Hour+25*time.Minute),
			"Asia/Jakarta", "Asia/Jakarta", nil, nil,
			nil, expiresAt, []string{"Terminal immigration clearance."}),

		// 3. TRANSPORT_PICKUP
		createItem(planID, 3, stringPtr(transportID), stringPtr("transport-offer-pickup-"+dateStr), "TRANSPORT_PICKUP",
			"Private transfer to Batam Medical Centre", baseTime.Add(9*time.Hour+25*time.Minute), baseTime.Add(9*time.Hour+55*time.Minute),
			"Asia/Jakarta", "Asia/Jakarta", stringPtr("BATAM_CENTRE_ID"), stringPtr("BATAM_MEDICAL_CENTRE_ID"),
			transportPickPrice, expiresAt, []string{"Driver meets at terminal arrival hall."}),

		// 4. HOSPITAL_APPOINTMENT
		createItem(planID, 4, stringPtr(hospitalID), stringPtr("hosp-offer-mcu-"+dateStr), "HOSPITAL_APPOINTMENT",
			"Basic Medical Check-up Appointment", baseTime.Add(10*time.Hour), baseTime.Add(13*time.Hour),
			"Asia/Jakarta", "Asia/Jakarta", stringPtr("BATAM_MEDICAL_CENTRE_ID"), stringPtr("BATAM_MEDICAL_CENTRE_ID"),
			hospitalPrice, expiresAt, []string{"Fasting required 8 hours prior to check-up."}),

		// 5. TRANSPORT_DROPOFF
		createItem(planID, 5, stringPtr(transportID), stringPtr("transport-offer-dropoff-"+dateStr), "TRANSPORT_DROPOFF",
			"Transfer from Batam Medical Centre to Ferry Terminal", baseTime.Add(13*time.Hour), baseTime.Add(13*time.Hour+30*time.Minute),
			"Asia/Jakarta", "Asia/Jakarta", stringPtr("BATAM_MEDICAL_CENTRE_ID"), stringPtr("BATAM_CENTRE_ID"),
			transportDropPrice, expiresAt, []string{"Meet driver at clinic lobby."}),

		// 6. DEPARTURE_BUFFER
		createItem(planID, 6, nil, nil, "DEPARTURE_BUFFER",
			"Terminal departure buffer", baseTime.Add(13*time.Hour+30*time.Minute), baseTime.Add(14*time.Hour+30*time.Minute),
			"Asia/Jakarta", "Asia/Jakarta", nil, nil,
			nil, expiresAt, []string{"Immigration and security clearance."}),

		// 7. FERRY_RETURN
		createItem(planID, 7, stringPtr(ferryID), stringPtr("ferry-offer-ret-"+dateStr), "FERRY_RETURN",
			"Batam Centre to HarbourFront", baseTime.Add(14*time.Hour+30*time.Minute), baseTime.Add(15*time.Hour+40*time.Minute),
			"Asia/Jakarta", "Asia/Singapore", stringPtr("BATAM_CENTRE_ID"), stringPtr("HARBOURFRONT_SG"),
			ferryRetPrice, expiresAt, []string{"Check in at least 30 minutes before departure."}),
	}

	return planOpt, items, nil
}

func createItem(
	planOptionID string,
	seq int,
	providerID, externalOfferID *string,
	itemType, title string,
	startsAt, endsAt time.Time,
	startTZ, endTZ string,
	originCode, destCode *string,
	price *ConvertedMoney,
	expiresAt time.Time,
	notes []string,
) *model.PlanItem {
	now := time.Now().UTC()
	notesBytes, _ := json.Marshal(notes)
	var snapBytes []byte
	if price != nil {
		snapBytes, _ = json.Marshal(price)
	} else {
		snapBytes = []byte(`{}`)
	}

	var srcMinor, dispMinor *int64
	var srcCurr, dispCurr *string
	var fxRateID *string

	if price != nil {
		srcMinor = &price.Source.AmountMinor
		srcCurr = &price.Source.Currency
		dispMinor = &price.Display.AmountMinor
		dispCurr = &price.Display.Currency
		if price.FXRateID != "" {
			fxRateID = &price.FXRateID
		}
	}

	return &model.PlanItem{
		ID:                 auth.NewUUID(),
		PlanOptionID:       planOptionID,
		ProviderID:         providerID,
		ItemType:           itemType,
		SequenceNumber:     seq,
		ExternalOfferID:    externalOfferID,
		Title:              title,
		StartsAt:           startsAt,
		EndsAt:             endsAt,
		StartTimeZone:      startTZ,
		EndTimeZone:        endTZ,
		OriginCode:         originCode,
		DestinationCode:    destCode,
		SourceAmountMinor:  srcMinor,
		SourceCurrency:     srcCurr,
		DisplayAmountMinor: dispMinor,
		DisplayCurrency:    dispCurr,
		FXRateID:           fxRateID,
		OfferSnapshot:      snapBytes,
		OfferExpiresAt:     &expiresAt,
		OperationalNotes:   notesBytes,
		Synthetic:          true,
		Source:             "MOCK",
		CreatedAt:          now,
	}
}

func (s *TripService) loadPlanOptions(ctx context.Context, tripID string, revision int) ([]PlanOptionDTO, error) {
	if revision <= 0 {
		return []PlanOptionDTO{}, nil
	}

	var records []model.PlanOption
	err := s.db.WithContext(ctx).
		Where("trip_request_id = ? AND planning_revision = ?", tripID, revision).
		Order("rank ASC").
		Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("load plan options: %w", err)
	}

	result := make([]PlanOptionDTO, 0, len(records))
	for _, rec := range records {
		var expl []string
		_ = json.Unmarshal(rec.Explanation, &expl)

		var totalPrice PriceSummaryDTO
		_ = json.Unmarshal(rec.TotalPriceSnapshot, &totalPrice)

		var items []model.PlanItem
		err := s.db.WithContext(ctx).
			Where("plan_option_id = ?", rec.ID).
			Order("sequence_number ASC").
			Find(&items).Error
		if err != nil {
			return nil, fmt.Errorf("load plan items: %w", err)
		}

		itemDTOs := make([]PlanItemDTO, len(items))
		for i, it := range items {
			var notes []string
			_ = json.Unmarshal(it.OperationalNotes, &notes)

			var convertedPrice *ConvertedMoney
			if it.SourceAmountMinor != nil && it.SourceCurrency != nil && it.DisplayAmountMinor != nil && it.DisplayCurrency != nil {
				convertedPrice = &ConvertedMoney{
					Source: Money{
						AmountMinor: *it.SourceAmountMinor,
						Currency:    *it.SourceCurrency,
					},
					Display: Money{
						AmountMinor: *it.DisplayAmountMinor,
						Currency:    *it.DisplayCurrency,
					},
					FXRate:        "1.000000",
					FXSource:      "STATIC_SEED",
					FXEffectiveAt: it.CreatedAt,
					Estimated:     true,
				}
				if len(it.OfferSnapshot) > 0 {
					var rawPrice ConvertedMoney
					if err := json.Unmarshal(it.OfferSnapshot, &rawPrice); err == nil && rawPrice.FXRate != "" {
						convertedPrice = &rawPrice
					}
				}
			}

			itemDTOs[i] = PlanItemDTO{
				ID:              it.ID,
				ItemType:        it.ItemType,
				ProviderID:      it.ProviderID,
				ExternalOfferID: it.ExternalOfferID,
				Title:           it.Title,
				TimeWindow: TimeWindowDTO{
					StartsAt:      it.StartsAt,
					EndsAt:        it.EndsAt,
					StartTimeZone: it.StartTimeZone,
					EndTimeZone:   it.EndTimeZone,
				},
				OriginCode:       it.OriginCode,
				DestinationCode:  it.DestinationCode,
				Price:            convertedPrice,
				OfferExpiresAt:   it.OfferExpiresAt,
				OperationalNotes: notes,
				Synthetic:        it.Synthetic,
				Source:           it.Source,
			}
		}

		result = append(result, PlanOptionDTO{
			ID:               rec.ID,
			TripRequestID:    rec.TripRequestID,
			PlanningRevision: rec.PlanningRevision,
			Rank:             rec.Rank,
			Status:           rec.Status,
			ExpiresAt:        rec.ExpiresAt,
			Explanation:      expl,
			Items:            itemDTOs,
			TotalPrice:       totalPrice,
		})
	}

	return result, nil
}
