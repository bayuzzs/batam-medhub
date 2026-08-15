#!/usr/bin/env bash
# Consolidated contract and isolation smoke verification suite for all four Batam MedHub provider services.
set -Eeuo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROVIDERS_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "================================================================="
echo "  Batam MedHub: Consolidated Provider Smoke Verification Suite   "
echo "================================================================="

# Optional reset before running smoke tests if --reset is passed
if [[ "${1:-}" == "--reset" ]]; then
  echo "--> Running deterministic database reset first..."
  "$SCRIPT_DIR/reset.sh"
fi

cd "$PROVIDERS_DIR"

HOSPITAL_URL="http://localhost:8081"
FERRY_URL="http://localhost:8082"
HOTEL_URL="http://localhost:8083"
TRANSPORT_URL="http://localhost:8084"

HOSPITAL_SECRET="hospital_dev_secret"
FERRY_SECRET="ferry_dev_secret"
HOTEL_SECRET="hotel_dev_secret"
TRANSPORT_SECRET="transport_dev_secret"

assert_equal() {
  local expected="$1"
  local actual="$2"
  local msg="$3"
  if [[ "$expected" != "$actual" ]]; then
    echo "FAILED: $msg (expected: '$expected', got: '$actual')" >&2
    exit 1
  fi
}

echo ""
echo "=== SECTION 1: Health & Request ID Tracking ==="
for svc in "HOSPITAL:$HOSPITAL_URL" "FERRY:$FERRY_URL" "HOTEL:$HOTEL_URL" "TRANSPORT:$TRANSPORT_URL"; do
  TYPE="${svc%%:*}"
  URL="${svc#*:}"
  echo "--> Checking $TYPE health ($URL/healthz)..."
  RESP=$(curl -sf "$URL/healthz")
  echo "$RESP" | grep -q '"status":"UP"' || { echo "FAILED: $TYPE status not UP"; exit 1; }
  echo "$RESP" | grep -q "\"provider_type\":\"$TYPE\"" || { echo "FAILED: $TYPE provider_type mismatch"; exit 1; }
  echo "$RESP" | grep -q '"database_status":"UP"' || { echo "FAILED: $TYPE database_status not UP"; exit 1; }

  # Test Request ID echo
  REQ_ID_RESP=$(curl -si -H "X-Request-ID: test-trace-req-001" "$URL/healthz")
  echo "$REQ_ID_RESP" | grep -qi "X-Request-Id: test-trace-req-001" || { echo "FAILED: $TYPE request ID not echoed"; exit 1; }

  # Test Request ID replacement on invalid format
  REQ_ID_BAD=$(curl -si -H "X-Request-ID: bad invalid req id with spaces" "$URL/healthz")
  echo "$REQ_ID_BAD" | grep -qi "X-Request-Id: req-" || { echo "FAILED: $TYPE invalid request ID not replaced"; exit 1; }
done
echo "PASSED: Health checks and Request ID tracking verified across all 4 providers."

echo ""
echo "=== SECTION 2: Authentication & Authorization (401 & 403) ==="
for svc in "HOSPITAL:$HOSPITAL_URL" "FERRY:$FERRY_URL" "HOTEL:$HOTEL_URL" "TRANSPORT:$TRANSPORT_URL"; do
  TYPE="${svc%%:*}"
  URL="${svc#*:}"
  echo "--> Testing 401 Unauthorized on $TYPE..."
  STATUS_401=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$URL/v1/offers/search" \
    -H "Content-Type: application/json" \
    -d "{\"provider_type\":\"$TYPE\"}")
  assert_equal "401" "$STATUS_401" "$TYPE missing auth header did not return 401"

  STATUS_BAD_KEY=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$URL/v1/offers/search" \
    -H "Content-Type: application/json" \
    -H "X-Integration-Key: wrong_secret_key" \
    -d "{\"provider_type\":\"$TYPE\"}")
  assert_equal "401" "$STATUS_BAD_KEY" "$TYPE invalid secret key did not return 401"
done
echo "PASSED: Authentication (401) verified across all 4 providers."

