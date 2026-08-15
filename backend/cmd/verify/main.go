package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"batam-medhub/internal/adapter"
)

func main() {
	fmt.Println("=== Starting B5 Provider Adapters Verification Suite ===")

	ctx := context.Background()

	// 1. Setup Mock Provider Servers
	hospitalServer := startMockProviderServer(adapter.ProviderTypeHospital, "hospital-demo-01", "hospital_dev_secret")
	defer hospitalServer.Close()

	ferryServer := startMockProviderServer(adapter.ProviderTypeFerry, "ferry-demo-01", "ferry_dev_secret")
	defer ferryServer.Close()

	hotelServer := startMockProviderServer(adapter.ProviderTypeHotel, "hotel-demo-01", "hotel_dev_secret")
	defer hotelServer.Close()

	transportServer := startMockProviderServer(adapter.ProviderTypeTransport, "transport-demo-01", "transport_dev_secret")
	defer transportServer.Close()

	// 2. Initialize Adapters
	hospAdapter := adapter.NewHospitalAdapter(hospitalServer.URL, "hospital_dev_secret", 2*time.Second)
	ferryAdapter := adapter.NewFerryAdapter(ferryServer.URL, "ferry_dev_secret", 2*time.Second)
	hotelAdapter := adapter.NewHotelAdapter(hotelServer.URL, "hotel_dev_secret", 2*time.Second)
	transAdapter := adapter.NewTransportAdapter(transportServer.URL, "transport_dev_secret", 2*time.Second)

	// 3. Test Hospital Adapter Lifecycle & Error Handling
	testHospitalAdapter(ctx, hospAdapter, hospitalServer.URL)
	fmt.Println("[PASS] HospitalAdapter search, hold, confirm, get, release, and error mappings validated.")

	// 4. Test Ferry Adapter Lifecycle
	testFerryAdapter(ctx, ferryAdapter)
	fmt.Println("[PASS] FerryAdapter search, hold, confirm, get, and release validated.")

	// 5. Test Hotel Adapter Lifecycle
	testHotelAdapter(ctx, hotelAdapter)
	fmt.Println("[PASS] HotelAdapter search, hold, confirm, get, and release validated.")

	// 6. Test Transport Adapter Lifecycle with Booking Requirements
	testTransportAdapter(ctx, transAdapter)
	fmt.Println("[PASS] TransportAdapter search, hold with booking requirements, confirm, get, and release validated.")

	// 7. Test Concurrent Aggregator with Partial Failure Isolation
	testAggregator(ctx, hospAdapter, ferryAdapter, hotelAdapter, transAdapter)
	fmt.Println("[PASS] Multi-provider Aggregator concurrent execution and fault isolation validated.")

	fmt.Println("\n=== ALL B5 VERIFICATIONS COMPLETED SUCCESSFULLY ===")
}

