package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"batam-medhub/internal/model"

	"gorm.io/gorm"
)

var (
	ErrJourneyNotFound          = errors.New("journey not found")
	ErrItineraryVersionNotFound = errors.New("itinerary version not found")
)

// JourneyDTO represents the core API view of a Journey.
type JourneyDTO struct {
	ID                     string    `json:"id"`
	TripRequestID          string    `json:"trip_request_id"`
	Status                 string    `json:"status"`
	ActiveItineraryVersion int       `json:"active_itinerary_version"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// ItineraryItemDTO represents an item inside a confirmed itinerary.
type ItineraryItemDTO struct {
	ID                    string          `json:"id"`
	ItemType              string          `json:"item_type"`
	ProviderID            *string         `json:"provider_id"`
	ExternalReservationID *string         `json:"external_reservation_id"`
	Title                 string          `json:"title"`
	Status                string          `json:"status"`
	TimeWindow            TimeWindowDTO   `json:"time_window"`
	OriginCode            *string         `json:"origin_code"`
	DestinationCode       *string         `json:"destination_code"`
	Price                 *ConvertedMoney `json:"price"`
	OperationalNotes      []string        `json:"operational_notes"`
	Synthetic             bool            `json:"synthetic"`
	Source                string          `json:"source"`
}

// ItineraryVersionDTO represents an immutable itinerary version.
type ItineraryVersionDTO struct {
	ID                  string             `json:"id"`
	JourneyID           string             `json:"journey_id"`
	Version             int                `json:"version"`
	Status              string             `json:"status"`
	BasedOnDisruptionID *string            `json:"based_on_disruption_id"`
	TotalPrice          PriceSummaryDTO    `json:"total_price"`
	Items               []ItineraryItemDTO `json:"items"`
	CreatedAt           time.Time          `json:"created_at"`
}

// ItineraryVersionSummaryDTO represents a compact summary of an itinerary version.
type ItineraryVersionSummaryDTO struct {
	ID                  string    `json:"id"`
	Version             int       `json:"version"`
	Status              string    `json:"status"`
	BasedOnDisruptionID *string   `json:"based_on_disruption_id"`
	CreatedAt           time.Time `json:"created_at"`
}

// JourneyDetail wraps a journey along with its active and historical itinerary versions.
type JourneyDetail struct {
	Journey           JourneyDTO                   `json:"journey"`
	ActiveItinerary   ItineraryVersionDTO          `json:"active_itinerary"`
	ItineraryVersions []ItineraryVersionSummaryDTO `json:"itinerary_versions"`
}

// JourneyService provides access to journeys and immutable itinerary versions.
type JourneyService struct {
	db *gorm.DB
}

// NewJourneyService constructs a JourneyService.
func NewJourneyService(db *gorm.DB) *JourneyService {
	return &JourneyService{db: db}
}

// GetJourneyItinerary returns the active journey itinerary and all version summaries.
func (s *JourneyService) GetJourneyItinerary(ctx context.Context, patientID, journeyID string) (*JourneyDetail, error) {
	var j model.Journey
	err := s.db.WithContext(ctx).
		Where("id = ? AND patient_id = ?", journeyID, patientID).
		First(&j).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrJourneyNotFound
		}
		return nil, fmt.Errorf("lookup journey: %w", err)
	}

	var activeVersion model.ItineraryVersion
	err = s.db.WithContext(ctx).
		Where("journey_id = ? AND version_number = ?", j.ID, j.CurrentVersionNumber).
		First(&activeVersion).Error
	if err != nil {
		return nil, fmt.Errorf("lookup active itinerary version: %w", err)
	}

	var items []model.ItineraryItem
	err = s.db.WithContext(ctx).
		Where("itinerary_version_id = ?", activeVersion.ID).
		Order("sequence_number ASC").
		Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("load itinerary items: %w", err)
	}

	var allVersions []model.ItineraryVersion
	err = s.db.WithContext(ctx).
		Where("journey_id = ?", j.ID).
		Order("version_number ASC").
		Find(&allVersions).Error
	if err != nil {
		return nil, fmt.Errorf("load itinerary version history: %w", err)
	}

	providerMap := s.loadProviderCodeMap(ctx)
	itItemDTOs := s.mapItineraryItems(items, providerMap)

	var totalPrice PriceSummaryDTO
	_ = json.Unmarshal(activeVersion.TotalPriceSnapshot, &totalPrice)

	activeItineraryDTO := ItineraryVersionDTO{
		ID:                  activeVersion.ID,
		JourneyID:           activeVersion.JourneyID,
		Version:             activeVersion.VersionNumber,
		Status:              activeVersion.Status,
		BasedOnDisruptionID: activeVersion.SourceDisruptionID,
		TotalPrice:          totalPrice,
		Items:               itItemDTOs,
		CreatedAt:           activeVersion.CreatedAt,
	}

	versionSummaries := make([]ItineraryVersionSummaryDTO, len(allVersions))
	for i, v := range allVersions {
		versionSummaries[i] = ItineraryVersionSummaryDTO{
			ID:                  v.ID,
			Version:             v.VersionNumber,
			Status:              v.Status,
			BasedOnDisruptionID: v.SourceDisruptionID,
			CreatedAt:           v.CreatedAt,
		}
	}

	return &JourneyDetail{
		Journey: JourneyDTO{
			ID:                     j.ID,
			TripRequestID:          j.TripRequestID,
			Status:                 j.Status,
			ActiveItineraryVersion: j.CurrentVersionNumber,
			CreatedAt:              j.CreatedAt,
			UpdatedAt:              j.UpdatedAt,
		},
		ActiveItinerary:   activeItineraryDTO,
		ItineraryVersions: versionSummaries,
	}, nil
}

// GetJourneyItineraryVersion returns a specific immutable itinerary version for a patient's journey.
func (s *JourneyService) GetJourneyItineraryVersion(ctx context.Context, patientID, journeyID string, versionNumber int) (*ItineraryVersionDTO, error) {
	var j model.Journey
	err := s.db.WithContext(ctx).
		Where("id = ? AND patient_id = ?", journeyID, patientID).
		First(&j).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrJourneyNotFound
		}
		return nil, fmt.Errorf("lookup journey: %w", err)
	}

	var version model.ItineraryVersion
	err = s.db.WithContext(ctx).
		Where("journey_id = ? AND version_number = ?", j.ID, versionNumber).
		First(&version).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrItineraryVersionNotFound
		}
		return nil, fmt.Errorf("lookup itinerary version: %w", err)
	}

	var items []model.ItineraryItem
	err = s.db.WithContext(ctx).
		Where("itinerary_version_id = ?", version.ID).
		Order("sequence_number ASC").
		Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("load itinerary items: %w", err)
	}

	providerMap := s.loadProviderCodeMap(ctx)
	itItemDTOs := s.mapItineraryItems(items, providerMap)

	var totalPrice PriceSummaryDTO
	_ = json.Unmarshal(version.TotalPriceSnapshot, &totalPrice)

	return &ItineraryVersionDTO{
		ID:                  version.ID,
		JourneyID:           version.JourneyID,
		Version:             version.VersionNumber,
		Status:              version.Status,
		BasedOnDisruptionID: version.SourceDisruptionID,
		TotalPrice:          totalPrice,
		Items:               itItemDTOs,
		CreatedAt:           version.CreatedAt,
	}, nil
}

// ListJourneys lists all journeys owned by the patient.
func (s *JourneyService) ListJourneys(ctx context.Context, patientID string) ([]JourneyDTO, error) {
	var list []model.Journey
	err := s.db.WithContext(ctx).
		Where("patient_id = ?", patientID).
		Order("created_at DESC").
		Find(&list).Error
	if err != nil {
		return nil, fmt.Errorf("list journeys: %w", err)
	}

	result := make([]JourneyDTO, len(list))
	for i, j := range list {
		result[i] = JourneyDTO{
			ID:                     j.ID,
			TripRequestID:          j.TripRequestID,
			Status:                 j.Status,
			ActiveItineraryVersion: j.CurrentVersionNumber,
			CreatedAt:              j.CreatedAt,
			UpdatedAt:              j.UpdatedAt,
		}
	}
	return result, nil
}

func (s *JourneyService) loadProviderCodeMap(ctx context.Context) map[string]string {
	var providers []model.Provider
	_ = s.db.WithContext(ctx).Find(&providers).Error
	m := make(map[string]string, len(providers))
	for _, p := range providers {
		m[p.ID] = p.Code
	}
	return m
}

func (s *JourneyService) mapItineraryItems(items []model.ItineraryItem, providerMap map[string]string) []ItineraryItemDTO {
	result := make([]ItineraryItemDTO, len(items))
	for i, it := range items {
		var convertedPrice *ConvertedMoney
		if it.SourceAmountMinor != nil && it.SourceCurrency != nil && it.DisplayAmountMinor != nil && it.DisplayCurrency != nil {
			fxVal := "1.000000"
			if it.FXRateValue != nil {
				fxVal = *it.FXRateValue
			}
			fxSrc := "DEMO_STATIC_2026_08"
			if it.FXSource != nil {
				fxSrc = *it.FXSource
			}
			fxEff := it.CreatedAt
			if it.FXEffectiveAt != nil {
				fxEff = *it.FXEffectiveAt
			}

			convertedPrice = &ConvertedMoney{
				Source: Money{
					AmountMinor: *it.SourceAmountMinor,
					Currency:    *it.SourceCurrency,
				},
				Display: Money{
					AmountMinor: *it.DisplayAmountMinor,
					Currency:    *it.DisplayCurrency,
				},
				FXRate:        fxVal,
				FXSource:      fxSrc,
				FXEffectiveAt: fxEff,
				Estimated:     true,
			}
		}

		var providerCode *string
		if it.ProviderID != nil {
			if code, ok := providerMap[*it.ProviderID]; ok {
				providerCode = &code
			} else {
				providerCode = it.ProviderID
			}
		}

		notes := []string{}
		if len(it.Snapshot) > 0 {
			var snapData struct {
				OperationalNotes []string `json:"operational_notes"`
			}
			if err := json.Unmarshal(it.Snapshot, &snapData); err == nil && snapData.OperationalNotes != nil {
				notes = snapData.OperationalNotes
			}
		}

		result[i] = ItineraryItemDTO{
			ID:                    it.ID,
			ItemType:              it.ItemType,
			ProviderID:            providerCode,
			ExternalReservationID: it.ExternalReservationID,
			Title:                 it.Title,
			Status:                it.Status,
			TimeWindow: TimeWindowDTO{
				StartsAt:      it.StartsAt,
				EndsAt:        it.EndsAt,
				StartTimeZone: it.StartTimeZone,
				EndTimeZone:   it.EndTimeZone,
			},
			OriginCode:       it.OriginCode,
			DestinationCode:  it.DestinationCode,
			Price:            convertedPrice,
			OperationalNotes: notes,
			Synthetic:        true,
			Source:           "MOCK",
		}
	}
	return result
}