echo ""
echo "=== SECTION 3: Hospital Provider (:8081) Lifecycle & Contract Verification ==="
echo "--> 3.1 Search MCU_BASIC offers..."
HOSP_SEARCH=$(curl -sf -X POST "$HOSPITAL_URL/v1/offers/search" \
  -H "Content-Type: application/json" \
  -H "X-Integration-Key: $HOSPITAL_SECRET" \
  -d '{
    "provider_type": "HOSPITAL",
    "criteria": {
      "service_code": "MCU_BASIC",
      "patient_count": 1,
      "appointment_window": {
        "starts_at": "2026-08-22T01:00:00Z",
        "ends_at": "2026-08-22T08:00:00Z",
        "start_time_zone": "Asia/Jakarta",
        "end_time_zone": "Asia/Jakarta"
      },
      "accessibility": []
    }
  }')
echo "$HOSP_SEARCH" | grep -q "hospital-offer-mcu-basic-20260822-1000" || { echo "FAILED: Canonical 10:00 MCU_BASIC offer not found"; exit 1; }
echo "$HOSP_SEARCH" | grep -q '"amount_minor":150000000' || { echo "FAILED: 10:00 MCU_BASIC price mismatch"; exit 1; }
echo "Search OK"

echo "--> 3.2 Concurrency Idempotency (10 parallel identical hold requests)..."
rm -f /tmp/hosp_conc_*.json /tmp/hosp_conc_*.code
for i in {1..10}; do
  curl -s -w "%{http_code}" -o /tmp/hosp_conc_$i.json -X POST "$HOSPITAL_URL/v1/holds" \
    -H "Content-Type: application/json" \
    -H "X-Integration-Key: $HOSPITAL_SECRET" \
    -H "Idempotency-Key: idem-hospital-conc-smoke-01" \
    -d '{
      "provider_id": "hospital-demo-01",
      "provider_type": "HOSPITAL",
      "offer_id": "hospital-offer-mcu-basic-20260822-1000",
      "units": 1,
      "expected_unit_price": {
        "amount_minor": 150000000,
        "currency": "IDR"
      },
      "client_reference": "journey-hosp-conc-smoke-01"
    }' > /tmp/hosp_conc_$i.code &
done
wait

HOSP_HOLD_ID=""
for i in {1..10}; do
  CODE=$(cat /tmp/hosp_conc_$i.code)
  assert_equal "201" "$CODE" "Hospital concurrent hold request $i failed"
  HID=$(grep -o '"hold_id":"[^"]*' /tmp/hosp_conc_$i.json | cut -d'"' -f4)
  if [[ -z "$HOSP_HOLD_ID" ]]; then HOSP_HOLD_ID="$HID"; fi
  assert_equal "$HOSP_HOLD_ID" "$HID" "Hospital concurrent hold ID mismatch"
done
echo "Concurrent hold OK: $HOSP_HOLD_ID"

echo "--> 3.3 Idempotency Conflict (409) & Identity Mismatch (403)..."
HOSP_IDEMP_409=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$HOSPITAL_URL/v1/holds" \
  -H "Content-Type: application/json" \
  -H "X-Integration-Key: $HOSPITAL_SECRET" \
  -H "Idempotency-Key: idem-hospital-conc-smoke-01" \
  -d '{
    "provider_id": "hospital-demo-01",
    "provider_type": "HOSPITAL",
    "offer_id": "hospital-offer-mcu-basic-20260822-1300",
    "units": 1,
    "expected_unit_price": {"amount_minor": 200000000, "currency": "IDR"},
    "client_reference": "journey-diff"
  }')
assert_equal "409" "$HOSP_IDEMP_409" "Hospital idempotency conflict expected 409"

HOSP_403=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$HOSPITAL_URL/v1/holds" \
  -H "Content-Type: application/json" \
  -H "X-Integration-Key: $HOSPITAL_SECRET" \
  -H "Idempotency-Key: idem-hospital-bad-id" \
  -d '{
    "provider_id": "wrong-provider-id",
    "provider_type": "HOSPITAL",
    "offer_id": "hospital-offer-mcu-basic-20260822-1000",
    "units": 1,
    "expected_unit_price": {"amount_minor": 150000000, "currency": "IDR"},
    "client_reference": "journey-bad-id"
  }')
