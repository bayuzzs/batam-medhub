# Batam MedHub — Core Backend Orchestrator

The Batam MedHub Core Backend is a high-reliability medical journey orchestration engine written in Go (using Gin, GORM, and PostgreSQL). It coordinates cross-border patient trips from Singapore to medical facilities in Batam, handling authentication, AI-driven unstructured intent extraction, multi-provider availability aggregation, distributed transactional booking sagas with automatic compensation, and real-time disruption ingestion with versioned itinerary recovery.

---

## 1. Architecture Overview & Bounded Contexts

The backend operates as the journey orchestrator across four independent mock provider domains (`hospital`, `ferry`, `hotel`, and `transport`). Services interact exclusively via HTTP APIs and never share databases.

```
                  +-----------------------------------+
                  |        Mobile Patient App         |
                  +-----------------+-----------------+
                                    | Bearer JWT
                                    v
+-------------------------------------------------------------------------+
|                       Core Backend Orchestrator                         |
|                                                                         |
|  +-------------------+  +-------------------+  +---------------------+  |
|  | Auth & Profile    |  | Intent Extractor  |  | Plan Generation     |  |
|  | (Sessions & JWT)  |  | (Workers AI / LLM)|  | (Multi-Provider Agg)|  |
|  +-------------------+  +-------------------+  +---------------------+  |
|                                                                         |
|  +-------------------+  +-------------------+  +---------------------+  |
|  | Booking Saga      |  | Journey & Version |  | Disruption Recovery |  |
|  | (Holds & Confirms)|  | (Immutable V1->V2)|  | (Event Ingest / Opt)|  |
|  +-------------------+  +-------------------+  +---------------------+  |
+-----------------------------------+-------------------------------------+
                                    | Provider Adapter HTTP Calls
           +------------------------+------------------------+
           |                        |                        |
           v                        v                        v
+--------------------+   +--------------------+   +--------------------+
|  Hospital Provider |   |   Ferry Provider   |   | Transport Provider |
|     (MCU/Cardio)   |   |   (Sg <-> Batam)   |   |  (Pickup/Dropoff)  |
+--------------------+   +--------------------+   +--------------------+
```

### Key Bounded Contexts:
- **Auth & Session Management**: Argon2id password hashing, rotating refresh tokens, HS256 access JWTs scoped to patient identities and preferred display currency.
- **Intent Extraction (`internal/ai`)**: Ingests unstructured patient requests, validates candidate medical service codes against reference catalogs, applies clinical emergency guardrails, and falls back to deterministic rule extraction on network/AI timeout.
- **Plan Generation & Aggregator (`internal/adapter`)**: Aggregates real-time availability across hospital appointment slots, outbound/return ferry departures, and private ground transfers; computes transit buffers and display currency conversions.
- **Booking Saga Orchestrator (`internal/service/saga.go`)**: Coordinates two-phase holds and sequential confirmations across providers. On any provider failure, automatically triggers compensation releases for all held/confirmed items.
- **Disruption & Recovery Engine (`internal/service/disruption.go`)**: Ingests provider-authenticated event webhooks, deduplicates requests using SHA-256 fingerprints, generates multi-provider recovery options, activates immutable Itinerary Version 2, and releases superseded reservations.
- **Demo Management (`internal/service/demo.go`)**: Exposes contracted `POST /v1/demo/reset` for cleaning dynamic orchestration state while preserving reference data.

---

## 2. Environment Variables Reference

