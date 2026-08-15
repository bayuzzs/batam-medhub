CREATE TABLE room_types (
    id uuid PRIMARY KEY,
    code varchar(64) NOT NULL UNIQUE,
    property_name varchar(160) NOT NULL,
    name varchar(80) NOT NULL,
    max_guests_per_room integer NOT NULL CHECK (max_guests_per_room > 0),
    accessibility jsonb NOT NULL DEFAULT '[]'::jsonb
        CHECK (jsonb_typeof(accessibility) = 'array'),
    active boolean NOT NULL DEFAULT true,
    synthetic boolean NOT NULL DEFAULT true CHECK (synthetic),
    source varchar(16) NOT NULL DEFAULT 'MOCK' CHECK (source = 'MOCK'),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE room_inventory_days (
    id uuid PRIMARY KEY,
    room_type_id uuid NOT NULL REFERENCES room_types(id) ON DELETE RESTRICT,
    stay_date date NOT NULL,
    rooms_total integer NOT NULL CHECK (rooms_total > 0),
    price_amount_minor bigint NOT NULL CHECK (price_amount_minor >= 0),
    price_currency char(3) NOT NULL CHECK (price_currency ~ '^[A-Z]{3}$'),
    status varchar(16) NOT NULL DEFAULT 'AVAILABLE'
        CHECK (status IN ('AVAILABLE', 'UNAVAILABLE')),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (room_type_id, stay_date),
    UNIQUE (id, room_type_id)
);

-- Search returns persisted opaque offer IDs. The core must pass the same value
-- back when creating a hold; it never reconstructs an ID from stay details.
CREATE TABLE hotel_offers (
    id uuid PRIMARY KEY,
    external_offer_id varchar(128) NOT NULL UNIQUE
        CHECK (external_offer_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$'),
    external_inventory_id varchar(128) NOT NULL UNIQUE
        CHECK (external_inventory_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$'),
    room_type_id uuid NOT NULL REFERENCES room_types(id) ON DELETE RESTRICT,
    check_in_date date NOT NULL,
    check_out_date date NOT NULL,
    service_starts_at timestamptz NOT NULL,
    service_ends_at timestamptz NOT NULL,
    start_time_zone varchar(64) NOT NULL,
    end_time_zone varchar(64) NOT NULL,
    valid_until timestamptz NOT NULL,
    status varchar(16) NOT NULL DEFAULT 'AVAILABLE'
        CHECK (status IN ('AVAILABLE', 'UNAVAILABLE', 'EXPIRED')),
    synthetic boolean NOT NULL DEFAULT true CHECK (synthetic),
    source varchar(16) NOT NULL DEFAULT 'MOCK' CHECK (source = 'MOCK'),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (check_out_date > check_in_date),
    CHECK (service_ends_at > service_starts_at),
    CHECK (start_time_zone ~ '^[A-Za-z_]+(/[A-Za-z0-9._+-]+)+$'),
    CHECK (end_time_zone ~ '^[A-Za-z_]+(/[A-Za-z0-9._+-]+)+$'),
    UNIQUE (id, room_type_id),
    UNIQUE (id, external_offer_id)
);

CREATE TABLE hotel_offer_nights (
    offer_pk uuid NOT NULL,
    inventory_day_id uuid NOT NULL,
    room_type_id uuid NOT NULL,
    PRIMARY KEY (offer_pk, inventory_day_id),
    FOREIGN KEY (offer_pk, room_type_id)
        REFERENCES hotel_offers(id, room_type_id) ON DELETE CASCADE,
    FOREIGN KEY (inventory_day_id, room_type_id)
        REFERENCES room_inventory_days(id, room_type_id) ON DELETE RESTRICT
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
    room_type_id uuid NOT NULL,
    client_reference varchar(128) NOT NULL UNIQUE
        CHECK (client_reference ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$'),
    check_in_date date NOT NULL,
    check_out_date date NOT NULL,
    time_zone varchar(64) NOT NULL,
    room_count integer NOT NULL CHECK (room_count > 0),
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
    CHECK (check_out_date > check_in_date),
    CHECK (time_zone ~ '^[A-Za-z_]+(/[A-Za-z0-9._+-]+)+$'),
    CHECK (service_ends_at > service_starts_at),
    CHECK (start_time_zone ~ '^[A-Za-z_]+(/[A-Za-z0-9._+-]+)+$'),
    CHECK (end_time_zone ~ '^[A-Za-z_]+(/[A-Za-z0-9._+-]+)+$'),
    CHECK (expires_at > created_at),
    CHECK (total_price_amount_minor = unit_price_amount_minor * room_count),
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
    FOREIGN KEY (offer_pk, room_type_id)
        REFERENCES hotel_offers(id, room_type_id) ON DELETE RESTRICT,
    FOREIGN KEY (offer_pk, offer_id)
        REFERENCES hotel_offers(id, external_offer_id) ON DELETE RESTRICT
);

CREATE TABLE hold_nights (
    hold_id uuid NOT NULL,
    offer_pk uuid NOT NULL,
    inventory_day_id uuid NOT NULL,
    room_count integer NOT NULL CHECK (room_count > 0),
    PRIMARY KEY (hold_id, inventory_day_id),
    FOREIGN KEY (hold_id, offer_pk)
        REFERENCES holds(id, offer_pk) ON DELETE CASCADE,
    FOREIGN KEY (offer_pk, inventory_day_id)
        REFERENCES hotel_offer_nights(offer_pk, inventory_day_id) ON DELETE RESTRICT
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
    room_count integer NOT NULL CHECK (room_count > 0),
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

CREATE INDEX room_inventory_days_search_idx
    ON room_inventory_days (stay_date, room_type_id)
    WHERE status = 'AVAILABLE';
CREATE INDEX hotel_offers_search_idx
    ON hotel_offers (check_in_date, check_out_date, room_type_id, valid_until)
    WHERE status = 'AVAILABLE';
CREATE INDEX hotel_offer_nights_inventory_idx
    ON hotel_offer_nights (inventory_day_id, offer_pk);
CREATE INDEX holds_offer_status_expiry_idx
    ON holds (offer_pk, status, expires_at);
CREATE INDEX hold_nights_inventory_idx
    ON hold_nights (inventory_day_id, hold_id);
CREATE INDEX reservations_offer_status_idx
    ON reservations (offer_pk, status);
CREATE INDEX idempotency_records_expiry_idx
    ON idempotency_records (expires_at)
    WHERE expires_at IS NOT NULL;