assert_equal "403" "$HOSP_403" "Hospital identity mismatch expected 403"

echo "--> 3.4 Confirm Hold, Idempotent Retry & Lookup..."
HOSP_CONFIRM=$(curl -si -X POST "$HOSPITAL_URL/v1/holds/$HOSP_HOLD_ID/confirm" \
  -H "Content-Type: application/json" \
  -H "X-Integration-Key: $HOSPITAL_SECRET" \
  -H "Idempotency-Key: idem-hospital-confirm-smoke-01")
echo "$HOSP_CONFIRM" | grep -q "HTTP/1.1 201" || { echo "FAILED: Hospital confirm status not 201"; exit 1; }
echo "$HOSP_CONFIRM" | grep -qi "Location: /v1/reservations/" || { echo "FAILED: Hospital confirm Location header missing"; exit 1; }
HOSP_RES_ID=$(echo "$HOSP_CONFIRM" | grep -o '"reservation_id":"[^"]*' | cut -d'"' -f4)

HOSP_RETRY_CONFIRM=$(curl -si -X POST "$HOSPITAL_URL/v1/holds/$HOSP_HOLD_ID/confirm" \
  -H "Content-Type: application/json" \
  -H "X-Integration-Key: $HOSPITAL_SECRET" \
  -H "Idempotency-Key: idem-hospital-confirm-smoke-01")
echo "$HOSP_RETRY_CONFIRM" | grep -qi "Idempotency-Replayed: true" || { echo "FAILED: Hospital confirm retry replay header missing"; exit 1; }

HOSP_LOOKUP=$(curl -sf "$HOSPITAL_URL/v1/reservations/$HOSP_RES_ID" -H "X-Integration-Key: $HOSPITAL_SECRET")
echo "$HOSP_LOOKUP" | grep -q "\"hold_id\":\"$HOSP_HOLD_ID\"" || { echo "FAILED: Hospital reservation lookup hold_id mismatch"; exit 1; }
echo "$HOSP_LOOKUP" | grep -q '"status":"CONFIRMED"' || { echo "FAILED: Hospital reservation status not CONFIRMED"; exit 1; }

echo "--> 3.5 Release Reservation & Direct Hold Release..."
HOSP_REL_RES=$(curl -si -X POST "$HOSPITAL_URL/v1/reservations/$HOSP_RES_ID/release" \
  -H "Content-Type: application/json" \
  -H "X-Integration-Key: $HOSPITAL_SECRET" \
  -H "Idempotency-Key: idem-hosp-rel-res-smoke")
echo "$HOSP_REL_RES" | grep -q "HTTP/1.1 200" || { echo "FAILED: Hospital release reservation status not 200"; exit 1; }
echo "$HOSP_REL_RES" | grep -q '"status":"RELEASED"' || { echo "FAILED: Hospital release reservation status not RELEASED"; exit 1; }

HOSP_HOLD2=$(curl -sf -X POST "$HOSPITAL_URL/v1/holds" \
  -H "Content-Type: application/json" \
  -H "X-Integration-Key: $HOSPITAL_SECRET" \
  -H "Idempotency-Key: idem-hosp-hold-smoke-02" \
  -d '{
    "provider_id": "hospital-demo-01",
    "provider_type": "HOSPITAL",
    "offer_id": "hospital-offer-mcu-basic-20260822-1300",
    "units": 1,
    "expected_unit_price": {"amount_minor": 200000000, "currency": "IDR"},
    "client_reference": "journey-hosp-smoke-02"
  }')
HOSP_HOLD2_ID=$(echo "$HOSP_HOLD2" | grep -o '"hold_id":"[^"]*' | cut -d'"' -f4)

HOSP_REL_HOLD=$(curl -si -X POST "$HOSPITAL_URL/v1/holds/$HOSP_HOLD2_ID/release" \
  -H "Content-Type: application/json" \
  -H "X-Integration-Key: $HOSPITAL_SECRET" \
  -H "Idempotency-Key: idem-hosp-rel-hold-smoke")