| Variable | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `HTTP_ADDR` | string | `:8080` | Server listen address and port. |
| `DATABASE_URL` | string | *required* | PostgreSQL connection URI for `core_db`. |
| `JWT_SIGNING_SECRET` | string | *required* | HS256 secret key for patient access tokens (min 32 chars). |
| `JWT_ISSUER` | string | `batam-medhub` | Expected JWT issuer claim. |
| `JWT_AUDIENCE` | string | `batam-medhub-mobile` | Expected JWT audience claim. |
| `HOSPITAL_BASE_URL` | string | `http://localhost:8081` | Base URL for mock hospital service. |
| `HOSPITAL_INTEGRATION_KEY` | string | `hospital_dev_secret` | Shared secret key for hospital provider calls. |
| `FERRY_BASE_URL` | string | `http://localhost:8082` | Base URL for mock ferry service. |
| `FERRY_INTEGRATION_KEY` | string | `ferry_dev_secret` | Shared secret key for ferry provider calls. |
| `HOTEL_BASE_URL` | string | `http://localhost:8083` | Base URL for mock hotel service. |
| `HOTEL_INTEGRATION_KEY` | string | `hotel_dev_secret` | Shared secret key for hotel provider calls. |
| `TRANSPORT_BASE_URL` | string | `http://localhost:8084` | Base URL for mock transport service. |
| `TRANSPORT_INTEGRATION_KEY`| string | `transport_dev_secret`| Shared secret key for transport provider calls. |
| `PROVIDER_HTTP_TIMEOUT_SECONDS`| int | `5` | Provider HTTP client timeout in seconds. |
| `CLOUDFLARE_ACCOUNT_ID` | string | `""` | Optional Cloudflare account identifier. |
| `CLOUDFLARE_API_TOKEN` | string | `""` | Optional Cloudflare Workers AI API token. |
| `CLOUDFLARE_AI_MODEL` | string | `@cf/meta/llama-3.1-8b-instruct` | Cloudflare Workers AI model tag. |
| `CLOUDFLARE_AI_BASE_URL` | string | `https://api.cloudflare.com/client/v4` | Cloudflare API base URL. |
| `CLOUDFLARE_AI_TIMEOUT_SECONDS` | int | `15` | Timeout for Workers AI inference. |
| `DEMO_SECRET` | string | `demo_dev_secret` | Secret key for authenticating `/v1/demo/reset`. |

---

## 3. Database Bootstrap and Startup Instructions

### Prerequisites
- Go 1.22+
- Docker & Docker Compose (for running PostgreSQL and mock providers)

### 1. Start Provider Infrastructure and PostgreSQL
From repository root:
```bash
docker compose -f ./providers/docker-compose.yml up -d
```

### 2. Apply Database Migrations
Run the migration CLI to apply schema and reference data:
```bash
cd backend
export DATABASE_URL="postgres://provider_admin:provider_admin_dev_password@localhost:5432/core_db?sslmode=disable"
go run ./cmd/migrate up
```

### 3. Run Backend API Server
```bash
cd backend
export DATABASE_URL="postgres://provider_admin:provider_admin_dev_password@localhost:5432/core_db?sslmode=disable"
export JWT_SIGNING_SECRET="12345678901234567890123456789012"
export DEMO_SECRET="demo_dev_secret"
go run ./cmd/api
```

### 4. Run Automated Verification Suite
```bash
cd backend
go run ./cmd/verify
```

---

## 4. End-to-End Smoke Test Commands

The following `curl` flow demonstrates the full journey lifecycle:

### Step 1: Demo Reset (Optional cleanup)
```bash
curl -s -X POST http://localhost:8080/v1/demo/reset \
  -H "X-Demo-Key: demo_dev_secret" \
  -H "Content-Type: application/json" \
  -d '{
    "scenario": "DEFAULT",
    "confirm": true
  }'
```

### Step 2: Register Synthetic Patient
```bash
curl -s -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Eleanor Vance",
    "email": "eleanor.vance@example.com",
    "password": "Password123!",
    "preferred_currency": "SGD",
    "nationality": "SG"
  }'
```
*Extract `token` from the response for subsequent requests:*
```bash
TOKEN="<access_token_here>"
```

