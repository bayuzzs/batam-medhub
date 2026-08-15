CREATE TABLE medical_services (
    id UUID PRIMARY KEY,
    code VARCHAR(64) NOT NULL UNIQUE,
    name VARCHAR(160) NOT NULL,
    default_duration_minutes INTEGER NOT NULL,
    patient_requestable BOOLEAN NOT NULL DEFAULT TRUE,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    synthetic BOOLEAN NOT NULL DEFAULT TRUE,
    source VARCHAR(16) NOT NULL DEFAULT 'MOCK',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT medical_services_code_format CHECK (code ~ '^[A-Z][A-Z0-9_]*$'),
    CONSTRAINT medical_services_duration_positive CHECK (default_duration_minutes > 0),
    CONSTRAINT medical_services_synthetic_only CHECK (synthetic AND source = 'MOCK')
);

CREATE TABLE appointment_slots (
    id UUID PRIMARY KEY,
    medical_service_id UUID NOT NULL REFERENCES medical_services(id) ON DELETE RESTRICT,
    external_slot_id VARCHAR(128) NOT NULL UNIQUE,
    offer_id VARCHAR(128) NOT NULL UNIQUE,
    facility_name VARCHAR(160) NOT NULL,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    start_time_zone VARCHAR(64) NOT NULL,
    end_time_zone VARCHAR(64) NOT NULL,
    duration_minutes INTEGER NOT NULL,
    capacity_total INTEGER NOT NULL,
    price_amount_minor BIGINT NOT NULL,
    price_currency CHAR(3) NOT NULL,
    valid_until TIMESTAMPTZ NOT NULL,
    languages JSONB NOT NULL DEFAULT '[]'::jsonb,
    accessibility JSONB NOT NULL DEFAULT '[]'::jsonb,
    status VARCHAR(16) NOT NULL DEFAULT 'AVAILABLE',
    synthetic BOOLEAN NOT NULL DEFAULT TRUE,
    source VARCHAR(16) NOT NULL DEFAULT 'MOCK',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT appointment_slots_id_offer_unique UNIQUE (id, offer_id),
    CONSTRAINT appointment_slots_external_slot_id_format CHECK (
        external_slot_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$'
    ),
    CONSTRAINT appointment_slots_offer_id_format CHECK (
        offer_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$'
    ),
    CONSTRAINT appointment_slots_window_valid CHECK (starts_at < ends_at),
    CONSTRAINT appointment_slots_duration_matches_window CHECK (
        ends_at = starts_at + duration_minutes * INTERVAL '1 minute'
    ),
    CONSTRAINT appointment_slots_duration_positive CHECK (duration_minutes > 0),
    CONSTRAINT appointment_slots_capacity_positive CHECK (capacity_total > 0),
    CONSTRAINT appointment_slots_price_nonnegative CHECK (price_amount_minor >= 0),
    CONSTRAINT appointment_slots_currency_format CHECK (price_currency ~ '^[A-Z]{3}$'),
    CONSTRAINT appointment_slots_offer_expiry_valid CHECK (valid_until <= starts_at),
    CONSTRAINT appointment_slots_time_zone_format CHECK (
        start_time_zone ~ '^[A-Za-z_]+(/[A-Za-z0-9._+-]+)+$'
        AND end_time_zone ~ '^[A-Za-z_]+(/[A-Za-z0-9._+-]+)+$'
    ),
    CONSTRAINT appointment_slots_languages_array CHECK (jsonb_typeof(languages) = 'array'),
    CONSTRAINT appointment_slots_accessibility_array CHECK (jsonb_typeof(accessibility) = 'array'),
    CONSTRAINT appointment_slots_status_valid CHECK (status IN ('AVAILABLE', 'UNAVAILABLE', 'CANCELLED')),
    CONSTRAINT appointment_slots_synthetic_only CHECK (synthetic AND source = 'MOCK')
);

CREATE INDEX appointment_slots_search_idx
    ON appointment_slots (medical_service_id, status, starts_at, ends_at);
CREATE INDEX appointment_slots_valid_until_idx ON appointment_slots (valid_until);

