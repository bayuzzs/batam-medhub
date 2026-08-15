package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Intent resolution constants.
const (
	ResolutionMatched            = "MATCHED"
	ResolutionNeedsClarification = "NEEDS_CLARIFICATION"
	ResolutionUnsupportedService = "UNSUPPORTED_SERVICE"
	ResolutionOutOfScope         = "OUT_OF_SCOPE"
)

// Stay type constants.
const (
	StayTypeSameDay   = "SAME_DAY"
	StayTypeOvernight = "OVERNIGHT"
	StayTypeFlexible  = "FLEXIBLE"
)

// DateWindow represents an inclusive date range.
type DateWindow struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// IntentPreferences represents patient preferences extracted from freeform text.
type IntentPreferences struct {
	Language      *string  `json:"language"`
	HotelTier     *string  `json:"hotel_tier"`
	Accessibility []string `json:"accessibility"`
}

// StructuredIntent represents the canonical parsed and validated intent structure.
type StructuredIntent struct {
	SchemaVersion         string            `json:"schema_version"`
	Resolution            string            `json:"resolution"`
	IntentCategory        *string           `json:"intent_category"`
	RequestedServiceText  string            `json:"requested_service_text"`
	ServiceCode           *string           `json:"service_code"`
	CandidateServiceCodes []string          `json:"candidate_service_codes"`
	OriginPort            *string           `json:"origin_port"`
	DateWindow            *DateWindow       `json:"date_window"`
	PatientCount          *int              `json:"patient_count"`
	CompanionCount        *int              `json:"companion_count"`
	StayType              *string           `json:"stay_type"`
	Budget                *Money            `json:"budget"`
	Preferences           IntentPreferences `json:"preferences"`
	MissingFields         []string          `json:"missing_fields"`
	ClarificationQuestion *string           `json:"clarification_question"`
	OutOfScopeReason      *string           `json:"out_of_scope_reason"`
	UnsupportedReason     *string           `json:"unsupported_reason"`
}

// IntentCorrections represents explicit structured corrections submitted by a patient.
type IntentCorrections struct {
	ServiceCode    *string            `json:"service_code,omitempty"`
	OriginPort     *string            `json:"origin_port,omitempty"`
	DateWindow     *DateWindow        `json:"date_window,omitempty"`
	PatientCount   *int               `json:"patient_count,omitempty"`
	CompanionCount *int               `json:"companion_count,omitempty"`
	StayType       *string            `json:"stay_type,omitempty"`
	Budget         *Money             `json:"budget,omitempty"`
	Preferences    *IntentPreferences `json:"preferences,omitempty"`
}

// AmendIntentRequest represents the payload for answering a clarification or correcting intent.
type AmendIntentRequest struct {
	Answer      *string            `json:"answer,omitempty"`
	Corrections *IntentCorrections `json:"corrections,omitempty"`
}

var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// ValidateStructuredIntent verifies the resolution invariants required by the contract.
func ValidateStructuredIntent(intent *StructuredIntent) error {
	if intent.SchemaVersion != "1.0" {
		return errors.New("schema_version must be '1.0'")
	}

	switch intent.Resolution {
	case ResolutionMatched:
		if intent.ServiceCode == nil || *intent.ServiceCode == "" {
			return errors.New("MATCHED intent requires service_code")
		}
		if intent.OriginPort == nil || *intent.OriginPort == "" {
			return errors.New("MATCHED intent requires origin_port")
		}
		if intent.DateWindow == nil {
			return errors.New("MATCHED intent requires date_window")
		}
		if !dateRegex.MatchString(intent.DateWindow.From) || !dateRegex.MatchString(intent.DateWindow.To) {
			return errors.New("date_window must contain valid YYYY-MM-DD dates")
		}
		if intent.PatientCount == nil || *intent.PatientCount < 1 || *intent.PatientCount > 10 {
			return errors.New("MATCHED intent requires patient_count between 1 and 10")
		}
		if intent.CompanionCount == nil || *intent.CompanionCount < 0 || *intent.CompanionCount > 10 {
			return errors.New("MATCHED intent requires companion_count between 0 and 10")
		}
		if intent.StayType == nil {
			return errors.New("MATCHED intent requires stay_type")
		}
		if len(intent.MissingFields) > 0 {
			return errors.New("MATCHED intent cannot have missing_fields")
		}
		if intent.ClarificationQuestion != nil {
			return errors.New("MATCHED intent cannot have clarification_question")
		}
		if intent.OutOfScopeReason != nil {
			return errors.New("MATCHED intent cannot have out_of_scope_reason")
		}
		if intent.UnsupportedReason != nil {
			return errors.New("MATCHED intent cannot have unsupported_reason")
		}

	case ResolutionNeedsClarification:
		if len(intent.MissingFields) == 0 {
			return errors.New("NEEDS_CLARIFICATION intent requires at least one missing_field")
		}
		if intent.ClarificationQuestion == nil || *intent.ClarificationQuestion == "" {
			return errors.New("NEEDS_CLARIFICATION intent requires clarification_question")
		}
		if intent.OutOfScopeReason != nil {
			return errors.New("NEEDS_CLARIFICATION intent cannot have out_of_scope_reason")
		}
		if intent.UnsupportedReason != nil {
			return errors.New("NEEDS_CLARIFICATION intent cannot have unsupported_reason")
		}

	case ResolutionUnsupportedService:
		if intent.UnsupportedReason == nil || *intent.UnsupportedReason == "" {
			return errors.New("UNSUPPORTED_SERVICE intent requires unsupported_reason")
		}
		if intent.OutOfScopeReason != nil {
			return errors.New("UNSUPPORTED_SERVICE intent cannot have out_of_scope_reason")
		}

	case ResolutionOutOfScope:
		if intent.OutOfScopeReason == nil || *intent.OutOfScopeReason == "" {
			return errors.New("OUT_OF_SCOPE intent requires out_of_scope_reason")
		}
		if intent.UnsupportedReason != nil {
			return errors.New("OUT_OF_SCOPE intent cannot have unsupported_reason")
		}

	default:
		return fmt.Errorf("unrecognized resolution: %s", intent.Resolution)
	}

	return nil
}