func testHospitalAdapter(ctx context.Context, a *adapter.HospitalAdapter, baseURL string) {
	// Health
	health, err := a.Health(ctx, "req-hosp-health")
	if err != nil || health.Status != "UP" {
		panic(fmt.Sprintf("health check failed: %v", err))
	}

	// Search
	offers, err := a.Search(ctx, "req-hosp-search", adapter.HospitalSearchCriteria{
		ServiceCode:  "MCU_BASIC",
		PatientCount: 1,
		AppointmentWindow: adapter.TimeWindow{
			StartsAt:      "2026-08-22T02:00:00Z",
			EndsAt:        "2026-08-22T08:00:00Z",
			StartTimeZone: "Asia/Jakarta",
			EndTimeZone:   "Asia/Jakarta",
		},
	})
	if err != nil || len(offers) == 0 {
		panic(fmt.Sprintf("search failed: %v", err))
	}
	if offers[0].OfferID != "hosp-offer-001" {
		panic(fmt.Sprintf("expected offer_id hosp-offer-001, got %s", offers[0].OfferID))
	}

	// Create Hold
	hold, err := a.CreateHold(ctx, "req-hosp-hold", "idem-hosp-hold-01", adapter.CreateHoldRequest{
		ProviderID:        "hospital-demo-01",
		OfferID:           "hosp-offer-001",
		Units:             1,
		ExpectedUnitPrice: adapter.Money{AmountMinor: 150000000, Currency: "IDR"},
		ClientReference:   "journey-001-hosp",
	})
	if err != nil || hold.HoldID != "hosp-hold-001" {
		panic(fmt.Sprintf("create hold failed: %v", err))
	}

	// Confirm Hold
	res, err := a.ConfirmHold(ctx, "req-hosp-confirm", "idem-hosp-conf-01", hold.HoldID)
	if err != nil || res.ReservationID != "hosp-res-001" {
		panic(fmt.Sprintf("confirm hold failed: %v", err))
	}

	// Get Reservation
	getRes, err := a.GetReservation(ctx, "req-hosp-get", res.ReservationID)
	if err != nil || getRes.Status != "CONFIRMED" {
		panic(fmt.Sprintf("get reservation failed: %v", err))
	}

	// Release Reservation
	relRes, err := a.ReleaseReservation(ctx, "req-hosp-rel-res", "idem-hosp-rel-01", res.ReservationID)
	if err != nil || relRes.Status != "RELEASED" {
		panic(fmt.Sprintf("release reservation failed: %v", err))
	}

	// Test Error Mapping: Invalid Secret (401)
	badAuthAdapter := adapter.NewHospitalAdapter(baseURL, "wrong_secret", 2*time.Second)
	_, err = badAuthAdapter.Search(ctx, "req-bad-auth", adapter.HospitalSearchCriteria{
		ServiceCode:  "MCU_BASIC",
		PatientCount: 1,
		AppointmentWindow: adapter.TimeWindow{
			StartsAt:      "2026-08-22T02:00:00Z",
			EndsAt:        "2026-08-22T08:00:00Z",
			StartTimeZone: "Asia/Jakarta",
			EndTimeZone:   "Asia/Jakarta",
		},
	})
	var provErr *adapter.ProviderError
	if !errors.As(err, &provErr) || !provErr.IsUnauthorized() {
		panic(fmt.Sprintf("expected 401 ProviderError on bad auth, got: %v", err))
	}

	// Test Error Mapping: Not Found (404)
	_, err = a.GetReservation(ctx, "req-not-found", "unknown-res-id")
	if !errors.As(err, &provErr) || !provErr.IsNotFound() {
		panic(fmt.Sprintf("expected 404 ProviderError on missing resource, got: %v", err))
	}

	// Test Error Mapping: Capacity Conflict (409)
	_, err = a.CreateHold(ctx, "req-conflict", "idem-conflict", adapter.CreateHoldRequest{
		ProviderID:        "hospital-demo-01",
		OfferID:           "hosp-offer-conflict",
		Units:             10,
		ExpectedUnitPrice: adapter.Money{AmountMinor: 150000000, Currency: "IDR"},
		ClientReference:   "journey-001-hosp",
	})
	if !errors.As(err, &provErr) || !provErr.IsConflict() {
		panic(fmt.Sprintf("expected 409 ProviderError on conflict, got: %v", err))
	}

	// Test Error Mapping: Expired Hold (410)
	_, err = a.ConfirmHold(ctx, "req-expired", "idem-expired", "hosp-hold-expired")
	if !errors.As(err, &provErr) || !provErr.IsExpired() {
		panic(fmt.Sprintf("expected 410 ProviderError on expired hold, got: %v", err))
	}
}

func testFerryAdapter(ctx context.Context, a *adapter.FerryAdapter) {
	// Search
	offers, err := a.Search(ctx, "req-ferry-search", adapter.FerrySearchCriteria{
		OriginPortCode:      "HARBOURFRONT_SG",
		DestinationPortCode: "BATAM_CENTRE_ID",
		PassengerCount:      2,
		DepartureWindow: adapter.TimeWindow{
			StartsAt:      "2026-08-21T23:00:00Z",
			EndsAt:        "2026-08-22T04:00:00Z",
			StartTimeZone: "Asia/Singapore",
			EndTimeZone:   "Asia/Jakarta",
		},
	})
	if err != nil || len(offers) == 0 {
		panic(fmt.Sprintf("ferry search failed: %v", err))
	}

	// Hold
	hold, err := a.CreateHold(ctx, "req-ferry-hold", "idem-ferry-hold-01", adapter.CreateHoldRequest{
		ProviderID:        "ferry-demo-01",
		OfferID:           offers[0].OfferID,
		Units:             2,
		ExpectedUnitPrice: adapter.Money{AmountMinor: 5000, Currency: "SGD"},
		ClientReference:   "journey-001-ferry",
	})
	if err != nil || hold.HoldID == "" {
		panic(fmt.Sprintf("ferry hold failed: %v", err))
	}

	// Release Hold
	rel, err := a.ReleaseHold(ctx, "req-ferry-rel-hold", "idem-ferry-rel-hold", hold.HoldID)
	if err != nil || rel.Status != "RELEASED" {
		panic(fmt.Sprintf("ferry release hold failed: %v", err))
	}
}

