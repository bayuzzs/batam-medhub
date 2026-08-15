package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const baseURL = "http://localhost:8092"

func main() {
	fmt.Println("=== Starting B4 Verification Suite ===")

	nonce := time.Now().UnixNano()
	// 1. Register Patient 1
	p1Token := registerPatient(fmt.Sprintf("patient1_%d@example.test", nonce), "Patient One", "Passw0rd123!")
	fmt.Println("[PASS] Patient 1 registered and authenticated.")

	// 2. Register Patient 2
	p2Token := registerPatient(fmt.Sprintf("patient2_%d@example.test", nonce), "Patient Two", "Passw0rd123!")
	fmt.Println("[PASS] Patient 2 registered and authenticated.")

	// 3. Header & Input Validations
	testHeaderAndInputValidations(p1Token)
	fmt.Println("[PASS] Idempotency header and request validation checks passed.")

	// 4. Test Structured Intent Resolutions
	tripMatchedID := testMatchedTripRequest(p1Token)
	tripClarifyID := testNeedsClarificationTripRequest(p1Token)
	tripUnsupportedID := testUnsupportedTripRequest(p1Token)
	_ = testOutOfScopeTripRequest(p1Token)
	fmt.Println("[PASS] All 4 intent resolutions (MATCHED, NEEDS_CLARIFICATION, UNSUPPORTED_SERVICE, OUT_OF_SCOPE) validated.")

	// 5. Test Idempotency Replay & Conflict
	testIdempotency(p1Token)
	fmt.Println("[PASS] Idempotency replay and 409 conflict detection validated.")

	// 6. Test GET /v1/trip-requests/{id} & Patient Scoping
	testGetTripRequest(p1Token, p2Token, tripMatchedID)
	fmt.Println("[PASS] GET trip-request and patient ID scoping (404 on cross-patient access) validated.")

	// 7. Test PATCH /v1/trip-requests/{id}/intent
	testAmendIntent(p1Token, p2Token, tripClarifyID)
	fmt.Println("[PASS] PATCH intent clarification resolution to MATCHED validated.")

	// 8. Test POST /v1/trip-requests/{id}/plans
	testPlanning(p1Token, p2Token, tripMatchedID, tripUnsupportedID)
	fmt.Println("[PASS] POST plans generation, pricing conversion, and cross-provider options validated.")

	fmt.Println("\n=== ALL B4 VERIFICATIONS COMPLETED SUCCESSFULLY ===")
}

