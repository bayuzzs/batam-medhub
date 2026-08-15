# Batam MedHub Logical ERD

Status: contract draft v0.1
Scope: first hackathon vertical slice

## Reading this document

These are logical models used to align OpenAPI and migrations. They intentionally show only identity, ownership, important state, and money/time fields. SQL migrations remain the physical-schema authority for exact PostgreSQL types, indexes, checks, and deletion behavior.

All internal IDs are UUIDs. Provider-owned external IDs remain opaque strings. All timestamps represent UTC instants unless a separate IANA `time_zone` is retained for local display.

## Core PostgreSQL

The core database is the Batam MedHub system of record for planning, orchestration, journeys, and disruption recovery. It stores provider snapshots and external references, not live provider inventory.

### Reference data and planning

```mermaid
erDiagram
    direction LR

    PROVIDER ||--o{ PROVIDER_CREDENTIAL : authenticates_with
    PROVIDER ||--o{ PROVIDER_CAPABILITY : publishes
    MEDICAL_SERVICE ||--o{ PROVIDER_CAPABILITY : maps_to
    MEDICAL_SERVICE o|..o{ TRIP_REQUEST : resolves
    TRIP_REQUEST ||--o{ PLAN_OPTION : produces
    PLAN_OPTION ||--|{ PLAN_ITEM : contains
    PROVIDER o|..o{ PLAN_ITEM : supplies
    FX_RATE o|..o{ PLAN_ITEM : converts

    PROVIDER {
        uuid id PK
        string provider_type
        string code UK
        string display_name
        string status
        bool synthetic
    }

    PROVIDER_CREDENTIAL {
        uuid id PK
        uuid provider_id FK
        string key_prefix UK
        string secret_hash
        string status
        datetime expires_at
    }

    MEDICAL_SERVICE {
        uuid id PK
        string code UK
        string name
        int default_duration_minutes
        bool active
        bool synthetic
    }

    PROVIDER_CAPABILITY {
        uuid id PK
        uuid provider_id FK
        uuid medical_service_id FK
        string external_service_id
        bool active
    }

    TRIP_REQUEST {
        uuid id PK
        string patient_subject
        string status
        uuid medical_service_id FK
        string requested_service_text
        json structured_intent
        string reference_currency
        int planning_revision
        datetime created_at
    }

    PLAN_OPTION {
        uuid id PK
        uuid trip_request_id FK
        int planning_revision
        int rank
        string status
        json explanation
        datetime expires_at
        datetime created_at
    }

    PLAN_ITEM {
        uuid id PK
        uuid plan_option_id FK
        uuid provider_id FK
        string item_type
        string external_offer_id
        datetime starts_at
        datetime ends_at
        string start_time_zone
        string end_time_zone
        int source_amount_minor
        string source_currency
        int display_amount_minor
        string display_currency
        uuid fx_rate_id FK
        json offer_snapshot
        datetime offer_expires_at
    }

    FX_RATE {
        uuid id PK
        string base_currency
        string quote_currency
        decimal rate
        string source
        datetime effective_at
    }

    IDEMPOTENCY_RECORD {
        uuid id PK
        string auth_scope
        string operation
        string idempotency_key
        string request_fingerprint
        int response_status
        json response_body
        datetime expires_at
    }
```

Important constraints not expressible as Mermaid cardinalities:

- `(provider_id, medical_service_id, external_service_id)` is unique for an active capability mapping.
- Only a provider with type `HOSPITAL` may own a medical-service capability.
- `(trip_request_id, planning_revision, rank)` is unique, and one planning revision exposes at most two options.
- A `PlanItem.provider_id` is nullable only for non-bookable operational steps such as immigration or safety buffers.
- A `PlanItem` copies conversion metadata into its snapshot; `fx_rate_id` preserves provenance but does not make a historical price dependent on a mutable rate.
- `(auth_scope, operation, idempotency_key)` is unique. The same key with another fingerprint is rejected.

### Booking, journey, and itinerary versioning