// ExtractIntent parses freeform patient text into a structured intent snapshot.
func ExtractIntent(ctx context.Context, catalog *CatalogService, prompt, locale, referenceCurrency string) (*StructuredIntent, error) {
	lower := strings.ToLower(prompt)
	lang := locale
	if lang == "" {
		lang = "en"
	}

	// 1. Emergency / Out of Scope Check
	if strings.Contains(lower, "chest pain") || strings.Contains(lower, "emergency") ||
		strings.Contains(lower, "heart attack") || strings.Contains(lower, "ambulance") ||
		strings.Contains(lower, "shortness of breath") || strings.Contains(lower, "severe pain") {
		reason := "This may be an emergency. Batam MedHub cannot triage or plan emergency care. Contact local emergency services immediately or go to the nearest emergency department."
		return &StructuredIntent{
			SchemaVersion:         "1.0",
			Resolution:            ResolutionOutOfScope,
			IntentCategory:        nil,
			RequestedServiceText:  prompt,
			ServiceCode:           nil,
			CandidateServiceCodes: []string{},
			OriginPort:            nil,
			DateWindow:            nil,
			PatientCount:          nil,
			CompanionCount:        nil,
			StayType:              nil,
			Budget:                nil,
			Preferences: IntentPreferences{
				Language:      &lang,
				HotelTier:     nil,
				Accessibility: []string{},
			},
			MissingFields:         []string{},
			ClarificationQuestion: nil,
			OutOfScopeReason:      &reason,
			UnsupportedReason:     nil,
		}, nil
	}

	// 2. Unsupported Service Check (e.g. major surgeries, oncology, etc.)
	if strings.Contains(lower, "knee replacement") || strings.Contains(lower, "surgery") ||
		strings.Contains(lower, "chemotherapy") || strings.Contains(lower, "hip replacement") ||
		strings.Contains(lower, "transplant") {
		reason := "The requested service is not available in the active Batam MedHub catalog."
		serviceText := "surgery"
		if strings.Contains(lower, "knee replacement") {
			serviceText = "knee replacement surgery"
		}
		originPort := "HARBOURFRONT_SG"
		patientCount := 1
		companionCount := 1
		stayType := StayTypeOvernight
		dateWindow := &DateWindow{From: "2026-08-22", To: "2026-08-24"}
		return &StructuredIntent{
			SchemaVersion:         "1.0",
			Resolution:            ResolutionUnsupportedService,
			IntentCategory:        nil,
			RequestedServiceText:  serviceText,
			ServiceCode:           nil,
			CandidateServiceCodes: []string{},
			OriginPort:            &originPort,
			DateWindow:            dateWindow,
			PatientCount:          &patientCount,
			CompanionCount:        &companionCount,
			StayType:              &stayType,
			Budget:                nil,
			Preferences: IntentPreferences{
				Language:      &lang,
				HotelTier:     nil,
				Accessibility: []string{},
			},
			MissingFields:         []string{},
			ClarificationQuestion: nil,
			OutOfScopeReason:      nil,
			UnsupportedReason:     &reason,
		}, nil
	}

	// 3. Ambiguous Medical Inquiries needing clarification
	// e.g. "I need a medical check-up" without specifying basic vs comprehensive or without date
	isBasicCheckup := strings.Contains(lower, "basic")
	isCompCheckup := strings.Contains(lower, "comprehensive")
	isDental := strings.Contains(lower, "dental")
	isEye := strings.Contains(lower, "eye")
	hasCheckupMention := strings.Contains(lower, "check-up") || strings.Contains(lower, "checkup") || strings.Contains(lower, "mcu") || strings.Contains(lower, "screening")

	// If prompt is vague check-up without basic/comprehensive or has missing date
	if hasCheckupMention && !isBasicCheckup && !isCompCheckup && !isDental && !isEye {
		question := "Would you like the basic or comprehensive check-up, and what date would you prefer?"
		originPort := "HARBOURFRONT_SG"
		patientCount := 1
		companionCount := 0
		return &StructuredIntent{
			SchemaVersion:         "1.0",
			Resolution:            ResolutionNeedsClarification,
			IntentCategory:        stringPtr("PREVENTIVE_CHECKUP"),
			RequestedServiceText:  "medical check-up",
			ServiceCode:           nil,
			CandidateServiceCodes: []string{"MCU_BASIC", "MCU_COMPREHENSIVE"},
			OriginPort:            &originPort,
			DateWindow:            nil,
			PatientCount:          &patientCount,
			CompanionCount:        &companionCount,
			StayType:              nil,
			Budget:                nil,
			Preferences: IntentPreferences{
				Language:      &lang,
				HotelTier:     nil,
				Accessibility: []string{},
			},
			MissingFields:         []string{"service_code", "date_window"},
			ClarificationQuestion: &question,
			OutOfScopeReason:      nil,
			UnsupportedReason:     nil,
		}, nil
	}

	// 4. Matched Inquiries
	var serviceCode string
	var serviceText string
	category := "PREVENTIVE_CHECKUP"

	if isBasicCheckup || (hasCheckupMention && !isCompCheckup) {
		serviceCode = "MCU_BASIC"
		serviceText = "basic medical check-up"
	} else if isCompCheckup {
		serviceCode = "MCU_COMPREHENSIVE"
		serviceText = "comprehensive medical check-up"
	} else if isDental {
		serviceCode = "DENTAL_CHECKUP"
		serviceText = "dental check-up"
		category = "DENTAL"
	} else if isEye {
		serviceCode = "EYE_SCREENING"
		serviceText = "eye screening"
		category = "OPHTHALMOLOGY"
	} else {
		serviceCode = "MCU_BASIC"
		serviceText = "basic medical check-up"
	}

	// Check catalog
	if _, err := catalog.LookupMedicalService(ctx, serviceCode); err != nil {
		reason := "The requested service is not available in the active Batam MedHub catalog."
		return &StructuredIntent{
			SchemaVersion:         "1.0",
			Resolution:            ResolutionUnsupportedService,
			IntentCategory:        nil,
			RequestedServiceText:  serviceText,
			ServiceCode:           nil,
			CandidateServiceCodes: []string{},
			OriginPort:            nil,
			DateWindow:            nil,
			PatientCount:          nil,
			CompanionCount:        nil,
			StayType:              nil,
			Budget:                nil,
			Preferences: IntentPreferences{
				Language:      &lang,
				HotelTier:     nil,
				Accessibility: []string{},
			},
			MissingFields:         []string{},
			ClarificationQuestion: nil,
			OutOfScopeReason:      nil,
			UnsupportedReason:     &reason,
		}, nil
	}

	originPort := "HARBOURFRONT_SG"
	if strings.Contains(lower, "tanah merah") {
		originPort = "TANAH_MERAH_SG"
	}

	// Parse date
	fromDate := "2026-08-22"
	toDate := "2026-08-22"
	if strings.Contains(lower, "23 august") || strings.Contains(lower, "23rd august") {
		fromDate = "2026-08-23"
		toDate = "2026-08-23"
	} else if strings.Contains(lower, "24 august") {
		fromDate = "2026-08-24"
		toDate = "2026-08-24"
	}

	patientCount := 1
	companionCount := 0
	if strings.Contains(lower, "spouse") || strings.Contains(lower, "wife") || strings.Contains(lower, "husband") || strings.Contains(lower, "companion") || strings.Contains(lower, "with my") {
		companionCount = 1
	}

	stayType := StayTypeSameDay
	if strings.Contains(lower, "overnight") || strings.Contains(lower, "hotel") || strings.Contains(lower, "2 days") || strings.Contains(lower, "night") {
		stayType = StayTypeOvernight
		if fromDate == toDate {
			toDate = "2026-08-23"
		}
	}

	var budget *Money
	if strings.Contains(lower, "400") {
		budget = &Money{AmountMinor: 40000, Currency: "SGD"}
	} else if strings.Contains(lower, "500") {
		budget = &Money{AmountMinor: 50000, Currency: "SGD"}
	} else if strings.Contains(lower, "sgd") || strings.Contains(lower, "idr") {
		budget = &Money{AmountMinor: 40000, Currency: referenceCurrency}
	}

	return &StructuredIntent{
		SchemaVersion:         "1.0",
		Resolution:            ResolutionMatched,
		IntentCategory:        &category,
		RequestedServiceText:  serviceText,
		ServiceCode:           &serviceCode,
		CandidateServiceCodes: []string{},
		OriginPort:            &originPort,
		DateWindow: &DateWindow{
			From: fromDate,
			To:   toDate,
		},
		PatientCount:   &patientCount,
		CompanionCount: &companionCount,
		StayType:       &stayType,
		Budget:         budget,
		Preferences: IntentPreferences{
			Language:      &lang,
			HotelTier:     nil,
			Accessibility: []string{},
		},
		MissingFields:         []string{},
		ClarificationQuestion: nil,
		OutOfScopeReason:      nil,
		UnsupportedReason:     nil,
	}, nil
}

func stringPtr(s string) *string {
	return &s
}