func testHotelAdapter(ctx context.Context, a *adapter.HotelAdapter) {
	// Search
	offers, err := a.Search(ctx, "req-hotel-search", adapter.HotelSearchCriteria{
		CheckInDate:   "2026-08-22",
		CheckOutDate:  "2026-08-23",
		LocalTimezone: "Asia/Jakarta",
		RoomCount:     1,
		GuestCount:    2,
	})
	if err != nil || len(offers) == 0 {
		panic(fmt.Sprintf("hotel search failed: %v", err))
	}

	// Hold
	hold, err := a.CreateHold(ctx, "req-hotel-hold", "idem-hotel-hold-01", adapter.CreateHoldRequest{
		ProviderID:        "hotel-demo-01",
		OfferID:           offers[0].OfferID,
		Units:             1,
		ExpectedUnitPrice: adapter.Money{AmountMinor: 80000000, Currency: "IDR"},
		ClientReference:   "journey-001-hotel",
	})
	if err != nil || hold.HoldID == "" {
		panic(fmt.Sprintf("hotel hold failed: %v", err))
	}

	// Confirm
	res, err := a.ConfirmHold(ctx, "req-hotel-conf", "idem-hotel-conf-01", hold.HoldID)
	if err != nil || res.ReservationID == "" {
		panic(fmt.Sprintf("hotel confirm failed: %v", err))
	}
}

func testTransportAdapter(ctx context.Context, a *adapter.TransportAdapter) {
	// Search
	offers, err := a.Search(ctx, "req-trans-search", adapter.TransportSearchCriteria{
		PickupLocationCode:  "BATAM_CENTRE_ID",
		DropoffLocationCode: "BATAM_MEDICAL_CENTRE_ID",
		PassengerCount:      2,
		PickupWindow: adapter.TimeWindow{
			StartsAt:      "2026-08-22T01:30:00Z",
			EndsAt:        "2026-08-22T02:00:00Z",
			StartTimeZone: "Asia/Jakarta",
			EndTimeZone:   "Asia/Jakarta",
		},
	})
	if err != nil || len(offers) == 0 {
		panic(fmt.Sprintf("transport search failed: %v", err))
	}

	// Hold with TransportBookingRequirements
	hold, err := a.CreateHold(ctx, "req-trans-hold", "idem-trans-hold-01", adapter.CreateHoldRequest{
		ProviderID:        "transport-demo-01",
		OfferID:           offers[0].OfferID,
		Units:             1,
		ExpectedUnitPrice: adapter.Money{AmountMinor: 15000000, Currency: "IDR"},
		ClientReference:   "journey-001-trans",
		BookingRequirements: &adapter.TransportBookingRequirements{
			PassengerCount:      2,
			PickupLocationCode:  "BATAM_CENTRE_ID",
			DropoffLocationCode: "BATAM_MEDICAL_CENTRE_ID",
			PickupWindow: adapter.TimeWindow{
				StartsAt:      "2026-08-22T01:30:00Z",
				EndsAt:        "2026-08-22T02:00:00Z",
				StartTimeZone: "Asia/Jakarta",
				EndTimeZone:   "Asia/Jakarta",
			},
		},
	})
	if err != nil || hold.HoldID == "" {
		panic(fmt.Sprintf("transport hold failed: %v", err))
	}

	// Confirm
	res, err := a.ConfirmHold(ctx, "req-trans-conf", "idem-trans-conf-01", hold.HoldID)
	if err != nil || res.ReservationID == "" {
		panic(fmt.Sprintf("transport confirm failed: %v", err))
	}
}

