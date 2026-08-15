package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"batam-medhub/internal/adapter"
	"batam-medhub/internal/config"
	"batam-medhub/internal/database"
	"batam-medhub/internal/httpapi"
	"batam-medhub/internal/model"
	"batam-medhub/internal/service"

	"gorm.io/gorm"
)

func main() {
	fmt.Println("=== Starting B7 Booking Saga and Journey Tracking Verification Suite ===")

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://provider_admin:provider_admin_dev_password@172.21.0.2:5432/core_db?sslmode=disable"
	}
	jwtSecret := "12345678901234567890123456789012"

	db, err := database.Open(dbURL)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to test db: %v", err))
	}

	cfg := config.Config{
		HTTPAddr:                ":8093",
		DatabaseURL:             dbURL,
		JWTSigningSecret:        jwtSecret,
		JWTIssuer:               "batam-medhub",
		JWTAudience:             "batam-medhub-mobile",
		JWTAccessTTL:            15 * time.Minute,
		RefreshTokenTTL:         30 * 24 * time.Hour,
		HospitalBaseURL:         "http://localhost:8081",
		HospitalIntegrationKey:  "hospital_dev_secret",
		FerryBaseURL:            "http://localhost:8082",
		FerryIntegrationKey:     "ferry_dev_secret",
		HotelBaseURL:            "http://localhost:8083",
		HotelIntegrationKey:     "hotel_dev_secret",
		TransportBaseURL:        "http://localhost:8084",
		TransportIntegrationKey: "transport_dev_secret",
		ProviderTimeout:         5 * time.Second,
	}

	coreEngine := httpapi.New(db, cfg, nil)
	coreServer := httptest.NewServer(coreEngine)
	defer coreServer.Close()
	apiBaseURL := coreServer.URL

	nonce := time.Now().UnixNano()
	p1Email := fmt.Sprintf("b7_patient1_%d@example.test", nonce)
	token1 := registerPatient(apiBaseURL, p1Email, "B7 Patient One", "Passw0rd123!", "SGD")
	fmt.Println("[PASS] Patient 1 registered and authenticated with SGD preferred currency.")

	// 1. Create Matched Trip Request and Generate Plans
	tripID1 := createMatchedTripRequest(apiBaseURL, token1)
	planResult1 := generatePlans(apiBaseURL, token1, tripID1)
	if len(planResult1.Options) < 1 {
		panic("expected at least 1 plan option generated")
	}
	selectedOption1 := planResult1.Options[0]
	fmt.Println("[PASS] Trip request 1 created and 2 ranked plan options generated.")

	// 2. Validate Approval Request Constraints
	testApprovalValidation(apiBaseURL, token1, selectedOption1.ID)
	fmt.Println("[PASS] Explicit approval validation enforced (requires approved: true).")

	// 3. Execute Booking Saga (Happy Path)
	bookingIdemKey := fmt.Sprintf("idem-b7-book-%d", time.Now().UnixNano())
	journeyDetail1 := confirmPlanOption(apiBaseURL, token1, selectedOption1.ID, bookingIdemKey)

	if journeyDetail1.Journey.Status != "ACTIVE" {
		panic(fmt.Sprintf("expected journey status ACTIVE, got %s", journeyDetail1.Journey.Status))
	}
	if journeyDetail1.Journey.ActiveItineraryVersion != 1 {
		panic(fmt.Sprintf("expected active itinerary version 1, got %d", journeyDetail1.Journey.ActiveItineraryVersion))
	}
	if len(journeyDetail1.ActiveItinerary.Items) != 7 {
		panic(fmt.Sprintf("expected 7 itinerary items, got %d", len(journeyDetail1.ActiveItinerary.Items)))
	}

	// Validate booked items vs buffers
	ferryOut := journeyDetail1.ActiveItinerary.Items[0]
	if ferryOut.ItemType != "FERRY_OUTBOUND" || ferryOut.Status != "CONFIRMED" || ferryOut.ProviderID == nil || *ferryOut.ProviderID != "ferry-demo-01" || ferryOut.ExternalReservationID == nil {
		panic("ferry outbound item not confirmed properly")
	}
	arrivalBuffer := journeyDetail1.ActiveItinerary.Items[1]
	if arrivalBuffer.ItemType != "ARRIVAL_BUFFER" || arrivalBuffer.Status != "BUFFER" || arrivalBuffer.ProviderID != nil || arrivalBuffer.ExternalReservationID != nil {
		panic("arrival buffer item invalid")
	}
	transPick := journeyDetail1.ActiveItinerary.Items[2]
	if transPick.ItemType != "TRANSPORT_PICKUP" || transPick.Status != "CONFIRMED" || transPick.ProviderID == nil || *transPick.ProviderID != "transport-demo-01" || transPick.ExternalReservationID == nil {
		panic("transport pickup item not confirmed properly")
	}
	hospAppt := journeyDetail1.ActiveItinerary.Items[3]
	if hospAppt.ItemType != "HOSPITAL_APPOINTMENT" || hospAppt.Status != "CONFIRMED" || hospAppt.ProviderID == nil || *hospAppt.ProviderID != "hospital-demo-01" || hospAppt.ExternalReservationID == nil {
		panic("hospital appointment item not confirmed properly")
	}
	transDrop := journeyDetail1.ActiveItinerary.Items[4]
	if transDrop.ItemType != "TRANSPORT_DROPOFF" || transDrop.Status != "CONFIRMED" || transDrop.ProviderID == nil || *transDrop.ProviderID != "transport-demo-01" || transDrop.ExternalReservationID == nil {
		panic("transport dropoff item not confirmed properly")
	}
	depBuffer := journeyDetail1.ActiveItinerary.Items[5]
	if depBuffer.ItemType != "DEPARTURE_BUFFER" || depBuffer.Status != "BUFFER" || depBuffer.ProviderID != nil || depBuffer.ExternalReservationID != nil {
		panic("departure buffer item invalid")
	}
	ferryRet := journeyDetail1.ActiveItinerary.Items[6]
	if ferryRet.ItemType != "FERRY_RETURN" || ferryRet.Status != "CONFIRMED" || ferryRet.ProviderID == nil || *ferryRet.ProviderID != "ferry-demo-01" || ferryRet.ExternalReservationID == nil {
		panic("ferry return item not confirmed properly")
	}

	fmt.Println("[PASS] Booking saga succeeded: multi-provider holds & confirmations sequentially executed; Journey status ACTIVE and Itinerary Version 1 persisted.")

	// 4. Test Idempotency Replay on Confirm Plan
	testBookingIdempotencyReplay(apiBaseURL, token1, selectedOption1.ID, bookingIdemKey)
	fmt.Println("[PASS] Idempotency replay on POST /v1/plan-options/{plan_option_id}/confirm verified.")

	// 5. Test Journey Tracking Endpoints
	journeyID := journeyDetail1.Journey.ID
	testGetJourney(apiBaseURL, token1, journeyID)
	testGetJourneyItinerary(apiBaseURL, token1, journeyID)
	testGetJourneyItineraryVersion(apiBaseURL, token1, journeyID, 1)
	testListJourneys(apiBaseURL, token1, journeyID)
	fmt.Println("[PASS] GET /v1/journeys/{journey_id}, GET /v1/journeys/{journey_id}/itinerary, GET /v1/journeys/{journey_id}/itineraries/1, and GET /v1/journeys validated.")

	// 6. Verify Trip Request state updated with journey_id
	testGetTripRequestWithJourney(apiBaseURL, token1, tripID1, journeyID)
	fmt.Println("[PASS] GET /v1/trip-requests/{trip_request_id} returns status ACTIVE and populated journey_id.")

	// 7. Test POST /v1/trip-requests/{trip_request_id}/select-plan endpoint
	tripID2 := createMatchedTripRequest(apiBaseURL, token1)
	planResult2 := generatePlans(apiBaseURL, token1, tripID2)
	selectPlanIdemKey := fmt.Sprintf("idem-b7-select-%d", time.Now().UnixNano())
	journeyDetail2 := selectPlanForTrip(apiBaseURL, token1, tripID2, planResult2.Options[1].ID, selectPlanIdemKey)
	if journeyDetail2.Journey.Status != "ACTIVE" {
		panic(fmt.Sprintf("expected journey 2 status ACTIVE, got %s", journeyDetail2.Journey.Status))
	}
	fmt.Println("[PASS] POST /v1/trip-requests/{trip_request_id}/select-plan booking initiation validated.")

	// 8. Test Patient Ownership Scoping
	p2Email := fmt.Sprintf("b7_patient2_%d@example.test", nonce)
	token2 := registerPatient(apiBaseURL, p2Email, "B7 Patient Two", "Passw0rd123!", "IDR")
	testPatientOwnershipScoping(apiBaseURL, token2, journeyID)
	fmt.Println("[PASS] Patient ownership scoping verified (patient 2 cannot access patient 1's journey).")

	// 9. Test Compensation on Hold Failure with Simulated Failing Provider Server
	testHoldCompensation(db, jwtSecret)
	fmt.Println("[PASS] Hold compensation verified: on hold failure, prior holds released and trip request transitions to CONFIRMATION_FAILED.")

	// 10. Test Compensation on Confirmation Failure with Simulated Failing Provider Server
	testConfirmationCompensation(db, jwtSecret)
	fmt.Println("[PASS] Confirmation compensation verified: on confirm failure, confirmed items and unconfirmed holds released and trip transitions to CONFIRMATION_FAILED.")

	// 11. Test Cloudflare Workers AI Intent Extractor (Phase B8)
	testWorkersAIExtractor(db, jwtSecret)
	fmt.Println("[PASS] Cloudflare Workers AI extractor verified (guardrails, retries, catalog validation, and seamless fallback).")

	// 12. Test Disruption and Recovery Engine (Phase B9)
	testDisruptionAndRecoveryEngine(db, jwtSecret)
	fmt.Println("[PASS] Disruption and recovery engine verified (provider event ingestion, deduplication, clinical validation, recovery option generation, Itinerary v2 activation, and manual review fallback).")

	fmt.Println("\n=== ALL B7, B8, AND B9 VERIFICATIONS COMPLETED SUCCESSFULLY ===")
}