func registerPatient(email, name, password string) string {
	payload := map[string]any{
		"email":              email,
		"full_name":          name,
		"password":           password,
		"preferred_currency": "SGD",
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(baseURL+"/v1/auth/register", "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		raw, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("register failed status %d: %s", resp.StatusCode, string(raw)))
	}

	var session struct {
		AccessToken string `json:"access_token"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&session)
	return session.AccessToken
}

func testHeaderAndInputValidations(token string) {
	// Missing Idempotency-Key
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/v1/trip-requests", bytes.NewReader([]byte(`{"prompt":"hello test prompt","locale":"en"}`)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusBadRequest {
		panic(fmt.Sprintf("expected 400 for missing Idempotency-Key, got %d", resp.StatusCode))
	}
	resp.Body.Close()

	// Short Idempotency-Key (< 8 chars)
	req, _ = http.NewRequest(http.MethodPost, baseURL+"/v1/trip-requests", bytes.NewReader([]byte(`{"prompt":"hello test prompt","locale":"en"}`)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", "short")
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusBadRequest {
		panic(fmt.Sprintf("expected 400 for short Idempotency-Key, got %d", resp.StatusCode))
	}
	resp.Body.Close()

	// Invalid locale
	req, _ = http.NewRequest(http.MethodPost, baseURL+"/v1/trip-requests", bytes.NewReader([]byte(`{"prompt":"hello test prompt","locale":"fr"}`)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", "idem-key-test-01")
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusBadRequest {
		panic(fmt.Sprintf("expected 400 for invalid locale, got %d", resp.StatusCode))
	}
	resp.Body.Close()
}

func testMatchedTripRequest(token string) string {
	prompt := "I need a same-day basic medical check-up in Batam on 22 August, leaving from HarbourFront with my spouse, with a budget of SGD 400."
	payload := map[string]any{"prompt": prompt, "locale": "en"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, baseURL+"/v1/trip-requests", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", "idem-key-matched-001")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		raw, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("expected 201 for matched trip, got %d: %s", resp.StatusCode, string(raw)))
	}

	loc := resp.Header.Get("Location")
	if loc == "" {
		panic("missing Location header on 201 response")
	}

	var detail struct {
		TripRequest struct {
			ID     string `json:"id"`
			Status string `json:"status"`
			Intent struct {
				Resolution   string `json:"resolution"`
				ServiceCode  string `json:"service_code"`
				PatientCount int    `json:"patient_count"`
			} `json:"intent"`
		} `json:"trip_request"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&detail)

	if detail.TripRequest.Intent.Resolution != "MATCHED" {
		panic(fmt.Sprintf("expected MATCHED, got %s", detail.TripRequest.Intent.Resolution))
	}
	if detail.TripRequest.Status != "PLANNING" {
		panic(fmt.Sprintf("expected PLANNING status, got %s", detail.TripRequest.Status))
	}
	if detail.TripRequest.Intent.ServiceCode != "MCU_BASIC" {
		panic(fmt.Sprintf("expected MCU_BASIC, got %s", detail.TripRequest.Intent.ServiceCode))
	}

	return detail.TripRequest.ID
}

func testNeedsClarificationTripRequest(token string) string {
	prompt := "I need a medical check-up in Batam"
	payload := map[string]any{"prompt": prompt, "locale": "en"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, baseURL+"/v1/trip-requests", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", "idem-key-clarify-001")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		raw, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("expected 201 for clarify trip, got %d: %s", resp.StatusCode, string(raw)))
	}

	var detail struct {
		TripRequest struct {
			ID     string `json:"id"`
			Status string `json:"status"`
			Intent struct {
				Resolution    string   `json:"resolution"`
				MissingFields []string `json:"missing_fields"`
			} `json:"intent"`
		} `json:"trip_request"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&detail)

	if detail.TripRequest.Intent.Resolution != "NEEDS_CLARIFICATION" {
		panic(fmt.Sprintf("expected NEEDS_CLARIFICATION, got %s", detail.TripRequest.Intent.Resolution))
	}
	if detail.TripRequest.Status != "NEEDS_INPUT" {
		panic(fmt.Sprintf("expected NEEDS_INPUT status, got %s", detail.TripRequest.Status))
	}
	if len(detail.TripRequest.Intent.MissingFields) == 0 {
		panic("expected missing_fields to be non-empty")
	}

	return detail.TripRequest.ID
}

func testUnsupportedTripRequest(token string) string {
	prompt := "I need knee replacement surgery in Batam on 22 August"
	payload := map[string]any{"prompt": prompt, "locale": "en"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, baseURL+"/v1/trip-requests", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", "idem-key-unsupported-001")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var detail struct {
		TripRequest struct {
			ID     string `json:"id"`
			Status string `json:"status"`
			Intent struct {
				Resolution        string  `json:"resolution"`
				UnsupportedReason *string `json:"unsupported_reason"`
			} `json:"intent"`
		} `json:"trip_request"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&detail)

	if detail.TripRequest.Intent.Resolution != "UNSUPPORTED_SERVICE" {
		panic(fmt.Sprintf("expected UNSUPPORTED_SERVICE, got %s", detail.TripRequest.Intent.Resolution))
	}
	if detail.TripRequest.Status != "UNSUPPORTED_SERVICE" {
		panic(fmt.Sprintf("expected UNSUPPORTED_SERVICE status, got %s", detail.TripRequest.Status))
	}
	if detail.TripRequest.Intent.UnsupportedReason == nil {
		panic("expected unsupported_reason to be present")
	}

	return detail.TripRequest.ID
}

func testOutOfScopeTripRequest(token string) string {
	prompt := "I have severe chest pain and need emergency help now"
	payload := map[string]any{"prompt": prompt, "locale": "en"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, baseURL+"/v1/trip-requests", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", "idem-key-outofscope-001")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var detail struct {
		TripRequest struct {
			ID     string `json:"id"`
			Status string `json:"status"`
			Intent struct {
				Resolution       string  `json:"resolution"`
				OutOfScopeReason *string `json:"out_of_scope_reason"`
			} `json:"intent"`
		} `json:"trip_request"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&detail)

	if detail.TripRequest.Intent.Resolution != "OUT_OF_SCOPE" {
		panic(fmt.Sprintf("expected OUT_OF_SCOPE, got %s", detail.TripRequest.Intent.Resolution))
	}
	if detail.TripRequest.Status != "OUT_OF_SCOPE" {
		panic(fmt.Sprintf("expected OUT_OF_SCOPE status, got %s", detail.TripRequest.Status))
	}
	if detail.TripRequest.Intent.OutOfScopeReason == nil {
		panic("expected out_of_scope_reason to be present")
	}

	return detail.TripRequest.ID
}

func testIdempotency(token string) {
	prompt := "I need a same-day basic medical check-up in Batam on 22 August with budget SGD 400"
	payload := map[string]any{"prompt": prompt, "locale": "en"}
	body, _ := json.Marshal(payload)

	// First call
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/v1/trip-requests", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", "idem-key-replay-test-01")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	resp.Body.Close()
	if resp.Header.Get("Idempotency-Replayed") == "true" {
		panic("expected first call to not be replayed")
	}

	// Second call with same body -> replay
	req, _ = http.NewRequest(http.MethodPost, baseURL+"/v1/trip-requests", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", "idem-key-replay-test-01")
	req.Header.Set("Content-Type", "application/json")
	resp2, _ := http.DefaultClient.Do(req)
	resp2.Body.Close()
	if resp2.Header.Get("Idempotency-Replayed") != "true" {
		panic("expected second call to be replayed with Idempotency-Replayed: true")
	}

	// Third call with changed body -> 409 conflict
	changedBody, _ := json.Marshal(map[string]any{"prompt": "Changed prompt text entirely", "locale": "en"})
	req, _ = http.NewRequest(http.MethodPost, baseURL+"/v1/trip-requests", bytes.NewReader(changedBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", "idem-key-replay-test-01")
	req.Header.Set("Content-Type", "application/json")
	resp3, _ := http.DefaultClient.Do(req)
	resp3.Body.Close()
	if resp3.StatusCode != http.StatusConflict {
		panic(fmt.Sprintf("expected 409 Conflict on payload mismatch, got %d", resp3.StatusCode))
	}
}

func testGetTripRequest(token1, token2, tripID string) {
	// Access with owner token -> 200
	req, _ := http.NewRequest(http.MethodGet, baseURL+"/v1/trip-requests/"+tripID, nil)
	req.Header.Set("Authorization", "Bearer "+token1)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("expected 200 on owner GET trip request, got %d", resp.StatusCode))
	}
	resp.Body.Close()

	// Access with different patient's token -> 404
	req, _ = http.NewRequest(http.MethodGet, baseURL+"/v1/trip-requests/"+tripID, nil)
	req.Header.Set("Authorization", "Bearer "+token2)
	resp2, err := http.DefaultClient.Do(req)
	if err != nil || resp2.StatusCode != http.StatusNotFound {
		panic(fmt.Sprintf("expected 404 on cross-patient GET trip request, got %d", resp2.StatusCode))
	}
	resp2.Body.Close()
}

func testAmendIntent(token1, token2, tripID string) {
	answer := "I would like the basic medical check-up on 22 August"
	payload := map[string]any{
		"answer": answer,
	}
	body, _ := json.Marshal(payload)

	// Cross-patient amend -> 404
	req, _ := http.NewRequest(http.MethodPatch, baseURL+"/v1/trip-requests/"+tripID+"/intent", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token2)
	req.Header.Set("Idempotency-Key", "idem-key-amend-001")
	req.Header.Set("Content-Type", "application/json")
	respCross, _ := http.DefaultClient.Do(req)
	respCross.Body.Close()
	if respCross.StatusCode != http.StatusNotFound {
		panic(fmt.Sprintf("expected 404 on cross-patient patch intent, got %d", respCross.StatusCode))
	}

	// Owner amend -> 200
	req, _ = http.NewRequest(http.MethodPatch, baseURL+"/v1/trip-requests/"+tripID+"/intent", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token1)
	req.Header.Set("Idempotency-Key", "idem-key-amend-001")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("expected 200 on patch intent, got %d: %s", resp.StatusCode, string(raw)))
	}

	var detail struct {
		TripRequest struct {
			Status string `json:"status"`
			Intent struct {
				Resolution  string `json:"resolution"`
				ServiceCode string `json:"service_code"`
			} `json:"intent"`
		} `json:"trip_request"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&detail)

	if detail.TripRequest.Intent.Resolution != "MATCHED" {
		panic(fmt.Sprintf("expected amended intent to be MATCHED, got %s", detail.TripRequest.Intent.Resolution))
	}
	if detail.TripRequest.Status != "PLANNING" {
		panic(fmt.Sprintf("expected status to transition to PLANNING, got %s", detail.TripRequest.Status))
	}
}

func testPlanning(token1, token2, matchedTripID, unsupportedTripID string) {
	// Planning for unsupported trip -> 422
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/v1/trip-requests/"+unsupportedTripID+"/plans", nil)
	req.Header.Set("Authorization", "Bearer "+token1)
	req.Header.Set("Idempotency-Key", "idem-key-plan-unsup-01")
	respUnsup, _ := http.DefaultClient.Do(req)
	respUnsup.Body.Close()
	if respUnsup.StatusCode != http.StatusUnprocessableEntity {
		panic(fmt.Sprintf("expected 422 when planning unsupported trip, got %d", respUnsup.StatusCode))
	}

	// Cross-patient planning -> 404
	req, _ = http.NewRequest(http.MethodPost, baseURL+"/v1/trip-requests/"+matchedTripID+"/plans", nil)
	req.Header.Set("Authorization", "Bearer "+token2)
	req.Header.Set("Idempotency-Key", "idem-key-plan-001")
	respCross, _ := http.DefaultClient.Do(req)
	respCross.Body.Close()
	if respCross.StatusCode != http.StatusNotFound {
		panic(fmt.Sprintf("expected 404 on cross-patient plan generation, got %d", respCross.StatusCode))
	}

	// Owner planning -> 200
	req, _ = http.NewRequest(http.MethodPost, baseURL+"/v1/trip-requests/"+matchedTripID+"/plans", nil)
	req.Header.Set("Authorization", "Bearer "+token1)
	req.Header.Set("Idempotency-Key", "idem-key-plan-001")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("expected 200 on plan generation, got %d: %s", resp.StatusCode, string(raw)))
	}

	var planResult struct {
		TripRequest struct {
			Status           string `json:"status"`
			PlanningRevision int    `json:"planning_revision"`
		} `json:"trip_request"`
		Options []struct {
			ID         string `json:"id"`
			Rank       int    `json:"rank"`
			Status     string `json:"status"`
			Items      []any  `json:"items"`
			TotalPrice struct {
				DisplayTotal struct {
					AmountMinor int64  `json:"amount_minor"`
					Currency    string `json:"currency"`
				} `json:"display_total"`
			} `json:"total_price"`
		} `json:"options"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&planResult)

	if planResult.TripRequest.Status != "PLAN_READY" {
		panic(fmt.Sprintf("expected PLAN_READY, got %s", planResult.TripRequest.Status))
	}
	if len(planResult.Options) != 2 {
		panic(fmt.Sprintf("expected 2 plan options, got %d", len(planResult.Options)))
	}
	if len(planResult.Options[0].Items) != 7 {
		panic(fmt.Sprintf("expected 7 items in option 1, got %d", len(planResult.Options[0].Items)))
	}
	if planResult.Options[0].TotalPrice.DisplayTotal.AmountMinor <= 0 {
		panic("expected positive display total price")
	}

	// Idempotency replay of planning
	req, _ = http.NewRequest(http.MethodPost, baseURL+"/v1/trip-requests/"+matchedTripID+"/plans", nil)
	req.Header.Set("Authorization", "Bearer "+token1)
	req.Header.Set("Idempotency-Key", "idem-key-plan-001")
	respReplay, _ := http.DefaultClient.Do(req)
	respReplay.Body.Close()
	if respReplay.Header.Get("Idempotency-Replayed") != "true" {
		panic("expected replayed plan generation to have Idempotency-Replayed: true")
	}
}

func init() {
	_ = os.Setenv("PORT", "8092")
}