echo "$HOSP_REL_HOLD" | grep -q "HTTP/1.1 200" || { echo "FAILED: Hospital release hold status not 200"; exit 1; }

# State conflict 409
HOSP_STATE_409=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$HOSPITAL_URL/v1/holds/$HOSP_HOLD2_ID/confirm" \
  -H "Content-Type: application/json" \
  -H "X-Integration-Key: $HOSPITAL_SECRET" \
  -H "Idempotency-Key: idem-hosp-confirm-released")
assert_equal "409" "$HOSP_STATE_409" "Hospital confirming released hold expected 409"
echo "PASSED: Hospital provider lifecycle completely verified."

echo ""
echo "=== SECTION 4: Ferry Provider (:8082) Lifecycle & Contract Verification ==="
echo "--> 4.1 Search Outbound & Return Sailings..."
FERRY_SEARCH_OUT=$(curl -sf -X POST "$FERRY_URL/v1/offers/search" \
  -H "Content-Type: application/json" \
  -H "X-Integration-Key: $FERRY_SECRET" \
  -d '{
    "provider_type": "FERRY",
    "criteria": {
      "origin_port_code": "HARBOURFRONT_SG",
      "destination_port_code": "BATAM_CENTRE_ID",
      "passenger_count": 2,
      "departure_window": {
        "starts_at": "2026-08-21T21:30:00Z",
        "ends_at": "2026-08-22T00:00:00Z",
        "start_time_zone": "Asia/Singapore",
        "end_time_zone": "Asia/Singapore"
      }
    }
  }')
echo "$FERRY_SEARCH_OUT" | grep -q "ferry-offer-hf-btm-20260822-0730" || { echo "FAILED: Outbound ferry offer not found"; exit 1; }
echo "$FERRY_SEARCH_OUT" | grep -q '"amount_minor":2500' || { echo "FAILED: Ferry unit price mismatch"; exit 1; }

echo "--> 4.2 Concurrency Idempotency (10 parallel identical hold requests)..."
rm -f /tmp/ferry_conc_*.json /tmp/ferry_conc_*.code
for i in {1..10}; do
  curl -s -w "%{http_code}" -o /tmp/ferry_conc_$i.json -X POST "$FERRY_URL/v1/holds" \
    -H "Content-Type: application/json" \
    -H "X-Integration-Key: $FERRY_SECRET" \
    -H "Idempotency-Key: idem-ferry-conc-smoke-01" \
    -d '{
      "provider_id": "ferry-demo-01",
      "provider_type": "FERRY",
      "offer_id": "ferry-offer-hf-btm-20260822-0730",
      "units": 2,
      "expected_unit_price": {"amount_minor": 2500, "currency": "SGD"},
      "client_reference": "journey-ferry-conc-smoke-01"
    }' > /tmp/ferry_conc_$i.code &
done
wait

FERRY_HOLD_ID=""
for i in {1..10}; do
  CODE=$(cat /tmp/ferry_conc_$i.code)
  assert_equal "201" "$CODE" "Ferry concurrent hold request $i failed"
  HID=$(grep -o '"hold_id":"[^"]*' /tmp/ferry_conc_$i.json | cut -d'"' -f4)
  if [[ -z "$FERRY_HOLD_ID" ]]; then FERRY_HOLD_ID="$HID"; fi
  assert_equal "$FERRY_HOLD_ID" "$HID" "Ferry concurrent hold ID mismatch"
done
echo "Concurrent hold OK: $FERRY_HOLD_ID"

echo "--> 4.3 Confirm Hold, Lookup, and Release..."
FERRY_CONFIRM=$(curl -si -X POST "$FERRY_URL/v1/holds/$FERRY_HOLD_ID/confirm" \
  -H "Content-Type: application/json" \
  -H "X-Integration-Key: $FERRY_SECRET" \
  -H "Idempotency-Key: idem-ferry-confirm-smoke-01")
