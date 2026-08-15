package ai

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"batam-medhub/internal/service"
)

// Extractor performs structured intent extraction using Cloudflare Workers AI with guardrails and deterministic fallback.
type Extractor struct {
	client  ModelClient
	catalog *service.CatalogService
	logger  *slog.Logger
}

// NewExtractor constructs a new Extractor instance.
func NewExtractor(client ModelClient, catalog *service.CatalogService, logger *slog.Logger) *Extractor {
	if logger == nil {
		logger = slog.Default()
	}
	return &Extractor{
		client:  client,
		catalog: catalog,
		logger:  logger,
	}
}

// ExtractIntent parses freeform patient inquiries into a structured intent snapshot with AI, guardrails, and deterministic fallback.
func (e *Extractor) ExtractIntent(ctx context.Context, prompt, locale, referenceCurrency string) (*service.StructuredIntent, error) {
	prompt = strings.TrimSpace(prompt)
	lower := strings.ToLower(prompt)

	// 1. Pre-extraction safety guardrail for medical emergencies
	if isEmergencyPrompt(lower) {
		e.logger.InfoContext(ctx, "emergency safety triage triggered", "prompt", prompt)
		reason := "This may be an emergency. Batam MedHub cannot triage or plan emergency care. Contact local emergency services immediately or go to the nearest emergency department."
		lang := locale
		if lang == "" {
			lang = "en"
		}
		return &service.StructuredIntent{
			SchemaVersion:         "1.0",
			Resolution:            service.ResolutionOutOfScope,
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
			Preferences: service.IntentPreferences{
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

	// 2. Fall back to deterministic rules if AI client is unconfigured
	if e.client == nil || !e.client.IsConfigured() {
		e.logger.DebugContext(ctx, "ai client unconfigured, using deterministic extractor")
		return service.ExtractIntentDeterministic(ctx, e.catalog, prompt, locale, referenceCurrency)
	}

	// 3. Fetch active catalog from PostgreSQL
	var catServices []service.MedicalServiceItem
	if e.catalog != nil {
		catResp, err := e.catalog.ListMedicalServices(ctx)
		if err != nil {
			e.logger.WarnContext(ctx, "failed to query medical catalog for prompt injection, falling back to deterministic extractor", "error", err)
			return service.ExtractIntentDeterministic(ctx, e.catalog, prompt, locale, referenceCurrency)
		}
		if catResp != nil {
			catServices = catResp.Services
		}
	}

	// 4. Construct AI prompts
	systemPrompt := BuildSystemPrompt(catServices)
	userPrompt := BuildUserPrompt(prompt, locale, referenceCurrency)

	messages := []ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	// 5. Model Inference (Attempt 1)
	respText, err := e.client.Infer(ctx, messages)
	if err != nil {
		e.logger.WarnContext(ctx, "ai inference failed, falling back to deterministic extractor", "error", err)
		return service.ExtractIntentDeterministic(ctx, e.catalog, prompt, locale, referenceCurrency)
	}

	// 6. Parse and validate candidate intent
	intent, parseErr := parseModelResponse(respText)
	if parseErr != nil {
		e.logger.WarnContext(ctx, "model returned malformed structured JSON, retrying once with corrective prompt", "error", parseErr, "raw_response", respText)

		// 7. Retry ONCE on malformed structured output
		retryMessages := append(messages,
			ChatMessage{Role: "assistant", Content: respText},
			ChatMessage{Role: "user", Content: "The previous output was not valid JSON adhering to the schema. Return ONLY valid, raw JSON with no backticks, markdown, or commentary."},
		)

		retryRespText, retryErr := e.client.Infer(ctx, retryMessages)
		if retryErr != nil {
			e.logger.WarnContext(ctx, "ai retry inference failed, falling back to deterministic extractor", "error", retryErr)
			return service.ExtractIntentDeterministic(ctx, e.catalog, prompt, locale, referenceCurrency)
		}

		intent, parseErr = parseModelResponse(retryRespText)
		if parseErr != nil {
			e.logger.WarnContext(ctx, "ai retry returned invalid JSON again, falling back to deterministic extractor", "error", parseErr, "raw_response", retryRespText)
			return service.ExtractIntentDeterministic(ctx, e.catalog, prompt, locale, referenceCurrency)
		}
	}

	// 8. Post-Validation and Database Catalog Verification
	if err := e.sanitizeAndVerifyIntent(ctx, intent, prompt, locale, referenceCurrency); err != nil {
		e.logger.WarnContext(ctx, "ai intent validation failed, falling back to deterministic extractor", "error", err)
		return service.ExtractIntentDeterministic(ctx, e.catalog, prompt, locale, referenceCurrency)
	}

	return intent, nil
}

func isEmergencyPrompt(lower string) bool {
	return strings.Contains(lower, "chest pain") ||
		strings.Contains(lower, "emergency") ||
		strings.Contains(lower, "heart attack") ||
		strings.Contains(lower, "ambulance") ||
		strings.Contains(lower, "shortness of breath") ||
		strings.Contains(lower, "severe pain") ||
		strings.Contains(lower, "stroke") ||
		strings.Contains(lower, "unconscious") ||
		strings.Contains(lower, "heavy bleeding")
}

func parseModelResponse(raw string) (*service.StructuredIntent, error) {
	cleaned := cleanJSONString(raw)
	var intent service.StructuredIntent
	if err := json.Unmarshal([]byte(cleaned), &intent); err != nil {
		return nil, err
	}
	return &intent, nil
}

func cleanJSONString(raw string) string {
	raw = strings.TrimSpace(raw)
	// Strip markdown fences ```json ... ``` if model output them
	if strings.HasPrefix(raw, "```") {
		lines := strings.Split(raw, "\n")
		if len(lines) >= 2 {
			// drop first line
			lines = lines[1:]
			// drop last line if it's ```
			if len(lines) > 0 && strings.HasPrefix(strings.TrimSpace(lines[len(lines)-1]), "```") {
				lines = lines[:len(lines)-1]
			}
			raw = strings.Join(lines, "\n")
		}
	}
	return strings.TrimSpace(raw)
}

func (e *Extractor) sanitizeAndVerifyIntent(ctx context.Context, intent *service.StructuredIntent, prompt, locale, referenceCurrency string) error {
	intent.SchemaVersion = "1.0"
	if intent.RequestedServiceText == "" {
		intent.RequestedServiceText = prompt
	}
	if intent.CandidateServiceCodes == nil {
		intent.CandidateServiceCodes = []string{}
	}
	if intent.MissingFields == nil {
		intent.MissingFields = []string{}
	}
	if intent.Preferences.Accessibility == nil {
		intent.Preferences.Accessibility = []string{}
	}
	if intent.Preferences.Language == nil {
		lang := locale
		if lang == "" {
			lang = "en"
		}
		intent.Preferences.Language = &lang
	}

	// Verify catalog for MATCHED resolution
	if intent.Resolution == service.ResolutionMatched {
		if intent.ServiceCode == nil || *intent.ServiceCode == "" {
			// Convert to NEEDS_CLARIFICATION if service code was not resolved
			intent.Resolution = service.ResolutionNeedsClarification
			q := "Which medical check-up package would you like to book?"
			intent.ClarificationQuestion = &q
			intent.MissingFields = append(intent.MissingFields, "service_code")
			intent.CandidateServiceCodes = []string{"MCU_BASIC", "MCU_COMPREHENSIVE"}
		} else if e.catalog != nil {
			// Verify against database
			item, err := e.catalog.LookupMedicalService(ctx, *intent.ServiceCode)
			if err != nil {
				// The service code returned by the model is not in our database
				intent.Resolution = service.ResolutionUnsupportedService
				reason := "The requested service is not available in the active Batam MedHub catalog."
				intent.UnsupportedReason = &reason
				intent.ServiceCode = nil
			} else {
				if intent.IntentCategory == nil {
					intent.IntentCategory = &item.Category
				}
			}
		}
	}

	// Run standard contract invariant validation
	return service.ValidateStructuredIntent(intent)
}
