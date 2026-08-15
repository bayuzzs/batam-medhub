CREATE TABLE vehicles (
    id uuid PRIMARY KEY,
    vehicle_code varchar(64) NOT NULL UNIQUE,
    vehicle_type varchar(80) NOT NULL,
    passenger_capacity integer NOT NULL CHECK (passenger_capacity > 0),
    accessibility jsonb NOT NULL DEFAULT '[]'::jsonb
        CHECK (jsonb_typeof(accessibility) = 'array'),
    status varchar(16) NOT NULL DEFAULT 'ACTIVE'
        CHECK (status IN ('ACTIVE', 'MAINTENANCE', 'INACTIVE')),
    synthetic boolean NOT NULL DEFAULT true CHECK (synthetic),
    source varchar(16) NOT NULL DEFAULT 'MOCK' CHECK (source = 'MOCK'),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Offers carry the exact opaque IDs exposed by search. Individual vehicle
-- availability rows beneath each offer are the lockable assignment units.
CREATE TABLE transport_offers (
    id uuid PRIMARY KEY,
    external_offer_id varchar(128) NOT NULL UNIQUE
        CHECK (external_offer_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$'),
    external_availability_id varchar(128) NOT NULL UNIQUE
        CHECK (external_availability_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$'),
    vehicle_type varchar(80) NOT NULL,
    pickup_location_code varchar(64) NOT NULL
        CHECK (pickup_location_code ~ '^[A-Z][A-Z0-9_]*$'),
    dropoff_location_code varchar(64) NOT NULL
        CHECK (dropoff_location_code ~ '^[A-Z][A-Z0-9_]*$'),
    passenger_capacity integer NOT NULL CHECK (passenger_capacity > 0),
    accessibility jsonb NOT NULL DEFAULT '[]'::jsonb
        CHECK (jsonb_typeof(accessibility) = 'array'),
    service_starts_at timestamptz NOT NULL,
    service_ends_at timestamptz NOT NULL,
    start_time_zone varchar(64) NOT NULL,
    end_time_zone varchar(64) NOT NULL,
    price_amount_minor bigint NOT NULL CHECK (price_amount_minor >= 0),
    price_currency char(3) NOT NULL CHECK (price_currency ~ '^[A-Z]{3}$'),
    valid_until timestamptz NOT NULL,
    status varchar(16) NOT NULL DEFAULT 'AVAILABLE'
        CHECK (status IN ('AVAILABLE', 'UNAVAILABLE', 'EXPIRED')),
    synthetic boolean NOT NULL DEFAULT true CHECK (synthetic),
    source varchar(16) NOT NULL DEFAULT 'MOCK' CHECK (source = 'MOCK'),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (pickup_location_code <> dropoff_location_code),
    CHECK (service_ends_at > service_starts_at),
    CHECK (start_time_zone ~ '^[A-Za-z_]+(/[A-Za-z0-9._+-]+)+$'),
    CHECK (end_time_zone ~ '^[A-Za-z_]+(/[A-Za-z0-9._+-]+)+$'),
    UNIQUE (id, external_offer_id)
);

CREATE TABLE availability_slots (
    id uuid PRIMARY KEY,
    slot_code varchar(128) NOT NULL UNIQUE,
    offer_pk uuid NOT NULL REFERENCES transport_offers(id) ON DELETE RESTRICT,
    vehicle_id uuid NOT NULL REFERENCES vehicles(id) ON DELETE RESTRICT,
    starts_at timestamptz NOT NULL,
    ends_at timestamptz NOT NULL,
    time_zone varchar(64) NOT NULL,
    status varchar(16) NOT NULL DEFAULT 'AVAILABLE'
        CHECK (status IN ('AVAILABLE', 'UNAVAILABLE')),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (ends_at > starts_at),
    CHECK (time_zone ~ '^[A-Za-z_]+(/[A-Za-z0-9._+-]+)+$'),
    UNIQUE (offer_pk, vehicle_id),
    UNIQUE (id, offer_pk)
);

CREATE TABLE holds (
    id uuid PRIMARY KEY,
    external_hold_id varchar(128) NOT NULL UNIQUE
        CHECK (external_hold_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$'),
    external_reference varchar(128) NOT NULL UNIQUE
        CHECK (external_reference ~ '^[A-Za-z0-9][A-Za-z0-9._:/-]{2,127}$'),
    offer_pk uuid NOT NULL,
    offer_id varchar(128) NOT NULL
        CHECK (offer_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$'),
    client_reference varchar(128) NOT NULL UNIQUE
        CHECK (client_reference ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$'),
    units integer NOT NULL CHECK (units > 0),
    passenger_count integer NOT NULL CHECK (passenger_count > 0),
    pickup_location_code varchar(64) NOT NULL
        CHECK (pickup_location_code ~ '^[A-Z][A-Z0-9_]*$'),
    dropoff_location_code varchar(64) NOT NULL
        CHECK (dropoff_location_code ~ '^[A-Z][A-Z0-9_]*$'),
    requested_pickup_starts_at timestamptz NOT NULL,
    requested_pickup_ends_at timestamptz NOT NULL,
    requested_start_time_zone varchar(64) NOT NULL,
    requested_end_time_zone varchar(64) NOT NULL,
    accessibility jsonb NOT NULL DEFAULT '[]'::jsonb
        CHECK (jsonb_typeof(accessibility) = 'array'),
    unit_price_amount_minor bigint NOT NULL CHECK (unit_price_amount_minor >= 0),
    total_price_amount_minor bigint NOT NULL CHECK (total_price_amount_minor >= 0),
    price_currency char(3) NOT NULL CHECK (price_currency ~ '^[A-Z]{3}$'),
    service_starts_at timestamptz NOT NULL,
    service_ends_at timestamptz NOT NULL,
    start_time_zone varchar(64) NOT NULL,
    end_time_zone varchar(64) NOT NULL,
    status varchar(16) NOT NULL DEFAULT 'HELD'
        CHECK (status IN ('HELD', 'CONFIRMED', 'RELEASED', 'EXPIRED')),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at timestamptz NOT NULL,
    confirmed_at timestamptz,
    released_at timestamptz,
    expired_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (pickup_location_code <> dropoff_location_code),
    CHECK (requested_pickup_ends_at > requested_pickup_starts_at),
    CHECK (requested_start_time_zone ~ '^[A-Za-z_]+(/[A-Za-z0-9._+-]+)+$'),
    CHECK (requested_end_time_zone ~ '^[A-Za-z_]+(/[A-Za-z0-9._+-]+)+$'),
    CHECK (service_ends_at > service_starts_at),
    CHECK (start_time_zone ~ '^[A-Za-z_]+(/[A-Za-z0-9._+-]+)+$'),
    CHECK (end_time_zone ~ '^[A-Za-z_]+(/[A-Za-z0-9._+-]+)+$'),
    CHECK (expires_at > created_at),
    CHECK (total_price_amount_minor = unit_price_amount_minor * units),
    CHECK (
        (
            status = 'HELD'
            AND confirmed_at IS NULL
            AND released_at IS NULL
            AND expired_at IS NULL
        ) OR (
            status = 'CONFIRMED'
            AND confirmed_at IS NOT NULL
            AND released_at IS NULL
            AND expired_at IS NULL
        ) OR (
            status = 'RELEASED'
            AND released_at IS NOT NULL
            AND expired_at IS NULL
        ) OR (
            status = 'EXPIRED'
            AND confirmed_at IS NULL
            AND released_at IS NULL
            AND expired_at IS NOT NULL
        )
    ),
    CHECK (
        (confirmed_at IS NULL OR (confirmed_at >= created_at AND confirmed_at < expires_at))
        AND (released_at IS NULL OR released_at >= created_at)
        AND (expired_at IS NULL OR expired_at >= expires_at)
    ),
    UNIQUE (id, offer_pk),
    UNIQUE (id, offer_pk, offer_id),
    FOREIGN KEY (offer_pk, offer_id)
        REFERENCES transport_offers(id, external_offer_id) ON DELETE RESTRICT
);

CREATE TABLE hold_assignments (
    hold_id uuid NOT NULL,
    offer_pk uuid NOT NULL,
    availability_slot_id uuid NOT NULL,
    PRIMARY KEY (hold_id, availability_slot_id),
    FOREIGN KEY (hold_id, offer_pk)
        REFERENCES holds(id, offer_pk) ON DELETE CASCADE,
    FOREIGN KEY (availability_slot_id, offer_pk)
        REFERENCES availability_slots(id, offer_pk) ON DELETE RESTRICT
);

CREATE TABLE reservations (
    id uuid PRIMARY KEY,
    external_reservation_id varchar(128) NOT NULL UNIQUE
        CHECK (external_reservation_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$'),
    external_reference varchar(128) NOT NULL UNIQUE
        CHECK (external_reference ~ '^[A-Za-z0-9][A-Za-z0-9._:/-]{2,127}$'),
    hold_id uuid NOT NULL UNIQUE,
    offer_pk uuid NOT NULL,
    offer_id varchar(128) NOT NULL
        CHECK (offer_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$'),
    client_reference varchar(128) NOT NULL UNIQUE
        CHECK (client_reference ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$'),
    units integer NOT NULL CHECK (units > 0),
    passenger_count integer NOT NULL CHECK (passenger_count > 0),
    total_price_amount_minor bigint NOT NULL CHECK (total_price_amount_minor >= 0),
    price_currency char(3) NOT NULL CHECK (price_currency ~ '^[A-Z]{3}$'),
    service_starts_at timestamptz NOT NULL,
    service_ends_at timestamptz NOT NULL,
    start_time_zone varchar(64) NOT NULL,
    end_time_zone varchar(64) NOT NULL,
    status varchar(16) NOT NULL DEFAULT 'CONFIRMED'
        CHECK (status IN ('CONFIRMED', 'RELEASED')),
    confirmed_at timestamptz NOT NULL,
    released_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (service_ends_at > service_starts_at),
    CHECK (start_time_zone ~ '^[A-Za-z_]+(/[A-Za-z0-9._+-]+)+$'),
    CHECK (end_time_zone ~ '^[A-Za-z_]+(/[A-Za-z0-9._+-]+)+$'),
    CHECK (
        (status = 'CONFIRMED' AND released_at IS NULL)
        OR (status = 'RELEASED' AND released_at IS NOT NULL)
    ),
    FOREIGN KEY (hold_id, offer_pk, offer_id)
        REFERENCES holds(id, offer_pk, offer_id) ON DELETE RESTRICT
);

CREATE TABLE idempotency_records (
    id uuid PRIMARY KEY,
    client_scope varchar(128) NOT NULL,
    operation varchar(128) NOT NULL,
    idempotency_key varchar(128) NOT NULL
        CHECK (idempotency_key ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{7,127}$'),
    request_fingerprint char(64) NOT NULL
        CHECK (request_fingerprint ~ '^[0-9a-f]{64}$'),
    response_status integer NOT NULL CHECK (response_status BETWEEN 100 AND 599),
    response_body jsonb NOT NULL CHECK (jsonb_typeof(response_body) = 'object'),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at timestamptz,
    UNIQUE (client_scope, operation, idempotency_key)
);

CREATE INDEX vehicles_bookable_idx
    ON vehicles (vehicle_type, passenger_capacity)
    WHERE status = 'ACTIVE';
CREATE INDEX transport_offers_search_idx
    ON transport_offers (
        pickup_location_code,
        dropoff_location_code,
        service_starts_at,
        service_ends_at,
        valid_until
    ) WHERE status = 'AVAILABLE';
CREATE INDEX availability_slots_offer_status_idx
    ON availability_slots (offer_pk, status);
CREATE INDEX availability_slots_vehicle_window_idx
    ON availability_slots (vehicle_id, starts_at, ends_at)
    WHERE status = 'AVAILABLE';
CREATE INDEX holds_offer_status_expiry_idx
    ON holds (offer_pk, status, expires_at);
CREATE INDEX hold_assignments_slot_idx
    ON hold_assignments (availability_slot_id, hold_id);
CREATE INDEX reservations_offer_status_idx
    ON reservations (offer_pk, status);
CREATE INDEX idempotency_records_expiry_idx
    ON idempotency_records (expires_at)
    WHERE expires_at IS NOT NULL;