echo "$FERRY_CONFIRM" | grep -q "HTTP/1.1 201" || { echo "FAILED: Ferry confirm status not 201"; exit 1; }
FERRY_RES_ID=$(echo "$FERRY_CONFIRM" | grep -o '"reservation_id":"[^"]*' | cut -d'"' -f4)

FERRY_LOOKUP=$(curl -sf "$FERRY_URL/v1/reservations/$FERRY_RES_ID" -H "X-Integration-Key: $FERRY_SECRET")
echo "$FERRY_LOOKUP" | grep -q "\"hold_id\":\"$FERRY_HOLD_ID\"" || { echo "FAILED: Ferry lookup hold_id mismatch"; exit 1; }

FERRY_REL_RES=$(curl -si -X POST "$FERRY_URL/v1/reservations/$FERRY_RES_ID/release" \
  -H "Content-Type: application/json" \
  -H "X-Integration-Key: $FERRY_SECRET" \
  -H "Idempotency-Key: idem-ferry-rel-res-smoke")
echo "$FERRY_REL_RES" | grep -q "HTTP/1.1 200" || { echo "FAILED: Ferry release reservation status not 200"; exit 1; }
echo "PASSED: Ferry provider lifecycle completely verified."

echo ""
echo "=== SECTION 5: Hotel Provider (:8083) Lifecycle & Contract Verification ==="
echo "--> 5.1 Search Hotel Offers..."
HOTEL_SEARCH=$(curl -sf -X POST "$HOTEL_URL/v1/offers/search" \
  -H "Content-Type: application/json" \
  -H "X-Integration-Key: $HOTEL_SECRET" \
  -d '{
    "provider_type": "HOTEL",
    "criteria": {
      "check_in_date": "2026-08-22",
      "check_out_date": "2026-08-23",
      "local_timezone": "Asia/Jakarta",
      "room_count": 1,
      "guest_count": 2,
      "accessibility": []
    }
  }')
echo "$HOTEL_SEARCH" | grep -q "hotel-offer-batam-centre-20260822-1n" || { echo "FAILED: Hotel offer not found"; exit 1; }
echo "$HOTEL_SEARCH" | grep -q '"amount_minor":85000000' || { echo "FAILED: Hotel unit price mismatch"; exit 1; }

echo "--> 5.2 Concurrency Idempotency (10 parallel identical hold requests)..."
rm -f /tmp/hotel_conc_*.json /tmp/hotel_conc_*.code
for i in {1..10}; do
  curl -s -w "%{http_code}" -o /tmp/hotel_conc_$i.json -X POST "$HOTEL_URL/v1/holds" \
    -H "Content-Type: application/json" \
    -H "X-Integration-Key: $HOTEL_SECRET" \
    -H "Idempotency-Key: idem-hotel-conc-smoke-01" \
    -d '{
      "provider_id": "hotel-demo-01",
      "provider_type": "HOTEL",
      "offer_id": "hotel-offer-batam-centre-20260822-1n",
      "units": 1,
      "expected_unit_price": {"amount_minor": 85000000, "currency": "IDR"},
      "client_reference": "journey-hotel-conc-smoke-01"
    }' > /tmp/hotel_conc_$i.code &
done
wait

HOTEL_HOLD_ID=""
for i in {1..10}; do
  CODE=$(cat /tmp/hotel_conc_$i.code)
  assert_equal "201" "$CODE" "Hotel concurrent hold request $i failed"
  HID=$(grep -o '"hold_id":"[^"]*' /tmp/hotel_conc_$i.json | cut -d'"' -f4)
  if [[ -z "$HOTEL_HOLD_ID" ]]; then HOTEL_HOLD_ID="$HID"; fi
  assert_equal "$HOTEL_HOLD_ID" "$HID" "Hotel concurrent hold ID mismatch"
done
echo "Concurrent hold OK: $HOTEL_HOLD_ID"