```mermaid
erDiagram
    direction LR

    TRIP_REQUEST ||--o| JOURNEY : becomes
    TRIP_REQUEST ||--o{ RESERVATION : coordinates
    PLAN_ITEM ||--o| RESERVATION : books
    PROVIDER ||--o{ RESERVATION : owns_externally
    JOURNEY o|..o{ RESERVATION : includes
    JOURNEY ||--|{ ITINERARY_VERSION : versions
    ITINERARY_VERSION ||--|{ ITINERARY_ITEM : contains
    RESERVATION o|..o{ ITINERARY_ITEM : represented_by

    TRIP_REQUEST {
        uuid id PK
        string patient_subject
        string status
        uuid selected_plan_option_id FK
    }

    PLAN_ITEM {
        uuid id PK
        uuid plan_option_id FK
        uuid provider_id FK
        string item_type
        string external_offer_id
    }

    PROVIDER {
        uuid id PK
        string provider_type
        string code UK
    }

    RESERVATION {
        uuid id PK
        uuid trip_request_id FK
        uuid journey_id FK
        uuid plan_item_id FK
        uuid provider_id FK
        string status
        string external_hold_id
        string external_reservation_id
        datetime hold_expires_at
        json provider_snapshot
        datetime updated_at
    }

    JOURNEY {
        uuid id PK
        uuid trip_request_id FK, UK
        string patient_subject
        string status
        int current_version_number
        datetime activated_at
    }

    ITINERARY_VERSION {
        uuid id PK
        uuid journey_id FK
        int version_number
        string status
        string change_reason
        uuid source_disruption_id FK
        datetime activated_at
    }

    ITINERARY_ITEM {
        uuid id PK
        uuid itinerary_version_id FK
        uuid reservation_id FK
        string item_type
        int sequence_number
        datetime starts_at
        datetime ends_at
        string start_time_zone
        string end_time_zone
        int source_amount_minor
        string source_currency
        int display_amount_minor
        string display_currency
        json snapshot
    }
```

Important constraints:

- `TripRequest` has zero or one `Journey`; a journey is created only after all required reservations confirm.
- `Reservation.journey_id` is nullable while confirmation is in progress.
- `(journey_id, version_number)` and `(itinerary_version_id, sequence_number)` are unique.
- A partial unique constraint permits only one `ACTIVE` itinerary version per journey.
- Activated itinerary versions and their items are append-only. Recovery creates a successor version.
- A reservation can appear in multiple versions when an unchanged booking is carried forward.
- Bookable itinerary items reference a core reservation; operational buffer items do not.

### Provider events, disruption, and recovery

```mermaid
erDiagram
    direction LR

    PROVIDER ||--o{ PROVIDER_EVENT : submits
    PROVIDER_EVENT ||--o| DISRUPTION : may_create
    JOURNEY ||--o{ DISRUPTION : experiences
    ITINERARY_VERSION ||--o{ DISRUPTION : analyzed_against
    DISRUPTION ||--o{ RECOVERY_OPTION : proposes
    RECOVERY_OPTION ||--|{ RECOVERY_ITEM : contains
    ITINERARY_ITEM o|..o{ RECOVERY_ITEM : changes

    PROVIDER {
        uuid id PK
        string provider_type
        string code UK
    }

    PROVIDER_EVENT {
        uuid id PK
        uuid provider_id FK
        string external_event_id
        string request_fingerprint
        string event_type
        datetime occurred_at
        json actor_snapshot
        json event_payload
        string assessment_outcome
        datetime received_at
    }

    JOURNEY {
        uuid id PK
        string status
        int current_version_number
    }

    ITINERARY_VERSION {
        uuid id PK
        uuid journey_id FK
        int version_number
        string status
    }

    ITINERARY_ITEM {
        uuid id PK
        uuid itinerary_version_id FK
        string item_type
        string external_reference
    }

    DISRUPTION {
        uuid id PK
        uuid provider_event_id FK, UK
        uuid journey_id FK
        uuid analyzed_itinerary_version_id FK
        string status
        json impact_summary
        datetime detected_at
        datetime resolved_at
    }

    RECOVERY_OPTION {
        uuid id PK
        uuid disruption_id FK
        int analysis_revision
        int rank
        string status
        json explanation
        json price_delta
        int time_delta_minutes
        datetime expires_at
    }

    RECOVERY_ITEM {
        uuid id PK
        uuid recovery_option_id FK
        uuid old_itinerary_item_id FK
        string change_type
        json replacement_offer_snapshot
        json item_delta
    }
```