### Step 3: Create Trip Request (Natural Language Intent)
```bash
curl -s -X POST http://localhost:8080/v1/trip-requests \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "I need a comprehensive health screening in Batam on 22 August 2026 for 1 person.",
    "preferred_currency": "SGD"
  }'
```
*Extract `trip_request_id` (e.g. `fc007984-9bb4-4b82-9648-bc953d7c500b`):*
```bash
TRIP_ID="<trip_request_id_here>"
```

### Step 4: Generate Ranked Plan Options
```bash
curl -s -X POST "http://localhost:8080/v1/trip-requests/$TRIP_ID/plans" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{}'
```
*Extract `plan_option_id` (e.g. `dfac0dd9-df15-44cc-88c7-d633b98c1faa`):*
```bash
PLAN_OPTION_ID="<plan_option_id_here>"
```

### Step 5: Execute Multi-Provider Booking Saga
```bash
curl -s -X POST "http://localhost:8080/v1/plan-options/$PLAN_OPTION_ID/confirm" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Idempotency-Key: idem-booking-smoke-001" \
  -H "Content-Type: application/json" \
  -d '{
    "approved": true
  }'
```
*Returns active Journey with Itinerary Version 1 containing 6 items.*

*Extract `journey_id` (e.g. `1afe4661-3b84-4f21-95ca-69a6a9c56f68`) and hospital itinerary item ID:*
```bash
JOURNEY_ID="<journey_id_here>"
HOSPITAL_ITEM_ID="<hospital_itinerary_item_id_here>"
```

### Step 6: Ingest Hospital Disruption Event
```bash
curl -s -X POST http://localhost:8080/v1/provider/disruptions \
  -H "X-Provider-Key: hospital_dev_secret" \
  -H "Content-Type: application/json" \
  -d '{
    "external_event_id": "evt-hosp-smoke-001",
    "journey_id": "'"$JOURNEY_ID"'",
    "event_type": "HOSPITAL_ADDITIONAL_CARE_REQUESTED",
    "occurred_at": "2026-08-22T04:30:00Z",
    "target": {
      "itinerary_item_id": "'"$HOSPITAL_ITEM_ID"'"
    },
    "actor": {
      "actor_id": "dr-lee-tan",
      "name": "Dr Lee Tan",
      "role": "Cardiologist"
    },
    "details": {
      "reason": "Patient requires additional cardiac observation following exam.",
      "instruction_reference": "hospital-instruction://followup-observation/FO-20260822-0001",
      "replacement_time_window": {
        "starts_at": "2026-08-22T05:00:00Z",
        "ends_at": "2026-08-22T06:30:00Z",
        "start_time_zone": "Asia/Jakarta",
        "end_time_zone": "Asia/Jakarta"
      },
      "additional_service_code": "FOLLOWUP_OBSERVATION",
      "additional_duration_minutes": 90,
      "priority": "MEDIUM",
      "travel_clearance_status": "CLEARED"
    }
  }'
```
*Extract `disruption_id` from receipt:*
```bash
DISRUPTION_ID="<disruption_id_here>"
```

### Step 7: Retrieve Disruption and Recovery Options
```bash
curl -s -X GET "http://localhost:8080/v1/disruptions/$DISRUPTION_ID" \
  -H "Authorization: Bearer $TOKEN"
```
*Extract `recovery_option_id` (e.g. `8d0896b3-c633-405f-a701-65d8a72aae95`):*
```bash
RECOVERY_OPTION_ID="<recovery_option_id_here>"
```

### Step 8: Approve Recovery Option & Activate Itinerary Version 2
```bash
curl -s -X POST "http://localhost:8080/v1/recovery-options/$RECOVERY_OPTION_ID/approve" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Idempotency-Key: idem-recovery-smoke-001" \
  -H "Content-Type: application/json" \
  -d '{
    "approved": true
  }'
```
*Returns updated Journey with `active_itinerary_version: 2` (8 items), while Version 1 remains preserved as `SUPERSEDED`.*