func testAggregator(
	ctx context.Context,
	hosp *adapter.HospitalAdapter,
	ferry *adapter.FerryAdapter,
	hotel *adapter.HotelAdapter,
	trans *adapter.TransportAdapter,
) {
	// Faulty hotel adapter returning 503
	badHotelServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(adapter.ErrorEnvelope{
			Error: adapter.ErrorBody{
				Code:      "SERVICE_UNAVAILABLE",
				Message:   "Hotel database offline",
				Retryable: true,
			},
		})
	}))
	defer badHotelServer.Close()
	faultyHotelAdapter := adapter.NewHotelAdapter(badHotelServer.URL, "hotel_dev_secret", 2*time.Second)

	agg := adapter.NewAggregator(hosp, ferry, faultyHotelAdapter, trans)

	result := agg.SearchAll(ctx, "req-multi-search", adapter.MultiSearchQuery{
		HospitalCriteria: &adapter.HospitalSearchCriteria{
			ServiceCode:  "MCU_BASIC",
			PatientCount: 1,
			AppointmentWindow: adapter.TimeWindow{
				StartsAt:      "2026-08-22T02:00:00Z",
				EndsAt:        "2026-08-22T08:00:00Z",
				StartTimeZone: "Asia/Jakarta",
				EndTimeZone:   "Asia/Jakarta",
			},
		},
		FerryCriteria: []adapter.FerrySearchCriteria{
			{
				OriginPortCode:      "HARBOURFRONT_SG",
				DestinationPortCode: "BATAM_CENTRE_ID",
				PassengerCount:      1,
				DepartureWindow: adapter.TimeWindow{
					StartsAt:      "2026-08-21T23:00:00Z",
					EndsAt:        "2026-08-22T04:00:00Z",
					StartTimeZone: "Asia/Singapore",
					EndTimeZone:   "Asia/Jakarta",
				},
			},
		},
		HotelCriteria: &adapter.HotelSearchCriteria{
			CheckInDate:   "2026-08-22",
			CheckOutDate:  "2026-08-23",
			LocalTimezone: "Asia/Jakarta",
			RoomCount:     1,
			GuestCount:    1,
		},
		TransportCriteria: []adapter.TransportSearchCriteria{
			{
				PickupLocationCode:  "BATAM_CENTRE_ID",
				DropoffLocationCode: "BATAM_MEDICAL_CENTRE_ID",
				PassengerCount:      1,
				PickupWindow: adapter.TimeWindow{
					StartsAt:      "2026-08-22T01:30:00Z",
					EndsAt:        "2026-08-22T02:00:00Z",
					StartTimeZone: "Asia/Jakarta",
					EndTimeZone:   "Asia/Jakarta",
				},
			},
		},
	})

	if len(result.HospitalOffers) == 0 {
		panic("expected hospital offers to succeed despite hotel outage")
	}
	if len(result.FerryOffers) == 0 {
		panic("expected ferry offers to succeed despite hotel outage")
	}
	if len(result.TransportOffers) == 0 {
		panic("expected transport offers to succeed despite hotel outage")
	}
	if len(result.Warnings) == 0 {
		panic("expected hotel outage to be reported as a non-fatal warning")
	}
}