Important constraints:

- `(provider_id, external_event_id)` is unique. A replay returns the existing result only when `request_fingerprint` matches; a different body is a conflict and does not create another disruption.
- `ProviderEvent.event_type` must be compatible with the authenticated provider's type. Every referenced reservation or itinerary item must belong to both that provider and the event's journey.
- One provider event creates zero or one disruption. A valid event with no itinerary impact records `NO_ACTION` on the event assessment and creates no disruption resource.
- `(disruption_id, analysis_revision, rank)` is unique, and one revision exposes at most two recovery options.
- `RecoveryItem.change_type` is `ADDED`, `CHANGED`, `REMOVED`, or `UNCHANGED`.
- `old_itinerary_item_id` is nullable for an added item. A removed item has no replacement offer.
- Actor details are a provider-asserted JSON snapshot; there is no provider-user table in this slice.

## Provider PostgreSQL topology

| PostgreSQL server | Logical database | Application role | Direct access allowed |
|---|---|---|---|
| Provider PostgreSQL | `hospital_db` | Hospital service credential | `hospital_db` only |
| Provider PostgreSQL | `ferry_db` | Ferry service credential | `ferry_db` only |
| Provider PostgreSQL | `hotel_db` | Hotel service credential | `hotel_db` only |
| Provider PostgreSQL | `transport_db` | Transport service credential | `transport_db` only |

There are no relationships across these diagrams. Similar names such as `HOLD` and `RESERVATION` describe the shared provider protocol but are separate tables and records in each database.

### `hospital_db`

```mermaid
erDiagram
    direction LR

    MEDICAL_SERVICE ||--o{ APPOINTMENT_SLOT : offers
    APPOINTMENT_SLOT ||--o{ HOLD : allocates
    HOLD ||--o| RESERVATION : confirms_as

    MEDICAL_SERVICE {
        uuid id PK
        string code UK
        string name
        int duration_minutes
        int price_amount_minor
        string price_currency
        bool active
        bool synthetic
    }

    APPOINTMENT_SLOT {
        uuid id PK
        uuid medical_service_id FK
        datetime starts_at
        datetime ends_at
        string time_zone
        int capacity_total
        string status
    }

    HOLD {
        uuid id PK
        uuid appointment_slot_id FK
        string external_request_id UK
        int patient_count
        string status
        datetime expires_at
        datetime created_at
    }

    RESERVATION {
        uuid id PK
        uuid hold_id FK, UK
        string confirmation_code UK
        string status
        datetime confirmed_at
        datetime released_at
    }

    IDEMPOTENCY_RECORD {
        uuid id PK
        string client_scope
        string operation
        string idempotency_key
        string request_fingerprint
        int response_status
        json response_body
    }
```

Availability is `capacity_total` minus active, non-expired holds and confirmed reservations, calculated while locking the appointment slot row.

### `ferry_db`

```mermaid
erDiagram
    direction LR

    SAILING ||--o{ HOLD : allocates
    HOLD ||--o| RESERVATION : confirms_as

    SAILING {
        uuid id PK
        string sailing_code UK
        string origin_port
        string destination_port
        datetime departs_at
        datetime arrives_at
        string departure_time_zone
        string arrival_time_zone
        datetime check_in_cutoff_at
        int seat_capacity
        int price_amount_minor
        string price_currency
        string status
        bool synthetic
    }

    HOLD {
        uuid id PK
        uuid sailing_id FK
        string external_request_id UK
        int passenger_count
        string status
        datetime expires_at
        datetime created_at
    }

    RESERVATION {
        uuid id PK
        uuid hold_id FK, UK
        string confirmation_code UK
        int passenger_count
        string status
        datetime confirmed_at
        datetime released_at
    }

    IDEMPOTENCY_RECORD {
        uuid id PK
        string client_scope
        string operation
        string idempotency_key
        string request_fingerprint
        int response_status
        json response_body
    }
```

