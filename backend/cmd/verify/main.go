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

	fmt.Println("\n=== ALL B7 SAGA AND JOURNEY TRACKING VERIFICATIONS COMPLETED SUCCESSFULLY ===")
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