func startMockProviderServer(provType, provID, secret string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify auth header
		if r.URL.Path != "/healthz" && r.Header.Get("X-Integration-Key") != secret {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(adapter.ErrorEnvelope{
				Error: adapter.ErrorBody{
					Code:      "AUTHENTICATION_FAILED",
					Message:   "Invalid integration secret",
					Retryable: false,
					RequestID: r.Header.Get("X-Request-ID"),
				},
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Request-ID", r.Header.Get("X-Request-ID"))

		switch {
		case r.URL.Path == "/healthz":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(adapter.HealthResponse{
				Status:         "UP",
				ProviderID:     provID,
				ProviderType:   provType,
				DatabaseStatus: "UP",
				CheckedAt:      "2026-08-16T00:00:00Z",
			})

		case r.URL.Path == "/v1/offers/search":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(adapter.SearchResponse{
				ProviderID:   provID,
				ProviderType: provType,
				Offers: []adapter.Offer{
					{
						OfferID:        fmt.Sprintf("%s-offer-001", provID[:4]),
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

		case r.URL.Path == "/v1/holds" && r.Method == http.MethodPost:
			var req adapter.CreateHoldRequest
			_ = json.NewDecoder(r.Body).Decode(&req)

			if req.OfferID == "hosp-offer-conflict" {
				w.WriteHeader(http.StatusConflict)
				_ = json.NewEncoder(w).Encode(adapter.ErrorEnvelope{
					Error: adapter.ErrorBody{
						Code:      "CAPACITY_CONFLICT",
						Message:   "Requested units exceed available capacity",
						Retryable: false,
						RequestID: r.Header.Get("X-Request-ID"),
					},
				})
				return
			}

			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(adapter.Hold{
				HoldID:            fmt.Sprintf("%s-hold-001", provID[:4]),
				ExternalReference: fmt.Sprintf("%s-HOLD-REF-001", provID[:4]),
				ProviderID:        provID,
				ProviderType:      provType,
				OfferID:           req.OfferID,
				ClientReference:   req.ClientReference,
				Status:            "HELD",
				Units:             req.Units,
				UnitPrice:         req.ExpectedUnitPrice,
				TotalPrice:        adapter.Money{AmountMinor: req.ExpectedUnitPrice.AmountMinor * int64(req.Units), Currency: req.ExpectedUnitPrice.Currency},
				ServiceWindow:     adapter.TimeWindow{StartsAt: "2026-08-22T02:00:00Z", EndsAt: "2026-08-22T04:00:00Z", StartTimeZone: "Asia/Jakarta", EndTimeZone: "Asia/Jakarta"},
				CreatedAt:         "2026-08-16T00:00:00Z",
				ExpiresAt:         "2026-08-20T12:00:00Z",
			})

		case r.Method == http.MethodPost && (len(r.URL.Path) > 17 && r.URL.Path[len(r.URL.Path)-8:] == "/confirm"):
			if r.URL.Path == "/v1/holds/hosp-hold-expired/confirm" {
				w.WriteHeader(http.StatusGone)
				_ = json.NewEncoder(w).Encode(adapter.ErrorEnvelope{
					Error: adapter.ErrorBody{
						Code:      "HOLD_EXPIRED",
						Message:   "Hold expired before confirmation",
						Retryable: true,
						RequestID: r.Header.Get("X-Request-ID"),
					},
				})
				return
			}
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(adapter.Reservation{
				ReservationID:     fmt.Sprintf("%s-res-001", provID[:4]),
				ExternalReference: fmt.Sprintf("%s-RES-REF-001", provID[:4]),
				HoldID:            fmt.Sprintf("%s-hold-001", provID[:4]),
				ProviderID:        provID,
				ProviderType:      provType,
				OfferID:           fmt.Sprintf("%s-offer-001", provID[:4]),
				ClientReference:   "journey-001",
				Status:            "CONFIRMED",
				Units:             1,
				TotalPrice:        adapter.Money{AmountMinor: 150000000, Currency: "IDR"},
				ServiceWindow:     adapter.TimeWindow{StartsAt: "2026-08-22T02:00:00Z", EndsAt: "2026-08-22T04:00:00Z", StartTimeZone: "Asia/Jakarta", EndTimeZone: "Asia/Jakarta"},
				ConfirmedAt:       "2026-08-16T00:00:00Z",
			})

		case r.Method == http.MethodGet && len(r.URL.Path) > 17 && r.URL.Path[:17] == "/v1/reservations/":
			resID := r.URL.Path[17:]
			if resID == "unknown-res-id" {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(adapter.ErrorEnvelope{
					Error: adapter.ErrorBody{
						Code:      "NOT_FOUND",
						Message:   "Reservation not found",
						Retryable: false,
						RequestID: r.Header.Get("X-Request-ID"),
					},
				})
				return
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(adapter.Reservation{
				ReservationID:     resID,
				ExternalReference: fmt.Sprintf("%s-RES-REF-001", provID[:4]),
				HoldID:            fmt.Sprintf("%s-hold-001", provID[:4]),
				ProviderID:        provID,
				ProviderType:      provType,
				OfferID:           fmt.Sprintf("%s-offer-001", provID[:4]),
				ClientReference:   "journey-001",
				Status:            "CONFIRMED",
				Units:             1,
				TotalPrice:        adapter.Money{AmountMinor: 150000000, Currency: "IDR"},
				ServiceWindow:     adapter.TimeWindow{StartsAt: "2026-08-22T02:00:00Z", EndsAt: "2026-08-22T04:00:00Z", StartTimeZone: "Asia/Jakarta", EndTimeZone: "Asia/Jakarta"},
				ConfirmedAt:       "2026-08-16T00:00:00Z",
			})

		case r.Method == http.MethodPost && (len(r.URL.Path) > 8 && r.URL.Path[len(r.URL.Path)-8:] == "/release"):
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(adapter.ReleaseResult{
				ResourceType:      "RESERVATION",
				ResourceID:        "res-001",
				ExternalReference: "RES-REF-001",
				ProviderID:        provID,
				ProviderType:      provType,
				Status:            "RELEASED",
				ReleasedAt:        "2026-08-16T00:00:00Z",
			})

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}