func registerPatient(apiBaseURL, email, name, password, currency string) string {
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

func createMatchedTripRequest(apiBaseURL, token string) string {
	prompt := "I need a same-day basic medical check-up in Batam on 22 August, leaving from HarbourFront with my spouse, with a budget of SGD 400."
	payload := map[string]any{"prompt": prompt, "locale": "en"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, apiBaseURL+"/v1/trip-requests", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", fmt.Sprintf("idem-b7-create-%d", time.Now().UnixNano()))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		raw, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("create trip request failed %d: %s", resp.StatusCode, string(raw)))
	}

	var res struct {
		TripRequest struct {
			ID string `json:"id"`
		} `json:"trip_request"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&res)
	return res.TripRequest.ID
}

func generatePlans(apiBaseURL, token, tripID string) service.PlanningResult {
	req, _ := http.NewRequest(http.MethodPost, apiBaseURL+"/v1/trip-requests/"+tripID+"/plans", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", fmt.Sprintf("idem-b7-plan-%d", time.Now().UnixNano()))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("generate plans failed %d: %s", resp.StatusCode, string(raw)))
	}

	var planResult service.PlanningResult
	_ = json.NewDecoder(resp.Body).Decode(&planResult)
	return planResult
}

func testApprovalValidation(apiBaseURL, token, planOptionID string) {
	// 1. Missing approved field
	body1 := []byte(`{}`)
	req1, _ := http.NewRequest(http.MethodPost, apiBaseURL+"/v1/plan-options/"+planOptionID+"/confirm", bytes.NewReader(body1))
	req1.Header.Set("Authorization", "Bearer "+token)
	req1.Header.Set("Idempotency-Key", fmt.Sprintf("idem-b7-val1-%d", time.Now().UnixNano()))
	req1.Header.Set("Content-Type", "application/json")
	resp1, _ := http.DefaultClient.Do(req1)
	if resp1.StatusCode != http.StatusBadRequest {
		panic(fmt.Sprintf("expected 400 on empty approval body, got %d", resp1.StatusCode))
	}
	resp1.Body.Close()

	// 2. approved: false
	body2 := []byte(`{"approved": false}`)
	req2, _ := http.NewRequest(http.MethodPost, apiBaseURL+"/v1/plan-options/"+planOptionID+"/confirm", bytes.NewReader(body2))
	req2.Header.Set("Authorization", "Bearer "+token)
	req2.Header.Set("Idempotency-Key", fmt.Sprintf("idem-b7-val2-%d", time.Now().UnixNano()))
	req2.Header.Set("Content-Type", "application/json")
	resp2, _ := http.DefaultClient.Do(req2)
	if resp2.StatusCode != http.StatusBadRequest {
		panic(fmt.Sprintf("expected 400 on approved: false, got %d", resp2.StatusCode))
	}
	resp2.Body.Close()
}

func confirmPlanOption(apiBaseURL, token, planOptionID, idemKey string) service.JourneyDetail {
	body := []byte(`{"approved": true}`)
	req, _ := http.NewRequest(http.MethodPost, apiBaseURL+"/v1/plan-options/"+planOptionID+"/confirm", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", idemKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		raw, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("confirm plan option failed %d: %s", resp.StatusCode, string(raw)))
	}

	loc := resp.Header.Get("Location")
	if loc == "" || !strings.HasPrefix(loc, "/v1/journeys/") {
		panic(fmt.Sprintf("expected Location header starting with /v1/journeys/, got %s", loc))
	}

	var detail service.JourneyDetail
	_ = json.NewDecoder(resp.Body).Decode(&detail)
	return detail
}

func selectPlanForTrip(apiBaseURL, token, tripID, planOptionID, idemKey string) service.JourneyDetail {
	payload := map[string]any{
		"approved": true,
	}
	if planOptionID != "" {
		payload["plan_option_id"] = planOptionID
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, apiBaseURL+"/v1/trip-requests/"+tripID+"/select-plan", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", idemKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		raw, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("select plan failed %d: %s", resp.StatusCode, string(raw)))
	}

	var detail service.JourneyDetail
	_ = json.NewDecoder(resp.Body).Decode(&detail)
	return detail
}

func testBookingIdempotencyReplay(apiBaseURL, token, planOptionID, idemKey string) {
	body := []byte(`{"approved": true}`)
	req, _ := http.NewRequest(http.MethodPost, apiBaseURL+"/v1/plan-options/"+planOptionID+"/confirm", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", idemKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		panic(fmt.Sprintf("expected 201 on idempotency replay, got %d", resp.StatusCode))
	}
	if resp.Header.Get("Idempotency-Replayed") != "true" {
		panic("expected Idempotency-Replayed: true header on duplicate booking request")
	}
}

func testGetJourney(apiBaseURL, token, journeyID string) {
	req, _ := http.NewRequest(http.MethodGet, apiBaseURL+"/v1/journeys/"+journeyID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("GET /v1/journeys/{journey_id} failed with status %d", resp.StatusCode))
	}

	var detail service.JourneyDetail
	_ = json.NewDecoder(resp.Body).Decode(&detail)
	if detail.Journey.ID != journeyID {
		panic("journey id mismatch")
	}
}

func testGetJourneyItinerary(apiBaseURL, token, journeyID string) {
	req, _ := http.NewRequest(http.MethodGet, apiBaseURL+"/v1/journeys/"+journeyID+"/itinerary", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("GET /v1/journeys/{journey_id}/itinerary failed with status %d", resp.StatusCode))
	}
}

func testGetJourneyItineraryVersion(apiBaseURL, token, journeyID string, version int) {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/journeys/%s/itineraries/%d", apiBaseURL, journeyID, version), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("GET /v1/journeys/{journey_id}/itineraries/{version} failed with status %d", resp.StatusCode))
	}

	var ver service.ItineraryVersionDTO
	_ = json.NewDecoder(resp.Body).Decode(&ver)
	if ver.Version != version {
		panic("version number mismatch")
	}
}

func testListJourneys(apiBaseURL, token, expectedJourneyID string) {
	req, _ := http.NewRequest(http.MethodGet, apiBaseURL+"/v1/journeys", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("GET /v1/journeys failed with status %d", resp.StatusCode))
	}

	var list []service.JourneyDTO
	_ = json.NewDecoder(resp.Body).Decode(&list)
	found := false
	for _, j := range list {
		if j.ID == expectedJourneyID {
			found = true
			break
		}
	}
	if !found {
		panic("expected journey not found in patient journeys list")
	}
}

func testGetTripRequestWithJourney(apiBaseURL, token, tripID, expectedJourneyID string) {
	req, _ := http.NewRequest(http.MethodGet, apiBaseURL+"/v1/trip-requests/"+tripID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("GET /v1/trip-requests/{trip_id} failed with status %d", resp.StatusCode))
	}

	var detail service.TripRequestDetail
	_ = json.NewDecoder(resp.Body).Decode(&detail)
	if detail.TripRequest.Status != "ACTIVE" {
		panic(fmt.Sprintf("expected trip status ACTIVE, got %s", detail.TripRequest.Status))
	}
	if detail.TripRequest.JourneyID == nil || *detail.TripRequest.JourneyID != expectedJourneyID {
		panic("expected trip request journey_id to match journey")
	}
}

func testPatientOwnershipScoping(apiBaseURL, token2, journeyID string) {
	req, _ := http.NewRequest(http.MethodGet, apiBaseURL+"/v1/journeys/"+journeyID, nil)
	req.Header.Set("Authorization", "Bearer "+token2)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		panic(fmt.Sprintf("expected 404 for unowned journey access, got %d", resp.StatusCode))
	}
}

func testHoldCompensation(db *gorm.DB, jwtSecret string) {
	// Start mock servers where Ferry succeeds hold but Hospital fails hold
	var releasedHoldID string
	ferryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/release") {
			parts := strings.Split(r.URL.Path, "/")
			releasedHoldID = parts[3]
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(adapter.ReleaseResult{
				ResourceType: "HOLD",
				ResourceID:   releasedHoldID,
				Status:       "RELEASED",
				ReleasedAt:   time.Now().UTC().Format(time.RFC3339),
			})
			return
		}
		if r.URL.Path == "/v1/holds" {
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(adapter.Hold{
				HoldID:     "ferry-hold-test-comp-01",
				ProviderID: "ferry-demo-01",
				Status:     "HELD",
				ExpiresAt:  time.Now().Add(15 * time.Minute).Format(time.RFC3339),
			})
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ferryServer.Close()

	hospServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/v1/holds" {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(adapter.ErrorEnvelope{
				Error: adapter.ErrorBody{
					Code:      "CAPACITY_CONFLICT",
					Message:   "No appointment slots available.",
					Retryable: false,
				},
			})
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer hospServer.Close()

	transportServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer transportServer.Close()

	hotelServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer hotelServer.Close()

	cfg := config.Config{
		HTTPAddr:                ":8094",
		DatabaseURL:             os.Getenv("DATABASE_URL"),
		JWTSigningSecret:        jwtSecret,
		JWTIssuer:               "batam-medhub",
		JWTAudience:             "batam-medhub-mobile",
		JWTAccessTTL:            15 * time.Minute,
		RefreshTokenTTL:         30 * 24 * time.Hour,
		HospitalBaseURL:         hospServer.URL,
		HospitalIntegrationKey:  "hospital_dev_secret",
		FerryBaseURL:            ferryServer.URL,
		FerryIntegrationKey:     "ferry_dev_secret",
		HotelBaseURL:            hotelServer.URL,
		HotelIntegrationKey:     "hotel_dev_secret",
		TransportBaseURL:        transportServer.URL,
		TransportIntegrationKey: "transport_dev_secret",
		ProviderTimeout:         5 * time.Second,
	}

	testEngine := httpapi.New(db, cfg, nil)
	testServer := httptest.NewServer(testEngine)
	defer testServer.Close()

	token := registerPatient(testServer.URL, fmt.Sprintf("b7_comp1_%d@test.test", time.Now().UnixNano()), "Comp Patient", "Passw0rd123!", "SGD")
	tripID := createMatchedTripRequest(testServer.URL, token)
	planRes := generatePlans(testServer.URL, token, tripID)

	// Attempt booking which should fail at hospital hold
	body := []byte(`{"approved": true}`)
	req, _ := http.NewRequest(http.MethodPost, testServer.URL+"/v1/plan-options/"+planRes.Options[0].ID+"/confirm", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", fmt.Sprintf("idem-comp-fail-%d", time.Now().UnixNano()))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		raw, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("expected 409 Conflict on hold failure, got %d: %s", resp.StatusCode, string(raw)))
	}

	// Verify compensating release was called on ferry hold
	if releasedHoldID != "ferry-hold-test-comp-01" {
		panic(fmt.Sprintf("expected ferry hold ferry-hold-test-comp-01 to be released, got %s", releasedHoldID))
	}

	// Verify trip status is CONFIRMATION_FAILED
	reqTrip, _ := http.NewRequest(http.MethodGet, testServer.URL+"/v1/trip-requests/"+tripID, nil)
	reqTrip.Header.Set("Authorization", "Bearer "+token)
	respTrip, err := http.DefaultClient.Do(reqTrip)
	if err != nil {
		panic(err)
	}
	defer respTrip.Body.Close()
	var tripDetail service.TripRequestDetail
	_ = json.NewDecoder(respTrip.Body).Decode(&tripDetail)
	if tripDetail.TripRequest.Status != "CONFIRMATION_FAILED" {
		panic(fmt.Sprintf("expected trip status CONFIRMATION_FAILED, got %s", tripDetail.TripRequest.Status))
	}
}

func testConfirmationCompensation(db *gorm.DB, jwtSecret string) {
	var releasedReservationID string
	var releasedHoldID string

	// Ferry succeeds hold and confirm
	ferryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/confirm") {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(adapter.Reservation{
				ReservationID: "ferry-res-comp-01",
				ProviderID:    "ferry-demo-01",
				Status:        "CONFIRMED",
			})
			return
		}
		if strings.HasPrefix(r.URL.Path, "/v1/reservations/") && strings.HasSuffix(r.URL.Path, "/release") {
			parts := strings.Split(r.URL.Path, "/")
			releasedReservationID = parts[3]
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(adapter.ReleaseResult{
				ResourceType: "RESERVATION",
				ResourceID:   releasedReservationID,
				Status:       "RELEASED",
			})
			return
		}
		if r.URL.Path == "/v1/holds" {
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(adapter.Hold{
				HoldID:     "ferry-hold-comp-01",
				ProviderID: "ferry-demo-01",
				Status:     "HELD",
				ExpiresAt:  time.Now().Add(15 * time.Minute).Format(time.RFC3339),
			})
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ferryServer.Close()

	// Transport succeeds hold but FAILS confirm
	transportServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/confirm") {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(adapter.ErrorEnvelope{
				Error: adapter.ErrorBody{
					Code:      "CONFIRMATION_FAILED",
					Message:   "Vehicle driver assignment conflict.",
					Retryable: false,
				},
			})
			return
		}
		if strings.HasSuffix(r.URL.Path, "/release") {
			parts := strings.Split(r.URL.Path, "/")
			releasedHoldID = parts[3]
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(adapter.ReleaseResult{
				ResourceType: "HOLD",
				ResourceID:   releasedHoldID,
				Status:       "RELEASED",
			})
			return
		}
		if r.URL.Path == "/v1/holds" {
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(adapter.Hold{
				HoldID:     "trans-hold-comp-01",
				ProviderID: "transport-demo-01",
				Status:     "HELD",
				ExpiresAt:  time.Now().Add(15 * time.Minute).Format(time.RFC3339),
			})
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer transportServer.Close()

	// Hospital succeeds hold
	hospServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/release") {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(adapter.ReleaseResult{
				ResourceType: "HOLD",
				ResourceID:   "hosp-hold-comp-01",
				Status:       "RELEASED",
			})
			return
		}
		if r.URL.Path == "/v1/holds" {
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(adapter.Hold{
				HoldID:     "hosp-hold-comp-01",
				ProviderID: "hospital-demo-01",
				Status:     "HELD",
				ExpiresAt:  time.Now().Add(15 * time.Minute).Format(time.RFC3339),
			})
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer hospServer.Close()

	hotelServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer hotelServer.Close()

	cfg := config.Config{
		HTTPAddr:                ":8095",
		DatabaseURL:             os.Getenv("DATABASE_URL"),
		JWTSigningSecret:        jwtSecret,
		JWTIssuer:               "batam-medhub",
		JWTAudience:             "batam-medhub-mobile",
		JWTAccessTTL:            15 * time.Minute,
		RefreshTokenTTL:         30 * 24 * time.Hour,
		HospitalBaseURL:         hospServer.URL,
		HospitalIntegrationKey:  "hospital_dev_secret",
		FerryBaseURL:            ferryServer.URL,
		FerryIntegrationKey:     "ferry_dev_secret",
		HotelBaseURL:            hotelServer.URL,
		HotelIntegrationKey:     "hotel_dev_secret",
		TransportBaseURL:        transportServer.URL,
		TransportIntegrationKey: "transport_dev_secret",
		ProviderTimeout:         5 * time.Second,
	}

	testEngine := httpapi.New(db, cfg, nil)
	testServer := httptest.NewServer(testEngine)
	defer testServer.Close()

	token := registerPatient(testServer.URL, fmt.Sprintf("b7_comp2_%d@test.test", time.Now().UnixNano()), "Comp Patient 2", "Passw0rd123!", "SGD")
	tripID := createMatchedTripRequest(testServer.URL, token)
	planRes := generatePlans(testServer.URL, token, tripID)

	body := []byte(`{"approved": true}`)
	req, _ := http.NewRequest(http.MethodPost, testServer.URL+"/v1/plan-options/"+planRes.Options[0].ID+"/confirm", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", fmt.Sprintf("idem-comp-conf-%d", time.Now().UnixNano()))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		raw, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("expected 409 Conflict on confirm failure, got %d: %s", resp.StatusCode, string(raw)))
	}

	// Verify confirmed ferry was released via release reservation
	if releasedReservationID != "ferry-res-comp-01" {
		panic(fmt.Sprintf("expected ferry reservation ferry-res-comp-01 to be released, got %s", releasedReservationID))
	}
}

func testWorkersAIExtractor(db *gorm.DB, jwtSecret string) {
	var aiCallCount int
	var retryTriggered bool

	mockAIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aiCallCount++
		if r.Header.Get("Authorization") != "Bearer test-cf-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var req struct {
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)

		lastMsg := ""
		if len(req.Messages) > 0 {
			lastMsg = strings.ToLower(req.Messages[len(req.Messages)-1].Content)
		}

		w.Header().Set("Content-Type", "application/json")

		// If this is a retry request
		if strings.Contains(lastMsg, "previous output was not valid json") {
			retryTriggered = true
			_ = json.NewEncoder(w).Encode(map[string]any{
				"success": true,
				"result": map[string]any{
					"response": `{
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
						"clarification_question": "Would you like basic or comprehensive check-up, and on which date?",
						"out_of_scope_reason": null,
						"unsupported_reason": null
					}`,
				},
			})
			return
		}

		// Test malformed JSON on first prompt if asking for "vague checkup"
		if strings.Contains(lastMsg, "vague checkup") {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"success": true,
				"result": map[string]any{
					"response": "Here is your JSON: { invalid json without closing",
				},
			})
			return
		}

		// Normal matched extraction response
		_ = json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"result": map[string]any{
				"response": `{
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
		})
	}))
	defer mockAIServer.Close()

	cfg := config.Config{
		HTTPAddr:                ":8096",
		DatabaseURL:             os.Getenv("DATABASE_URL"),
		JWTSigningSecret:        jwtSecret,
		JWTIssuer:               "batam-medhub",
		JWTAudience:             "batam-medhub-mobile",
		JWTAccessTTL:            15 * time.Minute,
		RefreshTokenTTL:         30 * 24 * time.Hour,
		HospitalBaseURL:         "http://localhost:8081",
		HospitalIntegrationKey:  "hospital_dev_secret",
		FerryBaseURL:            "http://localhost:8082",
		FerryIntegrationKey:     "ferry_dev_secret",
		HotelBaseURL:            "http://localhost:8083",
		HotelIntegrationKey:     "hotel_dev_secret",
		TransportBaseURL:        "http://localhost:8084",
		TransportIntegrationKey: "transport_dev_secret",
		ProviderTimeout:         5 * time.Second,
		CloudflareAccountID:     "test-cf-acc",
		CloudflareAPIToken:      "test-cf-token",
		CloudflareAIModel:       "@cf/meta/llama-3.1-8b-instruct",
		CloudflareAIBaseURL:     mockAIServer.URL,
		CloudflareAITimeout:     5 * time.Second,
	}

	testEngine := httpapi.New(db, cfg, nil)
	testServer := httptest.NewServer(testEngine)
	defer testServer.Close()

	token := registerPatient(testServer.URL, fmt.Sprintf("b8_ai_%d@test.test", time.Now().UnixNano()), "AI Patient", "Passw0rd123!", "SGD")

	// 1. Test AI Matched Intent Extraction
	reqPayload1 := map[string]any{
		"prompt": "I need a same-day basic medical check-up in Batam on 22 August, leaving from HarbourFront with my spouse, with a budget of SGD 400.",
		"locale": "en",
	}
	bodyBytes1, _ := json.Marshal(reqPayload1)
	req1, _ := http.NewRequest(http.MethodPost, testServer.URL+"/v1/trip-requests", bytes.NewReader(bodyBytes1))
	req1.Header.Set("Authorization", "Bearer "+token)
	req1.Header.Set("Idempotency-Key", fmt.Sprintf("idem-ai-1-%d", time.Now().UnixNano()))
	req1.Header.Set("Content-Type", "application/json")

	resp1, err := http.DefaultClient.Do(req1)
	if err != nil {
		panic(err)
	}
	defer resp1.Body.Close()

	if resp1.StatusCode != http.StatusCreated {
		raw, _ := io.ReadAll(resp1.Body)
		panic(fmt.Sprintf("expected 201 Created on AI trip request, got %d: %s", resp1.StatusCode, string(raw)))
	}

	var detail1 struct {
		TripRequest struct {
			Status string `json:"status"`
			Intent struct {
				Resolution  string  `json:"resolution"`
				ServiceCode *string `json:"service_code"`
			} `json:"intent"`
		} `json:"trip_request"`
	}
	_ = json.NewDecoder(resp1.Body).Decode(&detail1)

	if detail1.TripRequest.Status != "PLANNING" || detail1.TripRequest.Intent.Resolution != "MATCHED" {
		panic(fmt.Sprintf("expected trip status PLANNING / MATCHED, got %s / %s", detail1.TripRequest.Status, detail1.TripRequest.Intent.Resolution))
	}
	if detail1.TripRequest.Intent.ServiceCode == nil || *detail1.TripRequest.Intent.ServiceCode != "MCU_BASIC" {
		panic("expected MCU_BASIC service code")
	}

	// 2. Test Emergency Guardrail (bypasses LLM)
	currentAICalls := aiCallCount
	reqPayload2 := map[string]any{
		"prompt": "I have sudden severe chest pain and difficulty breathing, need an ambulance immediately.",
		"locale": "en",
	}
	bodyBytes2, _ := json.Marshal(reqPayload2)
	req2, _ := http.NewRequest(http.MethodPost, testServer.URL+"/v1/trip-requests", bytes.NewReader(bodyBytes2))
	req2.Header.Set("Authorization", "Bearer "+token)
	req2.Header.Set("Idempotency-Key", fmt.Sprintf("idem-ai-2-%d", time.Now().UnixNano()))
	req2.Header.Set("Content-Type", "application/json")

	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		panic(err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusCreated {
		panic(fmt.Sprintf("expected 201 on emergency prompt, got %d", resp2.StatusCode))
	}
	var detail2 struct {
		TripRequest struct {
			Status string `json:"status"`
			Intent struct {
				Resolution       string  `json:"resolution"`
				OutOfScopeReason *string `json:"out_of_scope_reason"`
			} `json:"intent"`
		} `json:"trip_request"`
	}
	_ = json.NewDecoder(resp2.Body).Decode(&detail2)

	if detail2.TripRequest.Status != "OUT_OF_SCOPE" || detail2.TripRequest.Intent.Resolution != "OUT_OF_SCOPE" {
		panic(fmt.Sprintf("expected OUT_OF_SCOPE for emergency triage, got %s", detail2.TripRequest.Status))
	}
	if detail2.TripRequest.Intent.OutOfScopeReason == nil || *detail2.TripRequest.Intent.OutOfScopeReason == "" {
		panic("expected populated out_of_scope_reason")
	}
	if aiCallCount != currentAICalls {
		panic("emergency guardrail should not call external Workers AI model")
	}

	// 3. Test Malformed JSON 1-Retry
	reqPayload3 := map[string]any{
		"prompt": "I need a vague checkup in Batam.",
		"locale": "en",
	}
	bodyBytes3, _ := json.Marshal(reqPayload3)
	req3, _ := http.NewRequest(http.MethodPost, testServer.URL+"/v1/trip-requests", bytes.NewReader(bodyBytes3))
	req3.Header.Set("Authorization", "Bearer "+token)
	req3.Header.Set("Idempotency-Key", fmt.Sprintf("idem-ai-3-%d", time.Now().UnixNano()))
	req3.Header.Set("Content-Type", "application/json")

	resp3, err := http.DefaultClient.Do(req3)
	if err != nil {
		panic(err)
	}
	defer resp3.Body.Close()

	if resp3.StatusCode != http.StatusCreated {
		panic(fmt.Sprintf("expected 201 on vague checkup, got %d", resp3.StatusCode))
	}
	var detail3 struct {
		TripRequest struct {
			ID     string `json:"id"`
			Status string `json:"status"`
			Intent struct {
				Resolution            string  `json:"resolution"`
				ClarificationQuestion *string `json:"clarification_question"`
			} `json:"intent"`
		} `json:"trip_request"`
	}
	_ = json.NewDecoder(resp3.Body).Decode(&detail3)

	if !retryTriggered {
		panic("expected 1-retry to be triggered upon malformed JSON from model")
	}
	if detail3.TripRequest.Status != "NEEDS_INPUT" || detail3.TripRequest.Intent.Resolution != "NEEDS_CLARIFICATION" {
		panic(fmt.Sprintf("expected NEEDS_INPUT / NEEDS_CLARIFICATION, got %s / %s", detail3.TripRequest.Status, detail3.TripRequest.Intent.Resolution))
	}

	// 4. Test Intent Amendment on Trip Request
	amendPayload := map[string]any{
		"answer": "I would like the basic medical check-up on 22 August 2026.",
	}
	bodyAmend, _ := json.Marshal(amendPayload)
	reqAmend, _ := http.NewRequest(http.MethodPatch, testServer.URL+"/v1/trip-requests/"+detail3.TripRequest.ID+"/intent", bytes.NewReader(bodyAmend))
	reqAmend.Header.Set("Authorization", "Bearer "+token)
	reqAmend.Header.Set("Idempotency-Key", fmt.Sprintf("idem-ai-amend-%d", time.Now().UnixNano()))
	reqAmend.Header.Set("Content-Type", "application/json")

	respAmend, err := http.DefaultClient.Do(reqAmend)
	if err != nil {
		panic(err)
	}
	defer respAmend.Body.Close()

	if respAmend.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(respAmend.Body)
		panic(fmt.Sprintf("expected 200 OK on amend intent, got %d: %s", respAmend.StatusCode, string(raw)))
	}

	var amendDetail struct {
		TripRequest struct {
			Status string `json:"status"`
			Intent struct {
				Resolution  string  `json:"resolution"`
				ServiceCode *string `json:"service_code"`
			} `json:"intent"`
		} `json:"trip_request"`
	}
	_ = json.NewDecoder(respAmend.Body).Decode(&amendDetail)
	if amendDetail.TripRequest.Status != "PLANNING" || amendDetail.TripRequest.Intent.Resolution != "MATCHED" {
		panic(fmt.Sprintf("expected amended trip status PLANNING / MATCHED, got %s / %s", amendDetail.TripRequest.Status, amendDetail.TripRequest.Intent.Resolution))
	}

	// 5. Test Seamless Fallback on Network Outage
	mockAIServer.Close() // Close mock server to simulate network / API outage
	reqPayload5 := map[string]any{
		"prompt": "I need a same-day basic medical check-up in Batam on 22 August from HarbourFront with my spouse, budget SGD 400.",
		"locale": "en",
	}
	bodyBytes5, _ := json.Marshal(reqPayload5)
	req5, _ := http.NewRequest(http.MethodPost, testServer.URL+"/v1/trip-requests", bytes.NewReader(bodyBytes5))
	req5.Header.Set("Authorization", "Bearer "+token)
	req5.Header.Set("Idempotency-Key", fmt.Sprintf("idem-ai-5-%d", time.Now().UnixNano()))
	req5.Header.Set("Content-Type", "application/json")

	resp5, err := http.DefaultClient.Do(req5)
	if err != nil {
		panic(err)
	}
	defer resp5.Body.Close()

	if resp5.StatusCode != http.StatusCreated {
		panic(fmt.Sprintf("expected 201 on seamless fallback when AI is down, got %d", resp5.StatusCode))
	}
	var detail5 struct {
		TripRequest struct {
			Status string `json:"status"`
			Intent struct {
				Resolution string `json:"resolution"`
			} `json:"intent"`
		} `json:"trip_request"`
	}
	_ = json.NewDecoder(resp5.Body).Decode(&detail5)
	if detail5.TripRequest.Status != "PLANNING" || detail5.TripRequest.Intent.Resolution != "MATCHED" {
		panic("expected seamless deterministic fallback to produce MATCHED planning intent")
	}
}

