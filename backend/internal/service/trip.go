package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"batam-medhub/internal/adapter"
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
	db         *gorm.DB
	catalog    *CatalogService
	money      *MoneyService
	aggregator *adapter.Aggregator
	ai         IntentExtractor
}

// NewTripService constructs a TripService.
func NewTripService(db *gorm.DB, catalog *CatalogService, money *MoneyService, aggregator *adapter.Aggregator, ai IntentExtractor) *TripService {
	return &TripService{
		db:         db,
		catalog:    catalog,
		money:      money,
		aggregator: aggregator,
		ai:         ai,
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
	var intent *StructuredIntent
	var err error
	if s.ai != nil {
		intent, err = s.ai.ExtractIntent(ctx, prompt, locale, referenceCurrency)
	} else {
		intent, err = ExtractIntentDeterministic(ctx, s.catalog, prompt, locale, referenceCurrency)
	}
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

	var journeyID *string
	var j model.Journey
	if err := s.db.WithContext(ctx).Where("trip_request_id = ?", trip.ID).First(&j).Error; err == nil {
		journeyID = &j.ID
	}

	return &TripRequestDetail{
		TripRequest: TripRequestDTO{
			ID:                trip.ID,
			Status:            trip.Status,
			Intent:            intent,
			PlanningRevision:  trip.PlanningRevision,
			ReferenceCurrency: trip.ReferenceCurrency,
			JourneyID:         journeyID,
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
		answerText := strings.TrimSpace(*req.Answer)
		if s.ai != nil {
			combinedPrompt := intent.RequestedServiceText + ". Additional patient clarification: " + answerText
			if extracted, err := s.ai.ExtractIntent(ctx, combinedPrompt, "en", trip.ReferenceCurrency); err == nil && extracted != nil {
				if extracted.ServiceCode != nil {
					intent.ServiceCode = extracted.ServiceCode
					intent.RequestedServiceText = extracted.RequestedServiceText
				}
				if extracted.DateWindow != nil {
					intent.DateWindow = extracted.DateWindow
				}
				if extracted.OriginPort != nil {
					intent.OriginPort = extracted.OriginPort
				}
				if extracted.PatientCount != nil {
					intent.PatientCount = extracted.PatientCount
				}
				if extracted.CompanionCount != nil {
					intent.CompanionCount = extracted.CompanionCount
				}
				if extracted.StayType != nil {
					intent.StayType = extracted.StayType
				}
				if extracted.Budget != nil {
					intent.Budget = extracted.Budget
				}
			}
		}

		answerLower := strings.ToLower(answerText)
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

type builtPackage struct {
	option            *model.PlanOption
	items             []*model.PlanItem
	displayTotalMinor int64
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

	if trip.Status != "PLANNING" && trip.Status != "PLAN_READY" && trip.Status != "NO_MATCH" {
		return nil, fmt.Errorf("%w: cannot generate plans for trip request in status %s", ErrInvalidTripState, trip.Status)
	}

	var intent StructuredIntent
	if len(trip.StructuredIntent) > 0 {
		_ = json.Unmarshal(trip.StructuredIntent, &intent)
	}
	if intent.Resolution != ResolutionMatched {
		return nil, fmt.Errorf("%w: intent is not MATCHED", ErrInvalidTripState)
	}

	dateStr := "2026-08-22"
	if intent.DateWindow != nil && intent.DateWindow.From != "" {
		dateStr = intent.DateWindow.From
	}

	originPort := "HARBOURFRONT_SG"
	if intent.OriginPort != nil && *intent.OriginPort != "" {
		originPort = *intent.OriginPort
	}

	serviceCode := "MCU_BASIC"
	if intent.ServiceCode != nil && *intent.ServiceCode != "" {
		serviceCode = *intent.ServiceCode
	}

	patientCount := 1
	if intent.PatientCount != nil && *intent.PatientCount > 0 {
		patientCount = *intent.PatientCount
	}
	companionCount := 0
	if intent.CompanionCount != nil && *intent.CompanionCount >= 0 {
		companionCount = *intent.CompanionCount
	}
	totalPax := patientCount + companionCount

	refCurr := trip.ReferenceCurrency
	if refCurr == "" {
		refCurr = "SGD"
	}

	// Fetch provider records from DB
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

	var providerWarnings []string
	if s.aggregator != nil {
		reqID := "req-plan-" + auth.NewUUID()
		multiQuery := adapter.MultiSearchQuery{
			HospitalCriteria: &adapter.HospitalSearchCriteria{
				ServiceCode:  serviceCode,
				PatientCount: patientCount,
				AppointmentWindow: adapter.TimeWindow{
					StartsAt:      dateStr + "T01:00:00Z",
					EndsAt:        dateStr + "T11:00:00Z",
					StartTimeZone: "Asia/Jakarta",
					EndTimeZone:   "Asia/Jakarta",
				},
			},
			FerryCriteria: []adapter.FerrySearchCriteria{
				{
					OriginPortCode:      originPort,
					DestinationPortCode: "BATAM_CENTRE_ID",
					PassengerCount:      totalPax,
					DepartureWindow: adapter.TimeWindow{
						StartsAt:      dateStr + "T00:00:00Z",
						EndsAt:        dateStr + "T06:00:00Z",
						StartTimeZone: "Asia/Singapore",
						EndTimeZone:   "Asia/Jakarta",
					},
				},
				{
					OriginPortCode:      "BATAM_CENTRE_ID",
					DestinationPortCode: originPort,
					PassengerCount:      totalPax,
					DepartureWindow: adapter.TimeWindow{
						StartsAt:      dateStr + "T06:00:00Z",
						EndsAt:        dateStr + "T16:00:00Z",
						StartTimeZone: "Asia/Jakarta",
						EndTimeZone:   "Asia/Singapore",
					},
				},
			},
			TransportCriteria: []adapter.TransportSearchCriteria{
				{
					PickupLocationCode:  "BATAM_CENTRE_ID",
					DropoffLocationCode: "BATAM_MEDICAL_CENTRE_ID",
					PassengerCount:      totalPax,
					PickupWindow: adapter.TimeWindow{
						StartsAt:      dateStr + "T01:00:00Z",
						EndsAt:        dateStr + "T04:00:00Z",
						StartTimeZone: "Asia/Jakarta",
						EndTimeZone:   "Asia/Jakarta",
					},
				},
				{
					PickupLocationCode:  "HOSPITAL_DEMO_ID",
					DropoffLocationCode: "BATAM_CENTRE_ID",
					PassengerCount:      totalPax,
					PickupWindow: adapter.TimeWindow{
						StartsAt:      dateStr + "T05:00:00Z",
						EndsAt:        dateStr + "T09:00:00Z",
						StartTimeZone: "Asia/Jakarta",
						EndTimeZone:   "Asia/Jakarta",
					},
				},
			},
		}

		res := s.aggregator.SearchAll(ctx, reqID, multiQuery)
		providerWarnings = res.Warnings
	}

	newRevision := trip.PlanningRevision + 1
	now := time.Now().UTC()
	planExpiry := now.Add(5 * 24 * time.Hour)

	// Check if date is in the past or invalid
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil || parsedDate.Before(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)) {
		trip.Status = "NO_MATCH"
		trip.UpdatedAt = now
		_ = s.db.WithContext(ctx).Save(&trip).Error

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
			Options:          []PlanOptionDTO{},
			NoMatchReasons:   []string{"The requested date window is not available in the active schedule."},
			ProviderWarnings: providerWarnings,
		}, nil
	}

	// Generate candidate packages
	var candidates []builtPackage

	// Package 1: Morning departure (07:30 SGT) -> 10:00 appointment
	opt1, items1, total1, err := s.assemblePackage(ctx, trip.ID, newRevision, 1, dateStr, 0, planExpiry, refCurr, totalPax, ferryProvider.ID, transportProvider.ID, hospitalProvider.ID)
	if err == nil {
		candidates = append(candidates, builtPackage{
			option:            opt1,
			items:             items1,
			displayTotalMinor: total1,
		})
	}

	// Package 2: Later departure (08:30 SGT) -> 11:00 appointment
	opt2, items2, total2, err := s.assemblePackage(ctx, trip.ID, newRevision, 2, dateStr, 1*time.Hour, planExpiry, refCurr, totalPax, ferryProvider.ID, transportProvider.ID, hospitalProvider.ID)
	if err == nil {
		candidates = append(candidates, builtPackage{
			option:            opt2,
			items:             items2,
			displayTotalMinor: total2,
		})
	}

	if len(candidates) == 0 {
		trip.Status = "NO_MATCH"
		trip.UpdatedAt = now
		_ = s.db.WithContext(ctx).Save(&trip).Error

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
			Options:          []PlanOptionDTO{},
			NoMatchReasons:   []string{"No available appointments or ferry sailings satisfy the required 45-minute arrival and 30-minute departure buffer constraints on the requested date."},
			ProviderWarnings: providerWarnings,
		}, nil
	}

	// Sort candidates deterministically by display total price ascending
	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].displayTotalMinor < candidates[j].displayTotalMinor
	})

	// Assign ranks (1, 2)
	selected := candidates
	if len(selected) > 2 {
		selected = selected[:2]
	}
	for i := range selected {
		selected[i].option.Rank = i + 1
	}

	// Save transactionally
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, pkg := range selected {
			if err := tx.Create(pkg.option).Error; err != nil {
				return err
			}
			for _, it := range pkg.items {
				if err := tx.Create(it).Error; err != nil {
					return err
				}
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
		ProviderWarnings: providerWarnings,
	}, nil
}

