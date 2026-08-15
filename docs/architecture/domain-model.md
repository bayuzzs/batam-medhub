# Batam MedHub Domain Model

Status: contract draft v0.1
Scope: first hackathon vertical slice

## Purpose

This document defines the shared domain vocabulary and ownership boundaries used by the core API, provider protocol, database migrations, and implementation. It is intentionally smaller than the long-term product model.

The vertical slice turns a supported planned-care request into a feasible cross-provider journey, confirms it, then recovers that journey after a provider disruption.

## Bounded contexts and ownership

| Context | Owns | Does not own |
|---|---|---|
| Core orchestration | Supported-service catalog, trip intent, normalized offer snapshots, plan options, booking coordination, journeys, immutable itinerary versions, provider events, disruptions, recovery options, and static FX rates | Live provider inventory or a provider's operational booking state |
| Hospital provider | Medical-service listings, appointment slots, appointment holds, and confirmed appointments | The patient's cross-provider journey |
| Ferry provider | Sailings, check-in cutoffs, seat capacity, holds, and confirmed passenger reservations | Transfers, hotel stays, medical appointments, or journey feasibility |
| Hotel provider | Room types, nightly inventory, room holds, and confirmed stays | Medical or transport decisions |
| Internal transport provider | Vehicles, bookable availability windows, assignment holds, and confirmed assignments | The rest of the itinerary |
| Cloudflare Workers AI | Natural-language-to-structured-intent transformation | Catalog authority, availability, planning, booking, clinical decisions, or journey state |

The core backend and the providers communicate only through HTTP contracts. There are no cross-database reads, writes, or foreign keys.

## Deployment-level persistence boundary

- The core backend has its own PostgreSQL datastore.
- The four providers share one PostgreSQL server for deployment efficiency.
- That provider server contains four isolated logical databases: `hospital_db`, `ferry_db`, `hotel_db`, and `transport_db`.
- Each provider connects with a credential restricted to its database.
- Each database has independent `golang-migrate` history.
- GORM is used for queries and transactions; runtime `AutoMigrate` is not used.

Sharing one PostgreSQL process is an infrastructure choice, not a shared domain or transaction boundary.

## Core aggregates and entities

### Reference and integration data

| Entity | Meaning | Important identity or rule |
|---|---|---|
| `Provider` | A hospital, ferry, hotel, or internal transport integration known to the core. | Stable core ID and exactly one `provider_type`. |
| `ProviderCredential` | A rotatable credential used to authenticate a provider actor submitting a disruption. | Store only a verifier/hash and key metadata, never the plaintext secret. |
| `MedicalService` | An active catalog item that Batam MedHub is allowed to orchestrate. | Stable, unique service code such as `MCU_BASIC`. The catalog, not the LLM, decides support. |
| `ProviderCapability` | Maps a core medical-service code to a hospital's external service reference. | Unique per provider and external service reference. |
| `FxRate` | A static, timestamped reference exchange rate used for display conversion. | Unique by base currency, quote currency, and effective timestamp/source. |

### Planning aggregate

`TripRequest` is the aggregate root before confirmation.

| Entity | Meaning | Important identity or rule |
|---|---|---|
| `TripRequest` | A patient's planning request and validated structured intent. | Owned by the JWT subject. The first slice does not need a patient table. |
| `PlanOption` | One feasible, explainable combination returned for a trip request. | At most two active options per planning run. It has an expiry derived from its provider offers. |
| `PlanItem` | An immutable normalized snapshot of one provider offer or required non-bookable itinerary step within a plan. | Stores provider and external offer references, time window, source money, synthetic/source markers, and expiry. It is not live inventory. |

The raw user prompt is processed at the API/LLM boundary and is not required as a durable domain record. Persist the validated structured intent and only the minimum text needed for review, such as `requested_service_text`.

### Booking and journey aggregate

`Journey` is created only after all required provider confirmations succeed. It is the aggregate root for the confirmed cross-provider journey.

| Entity | Meaning | Important identity or rule |
|---|---|---|
| `Reservation` | The core's record of a provider hold or confirmed reservation created during the orchestration saga. | References the provider's external hold/reservation IDs. It never replaces the provider record. |
| `Journey` | The confirmed journey owned by Batam MedHub. | Links back to exactly one trip request and identifies exactly one active itinerary version. |
| `ItineraryVersion` | An immutable version of the complete journey. | Version number is unique within a journey. At most one version is `ACTIVE`. |
| `ItineraryItem` | An immutable snapshot of a medical appointment, sailing, stay, transfer, or operational buffer in an itinerary version. | Stores local time zone, UTC instants, source/display money snapshots, and optional reservation reference. |

Reservations may be recorded while a trip request is in `CONFIRMING`. Successful confirmation links them to the new journey and active itinerary in one core transaction. Provider operations cannot participate in that database transaction, so failures are handled through the booking saga and compensating releases.

### Disruption and recovery aggregate

`Disruption` coordinates impact analysis and recovery for one active journey.

| Entity | Meaning | Important identity or rule |
|---|---|---|
| `ProviderEvent` | A canonical fact submitted manually by an authenticated provider actor. | Unique by `(provider_id, external_event_id)` for deduplication, with a body fingerprint that must match on replay. Provider identity comes from the verified secret, not the body. |
| `Disruption` | A provider event that invalidates at least one active itinerary constraint or requires a safety hold. | References the active itinerary version that was analyzed. A provider event may produce no disruption. |
| `RecoveryOption` | One feasible set of logistical changes proposed for a disruption. | At most two active options per analysis run. It includes signed `time_delta_minutes` and price delta relative to the active itinerary and expires when any replacement offer expires. Positive time means a later finish; negative time means an earlier finish. |
| `RecoveryItem` | The added, changed, removed, or unchanged item delta shown to the patient. | References old itinerary items and/or replacement offer snapshots as applicable. |

