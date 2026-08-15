package ai

import (
	"context"
	"errors"
	"testing"

	"batam-medhub/internal/service"
)

type mockModelClient struct {
	configured bool
	responses  []string
	errs       []error
	calls      int
}

func (m *mockModelClient) IsConfigured() bool {
	return m.configured
}

func (m *mockModelClient) Infer(ctx context.Context, messages []ChatMessage) (string, error) {
	idx := m.calls
	m.calls++
	if idx < len(m.errs) && m.errs[idx] != nil {
		return "", m.errs[idx]
	}
	if idx < len(m.responses) {
		return m.responses[idx], nil
	}
	return "", errors.New("no more mock responses")
}

func TestExtractor_EmergencyGuardrail(t *testing.T) {
	mock := &mockModelClient{configured: true}
	extractor := NewExtractor(mock, nil, nil)

	intent, err := extractor.ExtractIntent(context.Background(), "I am having severe chest pain and shortness of breath, please help", "en", "SGD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if intent.Resolution != service.ResolutionOutOfScope {
		t.Errorf("expected OUT_OF_SCOPE, got %s", intent.Resolution)
	}
	if intent.OutOfScopeReason == nil || *intent.OutOfScopeReason == "" {
		t.Errorf("expected out_of_scope_reason populated")
	}
	if mock.calls != 0 {
		t.Errorf("emergency guardrail should not invoke LLM, called %d times", mock.calls)
	}
}

func TestExtractor_RetryOnMalformedJSON(t *testing.T) {
	mock := &mockModelClient{
		configured: true,
		responses: []string{
			"Here is the intent: { not valid json",
			`{
				"schema_version": "1.0",
				"resolution": "NEEDS_CLARIFICATION",
				"intent_category": "PREVENTIVE_CHECKUP",
				"requested_service_text": "medical check-up",
				"service_code": null,
				"candidate_service_codes": ["MCU_BASIC", "MCU_COMPREHENSIVE"],
				"origin_port": "HARBOURFRONT_SG",
				"date_window": null,
				"patient_count": 1,
				"companion_count": 0,
				"stay_type": null,
				"budget": null,
				"preferences": {"language": "en", "hotel_tier": null, "accessibility": []},
				"missing_fields": ["service_code", "date_window"],
				"clarification_question": "Would you like basic or comprehensive check-up?",
				"out_of_scope_reason": null,
				"unsupported_reason": null
			}`,
		},
	}

	extractor := NewExtractor(mock, nil, nil)
	intent, err := extractor.ExtractIntent(context.Background(), "I want to do a medical check-up in Batam", "en", "SGD")
	if err != nil {
		t.Fatalf("unexpected error on retry: %v", err)
	}

	if mock.calls != 2 {
		t.Errorf("expected exactly 2 LLM calls (1 initial + 1 retry), got %d", mock.calls)
	}

	if intent.Resolution != service.ResolutionNeedsClarification {
		t.Errorf("expected NEEDS_CLARIFICATION, got %s", intent.Resolution)
	}
}

func TestExtractor_FallbackOnUnconfigured(t *testing.T) {
	mock := &mockModelClient{configured: false}
	extractor := NewExtractor(mock, nil, nil)

	intent, err := extractor.ExtractIntent(context.Background(), "I need a same-day basic medical check-up in Batam on 22 August from HarbourFront with my spouse, budget SGD 400", "en", "SGD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if intent.Resolution != service.ResolutionMatched {
		t.Errorf("expected MATCHED from deterministic fallback, got %s", intent.Resolution)
	}
	if mock.calls != 0 {
		t.Errorf("expected 0 LLM calls when unconfigured, got %d", mock.calls)
	}
}

func TestExtractor_FallbackOnNetworkFailure(t *testing.T) {
	mock := &mockModelClient{
		configured: true,
		errs: []error{
			errors.New("connection reset by peer"),
		},
	}

	extractor := NewExtractor(mock, nil, nil)
	intent, err := extractor.ExtractIntent(context.Background(), "I need a same-day basic medical check-up in Batam on 22 August from HarbourFront with my spouse, budget SGD 400", "en", "SGD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if intent.Resolution != service.ResolutionMatched {
		t.Errorf("expected MATCHED from deterministic fallback on network error, got %s", intent.Resolution)
	}
	if mock.calls != 1 {
		t.Errorf("expected 1 LLM call before fallback, got %d", mock.calls)
	}
}

func TestExtractor_ValidAIResponse(t *testing.T) {
	mock := &mockModelClient{
		configured: true,
		responses: []string{
			`{
				"schema_version": "1.0",
				"resolution": "MATCHED",
				"intent_category": "PREVENTIVE_CHECKUP",
				"requested_service_text": "basic medical check-up",
				"service_code": "MCU_BASIC",
				"candidate_service_codes": [],
				"origin_port": "HARBOURFRONT_SG",
				"date_window": {
					"from": "2026-08-22",
					"to": "2026-08-22"
				},
				"patient_count": 1,
				"companion_count": 1,
				"stay_type": "SAME_DAY",
				"budget": {
					"amount_minor": 40000,
					"currency": "SGD"
				},
				"preferences": {
					"language": "en",
					"hotel_tier": null,
					"accessibility": []
				},
				"missing_fields": [],
				"clarification_question": null,
				"out_of_scope_reason": null,
				"unsupported_reason": null
			}`,
		},
	}

	extractor := NewExtractor(mock, nil, nil)
	intent, err := extractor.ExtractIntent(context.Background(), "I need a same-day basic medical check-up in Batam on 22 August from HarbourFront with my spouse, budget SGD 400", "en", "SGD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if intent.Resolution != service.ResolutionMatched {
		t.Errorf("expected MATCHED, got %s", intent.Resolution)
	}
	if intent.ServiceCode == nil || *intent.ServiceCode != "MCU_BASIC" {
		t.Errorf("expected MCU_BASIC service_code, got %v", intent.ServiceCode)
	}
	if intent.PatientCount == nil || *intent.PatientCount != 1 {
		t.Errorf("expected patient_count 1, got %v", intent.PatientCount)
	}
	if intent.CompanionCount == nil || *intent.CompanionCount != 1 {
		t.Errorf("expected companion_count 1, got %v", intent.CompanionCount)
	}
}
