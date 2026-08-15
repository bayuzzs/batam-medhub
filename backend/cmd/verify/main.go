package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"batam-medhub/internal/adapter"
)

const apiBaseURL = "http://localhost:8093"

func main() {
	fmt.Println("=== Starting B6 Planning and Ranking Verification Suite ===")

	// 1. Setup Mock Provider Servers
	hospitalServer := startMockProviderServer(adapter.ProviderTypeHospital, "hospital-demo-01", "hospital_dev_secret")
	defer hospitalServer.Close()

	ferryServer := startMockProviderServer(adapter.ProviderTypeFerry, "ferry-demo-01", "ferry_dev_secret")
	defer ferryServer.Close()

	hotelServer := startMockProviderServer(adapter.ProviderTypeHotel, "hotel-demo-01", "hotel_dev_secret")
	defer hotelServer.Close()

	transportServer := startMockProviderServer(adapter.ProviderTypeTransport, "transport-demo-01", "transport_dev_secret")
	defer transportServer.Close()

	fmt.Println("[PASS] Mock provider servers initialized.")

	nonce := time.Now().UnixNano()
	p1Email := fmt.Sprintf("b6_patient1_%d@example.test", nonce)
	token := registerPatient(p1Email, "B6 Patient", "Passw0rd123!", "SGD")
	fmt.Println("[PASS] Patient registered and authenticated with SGD preferred currency.")

	// 2. Create Matched Trip Request
	tripID := createMatchedTripRequest(token)
	fmt.Println("[PASS] Matched trip request created (status PLANNING).")

	// 3. Generate Plans (Happy Path: 2 Ranked Options, Buffer Invariants & Multi-Currency Pricing)
	testPlanGeneration(token, tripID)
	fmt.Println("[PASS] Plan generation, schedule buffer invariants, and 2-option ranking validated.")

	// 4. Test NO_MATCH Handling
	testNoMatchTrip(token)
	fmt.Println("[PASS] NO_MATCH status and descriptive reason handling validated.")

	// 5. Test Idempotency Replay
	testPlanIdempotency(token, tripID)
	fmt.Println("[PASS] Idempotency replay on plan generation validated.")

	fmt.Println("\n=== ALL B6 VERIFICATIONS COMPLETED SUCCESSFULLY ===")
}