echo "--> 5.3 Confirm Hold, Lookup, and Release..."
HOTEL_CONFIRM=$(curl -si -X POST "$HOTEL_URL/v1/holds/$HOTEL_HOLD_ID/confirm" \
  -H "Content-Type: application/json" \
  -H "X-Integration-Key: $HOTEL_SECRET" \
  -H "Idempotency-Key: idem-hotel-confirm-smoke-01")
echo "$HOTEL_CONFIRM" | grep -q "HTTP/1.1 201" || { echo "FAILED: Hotel confirm status not 201"; exit 1; }
HOTEL_RES_ID=$(echo "$HOTEL_CONFIRM" | grep -o '"reservation_id":"[^"]*' | cut -d'"' -f4)

HOTEL_LOOKUP=$(curl -sf "$HOTEL_URL/v1/reservations/$HOTEL_RES_ID" -H "X-Integration-Key: $HOTEL_SECRET")
echo "$HOTEL_LOOKUP" | grep -q "\"hold_id\":\"$HOTEL_HOLD_ID\"" || { echo "FAILED: Hotel lookup hold_id mismatch"; exit 1; }

HOTEL_REL_RES=$(curl -si -X POST "$HOTEL_URL/v1/reservations/$HOTEL_RES_ID/release" \
  -H "Content-Type: application/json" \
  -H "X-Integration-Key: $HOTEL_SECRET" \
  -H "Idempotency-Key: idem-hotel-rel-res-smoke")
echo "$HOTEL_REL_RES" | grep -q "HTTP/1.1 200" || { echo "FAILED: Hotel release reservation status not 200"; exit 1; }
echo "PASSED: Hotel provider lifecycle completely verified."

echo ""
echo "=== SECTION 6: Transport Provider (:8084) Lifecycle & Contract Verification ==="
echo "--> 6.1 Search Transport Offers..."
TRANSPORT_SEARCH=$(curl -sf -X POST "$TRANSPORT_URL/v1/offers/search" \
  -H "Content-Type: application/json" \
  -H "X-Integration-Key: $TRANSPORT_SECRET" \
  -d '{
    "provider_type": "TRANSPORT",
    "criteria": {
      "pickup_location_code": "BATAM_CENTRE_ID",
      "dropoff_location_code": "HOSPITAL_DEMO_ID",
      "passenger_count": 2,
      "pickup_window": {
        "starts_at": "2026-08-22T01:20:00Z",
        "ends_at": "2026-08-22T01:40:00Z",
        "start_time_zone": "Asia/Jakarta",
        "end_time_zone": "Asia/Jakarta"
      },
      "accessibility": []
    }
  }')
echo "$TRANSPORT_SEARCH" | grep -q "transport-offer-btm-hospital-20260822-0825" || { echo "FAILED: Transport offer not found"; exit 1; }
echo "$TRANSPORT_SEARCH" | grep -q '"amount_minor":15000000' || { echo "FAILED: Transport unit price mismatch"; exit 1; }

echo "--> 6.2 Concurrency Idempotency with Required Booking Requirements (10 parallel requests)..."
rm -f /tmp/transport_conc_*.json /tmp/transport_conc_*.code
for i in {1..10}; do
  curl -s -w "%{http_code}" -o /tmp/transport_conc_$i.json -X POST "$TRANSPORT_URL/v1/holds" \
    -H "Content-Type: application/json" \
    -H "X-Integration-Key: $TRANSPORT_SECRET" \
    -H "Idempotency-Key: idem-transport-conc-smoke-01" \
    -d '{
      "provider_id": "transport-demo-01",
      "provider_type": "TRANSPORT",
      "offer_id": "transport-offer-btm-hospital-20260822-0825",
      "units": 1,
      "expected_unit_price": {"amount_minor": 15000000, "currency": "IDR"},
      "client_reference": "journey-transport-conc-smoke-01",
      "booking_requirements": {
        "passenger_count": 2,
        "pickup_location_code": "BATAM_CENTRE_ID",
        "dropoff_location_code": "HOSPITAL_DEMO_ID",
        "pickup_window": {
          "starts_at": "2026-08-22T01:25:00Z",
          "ends_at": "2026-08-22T02:10:00Z",
          "start_time_zone": "Asia/Jakarta",
          "end_time_zone": "Asia/Jakarta"
        },
        "accessibility": []
      }
    }' > /tmp/transport_conc_$i.code &
