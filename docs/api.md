# Batam MedHub — API Documentation

Batam MedHub defines **two** OpenAPI 3.1 contracts:

1. **Core API** — the patient-facing backend API (auth, trip requests, plans,
   booking, itinerary, disruptions). → [`specs/openapi.yaml`](../specs/openapi.yaml)
2. **Provider API** — the protocol the backend uses to talk to the four mock
   providers (hospital, ferry, hotel, transport). → [`specs/provider-openapi.yaml`](../specs/provider-openapi.yaml)

Both contracts are the source of truth. Every contract change goes through the
control plane, and golden example payloads are validated against the schemas:

```bash
bash specs/validate.sh
```

---

## 1. Core API — endpoint overview

| Method | Path | Purpose |
| :--- | :--- | :--- |
| `GET` | `/healthz` | Liveness check. |
| `GET` | `/readyz` | Readiness check (DB reachable). |
| `POST` | `/v1/auth/register` | Create a patient account (email/password). |
| `POST` | `/v1/auth/login` | Login, returns access JWT + refresh token. |
| `POST` | `/v1/auth/refresh` | Rotate the refresh token (single-use). |
| `POST` | `/v1/auth/logout` | Idempotent session revocation. |
| `GET` / `PUT` | `/v1/profile` | Read / update the patient profile (currency preference). |
| `GET` | `/v1/medical-services` | The active, catalog-driven medical-service catalog. |
| `POST` | `/v1/trip-requests` | Submit a natural-language trip request (AI intent extraction). |
| `GET` | `/v1/trip-requests/{id}` | Fetch a trip request and its intent resolution. |
| `GET` | `/v1/trip-requests/{id}/intent` | Fetch the structured intent result. |
| `POST` | `/v1/trip-requests/{id}/plans` | Generate deterministic plan options (≤ 2). |
| `POST` | `/v1/plan-options/{id}/confirm` | Approve a plan → run the multi-provider booking saga. |
| `GET` | `/v1/journeys/{id}/itinerary` | Fetch the active itinerary. |
| `GET` | `/v1/journeys/{id}/itineraries/{version}` | Fetch a specific immutable itinerary version. |
| `POST` | `/v1/provider/disruptions` | Provider-authenticated disruption event ingestion. |
| `GET` | `/v1/disruptions/{id}` | Fetch a disruption and its recovery options. |
| `POST` | `/v1/recovery-options/{id}/approve` | Approve a recovery option → activate itinerary v2. |
| `POST` | `/v1/demo/reset` | Reset dynamic state to the golden demo (demo secret). |

### Auth model

- **Access tokens:** short-lived HS256 JWTs (patient identity + display
  currency).
- **Refresh tokens:** rotating opaque tokens; only SHA-256 hashes are stored;
  single-use with atomic rotation.
- Provider disruption ingestion is authenticated with the provider integration
  key (`X-Provider-Key`); demo reset uses `X-Demo-Key`.

### Example golden payloads (`specs/examples/core/`)

- `auth-register-request.json`, `auth-login-request.json`, `auth-session.json`
- `matched-trip-request.json`, `needs-clarification-trip-request.json`,
  `out-of-scope-trip-request.json`, `unsupported-trip-request.json`
- `medical-services.json`, `plan-result.json`
- `active-journey-v1.json`, `disruption-request.json`, `disruption-detail.json`,
  `recovery-approved-v2.json`

---

## 2. Provider API — endpoint overview

All four providers (hospital, ferry, hotel, transport) implement the same
lifecycle, distinguished by the resource they hold:

| Method | Path | Purpose |
| :--- | :--- | :--- |
| `GET` | `/healthz` | Liveness check. |
| `POST` | `/v1/offers/search` | Search available offers (slots/sailings/rooms/assignments). |
| `POST` | `/v1/holds` | Place a temporary hold on an offer. |
| `POST` | `/v1/holds/{id}/confirm` | Confirm a hold into a reservation. |
| `POST` | `/v1/holds/{id}/release` | Release a hold (compensation). |
| `GET` | `/v1/reservations/{id}` | Look up a confirmed reservation. |
| `POST` | `/v1/reservations/{id}/release` | Release a reservation. |

Backend ↔ provider calls are authenticated with the shared integration key
(`X-Integration-Key`).

### Example golden payloads (`specs/examples/provider/`)

- `search-hospital-request.json` / `search-hospital-response.json`
- `search-ferry-request.json` / `search-ferry-response.json`
- `search-hotel-request.json` / `search-hotel-response.json`
- `search-transport-request.json` / `search-transport-response.json`
- `create-hold-request.json` / `create-hold-response.json`
- `confirm-hold-response.json`, `release-hold-response.json`,
  `reservation-response.json`, `capacity-conflict-error.json`,
  `hold-expired-error.json`

---

## 3. Key flows to show judges

1. **Register → login → trip request** — the AI intent pipeline and the
   `matched` / `needs_clarification` / `unsupported` / `out_of_scope` branches.
2. **Plan generation** — deterministic constraint planning producing ≤ 2
   options with explainable rankings.
3. **Confirm** — the idempotent booking saga across all four providers.
4. **Disruption** — provider event ingestion, impact analysis, recovery options.
5. **Recovery approval** — itinerary v2 activation with v1 superseded.

See [demo-script.md](demo-script.md) for the full walking order and
`backend/README.md` for copy-paste `curl` commands.