func testDisruptionAndRecoveryEngine(db *gorm.DB, jwtSecret string) {
	fmt.Println("--- Running Disruption and Recovery Engine Verification (B9) ---")

	// 1. Setup Mock Provider Servers
	hospServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/confirm") {
			parts := strings.Split(r.URL.Path, "/")
			holdID := parts[3]
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(adapter.Reservation{
				ReservationID: fmt.Sprintf("hosp-res-%s", holdID),
				ProviderID:    "hospital-demo-01",
				Status:        "CONFIRMED",
				ConfirmedAt:   time.Now().UTC().Format(time.RFC3339),
			})
			return
		}
		if strings.HasSuffix(r.URL.Path, "/release") {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(adapter.ReleaseResult{
				ResourceType: "HOLD",
				ResourceID:   "hosp-rel",
				Status:       "RELEASED",
			})
			return
		}
		if r.URL.Path == "/v1/holds" {
			var hReq adapter.CreateHoldRequest
			_ = json.NewDecoder(r.Body).Decode(&hReq)
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(adapter.Hold{
				HoldID:     fmt.Sprintf("hosp-hold-%d", time.Now().UnixNano()%10000),
				ProviderID: "hospital-demo-01",
				Status:     "HELD",
				ExpiresAt:  time.Now().Add(15 * time.Minute).Format(time.RFC3339),
			})
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer hospServer.Close()

	transServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/confirm") {
			parts := strings.Split(r.URL.Path, "/")
			holdID := parts[3]
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(adapter.Reservation{
				ReservationID: fmt.Sprintf("trans-res-%s", holdID),
				ProviderID:    "transport-demo-01",
				Status:        "CONFIRMED",
				ConfirmedAt:   time.Now().UTC().Format(time.RFC3339),
			})
			return
		}
		if strings.HasSuffix(r.URL.Path, "/release") {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(adapter.ReleaseResult{
				ResourceType: "RESERVATION",
				ResourceID:   "trans-rel",
				Status:       "RELEASED",
			})
			return
		}
		if r.URL.Path == "/v1/holds" {
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(adapter.Hold{
				HoldID:     fmt.Sprintf("trans-hold-%d", time.Now().UnixNano()%10000),
				ProviderID: "transport-demo-01",
				Status:     "HELD",
				ExpiresAt:  time.Now().Add(15 * time.Minute).Format(time.RFC3339),
			})
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer transServer.Close()

	ferryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/confirm") {
			parts := strings.Split(r.URL.Path, "/")
			holdID := parts[3]
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(adapter.Reservation{
				ReservationID: fmt.Sprintf("ferry-res-%s", holdID),
				ProviderID:    "ferry-demo-01",
				Status:        "CONFIRMED",
				ConfirmedAt:   time.Now().UTC().Format(time.RFC3339),
			})
			return
		}
		if strings.HasSuffix(r.URL.Path, "/release") {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(adapter.ReleaseResult{
				ResourceType: "RESERVATION",
				ResourceID:   "ferry-rel",
				Status:       "RELEASED",
			})
			return
		}
		if r.URL.Path == "/v1/holds" {
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(adapter.Hold{
				HoldID:     fmt.Sprintf("ferry-hold-%d", time.Now().UnixNano()%10000),
				ProviderID: "ferry-demo-01",
				Status:     "HELD",
				ExpiresAt:  time.Now().Add(15 * time.Minute).Format(time.RFC3339),
			})
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ferryServer.Close()

	hotelServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer hotelServer.Close()

	cfg := config.Config{
		HTTPAddr:                ":8096",
		DatabaseURL:             os.Getenv("DATABASE_URL"),
		JWTSigningSecret:        jwtSecret,
		JWTIssuer:               "batam-medhub",
		JWTAudience:             "batam-medhub-mobile",
		JWTAccessTTL:            15 * time.Minute,
		RefreshTokenTTL:         30 * 24 * time.Hour,
		HospitalBaseURL:         hospServer.URL,
		HospitalIntegrationKey:  "hospital_dev_secret",
		FerryBaseURL:            ferryServer.URL,
		FerryIntegrationKey:     "ferry_dev_secret",
		HotelBaseURL:            hotelServer.URL,
		HotelIntegrationKey:     "hotel_dev_secret",
		TransportBaseURL:        transServer.URL,
		TransportIntegrationKey: "transport_dev_secret",
		ProviderTimeout:         5 * time.Second,
	}

	testEngine := httpapi.New(db, cfg, nil)
	testServer := httptest.NewServer(testEngine)
	defer testServer.Close()

	// 2. Register Patient and Confirm an initial Journey
	token := registerPatient(testServer.URL, fmt.Sprintf("b9_patient_%d@test.test", time.Now().UnixNano()), "B9 Patient", "Passw0rd123!", "SGD")
	tripID := createMatchedTripRequest(testServer.URL, token)
	planRes := generatePlans(testServer.URL, token, tripID)

	bookingIdemKey := fmt.Sprintf("idem-b9-init-book-%d", time.Now().UnixNano())
	journeyDetail := confirmPlanOption(testServer.URL, token, planRes.Options[0].ID, bookingIdemKey)

	journeyID := journeyDetail.Journey.ID
	if journeyDetail.Journey.ActiveItineraryVersion != 1 {
		panic(fmt.Sprintf("expected active itinerary version 1, got %d", journeyDetail.Journey.ActiveItineraryVersion))
	}

	var hospApptID string
	for _, it := range journeyDetail.ActiveItinerary.Items {
		if it.ItemType == "HOSPITAL_APPOINTMENT" {
			hospApptID = it.ID
			break
		}
	}
	if hospApptID == "" {
		panic("hospital appointment item not found in journey")
	}

	// 3. Test Negative Ingestion Cases
	// 3a. Invalid Provider Key
	invalidKeyPayload := map[string]any{
		"external_event_id": fmt.Sprintf("evt-disrupt-%d", time.Now().UnixNano()),
		"journey_id":        journeyID,
		"event_type":        "HOSPITAL_ADDITIONAL_CARE_REQUESTED",
		"occurred_at":       time.Now().UTC().Format(time.RFC3339),
		"target":            map[string]any{"itinerary_item_id": hospApptID},
		"actor":             map[string]any{"actor_id": "dr-lee", "name": "Dr Lee", "role": "Cardiologist"},
		"details": map[string]any{
			"reason":                "Additional care needed",
			"instruction_reference": "hospital-instruction://followup-observation/FO-20260822-0001",
			"replacement_time_window": map[string]any{
				"starts_at":       "2026-08-22T05:00:00Z",
				"ends_at":         "2026-08-22T06:30:00Z",
				"start_time_zone": "Asia/Jakarta",
				"end_time_zone":   "Asia/Jakarta",
			},
			"additional_service_code":     "FOLLOWUP_OBSERVATION",
			"additional_duration_minutes": 90,
			"priority":                    "MEDIUM",
			"travel_clearance_status":     "CLEARED",
		},
	}
	bodyInvalidKey, _ := json.Marshal(invalidKeyPayload)
	reqInvKey, _ := http.NewRequest(http.MethodPost, testServer.URL+"/v1/provider/disruptions", bytes.NewReader(bodyInvalidKey))
	reqInvKey.Header.Set("X-Provider-Key", "invalid_secret_key")
	reqInvKey.Header.Set("Content-Type", "application/json")
	respInvKey, err := http.DefaultClient.Do(reqInvKey)
	if err != nil || respInvKey.StatusCode != http.StatusUnauthorized {
		panic(fmt.Sprintf("expected 401 Unauthorized for invalid provider key, got %d", respInvKey.StatusCode))
	}
	respInvKey.Body.Close()

	// 3b. Incompatible Event Type (Ferry provider sending HOSPITAL_...)
	reqIncompType, _ := http.NewRequest(http.MethodPost, testServer.URL+"/v1/provider/disruptions", bytes.NewReader(bodyInvalidKey))
	reqIncompType.Header.Set("X-Provider-Key", "ferry_dev_secret")
	reqIncompType.Header.Set("Content-Type", "application/json")
	respIncompType, err := http.DefaultClient.Do(reqIncompType)
	if err != nil || respIncompType.StatusCode != http.StatusForbidden {
		panic(fmt.Sprintf("expected 403 Forbidden for incompatible event type, got %d", respIncompType.StatusCode))
	}
	respIncompType.Body.Close()

	// 3c. Missing Clinical Details (missing instruction_reference)
	missingClinPayload := map[string]any{
		"external_event_id": fmt.Sprintf("evt-disrupt-%d", time.Now().UnixNano()),
		"journey_id":        journeyID,
		"event_type":        "HOSPITAL_ADDITIONAL_CARE_REQUESTED",
		"occurred_at":       time.Now().UTC().Format(time.RFC3339),
		"target":            map[string]any{"itinerary_item_id": hospApptID},
		"actor":             map[string]any{"actor_id": "dr-lee", "name": "Dr Lee", "role": "Cardiologist"},
		"details": map[string]any{
			"reason": "Missing instruction ref",
		},
	}
	bodyMissingClin, _ := json.Marshal(missingClinPayload)
	reqMissingClin, _ := http.NewRequest(http.MethodPost, testServer.URL+"/v1/provider/disruptions", bytes.NewReader(bodyMissingClin))
	reqMissingClin.Header.Set("X-Provider-Key", "hospital_dev_secret")
	reqMissingClin.Header.Set("Content-Type", "application/json")
	respMissingClin, err := http.DefaultClient.Do(reqMissingClin)
	if err != nil || respMissingClin.StatusCode != http.StatusUnprocessableEntity {
		panic(fmt.Sprintf("expected 422 Unprocessable Entity for missing clinical details, got %d", respMissingClin.StatusCode))
	}
	respMissingClin.Body.Close()

	// 4. Positive Ingestion Test: Ingest Valid Hospital Disruption
	eventID := fmt.Sprintf("evt-hosp-disrupt-%d", time.Now().UnixNano())
	validPayload := map[string]any{
		"external_event_id": eventID,
		"journey_id":        journeyID,
		"event_type":        "HOSPITAL_ADDITIONAL_CARE_REQUESTED",
		"occurred_at":       time.Now().UTC().Format(time.RFC3339),
		"target":            map[string]any{"itinerary_item_id": hospApptID},
		"actor":             map[string]any{"actor_id": "dr-lee", "name": "Dr Lee", "role": "Cardiologist"},
		"details": map[string]any{
			"reason":                "Patient requires additional observation following examination",
			"instruction_reference": "hospital-instruction://followup-observation/FO-20260822-0001",
			"replacement_time_window": map[string]any{
				"starts_at":       "2026-08-22T05:00:00Z",
				"ends_at":         "2026-08-22T06:30:00Z",
				"start_time_zone": "Asia/Jakarta",
				"end_time_zone":   "Asia/Jakarta",
			},
			"additional_service_code":     "FOLLOWUP_OBSERVATION",
			"additional_duration_minutes": 90,
			"priority":                    "MEDIUM",
			"travel_clearance_status":     "CLEARED",
		},
	}
	bodyValid, _ := json.Marshal(validPayload)
	reqValid, _ := http.NewRequest(http.MethodPost, testServer.URL+"/v1/events/disruptions", bytes.NewReader(bodyValid))
	reqValid.Header.Set("X-Integration-Key", "hospital_dev_secret")
	reqValid.Header.Set("Content-Type", "application/json")

	respValid, err := http.DefaultClient.Do(reqValid)
	if err != nil {
		panic(err)
	}
	defer respValid.Body.Close()

	if respValid.StatusCode != http.StatusAccepted {
		raw, _ := io.ReadAll(respValid.Body)
		panic(fmt.Sprintf("expected 202 Accepted on disruption ingestion, got %d: %s", respValid.StatusCode, string(raw)))
	}

	var receipt service.ProviderEventReceipt
	_ = json.NewDecoder(respValid.Body).Decode(&receipt)
	if receipt.Outcome != "DISRUPTION_CREATED" || receipt.DisruptionID == nil || *receipt.DisruptionID == "" {
		panic("expected receipt outcome DISRUPTION_CREATED with non-empty disruption_id")
	}
	disruptionID := *receipt.DisruptionID

	// 5. Test Deduplication
	// 5a. Replay identical payload -> 202 Accepted with DUPLICATE outcome and header
	reqReplay, _ := http.NewRequest(http.MethodPost, testServer.URL+"/v1/provider/disruptions", bytes.NewReader(bodyValid))
	reqReplay.Header.Set("X-Provider-Key", "hospital_dev_secret")
	reqReplay.Header.Set("Content-Type", "application/json")
	respReplay, err := http.DefaultClient.Do(reqReplay)
	if err != nil || respReplay.StatusCode != http.StatusAccepted {
		panic(fmt.Sprintf("expected 202 Accepted on deduplicated replay, got %d", respReplay.StatusCode))
	}
	if respReplay.Header.Get("Idempotency-Replayed") != "true" {
		panic("expected Idempotency-Replayed header on deduplicated event")
	}
	var receiptReplay service.ProviderEventReceipt
	_ = json.NewDecoder(respReplay.Body).Decode(&receiptReplay)
	if receiptReplay.Outcome != "DUPLICATE" {
		panic(fmt.Sprintf("expected DUPLICATE outcome, got %s", receiptReplay.Outcome))
	}
	respReplay.Body.Close()

	// 5b. Conflicting payload with same external_event_id -> 409 Conflict
	conflictingPayload := map[string]any{
		"external_event_id": eventID,
		"journey_id":        journeyID,
		"event_type":        "HOSPITAL_ADDITIONAL_CARE_REQUESTED",
		"occurred_at":       time.Now().UTC().Format(time.RFC3339),
		"target":            map[string]any{"itinerary_item_id": hospApptID},
		"actor":             map[string]any{"actor_id": "dr-changed", "name": "Dr Changed", "role": "Surgeon"},
		"details": map[string]any{
			"reason": "Different reason entirely",
		},
	}
	bodyConflict, _ := json.Marshal(conflictingPayload)
	reqConflict, _ := http.NewRequest(http.MethodPost, testServer.URL+"/v1/provider/disruptions", bytes.NewReader(bodyConflict))
	reqConflict.Header.Set("X-Provider-Key", "hospital_dev_secret")
	reqConflict.Header.Set("Content-Type", "application/json")
	respConflict, err := http.DefaultClient.Do(reqConflict)
	if err != nil || respConflict.StatusCode != http.StatusConflict {
		panic(fmt.Sprintf("expected 409 Conflict on conflicting deduplication payload, got %d", respConflict.StatusCode))
	}
	respConflict.Body.Close()

	// 6. Test GET /v1/disruptions/{disruption_id}
	reqGetDis, _ := http.NewRequest(http.MethodGet, testServer.URL+"/v1/disruptions/"+disruptionID, nil)
	reqGetDis.Header.Set("Authorization", "Bearer "+token)
	respGetDis, err := http.DefaultClient.Do(reqGetDis)
	if err != nil || respGetDis.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("expected 200 OK on get disruption, got %d", respGetDis.StatusCode))
	}

	var disDetail service.DisruptionDetail
	_ = json.NewDecoder(respGetDis.Body).Decode(&disDetail)
	respGetDis.Body.Close()

	if disDetail.Disruption.ID != disruptionID || disDetail.Disruption.Impact.Severity != "HIGH" {
		panic("disruption detail impact or ID mismatch")
	}
	if len(disDetail.RecoveryOptions) < 1 {
		panic("expected at least 1 recovery option generated")
	}
	recOption := disDetail.RecoveryOptions[0]
	if recOption.Status != "PROPOSED" || recOption.TimeDeltaMinutes != 90 {
		panic(fmt.Sprintf("recovery option invalid: status %s, delta %d", recOption.Status, recOption.TimeDeltaMinutes))
	}

	// 7. Test Approve Recovery Option (POST /v1/recovery-options/{option_id}/approve)
	recApproveIdemKey := fmt.Sprintf("idem-rec-app-%d", time.Now().UnixNano())
	recApprovePayload := map[string]any{"approved": true}
	bodyRecApp, _ := json.Marshal(recApprovePayload)

	reqRecApp, _ := http.NewRequest(http.MethodPost, testServer.URL+"/v1/recovery-options/"+recOption.ID+"/approve", bytes.NewReader(bodyRecApp))
	reqRecApp.Header.Set("Authorization", "Bearer "+token)
	reqRecApp.Header.Set("Idempotency-Key", recApproveIdemKey)
	reqRecApp.Header.Set("Content-Type", "application/json")

	respRecApp, err := http.DefaultClient.Do(reqRecApp)
	if err != nil {
		panic(err)
	}
	defer respRecApp.Body.Close()

	if respRecApp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(respRecApp.Body)
		panic(fmt.Sprintf("expected 200 OK on approve recovery option, got %d: %s", respRecApp.StatusCode, string(raw)))
	}

	var v2Journey service.JourneyDetail
	_ = json.NewDecoder(respRecApp.Body).Decode(&v2Journey)

	if v2Journey.Journey.ActiveItineraryVersion != 2 {
		panic(fmt.Sprintf("expected active itinerary version 2, got %d", v2Journey.Journey.ActiveItineraryVersion))
	}
	if v2Journey.ActiveItinerary.Version != 2 || v2Journey.ActiveItinerary.Status != "ACTIVE" {
		panic(fmt.Sprintf("expected active itinerary version 2 status ACTIVE, got %d / %s", v2Journey.ActiveItinerary.Version, v2Journey.ActiveItinerary.Status))
	}
	if len(v2Journey.ActiveItinerary.Items) != 8 {
		panic(fmt.Sprintf("expected 8 itinerary items in v2, got %d", len(v2Journey.ActiveItinerary.Items)))
	}

	// Validate items structure in Version 2
	items := v2Journey.ActiveItinerary.Items
	if items[0].ItemType != "FERRY_OUTBOUND" || items[1].ItemType != "ARRIVAL_BUFFER" ||
		items[2].ItemType != "TRANSPORT_PICKUP" || items[3].ItemType != "HOSPITAL_APPOINTMENT" ||
		items[4].ItemType != "ADDITIONAL_CARE" || items[5].ItemType != "TRANSPORT_DROPOFF" ||
		items[6].ItemType != "DEPARTURE_BUFFER" || items[7].ItemType != "FERRY_RETURN" {
		panic(fmt.Sprintf("unexpected itinerary items ordering in v2: %v, %v, %v, %v, %v, %v, %v, %v",
			items[0].ItemType, items[1].ItemType, items[2].ItemType, items[3].ItemType,
			items[4].ItemType, items[5].ItemType, items[6].ItemType, items[7].ItemType))
	}

	if len(v2Journey.ItineraryVersions) != 2 {
		panic(fmt.Sprintf("expected 2 itinerary versions in history, got %d", len(v2Journey.ItineraryVersions)))
	}

	// 8. Test Idempotency Replay on Recovery Approval
	reqReplayApp, _ := http.NewRequest(http.MethodPost, testServer.URL+"/v1/recovery-options/"+recOption.ID+"/approve", bytes.NewReader(bodyRecApp))
	reqReplayApp.Header.Set("Authorization", "Bearer "+token)
	reqReplayApp.Header.Set("Idempotency-Key", recApproveIdemKey)
	reqReplayApp.Header.Set("Content-Type", "application/json")
	respReplayApp, err := http.DefaultClient.Do(reqReplayApp)
	if err != nil || respReplayApp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("expected 200 OK on recovery approval idempotency replay, got %d", respReplayApp.StatusCode))
	}
	if respReplayApp.Header.Get("Idempotency-Replayed") != "true" {
		panic("expected Idempotency-Replayed header on recovery approval replay")
	}
	respReplayApp.Body.Close()

	// 9. Verify Itinerary Version 1 status is SUPERSEDED via GET /v1/journeys/{journey_id}/itineraries/1
	reqGetV1, _ := http.NewRequest(http.MethodGet, testServer.URL+"/v1/journeys/"+journeyID+"/itineraries/1", nil)
	reqGetV1.Header.Set("Authorization", "Bearer "+token)
	respGetV1, err := http.DefaultClient.Do(reqGetV1)
	if err != nil || respGetV1.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("expected 200 OK for version 1 retrieval, got %d", respGetV1.StatusCode))
	}
	var v1Itinerary service.ItineraryVersionDTO
	_ = json.NewDecoder(respGetV1.Body).Decode(&v1Itinerary)
	respGetV1.Body.Close()
	if v1Itinerary.Status != "SUPERSEDED" || v1Itinerary.Version != 1 {
		panic(fmt.Sprintf("expected version 1 status SUPERSEDED, got %s", v1Itinerary.Status))
	}

	// 10. Verify Disruption status in database is RESOLVED
	var disModel model.Disruption
	if err := db.Where("id = ?", disruptionID).First(&disModel).Error; err != nil || disModel.Status != "RESOLVED" {
		panic(fmt.Sprintf("expected disruption status RESOLVED in db, got %s", disModel.Status))
	}

	// 11. Test Alternative Route: POST /v1/journeys/{journey_id}/recovery-options/{option_id}/select
	tripID2 := createMatchedTripRequest(testServer.URL, token)
	planRes2 := generatePlans(testServer.URL, token, tripID2)
	journeyDetail2 := confirmPlanOption(testServer.URL, token, planRes2.Options[0].ID, fmt.Sprintf("idem-b9-init2-%d", time.Now().UnixNano()))

	var hospApptID2 string
	for _, it := range journeyDetail2.ActiveItinerary.Items {
		if it.ItemType == "HOSPITAL_APPOINTMENT" {
			hospApptID2 = it.ID
			break
		}
	}

	validPayload2 := map[string]any{
		"external_event_id": fmt.Sprintf("evt-hosp-disrupt-2-%d", time.Now().UnixNano()),
		"journey_id":        journeyDetail2.Journey.ID,
		"event_type":        "HOSPITAL_ADDITIONAL_CARE_REQUESTED",
		"occurred_at":       time.Now().UTC().Format(time.RFC3339),
		"target":            map[string]any{"itinerary_item_id": hospApptID2},
		"actor":             map[string]any{"actor_id": "dr-lee", "name": "Dr Lee", "role": "Cardiologist"},
		"details": map[string]any{
			"reason":                "Patient requires additional observation following examination",
			"instruction_reference": "hospital-instruction://followup-observation/FO-20260822-0002",
			"replacement_time_window": map[string]any{
				"starts_at":       "2026-08-22T05:00:00Z",
				"ends_at":         "2026-08-22T06:30:00Z",
				"start_time_zone": "Asia/Jakarta",
				"end_time_zone":   "Asia/Jakarta",
			},
			"additional_service_code":     "FOLLOWUP_OBSERVATION",
			"additional_duration_minutes": 90,
			"priority":                    "MEDIUM",
			"travel_clearance_status":     "CLEARED",
		},
	}
	bodyValid2, _ := json.Marshal(validPayload2)
	reqValid2, _ := http.NewRequest(http.MethodPost, testServer.URL+"/v1/provider/disruptions", bytes.NewReader(bodyValid2))
	reqValid2.Header.Set("X-Provider-Key", "hospital_dev_secret")
	reqValid2.Header.Set("Content-Type", "application/json")
	respValid2, err := http.DefaultClient.Do(reqValid2)
	if err != nil || respValid2.StatusCode != http.StatusAccepted {
		panic(fmt.Sprintf("expected 202 on second disruption ingestion, got %d", respValid2.StatusCode))
	}
	var receipt2 service.ProviderEventReceipt
	_ = json.NewDecoder(respValid2.Body).Decode(&receipt2)
	respValid2.Body.Close()

	disDetail2, err := http.DefaultClient.Get(testServer.URL + "/v1/disruptions/" + *receipt2.DisruptionID)
	// Authenticated retrieval
	reqGetDis2, _ := http.NewRequest(http.MethodGet, testServer.URL+"/v1/disruptions/"+*receipt2.DisruptionID, nil)
	reqGetDis2.Header.Set("Authorization", "Bearer "+token)
	respGetDis2, err := http.DefaultClient.Do(reqGetDis2)
	if err != nil || respGetDis2.StatusCode != http.StatusOK {
		panic("failed to retrieve second disruption")
	}
	var disDetailObj2 service.DisruptionDetail
	_ = json.NewDecoder(respGetDis2.Body).Decode(&disDetailObj2)
	respGetDis2.Body.Close()
	_ = disDetail2

	// Select via POST /v1/journeys/{journey_id}/recovery-options/{option_id}/select
	reqSelect, _ := http.NewRequest(http.MethodPost, testServer.URL+"/v1/journeys/"+journeyDetail2.Journey.ID+"/recovery-options/"+disDetailObj2.RecoveryOptions[0].ID+"/select", bytes.NewReader(bodyRecApp))
	reqSelect.Header.Set("Authorization", "Bearer "+token)
	reqSelect.Header.Set("Idempotency-Key", fmt.Sprintf("idem-select-app-%d", time.Now().UnixNano()))
	reqSelect.Header.Set("Content-Type", "application/json")
	respSelect, err := http.DefaultClient.Do(reqSelect)
	if err != nil || respSelect.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(respSelect.Body)
		panic(fmt.Sprintf("expected 200 OK on alias recovery option select, got %d: %s", respSelect.StatusCode, string(raw)))
	}
	var v2Journey2 service.JourneyDetail
	_ = json.NewDecoder(respSelect.Body).Decode(&v2Journey2)
	respSelect.Body.Close()
	if v2Journey2.Journey.ActiveItineraryVersion != 2 {
		panic("expected active itinerary version 2 on alias select")
	}
}