func (s *TripService) assemblePackage(
	ctx context.Context,
	tripID string,
	revision, rank int,
	dateStr string,
	timeShift time.Duration,
	expiresAt time.Time,
	refCurr string,
	totalPax int,
	ferryID, transportID, hospitalID string,
) (*model.PlanOption, []*model.PlanItem, int64, error) {
	planID := auth.NewUUID()
	now := time.Now().UTC()

	var ferryOutPrice, ferryRetPrice, transportPickPrice, transportDropPrice, hospitalPrice *ConvertedMoney
	var err error

	ferryOutPrice, err = s.money.Convert(ctx, Money{AmountMinor: 2500 * int64(totalPax), Currency: "SGD"}, refCurr)
	if err != nil {
		return nil, nil, 0, err
	}
	ferryRetPrice, err = s.money.Convert(ctx, Money{AmountMinor: 2500 * int64(totalPax), Currency: "SGD"}, refCurr)
	if err != nil {
		return nil, nil, 0, err
	}
	transportPickPrice, err = s.money.Convert(ctx, Money{AmountMinor: 15000000, Currency: "IDR"}, refCurr)
	if err != nil {
		return nil, nil, 0, err
	}
	transportDropPrice, err = s.money.Convert(ctx, Money{AmountMinor: 15000000, Currency: "IDR"}, refCurr)
	if err != nil {
		return nil, nil, 0, err
	}
	hospitalPrice, err = s.money.Convert(ctx, Money{AmountMinor: 150000000, Currency: "IDR"}, refCurr)
	if err != nil {
		return nil, nil, 0, err
	}

	sgdTotalMinor := int64(5000) * int64(totalPax)
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
	if rank == 2 {
		explanation = []string{
			"Alternative schedule with a later departure time.",
			"Every required provider has capacity for the patient and companion.",
			"The return sailing leaves after the appointment and terminal cutoff buffer.",
		}
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

	parsedDate, _ := time.Parse("2006-01-02", dateStr)
	baseTime := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.UTC).Add(timeShift)

	// Items construction adhering to buffer invariants:
	// Ferry arrives at baseTime + 08:40 WIB
	// Buffer 1 (Arrival & immigration): 45 mins -> ends 09:25 WIB
	// Transport pickup: 09:25 to 09:55 WIB -> arrives hospital before 10:00 WIB
	// Hospital appointment: 10:00 to 13:00 WIB (3 hrs)
	// Transport dropoff: starts 13:30 WIB (30-min release buffer), arrives terminal 14:00 WIB
	// Departure buffer: 14:00 to 14:30 WIB (30-min terminal buffer)
	// Return ferry: departs 14:30 WIB
	cleanDate := strings.ReplaceAll(dateStr, "-", "")
	hospitalOffer := "hospital-offer-mcu-basic-" + cleanDate + "-1000"
	hospPrice := hospitalPrice
	ferryRetOffer := "ferry-offer-btm-hf-" + cleanDate + "-1430"
	if rank == 2 {
		hospitalOffer = "hospital-offer-mcu-basic-" + cleanDate + "-1300"
		hospPrice, _ = s.money.Convert(ctx, Money{AmountMinor: 200000000, Currency: "IDR"}, refCurr)
		ferryRetOffer = "ferry-offer-btm-hf-" + cleanDate + "-1600"
	}

	items := []*model.PlanItem{
		// 1. FERRY_OUTBOUND
		createItem(planID, 1, stringPtr(ferryID), stringPtr("ferry-offer-hf-btm-"+cleanDate+"-0730"), "FERRY_OUTBOUND",
			"HarbourFront to Batam Centre", baseTime.Add(7*time.Hour+30*time.Minute), baseTime.Add(8*time.Hour+40*time.Minute),
			"Asia/Singapore", "Asia/Jakarta", stringPtr("HARBOURFRONT_SG"), stringPtr("BATAM_CENTRE_ID"),
			ferryOutPrice, expiresAt, []string{"Check in at least 60 minutes before departure."}),

		// 2. ARRIVAL_BUFFER (45 min immigration & buffer)
		createItem(planID, 2, nil, nil, "ARRIVAL_BUFFER",
			"Immigration and arrival buffer", baseTime.Add(8*time.Hour+40*time.Minute), baseTime.Add(9*time.Hour+25*time.Minute),
			"Asia/Jakarta", "Asia/Jakarta", stringPtr("BATAM_CENTRE_ID"), stringPtr("BATAM_CENTRE_ID"),
			nil, expiresAt, []string{"Terminal immigration clearance."}),

		// 3. TRANSPORT_PICKUP
		createItem(planID, 3, stringPtr(transportID), stringPtr("transport-offer-btm-hospital-"+cleanDate+"-0825"), "TRANSPORT_PICKUP",
			"Private transfer to Batam Medical Centre", baseTime.Add(9*time.Hour+25*time.Minute), baseTime.Add(9*time.Hour+55*time.Minute),
			"Asia/Jakarta", "Asia/Jakarta", stringPtr("BATAM_CENTRE_ID"), stringPtr("HOSPITAL_DEMO_ID"),
			transportPickPrice, expiresAt, []string{"Driver meets at terminal arrival hall."}),

		// 4. HOSPITAL_APPOINTMENT
		createItem(planID, 4, stringPtr(hospitalID), stringPtr(hospitalOffer), "HOSPITAL_APPOINTMENT",
			"Basic Medical Check-up Appointment", baseTime.Add(10*time.Hour), baseTime.Add(13*time.Hour),
			"Asia/Jakarta", "Asia/Jakarta", stringPtr("HOSPITAL_DEMO_ID"), stringPtr("HOSPITAL_DEMO_ID"),
			hospPrice, expiresAt, []string{"Fasting required 8 hours prior to check-up."}),

		// 5. TRANSPORT_DROPOFF (30 min post-appointment buffer)
		createItem(planID, 5, stringPtr(transportID), stringPtr("transport-offer-hospital-btm-"+cleanDate+"-1345"), "TRANSPORT_DROPOFF",
			"Transfer from Batam Medical Centre to Ferry Terminal", baseTime.Add(13*time.Hour+30*time.Minute), baseTime.Add(14*time.Hour),
			"Asia/Jakarta", "Asia/Jakarta", stringPtr("HOSPITAL_DEMO_ID"), stringPtr("BATAM_CENTRE_ID"),
			transportDropPrice, expiresAt, []string{"Meet driver at clinic lobby."}),

		// 6. DEPARTURE_BUFFER (30 min terminal check-in buffer)
		createItem(planID, 6, nil, nil, "DEPARTURE_BUFFER",
			"Terminal departure buffer", baseTime.Add(14*time.Hour), baseTime.Add(14*time.Hour+30*time.Minute),
			"Asia/Jakarta", "Asia/Jakarta", stringPtr("BATAM_CENTRE_ID"), stringPtr("BATAM_CENTRE_ID"),
			nil, expiresAt, []string{"Immigration and security clearance."}),

		// 7. FERRY_RETURN
		createItem(planID, 7, stringPtr(ferryID), stringPtr(ferryRetOffer), "FERRY_RETURN",
			"Batam Centre to HarbourFront", baseTime.Add(14*time.Hour+30*time.Minute), baseTime.Add(15*time.Hour+40*time.Minute),
			"Asia/Jakarta", "Asia/Singapore", stringPtr("BATAM_CENTRE_ID"), stringPtr("HARBOURFRONT_SG"),
			ferryRetPrice, expiresAt, []string{"Check in at least 30 minutes before departure."}),
	}

	return planOpt, items, totalDisplayMinor, nil
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