The actor data in a provider event is an audit snapshot asserted by the provider (`actor_id`, `name`, and `role`). It is not an independently authenticated provider user in this slice.

`HOSPITAL_ADDITIONAL_CARE_REQUESTED` is accepted only when the hospital supplies the service code, instruction reference, explicit time window, duration, priority, travel-clearance status, and operational requirements. Batam MedHub and the LLM must not infer those clinical facts. `FOLLOWUP_OBSERVATION` is provider-only hospital recovery inventory for the demo, not an active patient-requestable core `MedicalService`; recovery must search, hold, and confirm the matching hospital offer through the normal provider protocol.

### Cross-cutting persistence

| Entity | Meaning | Important identity or rule |
|---|---|---|
| `IdempotencyRecord` | Stores the outcome of a mutating core request so retries return the original result. | Unique by authenticated scope, operation, and idempotency key. |

The stored request fingerprint must match on a retry. Reusing a key with a different request is a conflict.

## Provider-owned domain model

Each provider has different inventory but implements the same operational concepts:

| Concept | Definition |
|---|---|
| `Offer` | A time-bounded, priced representation of currently bookable provider inventory returned by search. |
| `Hold` | A temporary capacity allocation with a provider-generated ID and expiry. |
| `Reservation` | A provider-confirmed booking created from a valid hold. |
| `IdempotencyRecord` | The provider-local replay record for hold, confirm, and release requests. |

Provider-specific inventory is persisted as follows:

| Database | Inventory authority |
|---|---|
| `hospital_db` | Patient-requestable medical services, provider-only follow-up inventory, and appointment slots |
| `ferry_db` | Dated sailings, check-in cutoffs, and seat capacity |
| `hotel_db` | Room types and nightly room inventory |
| `transport_db` | Vehicles and bookable availability windows |

Search availability is derived inside the provider transaction from inventory minus non-expired holds and confirmed reservations. The core never calculates provider capacity from its snapshots.

## Canonical value objects

| Value object | Contract rule |
|---|---|
| `Money` | Integer `amount_minor` plus ISO 4217 `currency`; never floating point. |
| `ConvertedMoney` | Source money plus reference/display money, rate, rate source, and rate effective timestamp; explicitly marked estimated/reference. |
| `TimeWindow` | UTC `starts_at` and `ends_at`, plus `start_time_zone` and `end_time_zone` IANA names for local display. End must be after start. |
| `ExternalReference` | Provider ID plus opaque provider-owned external ID. Core must not parse meaning from it. |
| `OfferSnapshot` | Provider/external reference, item type, times, location, source price, expiry, source marker, and only the operational metadata needed by the planner. |
| `ActorSnapshot` | Provider-asserted actor ID, display name, and role stored with the event for audit context. |

## Business invariants

1. Only an active `MedicalService` can resolve an intent to `MATCHED` and enter planning.
2. LLM output is untrusted until schema, enum, date, range, and catalog validation pass.
3. Hard constraints are evaluated before scoring. An infeasible combination is never returned, even if it is cheaper or ranks higher.
4. A planning or recovery run returns no more than two feasible options.
5. A `PlanItem` or `ItineraryItem` is a snapshot; it must not be treated as current provider availability.
6. Every provider offer is revalidated before hold or confirmation.
7. All mutating core and provider operations are idempotent within an authenticated scope.
8. A journey becomes `ACTIVE` only after every required provider reservation is confirmed.
9. Confirmation failure releases all new holds that can be safely released. Confirmed cleanup sets the trip request to `CONFIRMATION_FAILED`; an uncertain compensation sets it to `MANUAL_REVIEW`, retains every external reference, and creates no active journey.
10. An itinerary version is immutable after activation. Recovery creates a new version; it never edits the old one.
11. A journey has at most one active itinerary version. Activating v2 and superseding v1 is one core database transaction after replacement confirmations succeed.
12. Replacement reservations are confirmed before superseded reservations are released.
13. A provider event is deduplicated by authenticated provider plus external event ID. Only a matching request-body fingerprint replays the original result; reuse with another body is a conflict.
14. A provider event is not automatically a disruption. If all active constraints still pass, record `NO_ACTION` and keep the itinerary unchanged.
15. A provider-reported clinical or travel hold stops automatic travel replanning until the provider clears it.
16. Provider identity and type are resolved from the verified secret. A body-supplied provider ID is ignored or rejected. The event type must be compatible with the authenticated provider type, and any target reservation or itinerary item must belong to both that provider and the stated journey.
17. Ferry, hotel, and transport providers receive only the minimum operational data; they do not receive raw medical prompts or medical records.
18. Store instants in UTC and retain IANA time zones. Store original provider money and any display conversion snapshot.
19. All demo provider, patient, inventory, price, and booking data is synthetic and visibly marked `synthetic: true` with source `MOCK`.
20. Active and superseded itinerary versions remain retrievable with their full immutable item contents.

## Transaction boundaries

- A provider hold, confirm, release, or capacity update is atomic only inside that provider's database.
- Creating a journey, linking core reservation records, activating its first itinerary, and completing the trip-request workflow is atomic inside the core database after external confirmations succeed.
- Activating a recovery itinerary and superseding the previous itinerary is atomic inside the core database after replacement confirmations succeed.
- HTTP calls across providers are coordinated as a saga. PostgreSQL transactions must never remain open while waiting on provider HTTP calls.

## Explicit exclusions

The first vertical slice has no medical records, diagnosis, treatment selection, emergency triage, payment, refund, insurance, live currency feed, provider dashboard, provider-user authentication, multilingual clinical translation, or real provider integration.
