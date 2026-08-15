# Provider Codex Worker Brief

## Copy-paste worker prompt

```text
You are the Batam MedHub provider-platform implementation worker. Read PROJECT_UNDERSTANDING.md, docs/architecture/, specs/provider-openapi.yaml, its provider examples, and this task completely before editing. You own only providers/**. Treat specs/**, backend/**, mobile/**, and docs/architecture/** as read-only. Implement four standalone headless Go/Gin provider applications—hospital, ferry, hotel, and internal transport—started by providers/docker-compose.yml. They share one PostgreSQL server but use four separate logical databases, credentials, migrations, and domain data. Use GORM for data access and golang-migrate SQL migrations; never use GORM AutoMigrate. Implement the provider OpenAPI contract exactly. Do not implement disruption callbacks or a provider UI: provider actors submit disruptions directly to the core backend. If the contract blocks implementation, report a precise contract-change request to the control-plane conversation instead of editing the specification. Preserve unrelated user changes. Automated tests are deferred, but each completed phase must build, migrate from empty databases, and leave a repeatable manual smoke command.
```

## Ownership

Write only inside `providers/**`. Do not modify backend, mobile, shared specifications, architecture documents, or root task briefs.

## Required runtime topology

`providers/docker-compose.yml` runs:

- hospital Go service;
- ferry Go service;
- hotel Go service;
- internal transport Go service; and
- one PostgreSQL service with one named volume.

The PostgreSQL initialization creates:

| Database | Owner credential |
|---|---|
| `hospital_db` | `hospital_user` |
| `ferry_db` | `ferry_user` |
| `hotel_db` | `hotel_user` |
| `transport_db` | `transport_user` |

Revoke public cross-database connection privileges and grant each application access only to its own database. No provider may access the core database or another provider database.

## Shared behavior allowed

Share only small infrastructure that is genuinely identical:

- configuration loading;
- Gin request ID, recovery, auth, and error middleware;
- integration-secret comparison;
- idempotency mechanics;
- health response shape; and
- common hold-state constants where the contract is identical.

Across all four providers, a hold must reference the exact opaque `offer_id` returned by search; services must not derive a different identifier from schedule or inventory fields.

For every hold, compare the body `provider_id` and `provider_type` with the identity configured for the authenticated receiving service. Return `403` on mismatch before selecting provider-specific behavior; a caller-supplied type must never bypass transport `booking_requirements` or another provider rule.

Provider schemas, search filters, inventory rules, repositories, migrations, and seed data remain provider-specific.

## Task order

### P1 — Compose and database bootstrap

- Create one provider Go module and four `cmd` entrypoints.
- Create the common Dockerfile and Compose topology.
- Create four databases, four owners, and least-privilege grants through an initialization script.
- Add health checks and startup dependency conditions.

Done when a fresh named volume creates all four isolated databases predictably.

### P2 — Independent migrations and seeds

Create separate `golang-migrate` histories for:

- hospital services (MCU Basic, MCU Comprehensive, Dental Check-up, and Eye Screening), provider-only `FOLLOWUP_OBSERVATION` recovery inventory, appointment slots, holds, reservations, and idempotency;
- ferry sailings, seat capacity, holds, reservations, and idempotency;
- hotel rooms or room types, per-night inventory, holds, reservations, and idempotency;
- transport vehicles, availability or assignments, holds, reservations, and idempotency.

Use deterministic English synthetic seed data. `FOLLOWUP_OBSERVATION` must be searchable, holdable, and confirmable by the backend for recovery but must not be presented as a patient-requestable core catalog service. Runtime GORM `AutoMigrate` must remain disabled.

The hospital golden seed must include these `MCU_BASIC` slots on 22 August 2026 (WIB): an intentionally infeasible 08:00 offer at IDR 1,200,000, the canonical selected 10:00 offer at IDR 1,500,000, and a later 13:00 comparison offer at IDR 2,000,000. Use the exact offer IDs and ISO-minor-unit prices from `specs/examples/provider/search-hospital-response.json`; the planner, not the provider, decides which complete cross-provider combinations are feasible.

### P3 — Hospital reference service

Implement the complete contracted flow before copying the pattern:

- health;
- search supported services and slots;
- transactional hold with expiry using the exact opaque `offer_id` returned by search;
- idempotent hold retry;
- confirm;
- reservation lookup; and
- release.

Use a database transaction and row-level locking or an equivalent PostgreSQL constraint so concurrent holds cannot oversell capacity. Lazy expiry checks using `expires_at` are sufficient for the hackathon.

### P4 — Ferry service

Implement the same reservation lifecycle with ferry-specific search criteria, check-in cutoff, departure/arrival time zones, and seat capacity.

### P5 — Hotel service

Implement the same lifecycle with date-range coverage and per-night room inventory. A hold must cover every required night atomically.

### P6 — Internal transport service

Implement the same lifecycle with pickup/drop-off windows, passenger capacity, and vehicle or driver availability. A transport hold must validate the contracted `booking_requirements`, including passenger count, pickup/drop-off route, pickup window, and accessibility needs, against the selected offer.

### P7 — Contract and isolation smoke flows

For every provider, document commands proving:

- search succeeds with valid credentials;
- hold succeeds from the exact searched `offer_id` and returns expiry;
- retrying the same idempotency key returns the same result;
- confirm and lookup succeed;
- release is idempotent;
- invalid secret returns `401`;
- capacity, price, or state conflicts return `409`;
- an expired offer or hold returns `410`; and
- its credential cannot connect to the other three databases.

### P8 — Deterministic reset and documentation

- Provide one reset/reseed procedure covering all four logical databases.
- Add an English README and `.env.example` under `providers/`.
- Document Compose startup, migrations, seeds, health, and smoke flows.

## Explicitly out of scope

- Provider dashboard UI or provider-user authentication.
- Disruption callbacks to core; actors call the core disruption route directly.
- Direct access to core journey, intent, FX, or patient data.
- Message brokers, service discovery, Kubernetes, and separate physical PostgreSQL clusters.
- Real provider integrations and production data.
- Broad automated test suites until the feature-complete phase is approved.
