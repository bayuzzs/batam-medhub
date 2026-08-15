# Backend Codex Worker Brief

## Copy-paste worker prompt

```text
You are the Batam MedHub backend implementation worker. Read PROJECT_UNDERSTANDING.md, docs/architecture/, specs/openapi.yaml, specs/provider-openapi.yaml, and this task completely before editing. You own only backend/**. Treat specs/**, providers/**, mobile/**, and docs/architecture/** as read-only. Implement the smallest robust feature slice in the task order below using Go, Gin, GORM, PostgreSQL, and golang-migrate SQL migrations. Never use GORM AutoMigrate. The core must communicate with providers only through HTTP. If the contract blocks implementation, stop the affected operation and report a precise contract-change request to the control-plane conversation; do not edit the specifications yourself. Preserve unrelated user changes. Automated tests are deferred, but every completed phase must build, migrate from an empty database, and leave a repeatable manual smoke command.
```

## Ownership

Write only inside `backend/**`. Do not modify mobile, providers, shared specifications, architecture documents, or root task briefs.

## Technical constraints

- Go with Gin, GORM, PostgreSQL, and `golang-migrate/migrate` SQL files.
- SQL migrations are the only schema authority; disable runtime `AutoMigrate`.
- Store UUIDs or another contract-compatible opaque identifier consistently.
- Store money as integer minor units using the currency's ISO 4217 exponent plus its currency code.
- Store instants in UTC and retain source IANA time zones where local schedule meaning matters.
- Validate all Workers AI output as untrusted input.
- Use provider-specific HTTP adapters; never query provider databases.
- Use synchronous orchestration for the hackathon; do not add a message broker or distributed-saga framework.
- Preserve immutable itinerary history.

## Task order

### B1 — Runtime foundation

- Create the Gin API entrypoint and bounded internal packages.
- Load validated environment configuration.
- Add request ID, recovery, structured error, and logging middleware.
- Add liveness and database-aware readiness operations matching OpenAPI.
- Connect GORM to the core PostgreSQL database.

Done when the service builds, starts with configuration, and reports health without business routes.

### B2 — Core schema and deterministic reference data

Create ordered SQL migrations for:

- providers and provider-key hashes;
- medical services and provider capability mappings;
- static FX rates;
- trip requests and extracted-intent snapshots;
- plan options and normalized offer snapshots;
- journeys, reservations, itinerary versions, and itinerary items;
- provider events, disruptions, and recovery options;
- idempotency records where required by the contract.

Add deterministic synthetic seeds for MCU Basic, MCU Comprehensive, Dental Check-up, Eye Screening, the provider registry, and static FX rates.

Done when migrations apply cleanly to an empty database and GORM models match them without AutoMigrate.

### B3 — Trust-boundary primitives

- Implement patient JWT verification and read the preferred-currency claim.
- Implement provider-secret verification using stored hashes.
- Implement request validation, common errors, and idempotency behavior.
- Implement source-money preservation and static reference-currency conversion.
- Implement catalog lookup and unsupported-service handling.

### B4 — Structured-intent and trip-request flow

- Persist a trip request before external model work.
- Implement strict structured-intent decoding and backend validation.
- Support `MATCHED`, `NEEDS_CLARIFICATION`, `UNSUPPORTED_SERVICE`, and `OUT_OF_SCOPE`.
- Permit patient corrections through the contracted operation.
- Initially exercise planning with golden structured-intent fixtures internally; do not expose an undocumented temporary public endpoint.

### B5 — Provider adapters

- Implement one typed HTTP adapter per provider category.
- Apply provider base URL, integration secret, request ID, idempotency key, and timeout.
- Query independent providers concurrently where useful.
- Normalize provider-specific responses into core offer snapshots while preserving each opaque provider `offer_id` exactly as returned.
- Preserve partial successes and explicit provider errors.

### B6 — Deterministic planner

- Anchor plans on a supported hospital appointment.
- Enforce ferry cutoff, arrival, immigration, transfer, medical buffer, non-overlap, capacity, accessibility, stay, and return constraints.
- Filter impossible combinations before scoring.
- Return at most two explainable options.
- Store original provider money and reference-currency display values.

### B7 — Confirmation saga and itinerary v1

- Require explicit patient approval and an idempotency key.
- Revalidate relevant offers.
- Hold all required provider resources using the exact selected provider `offer_id`; never derive or substitute an offer ID.
- Send contracted `booking_requirements` for transport holds, including passenger count, route, pickup window, and accessibility needs.
- Release earlier holds when a later hold or confirmation fails.
- Set the trip request to `MANUAL_REVIEW` and create no journey when any compensating release has an uncertain outcome; use `CONFIRMATION_FAILED` only when cleanup outcomes are known.
- Confirm all successful holds and store provider references.
- Create an `ACTIVE` journey with immutable itinerary version 1.
- Implement both active-itinerary lookup and full immutable version lookup at `GET /v1/journeys/{journey_id}/itineraries/{version}` so superseded contents remain queryable.

### B8 — Cloudflare Workers AI

- Call Workers AI directly from the backend using secret configuration.
- Supply only the active small catalog and extraction instructions needed by the contract.
- Retry malformed structured output once, then return a safe deterministic failure.
- Never allow the model to select treatment, diagnose, invent provider facts, plan constraints, or book resources.

### B9 — Disruption and recovery

- Authenticate the provider organization from its secret.
- Accept a provider-asserted actor snapshot without claiming individual authentication.
- Validate that the event type is compatible with the authenticated provider type and that every target reservation/item belongs to both that provider and the stated journey; return `403` on mismatch.
- Deduplicate by authenticated provider plus `external_event_id`: replay only when the request-body fingerprint is identical, and return `409` for a changed body.
- Require `instruction_reference` and all contracted hospital-authored facts for additional care.
- Persist the canonical event and compute itinerary impact.
- Produce at most two recovery options with signed `time_delta_minutes` and price deltas.
- Treat `FOLLOWUP_OBSERVATION` as provider-only recovery inventory, not a user-requestable core service; search it and use its exact provider offer ID through the normal hold/confirm flow.
- Confirm replacement resources before releasing superseded reservations.
- Create itinerary version 2 and retain version 1.
- Fall back to manual review when no safe recovery exists.

### B10 — Demo hardening

- Implement the contracted deterministic demo reset behavior.
- Add an English backend README and environment example inside `backend/`.
- Document exact migration, startup, health, happy-path, and disruption smoke commands.

## Explicitly out of scope

- Provider or patient UI.
- Provider database access.
- Payments, refunds, insurance, medical records, multilingual UI, and post-care.
- Vector RAG, autonomous booking agents, message brokers, Kubernetes, and production compliance claims.
- Broad automated test suites until the feature-complete phase is approved.