done
wait

TRANSPORT_HOLD_ID=""
for i in {1..10}; do
  CODE=$(cat /tmp/transport_conc_$i.code)
  assert_equal "201" "$CODE" "Transport concurrent hold request $i failed"
  HID=$(grep -o '"hold_id":"[^"]*' /tmp/transport_conc_$i.json | cut -d'"' -f4)
  if [[ -z "$TRANSPORT_HOLD_ID" ]]; then TRANSPORT_HOLD_ID="$HID"; fi
  assert_equal "$TRANSPORT_HOLD_ID" "$HID" "Transport concurrent hold ID mismatch"
done
echo "Concurrent hold OK: $TRANSPORT_HOLD_ID"

echo "--> 6.3 Confirm Hold, Lookup, and Release..."
TRANSPORT_CONFIRM=$(curl -si -X POST "$TRANSPORT_URL/v1/holds/$TRANSPORT_HOLD_ID/confirm" \
  -H "Content-Type: application/json" \
  -H "X-Integration-Key: $TRANSPORT_SECRET" \
  -H "Idempotency-Key: idem-transport-confirm-smoke-01")
echo "$TRANSPORT_CONFIRM" | grep -q "HTTP/1.1 201" || { echo "FAILED: Transport confirm status not 201"; exit 1; }
TRANSPORT_RES_ID=$(echo "$TRANSPORT_CONFIRM" | grep -o '"reservation_id":"[^"]*' | cut -d'"' -f4)

TRANSPORT_LOOKUP=$(curl -sf "$TRANSPORT_URL/v1/reservations/$TRANSPORT_RES_ID" -H "X-Integration-Key: $TRANSPORT_SECRET")
echo "$TRANSPORT_LOOKUP" | grep -q "\"hold_id\":\"$TRANSPORT_HOLD_ID\"" || { echo "FAILED: Transport lookup hold_id mismatch"; exit 1; }

TRANSPORT_REL_RES=$(curl -si -X POST "$TRANSPORT_URL/v1/reservations/$TRANSPORT_RES_ID/release" \
  -H "Content-Type: application/json" \
  -H "X-Integration-Key: $TRANSPORT_SECRET" \
  -H "Idempotency-Key: idem-transport-rel-res-smoke")
echo "$TRANSPORT_REL_RES" | grep -q "HTTP/1.1 200" || { echo "FAILED: Transport release reservation status not 200"; exit 1; }
echo "PASSED: Transport provider lifecycle completely verified."

echo ""
echo "=== SECTION 7: Database Isolation & Least-Privilege Verification ==="

verify_isolation() {
  local user="$1"
  local target_db="$2"
  local output
  output=$(docker compose exec -T postgres psql -U "$user" -d "$target_db" -c "SELECT 1;" 2>&1 || true)
  if ! echo "$output" | grep -q "permission denied"; then
    echo "FAILED: $user unexpectedly connected to $target_db! Output: $output" >&2
    exit 1
  fi
}

echo "--> 7.1 Testing hospital_user isolation..."
verify_isolation hospital_user ferry_db
verify_isolation hospital_user hotel_db
verify_isolation hospital_user transport_db

echo "--> 7.2 Testing ferry_user isolation..."
verify_isolation ferry_user hospital_db
verify_isolation ferry_user hotel_db
verify_isolation ferry_user transport_db

echo "--> 7.3 Testing hotel_user isolation..."
verify_isolation hotel_user hospital_db
verify_isolation hotel_user ferry_db
verify_isolation hotel_user transport_db

echo "--> 7.4 Testing transport_user isolation..."
verify_isolation transport_user hospital_db
verify_isolation transport_user ferry_db
verify_isolation transport_user hotel_db
echo "PASSED: All 4 provider database users are strictly isolated with zero cross-database connection privileges."

echo ""
echo "================================================================="
echo "  SUCCESS: All 4 Provider Services Passed Smoke Verification!    "
echo "================================================================="