func registerPatient(email, name, password, currency string) string {
	payload := map[string]any{
		"email":              email,
		"full_name":          name,
		"password":           password,
		"preferred_currency": currency,
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(apiBaseURL+"/v1/auth/register", "application/json", bytes.NewReader(body))
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

func createMatchedTripRequest(token string) string {
	prompt := "I need a same-day basic medical check-up in Batam on 22 August, leaving from HarbourFront with my spouse, with a budget of SGD 400."
	payload := map[string]any{"prompt": prompt, "locale": "en"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, apiBaseURL+"/v1/trip-requests", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", fmt.Sprintf("idem-b6-create-%d", time.Now().UnixNano()))
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

	var detail struct {
		TripRequest struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"trip_request"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&detail)
	return detail.TripRequest.ID
}

func testPlanGeneration(token, tripID string) {
	req, _ := http.NewRequest(http.MethodPost, apiBaseURL+"/v1/trip-requests/"+tripID+"/plans", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", "idem-b6-plan-generate-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("expected 200 on plan generation, got %d: %s", resp.StatusCode, string(raw)))
	}

	var result struct {
		TripRequest struct {
			ID               string `json:"id"`
			Status           string `json:"status"`
			PlanningRevision int    `json:"planning_revision"`
		} `json:"trip_request"`
		Options []struct {
			ID          string   `json:"id"`
			Rank        int      `json:"rank"`
			Status      string   `json:"status"`
			Explanation []string `json:"explanation"`
			TotalPrice  struct {
				SourceTotals []struct {
					AmountMinor int64  `json:"amount_minor"`
					Currency    string `json:"currency"`
				} `json:"source_totals"`
				DisplayTotal struct {
					AmountMinor int64  `json:"amount_minor"`
					Currency    string `json:"currency"`
				} `json:"display_total"`
				Estimated bool `json:"estimated"`
			} `json:"total_price"`
			Items []struct {
				ID         string `json:"id"`
				ItemType   string `json:"item_type"`
				TimeWindow struct {
					StartsAt string `json:"starts_at"`
					EndsAt   string `json:"ends_at"`
				} `json:"time_window"`
				Price *struct {
					Source struct {
						AmountMinor int64  `json:"amount_minor"`
						Currency    string `json:"currency"`
					} `json:"source"`
					Display struct {
						AmountMinor int64  `json:"amount_minor"`
						Currency    string `json:"currency"`
					} `json:"display"`
					FXRate string `json:"fx_rate"`
				} `json:"price"`
			} `json:"items"`
		} `json:"options"`
		NoMatchReasons   []string `json:"no_match_reasons"`
		ProviderWarnings []string `json:"provider_warnings"`
	}

	_ = json.NewDecoder(resp.Body).Decode(&result)

	// Check TripRequest Status
	if result.TripRequest.Status != "PLAN_READY" {
		panic(fmt.Sprintf("expected status PLAN_READY, got %s", result.TripRequest.Status))
	}
	if result.TripRequest.PlanningRevision != 1 {
		panic(fmt.Sprintf("expected planning_revision 1, got %d", result.TripRequest.PlanningRevision))
	}

	// Check Exactly 2 Ranked Options
	if len(result.Options) != 2 {
		panic(fmt.Sprintf("expected exactly 2 plan options, got %d", len(result.Options)))
	}
	if result.Options[0].Rank != 1 || result.Options[1].Rank != 2 {
		panic(fmt.Sprintf("expected options ranked [1, 2], got [%d, %d]", result.Options[0].Rank, result.Options[1].Rank))
	}

	// Verify Rank 1 is lowest total display cost
	if result.Options[0].TotalPrice.DisplayTotal.AmountMinor > result.Options[1].TotalPrice.DisplayTotal.AmountMinor {
		panic("expected Rank 1 display total price to be <= Rank 2 display total price")
	}

	// Check Plan Items and Schedule Buffer Invariants for Option 1
	items := result.Options[0].Items
	if len(items) != 7 {
		panic(fmt.Sprintf("expected 7 plan items in option 1, got %d", len(items)))
	}

	var ferryArr, appStart, appEnd, ferryDep time.Time
	for _, it := range items {
		switch it.ItemType {
		case "FERRY_OUTBOUND":
			ferryArr, _ = time.Parse(time.RFC3339, it.TimeWindow.EndsAt)
		case "HOSPITAL_APPOINTMENT":
			appStart, _ = time.Parse(time.RFC3339, it.TimeWindow.StartsAt)
			appEnd, _ = time.Parse(time.RFC3339, it.TimeWindow.EndsAt)
		case "FERRY_RETURN":
			ferryDep, _ = time.Parse(time.RFC3339, it.TimeWindow.StartsAt)
		}
	}

	// Invariant 1: At least 45 minutes between Batam arrival and appointment start
	if appStart.Sub(ferryArr) < 45*time.Minute {
		panic(fmt.Sprintf("buffer invariant violated: arrival to appointment buffer is %v (expected >= 45m)", appStart.Sub(ferryArr)))
	}

	// Invariant 2: At least 30 minutes release buffer plus transfer between appointment end and return departure
	if ferryDep.Sub(appEnd) < 30*time.Minute {
		panic(fmt.Sprintf("buffer invariant violated: appointment end to return departure buffer is %v (expected >= 30m)", ferryDep.Sub(appEnd)))
	}

	// Multi-currency price check
	if len(result.Options[0].TotalPrice.SourceTotals) < 2 {
		panic("expected mixed source currencies (SGD and IDR) in total_price.source_totals")
	}
	if result.Options[0].TotalPrice.DisplayTotal.Currency != "SGD" {
		panic(fmt.Sprintf("expected display currency SGD, got %s", result.Options[0].TotalPrice.DisplayTotal.Currency))
	}
}

func testNoMatchTrip(token string) {
	// Request with past impossible date
	prompt := "I need a basic check-up in Batam on 2020-01-01"
	payload := map[string]any{"prompt": prompt, "locale": "en"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, apiBaseURL+"/v1/trip-requests", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", fmt.Sprintf("idem-nomatch-create-%d", time.Now().UnixNano()))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		panic("failed to create trip for no match test")
	}
	var detail struct {
		TripRequest struct {
			ID string `json:"id"`
		} `json:"trip_request"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&detail)
	resp.Body.Close()

	// Amend with explicit past date to bypass intent parser default
	amendBody, _ := json.Marshal(map[string]any{
		"corrections": map[string]any{
			"service_code": "MCU_BASIC",
			"date_window": map[string]any{
				"from": "2020-01-01",
				"to":   "2020-01-01",
			},
		},
	})
	reqAmend, _ := http.NewRequest(http.MethodPatch, apiBaseURL+"/v1/trip-requests/"+detail.TripRequest.ID+"/intent", bytes.NewReader(amendBody))
	reqAmend.Header.Set("Authorization", "Bearer "+token)
	reqAmend.Header.Set("Idempotency-Key", fmt.Sprintf("idem-nomatch-amend-%d", time.Now().UnixNano()))
	reqAmend.Header.Set("Content-Type", "application/json")
	respAmend, _ := http.DefaultClient.Do(reqAmend)
	respAmend.Body.Close()

	// Plan should return 200 with NO_MATCH
	reqPlan, _ := http.NewRequest(http.MethodPost, apiBaseURL+"/v1/trip-requests/"+detail.TripRequest.ID+"/plans", nil)
	reqPlan.Header.Set("Authorization", "Bearer "+token)
	reqPlan.Header.Set("Idempotency-Key", fmt.Sprintf("idem-nomatch-plan-%d", time.Now().UnixNano()))
	respPlan, err := http.DefaultClient.Do(reqPlan)
	if err != nil {
		panic(err)
	}
	defer respPlan.Body.Close()

	if respPlan.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(respPlan.Body)
		panic(fmt.Sprintf("expected 200 for NO_MATCH planning, got %d: %s", respPlan.StatusCode, string(raw)))
	}

	var planResult struct {
		TripRequest struct {
			Status string `json:"status"`
		} `json:"trip_request"`
		Options        []any    `json:"options"`
		NoMatchReasons []string `json:"no_match_reasons"`
	}
	_ = json.NewDecoder(respPlan.Body).Decode(&planResult)

	if planResult.TripRequest.Status != "NO_MATCH" {
		panic(fmt.Sprintf("expected status NO_MATCH, got %s", planResult.TripRequest.Status))
	}
	if len(planResult.NoMatchReasons) == 0 {
		panic("expected no_match_reasons to be populated")
	}
}

func testPlanIdempotency(token, tripID string) {
	req, _ := http.NewRequest(http.MethodPost, apiBaseURL+"/v1/trip-requests/"+tripID+"/plans", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", "idem-b6-plan-generate-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("expected 200 on plan idempotency replay, got %d", resp.StatusCode))
	}
	if resp.Header.Get("Idempotency-Replayed") != "true" {
		panic("expected Idempotency-Replayed: true header on duplicate request")
	}
}

func startMockProviderServer(provType, provID, secret string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Request-ID", r.Header.Get("X-Request-ID"))

		if r.URL.Path == "/healthz" {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(adapter.HealthResponse{
				Status:         "UP",
				ProviderID:     provID,
				ProviderType:   provType,
				DatabaseStatus: "UP",
				CheckedAt:      "2026-08-16T00:00:00Z",
			})
			return
		}

		if r.Header.Get("X-Integration-Key") != secret {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(adapter.ErrorEnvelope{
				Error: adapter.ErrorBody{
					Code:      "AUTHENTICATION_FAILED",
					Message:   "Invalid secret",
					Retryable: false,
					RequestID: r.Header.Get("X-Request-ID"),
				},
			})
			return
		}

		switch r.URL.Path {
		case "/v1/offers/search":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(adapter.SearchResponse{
				ProviderID:   provID,
				ProviderType: provType,
				Offers: []adapter.Offer{
					{
						OfferID:        fmt.Sprintf("%s-offer-01", provID[:4]),
						ProviderID:     provID,
						ProviderType:   provType,
						Status:         "AVAILABLE",
						ServiceWindow:  adapter.TimeWindow{StartsAt: "2026-08-22T02:00:00Z", EndsAt: "2026-08-22T04:00:00Z", StartTimeZone: "Asia/Jakarta", EndTimeZone: "Asia/Jakarta"},
						AvailableUnits: 5,
						UnitPrice:      adapter.Money{AmountMinor: 150000000, Currency: "IDR"},
						ValidUntil:     "2026-08-20T12:00:00Z",
						Synthetic:      true,
						Source:         "MOCK",
						Details:        json.RawMessage(`{}`),
					},
				},
			})

		default:
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]any{"status": "OK"})
		}
	}))
}