Seat availability is calculated transactionally from sailing capacity, active holds, and confirmed reservations.

### `hotel_db`

```mermaid
erDiagram
    direction LR

    ROOM_TYPE ||--o{ ROOM_INVENTORY_DAY : inventories
    ROOM_TYPE ||--o{ HOLD : requested_for
    HOLD ||--|{ HOLD_NIGHT : covers
    ROOM_INVENTORY_DAY ||--o{ HOLD_NIGHT : allocates
    HOLD ||--o| RESERVATION : confirms_as

    ROOM_TYPE {
        uuid id PK
        string code UK
        string name
        int occupancy
        json accessibility
        bool active
        bool synthetic
    }

    ROOM_INVENTORY_DAY {
        uuid id PK
        uuid room_type_id FK
        date stay_date
        int rooms_total
        int price_amount_minor
        string price_currency
        string status
    }

    HOLD {
        uuid id PK
        uuid room_type_id FK
        string external_request_id UK
        date check_in_date
        date check_out_date
        string time_zone
        int room_count
        string status
        datetime expires_at
    }

    HOLD_NIGHT {
        uuid hold_id PK, FK
        uuid inventory_day_id PK, FK
        int room_count
    }

    RESERVATION {
        uuid id PK
        uuid hold_id FK, UK
        string confirmation_code UK
        string status
        datetime confirmed_at
        datetime released_at
    }

    IDEMPOTENCY_RECORD {
        uuid id PK
        string client_scope
        string operation
        string idempotency_key
        string request_fingerprint
        int response_status
        json response_body
    }
```

`(room_type_id, stay_date)` is unique. A hold succeeds only when every night in `[check_in_date, check_out_date)` has enough inventory in one transaction.

### `transport_db`

```mermaid
erDiagram
    direction LR

    VEHICLE ||--o{ AVAILABILITY_SLOT : exposes
    AVAILABILITY_SLOT ||--o{ HOLD : allocates
    HOLD ||--o| RESERVATION : confirms_as

    VEHICLE {
        uuid id PK
        string vehicle_code UK
        string vehicle_type
        int passenger_capacity
        json accessibility
        string status
        bool synthetic
    }

    AVAILABILITY_SLOT {
        uuid id PK
        uuid vehicle_id FK
        datetime starts_at
        datetime ends_at
        string time_zone
        int price_amount_minor
        string price_currency
        string status
    }

    HOLD {
        uuid id PK
        uuid availability_slot_id FK
        string external_request_id UK
        int passenger_count
        string pickup_location
        string dropoff_location
        string status
        datetime expires_at
    }

    RESERVATION {
        uuid id PK
        uuid hold_id FK, UK
        string confirmation_code UK
        string status
        datetime confirmed_at
        datetime released_at
    }

    IDEMPOTENCY_RECORD {
        uuid id PK
        string client_scope
        string operation
        string idempotency_key
        string request_fingerprint
        int response_status
        json response_body
    }
```

Availability windows for the same vehicle must not overlap once held or confirmed, and passenger/accessibility requirements must fit the vehicle.

## Provider-wide database invariants

1. `HOLD.status` is `HELD`, `CONFIRMED`, `RELEASED`, or `EXPIRED`; terminal states are not reopened.
2. `RESERVATION.status` is `CONFIRMED` or `RELEASED`.
3. Confirm accepts only a non-expired `HELD` record and creates at most one reservation per hold.
4. Release is idempotent for both a hold and a confirmed reservation.
5. Expired holds do not consume availability and are transitioned lazily during access or by a small cleanup job.
6. Each capacity check and hold creation locks the relevant inventory rows in one provider transaction.
7. `(client_scope, operation, idempotency_key)` is unique, and request fingerprints must match on replay.
8. Provider databases store no raw medical prompt, medical record, or cross-provider itinerary.