CREATE TABLE holds (
    id UUID PRIMARY KEY,
    external_hold_id VARCHAR(128) NOT NULL UNIQUE,
    external_reference VARCHAR(128) NOT NULL UNIQUE,
    appointment_slot_id UUID NOT NULL REFERENCES appointment_slots(id) ON DELETE RESTRICT,
    offer_id VARCHAR(128) NOT NULL,
    client_reference VARCHAR(128) NOT NULL UNIQUE,
    patient_count INTEGER NOT NULL,
    unit_price_amount_minor BIGINT NOT NULL,
    total_price_amount_minor BIGINT NOT NULL,
    price_currency CHAR(3) NOT NULL,
    service_starts_at TIMESTAMPTZ NOT NULL,
    service_ends_at TIMESTAMPTZ NOT NULL,
    start_time_zone VARCHAR(64) NOT NULL,
    end_time_zone VARCHAR(64) NOT NULL,
    status VARCHAR(16) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    confirmed_at TIMESTAMPTZ,
    released_at TIMESTAMPTZ,
    expired_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT holds_reservation_identity_unique UNIQUE (id, appointment_slot_id, offer_id),
    CONSTRAINT holds_offer_matches_slot FOREIGN KEY (appointment_slot_id, offer_id)
        REFERENCES appointment_slots(id, offer_id) ON DELETE RESTRICT,
    CONSTRAINT holds_external_hold_id_format CHECK (
        external_hold_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$'
    ),
    CONSTRAINT holds_external_reference_format CHECK (
        external_reference ~ '^[A-Za-z0-9][A-Za-z0-9._:/-]{2,127}$'
    ),
    CONSTRAINT holds_offer_id_format CHECK (
        offer_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$'
    ),
    CONSTRAINT holds_client_reference_format CHECK (
        client_reference ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$'
    ),
    CONSTRAINT holds_patient_count_positive CHECK (patient_count > 0),
    CONSTRAINT holds_price_nonnegative CHECK (
        unit_price_amount_minor >= 0
        AND total_price_amount_minor = unit_price_amount_minor * patient_count
    ),
    CONSTRAINT holds_currency_format CHECK (price_currency ~ '^[A-Z]{3}$'),
    CONSTRAINT holds_service_window_valid CHECK (service_starts_at < service_ends_at),
    CONSTRAINT holds_time_zone_format CHECK (
        start_time_zone ~ '^[A-Za-z_]+(/[A-Za-z0-9._+-]+)+$'
        AND end_time_zone ~ '^[A-Za-z_]+(/[A-Za-z0-9._+-]+)+$'
    ),
    CONSTRAINT holds_expiry_valid CHECK (created_at < expires_at),
    CONSTRAINT holds_status_valid CHECK (status IN ('HELD', 'CONFIRMED', 'RELEASED', 'EXPIRED')),
    CONSTRAINT holds_status_timestamps_valid CHECK (
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
    CONSTRAINT holds_event_timestamps_valid CHECK (
        (confirmed_at IS NULL OR (confirmed_at >= created_at AND confirmed_at < expires_at))
        AND (released_at IS NULL OR released_at >= created_at)
        AND (expired_at IS NULL OR expired_at >= expires_at)
    )
);

CREATE INDEX holds_capacity_idx
    ON holds (appointment_slot_id, status, expires_at);

CREATE TABLE reservations (
    id UUID PRIMARY KEY,
    external_reservation_id VARCHAR(128) NOT NULL UNIQUE,
    external_reference VARCHAR(128) NOT NULL UNIQUE,
    hold_id UUID NOT NULL UNIQUE,
    appointment_slot_id UUID NOT NULL REFERENCES appointment_slots(id) ON DELETE RESTRICT,
    offer_id VARCHAR(128) NOT NULL,
    client_reference VARCHAR(128) NOT NULL,
    patient_count INTEGER NOT NULL,
    total_price_amount_minor BIGINT NOT NULL,
    price_currency CHAR(3) NOT NULL,
    service_starts_at TIMESTAMPTZ NOT NULL,
    service_ends_at TIMESTAMPTZ NOT NULL,
    start_time_zone VARCHAR(64) NOT NULL,
    end_time_zone VARCHAR(64) NOT NULL,
    status VARCHAR(16) NOT NULL,
    confirmed_at TIMESTAMPTZ NOT NULL,
    released_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT reservations_match_hold FOREIGN KEY (hold_id, appointment_slot_id, offer_id)
        REFERENCES holds(id, appointment_slot_id, offer_id) ON DELETE RESTRICT,
    CONSTRAINT reservations_external_reservation_id_format CHECK (
        external_reservation_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$'
    ),
    CONSTRAINT reservations_external_reference_format CHECK (
        external_reference ~ '^[A-Za-z0-9][A-Za-z0-9._:/-]{2,127}$'
    ),
    CONSTRAINT reservations_offer_id_format CHECK (
        offer_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$'
    ),
    CONSTRAINT reservations_client_reference_format CHECK (
        client_reference ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$'
    ),
    CONSTRAINT reservations_patient_count_positive CHECK (patient_count > 0),
    CONSTRAINT reservations_price_nonnegative CHECK (total_price_amount_minor >= 0),
    CONSTRAINT reservations_currency_format CHECK (price_currency ~ '^[A-Z]{3}$'),
    CONSTRAINT reservations_service_window_valid CHECK (service_starts_at < service_ends_at),
    CONSTRAINT reservations_time_zone_format CHECK (
        start_time_zone ~ '^[A-Za-z_]+(/[A-Za-z0-9._+-]+)+$'
        AND end_time_zone ~ '^[A-Za-z_]+(/[A-Za-z0-9._+-]+)+$'
    ),
    CONSTRAINT reservations_status_valid CHECK (status IN ('CONFIRMED', 'RELEASED')),
    CONSTRAINT reservations_release_timestamp_valid CHECK (
        (status = 'CONFIRMED' AND released_at IS NULL)
        OR (status = 'RELEASED' AND released_at IS NOT NULL)
    )
);

CREATE INDEX reservations_slot_status_idx
    ON reservations (appointment_slot_id, status);
CREATE INDEX reservations_client_reference_idx ON reservations (client_reference);

CREATE TABLE idempotency_records (
    id UUID PRIMARY KEY,
    client_scope VARCHAR(128) NOT NULL,
    operation VARCHAR(160) NOT NULL,
    idempotency_key VARCHAR(128) NOT NULL,
    request_fingerprint CHAR(64) NOT NULL,
    response_status INTEGER NOT NULL,
    response_body JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT idempotency_records_scope_key UNIQUE (client_scope, operation, idempotency_key),
    CONSTRAINT idempotency_records_key_format CHECK (
        idempotency_key ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{7,127}$'
    ),
    CONSTRAINT idempotency_records_fingerprint_format CHECK (
        request_fingerprint ~ '^[0-9a-f]{64}$'
    ),
    CONSTRAINT idempotency_records_response_status_valid CHECK (response_status BETWEEN 100 AND 599),
    CONSTRAINT idempotency_records_response_body_object CHECK (jsonb_typeof(response_body) = 'object')
);

CREATE INDEX idempotency_records_created_at_idx ON idempotency_records (created_at);
