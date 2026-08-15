CREATE TABLE patients (
    id uuid PRIMARY KEY,
    email_normalized varchar(254) NOT NULL,
    password_hash char(60) NOT NULL,
    full_name varchar(120) NOT NULL,
    preferred_currency char(3) NOT NULL DEFAULT 'SGD',
    status varchar(32) NOT NULL DEFAULT 'ACTIVE',
    synthetic boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT patients_email_normalized_format CHECK (
        email_normalized = lower(btrim(email_normalized))
        AND char_length(email_normalized) BETWEEN 3 AND 254
    ),
    CONSTRAINT patients_full_name_length CHECK (char_length(btrim(full_name)) BETWEEN 2 AND 120),
    CONSTRAINT patients_password_hash_format CHECK (
        password_hash ~ '^\$2[aby]\$[0-9]{2}\$[./A-Za-z0-9]{53}$'
    ),
    CONSTRAINT patients_preferred_currency_check CHECK (preferred_currency IN ('SGD', 'IDR')),
    CONSTRAINT patients_status_check CHECK (status IN ('ACTIVE', 'DISABLED')),
    CONSTRAINT patients_synthetic_check CHECK (synthetic),
    CONSTRAINT patients_updated_after_created CHECK (updated_at >= created_at),
    CONSTRAINT patients_email_normalized_key UNIQUE (email_normalized)
);

CREATE TABLE auth_sessions (
    id uuid PRIMARY KEY,
    patient_id uuid NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    refresh_token_hash char(64) NOT NULL,
    replaced_by_session_id uuid REFERENCES auth_sessions(id) ON DELETE SET NULL,
    expires_at timestamptz NOT NULL,
    revoked_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    last_used_at timestamptz,
    CONSTRAINT auth_sessions_refresh_hash_format CHECK (refresh_token_hash ~ '^[0-9a-f]{64}$'),
    CONSTRAINT auth_sessions_expiry_check CHECK (expires_at > created_at),
    CONSTRAINT auth_sessions_revoked_check CHECK (revoked_at IS NULL OR revoked_at >= created_at),
    CONSTRAINT auth_sessions_last_used_check CHECK (last_used_at IS NULL OR last_used_at >= created_at),
    CONSTRAINT auth_sessions_no_self_replacement CHECK (replaced_by_session_id IS NULL OR replaced_by_session_id <> id),
    CONSTRAINT auth_sessions_refresh_token_hash_key UNIQUE (refresh_token_hash)
);

CREATE UNIQUE INDEX auth_sessions_replaced_by_session_id_key
    ON auth_sessions(replaced_by_session_id)
    WHERE replaced_by_session_id IS NOT NULL;
CREATE INDEX auth_sessions_patient_id_idx ON auth_sessions(patient_id);
CREATE INDEX auth_sessions_active_lookup_idx
    ON auth_sessions(patient_id, expires_at)
    WHERE revoked_at IS NULL;

CREATE TABLE providers (
    id uuid PRIMARY KEY,
    provider_type varchar(16) NOT NULL,
    code varchar(64) NOT NULL,
    display_name varchar(160) NOT NULL,
    status varchar(32) NOT NULL DEFAULT 'ACTIVE',
    synthetic boolean NOT NULL DEFAULT true,
    source varchar(16) NOT NULL DEFAULT 'MOCK',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT providers_type_check CHECK (provider_type IN ('HOSPITAL', 'FERRY', 'HOTEL', 'TRANSPORT')),
    CONSTRAINT providers_code_format CHECK (code ~ '^[a-z][a-z0-9]*(-[a-z0-9]+)*$'),
    CONSTRAINT providers_status_check CHECK (status IN ('ACTIVE', 'INACTIVE')),
    CONSTRAINT providers_synthetic_check CHECK (synthetic),
    CONSTRAINT providers_source_check CHECK (source = 'MOCK'),
    CONSTRAINT providers_updated_after_created CHECK (updated_at >= created_at),
    CONSTRAINT providers_code_key UNIQUE (code)
);

CREATE TABLE provider_credentials (
    id uuid PRIMARY KEY,
    provider_id uuid NOT NULL REFERENCES providers(id) ON DELETE CASCADE,
    key_prefix varchar(32) NOT NULL,
    secret_hash char(64) NOT NULL,
    hash_algorithm varchar(16) NOT NULL DEFAULT 'SHA256',
    status varchar(32) NOT NULL DEFAULT 'ACTIVE',
    expires_at timestamptz,
    synthetic boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT provider_credentials_key_prefix_length CHECK (char_length(key_prefix) BETWEEN 4 AND 32),
    CONSTRAINT provider_credentials_secret_hash_format CHECK (secret_hash ~ '^[0-9a-f]{64}$'),
    CONSTRAINT provider_credentials_hash_algorithm_check CHECK (hash_algorithm = 'SHA256'),
    CONSTRAINT provider_credentials_status_check CHECK (status IN ('ACTIVE', 'REVOKED')),
    CONSTRAINT provider_credentials_expiry_check CHECK (expires_at IS NULL OR expires_at > created_at),
    CONSTRAINT provider_credentials_synthetic_check CHECK (synthetic),
    CONSTRAINT provider_credentials_updated_after_created CHECK (updated_at >= created_at),
    CONSTRAINT provider_credentials_key_prefix_key UNIQUE (key_prefix),
    CONSTRAINT provider_credentials_secret_hash_key UNIQUE (secret_hash)
);

CREATE INDEX provider_credentials_provider_id_idx ON provider_credentials(provider_id);
CREATE INDEX provider_credentials_active_lookup_idx
    ON provider_credentials(key_prefix)
    WHERE status = 'ACTIVE';

CREATE TABLE medical_services (
    id uuid PRIMARY KEY,
    code varchar(64) NOT NULL,
    name varchar(160) NOT NULL,
    category varchar(64) NOT NULL,
    description text,
    default_duration_minutes integer NOT NULL,
    active boolean NOT NULL DEFAULT true,
    synthetic boolean NOT NULL DEFAULT true,
    source varchar(16) NOT NULL DEFAULT 'MOCK',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT medical_services_code_format CHECK (code ~ '^[A-Z][A-Z0-9_]{2,63}$'),
    CONSTRAINT medical_services_name_length CHECK (char_length(btrim(name)) BETWEEN 1 AND 160),
    CONSTRAINT medical_services_category_length CHECK (char_length(category) BETWEEN 1 AND 64),
    CONSTRAINT medical_services_duration_check CHECK (default_duration_minutes BETWEEN 1 AND 1440),
    CONSTRAINT medical_services_synthetic_check CHECK (synthetic),
    CONSTRAINT medical_services_source_check CHECK (source = 'MOCK'),
    CONSTRAINT medical_services_updated_after_created CHECK (updated_at >= created_at),
    CONSTRAINT medical_services_code_key UNIQUE (code)
);

CREATE TABLE provider_capabilities (
    id uuid PRIMARY KEY,
    provider_id uuid NOT NULL REFERENCES providers(id) ON DELETE RESTRICT,
    medical_service_id uuid NOT NULL REFERENCES medical_services(id) ON DELETE RESTRICT,
    external_service_id varchar(128) NOT NULL,
    active boolean NOT NULL DEFAULT true,
    synthetic boolean NOT NULL DEFAULT true,
    source varchar(16) NOT NULL DEFAULT 'MOCK',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT provider_capabilities_external_service_length CHECK (char_length(external_service_id) BETWEEN 1 AND 128),
    CONSTRAINT provider_capabilities_synthetic_check CHECK (synthetic),
    CONSTRAINT provider_capabilities_source_check CHECK (source = 'MOCK'),
    CONSTRAINT provider_capabilities_updated_after_created CHECK (updated_at >= created_at)
);

CREATE UNIQUE INDEX provider_capabilities_active_mapping_key
    ON provider_capabilities(provider_id, medical_service_id, external_service_id)
    WHERE active;
CREATE UNIQUE INDEX provider_capabilities_active_external_key
    ON provider_capabilities(provider_id, external_service_id)
    WHERE active;
CREATE INDEX provider_capabilities_medical_service_id_idx ON provider_capabilities(medical_service_id);

CREATE FUNCTION enforce_hospital_provider_capability() RETURNS trigger AS $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM providers
        WHERE id = NEW.provider_id AND provider_type = 'HOSPITAL'
    ) THEN
        RAISE EXCEPTION 'provider capabilities require a HOSPITAL provider';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER provider_capabilities_hospital_only
    BEFORE INSERT OR UPDATE OF provider_id ON provider_capabilities
    FOR EACH ROW EXECUTE FUNCTION enforce_hospital_provider_capability();

CREATE TABLE fx_rates (
    id uuid PRIMARY KEY,
    base_currency char(3) NOT NULL,
    quote_currency char(3) NOT NULL,
    rate numeric(24,12) NOT NULL,
    source varchar(64) NOT NULL,
    effective_at timestamptz NOT NULL,
    estimated boolean NOT NULL DEFAULT true,
    synthetic boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT fx_rates_base_currency_format CHECK (base_currency ~ '^[A-Z]{3}$'),
    CONSTRAINT fx_rates_quote_currency_format CHECK (quote_currency ~ '^[A-Z]{3}$'),
    CONSTRAINT fx_rates_rate_positive CHECK (rate > 0),
    CONSTRAINT fx_rates_source_length CHECK (char_length(source) BETWEEN 1 AND 64),
    CONSTRAINT fx_rates_estimated_check CHECK (estimated),
    CONSTRAINT fx_rates_synthetic_check CHECK (synthetic),
    CONSTRAINT fx_rates_reference_currencies_check CHECK (
        base_currency IN ('SGD', 'IDR') AND quote_currency IN ('SGD', 'IDR')
    ),
    CONSTRAINT fx_rates_identity_rate_check CHECK (base_currency <> quote_currency OR rate = 1),
    CONSTRAINT fx_rates_identity_key UNIQUE (base_currency, quote_currency, effective_at, source)
);

CREATE INDEX fx_rates_lookup_idx
    ON fx_rates(base_currency, quote_currency, effective_at DESC);

CREATE TABLE trip_requests (
    id uuid PRIMARY KEY,
    patient_id uuid NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
    status varchar(32) NOT NULL DEFAULT 'DRAFT',
    medical_service_id uuid REFERENCES medical_services(id) ON DELETE RESTRICT,
    requested_service_text varchar(300),
    structured_intent jsonb,
    reference_currency char(3) NOT NULL,
    planning_revision integer NOT NULL DEFAULT 0,
    selected_plan_option_id uuid,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT trip_requests_status_check CHECK (status IN (
        'DRAFT', 'PARSING_INTENT', 'NEEDS_INPUT', 'UNSUPPORTED_SERVICE', 'OUT_OF_SCOPE',
        'PLANNING', 'NO_MATCH', 'PLAN_READY', 'CONFIRMING', 'ACTIVE',
        'CONFIRMATION_FAILED', 'MANUAL_REVIEW'
    )),
    CONSTRAINT trip_requests_requested_service_length CHECK (
        requested_service_text IS NULL OR char_length(requested_service_text) <= 300
    ),
    CONSTRAINT trip_requests_structured_intent_object CHECK (
        structured_intent IS NULL OR jsonb_typeof(structured_intent) = 'object'
    ),
    CONSTRAINT trip_requests_reference_currency_check CHECK (reference_currency IN ('SGD', 'IDR')),
    CONSTRAINT trip_requests_planning_revision_check CHECK (planning_revision >= 0),
    CONSTRAINT trip_requests_updated_after_created CHECK (updated_at >= created_at),
    CONSTRAINT trip_requests_id_patient_key UNIQUE (id, patient_id)
);

CREATE INDEX trip_requests_patient_id_idx ON trip_requests(patient_id, created_at DESC);
CREATE INDEX trip_requests_status_idx ON trip_requests(status);

CREATE TABLE plan_options (
    id uuid PRIMARY KEY,
    trip_request_id uuid NOT NULL REFERENCES trip_requests(id) ON DELETE CASCADE,
    planning_revision integer NOT NULL,
    rank integer NOT NULL,
    status varchar(16) NOT NULL DEFAULT 'PROPOSED',
    explanation jsonb NOT NULL,
    total_price_snapshot jsonb NOT NULL,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT plan_options_planning_revision_check CHECK (planning_revision >= 1),
    CONSTRAINT plan_options_rank_check CHECK (rank BETWEEN 1 AND 2),
    CONSTRAINT plan_options_status_check CHECK (status IN ('PROPOSED', 'SELECTED', 'EXPIRED', 'CONFIRMED')),
    CONSTRAINT plan_options_explanation_array CHECK (jsonb_typeof(explanation) = 'array'),
    CONSTRAINT plan_options_total_price_object CHECK (jsonb_typeof(total_price_snapshot) = 'object'),
    CONSTRAINT plan_options_expiry_check CHECK (expires_at > created_at),
    CONSTRAINT plan_options_updated_after_created CHECK (updated_at >= created_at),
    CONSTRAINT plan_options_trip_revision_rank_key UNIQUE (trip_request_id, planning_revision, rank),
    CONSTRAINT plan_options_id_trip_key UNIQUE (id, trip_request_id)
);

CREATE UNIQUE INDEX plan_options_one_selection_per_revision_key
    ON plan_options(trip_request_id, planning_revision)
    WHERE status IN ('SELECTED', 'CONFIRMED');
CREATE INDEX plan_options_trip_request_id_idx ON plan_options(trip_request_id, planning_revision DESC);

ALTER TABLE trip_requests
    ADD CONSTRAINT trip_requests_selected_plan_option_fk
    FOREIGN KEY (selected_plan_option_id, id)
    REFERENCES plan_options(id, trip_request_id)
    DEFERRABLE INITIALLY DEFERRED;

CREATE TABLE plan_items (
    id uuid PRIMARY KEY,
    plan_option_id uuid NOT NULL REFERENCES plan_options(id) ON DELETE CASCADE,
    provider_id uuid REFERENCES providers(id) ON DELETE RESTRICT,
    item_type varchar(32) NOT NULL,
    sequence_number integer NOT NULL,
    external_offer_id varchar(128),
    title varchar(200) NOT NULL,
    starts_at timestamptz NOT NULL,
    ends_at timestamptz NOT NULL,
    start_time_zone varchar(64) NOT NULL,
    end_time_zone varchar(64) NOT NULL,
    origin_code varchar(64),
    destination_code varchar(64),
    source_amount_minor bigint,
    source_currency char(3),
    display_amount_minor bigint,
    display_currency char(3),
    fx_rate_id uuid REFERENCES fx_rates(id) ON DELETE RESTRICT,
    offer_snapshot jsonb NOT NULL,
    offer_expires_at timestamptz,
    operational_notes jsonb NOT NULL DEFAULT '[]'::jsonb,
    synthetic boolean NOT NULL DEFAULT true,
    source varchar(16) NOT NULL DEFAULT 'MOCK',
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT plan_items_item_type_check CHECK (item_type IN (
        'FERRY_OUTBOUND', 'ARRIVAL_BUFFER', 'TRANSPORT_PICKUP', 'HOSPITAL_APPOINTMENT',
        'ADDITIONAL_CARE', 'HOTEL_STAY', 'TRANSPORT_DROPOFF', 'DEPARTURE_BUFFER', 'FERRY_RETURN'
    )),
    CONSTRAINT plan_items_sequence_check CHECK (sequence_number >= 1),
    CONSTRAINT plan_items_title_length CHECK (char_length(btrim(title)) BETWEEN 1 AND 200),
    CONSTRAINT plan_items_time_window_check CHECK (ends_at > starts_at),
    CONSTRAINT plan_items_time_zone_length CHECK (
        char_length(start_time_zone) BETWEEN 3 AND 64 AND char_length(end_time_zone) BETWEEN 3 AND 64
    ),
    CONSTRAINT plan_items_provider_offer_pair CHECK (
        (provider_id IS NULL AND external_offer_id IS NULL)
        OR (provider_id IS NOT NULL AND external_offer_id IS NOT NULL)
    ),
    CONSTRAINT plan_items_provider_required_check CHECK (
        provider_id IS NOT NULL OR item_type IN ('ARRIVAL_BUFFER', 'DEPARTURE_BUFFER')
    ),
    CONSTRAINT plan_items_money_group_check CHECK (
        (source_amount_minor IS NULL AND source_currency IS NULL AND display_amount_minor IS NULL
            AND display_currency IS NULL AND fx_rate_id IS NULL)
        OR (source_amount_minor IS NOT NULL AND source_amount_minor >= 0
            AND source_currency IS NOT NULL AND source_currency ~ '^[A-Z]{3}$'
            AND display_amount_minor IS NOT NULL AND display_amount_minor >= 0
            AND display_currency IS NOT NULL AND display_currency ~ '^[A-Z]{3}$'
            AND fx_rate_id IS NOT NULL)
    ),
    CONSTRAINT plan_items_offer_snapshot_object CHECK (jsonb_typeof(offer_snapshot) = 'object'),
    CONSTRAINT plan_items_operational_notes_array CHECK (jsonb_typeof(operational_notes) = 'array'),
    CONSTRAINT plan_items_offer_expiry_check CHECK (
        offer_expires_at IS NULL OR offer_expires_at > created_at
    ),
    CONSTRAINT plan_items_synthetic_check CHECK (synthetic),
    CONSTRAINT plan_items_source_check CHECK (source = 'MOCK'),
    CONSTRAINT plan_items_plan_sequence_key UNIQUE (plan_option_id, sequence_number)
);

CREATE INDEX plan_items_plan_option_id_idx ON plan_items(plan_option_id, sequence_number);
CREATE INDEX plan_items_provider_offer_idx ON plan_items(provider_id, external_offer_id)
    WHERE provider_id IS NOT NULL;

CREATE TABLE journeys (
    id uuid PRIMARY KEY,
    trip_request_id uuid NOT NULL UNIQUE,
    patient_id uuid NOT NULL,
    status varchar(32) NOT NULL,
    current_version_number integer NOT NULL,
    activated_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT journeys_trip_patient_fk FOREIGN KEY (trip_request_id, patient_id)
        REFERENCES trip_requests(id, patient_id) ON DELETE RESTRICT,
    CONSTRAINT journeys_status_check CHECK (status IN ('ACTIVE', 'MANUAL_REVIEW')),
    CONSTRAINT journeys_current_version_check CHECK (current_version_number >= 1),
    CONSTRAINT journeys_activation_check CHECK (activated_at >= created_at),
    CONSTRAINT journeys_updated_after_created CHECK (updated_at >= created_at),
    CONSTRAINT journeys_id_trip_key UNIQUE (id, trip_request_id)
);

CREATE INDEX journeys_patient_id_idx ON journeys(patient_id, created_at DESC);

CREATE TABLE reservations (
    id uuid PRIMARY KEY,
    trip_request_id uuid NOT NULL REFERENCES trip_requests(id) ON DELETE RESTRICT,
    journey_id uuid,
    plan_item_id uuid REFERENCES plan_items(id) ON DELETE RESTRICT,
    provider_id uuid NOT NULL REFERENCES providers(id) ON DELETE RESTRICT,
    status varchar(16) NOT NULL DEFAULT 'PENDING',
    external_offer_id varchar(128) NOT NULL,
    external_hold_id varchar(128),
    external_reservation_id varchar(128),
    hold_expires_at timestamptz,
    provider_snapshot jsonb NOT NULL DEFAULT '{}'::jsonb,
    cleanup_required boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT reservations_journey_trip_fk FOREIGN KEY (journey_id, trip_request_id)
        REFERENCES journeys(id, trip_request_id) ON DELETE RESTRICT DEFERRABLE INITIALLY DEFERRED,
    CONSTRAINT reservations_status_check CHECK (status IN ('PENDING', 'HELD', 'CONFIRMED', 'RELEASED', 'FAILED')),
    CONSTRAINT reservations_external_offer_length CHECK (char_length(external_offer_id) BETWEEN 1 AND 128),
    CONSTRAINT reservations_hold_state_check CHECK (
        status NOT IN ('HELD', 'CONFIRMED') OR external_hold_id IS NOT NULL
    ),
    CONSTRAINT reservations_confirmation_state_check CHECK (
        status <> 'CONFIRMED' OR external_reservation_id IS NOT NULL
    ),
    CONSTRAINT reservations_hold_expiry_check CHECK (
        hold_expires_at IS NULL OR hold_expires_at > created_at
    ),
    CONSTRAINT reservations_provider_snapshot_object CHECK (jsonb_typeof(provider_snapshot) = 'object'),
    CONSTRAINT reservations_updated_after_created CHECK (updated_at >= created_at)
);

CREATE UNIQUE INDEX reservations_plan_item_key ON reservations(plan_item_id) WHERE plan_item_id IS NOT NULL;
CREATE UNIQUE INDEX reservations_provider_hold_key
    ON reservations(provider_id, external_hold_id) WHERE external_hold_id IS NOT NULL;
CREATE UNIQUE INDEX reservations_provider_reservation_key
    ON reservations(provider_id, external_reservation_id) WHERE external_reservation_id IS NOT NULL;
CREATE INDEX reservations_trip_request_id_idx ON reservations(trip_request_id);
CREATE INDEX reservations_journey_id_idx ON reservations(journey_id) WHERE journey_id IS NOT NULL;

CREATE TABLE itinerary_versions (
    id uuid PRIMARY KEY,
    journey_id uuid NOT NULL REFERENCES journeys(id) ON DELETE CASCADE,
    version_number integer NOT NULL,
    status varchar(16) NOT NULL DEFAULT 'DRAFT',
    change_reason varchar(500) NOT NULL,
    source_disruption_id uuid,
    total_price_snapshot jsonb NOT NULL,
    activated_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT itinerary_versions_version_check CHECK (version_number >= 1),
    CONSTRAINT itinerary_versions_status_check CHECK (status IN ('DRAFT', 'ACTIVE', 'SUPERSEDED', 'ABANDONED')),
    CONSTRAINT itinerary_versions_change_reason_length CHECK (char_length(btrim(change_reason)) BETWEEN 1 AND 500),
    CONSTRAINT itinerary_versions_total_price_object CHECK (jsonb_typeof(total_price_snapshot) = 'object'),
    CONSTRAINT itinerary_versions_activation_check CHECK (
        (status IN ('ACTIVE', 'SUPERSEDED') AND activated_at IS NOT NULL)
        OR (status IN ('DRAFT', 'ABANDONED'))
    ),
    CONSTRAINT itinerary_versions_journey_version_key UNIQUE (journey_id, version_number),
    CONSTRAINT itinerary_versions_id_journey_key UNIQUE (id, journey_id)
);

CREATE UNIQUE INDEX itinerary_versions_one_active_key
    ON itinerary_versions(journey_id) WHERE status = 'ACTIVE';
CREATE INDEX itinerary_versions_journey_id_idx ON itinerary_versions(journey_id, version_number);

ALTER TABLE journeys
    ADD CONSTRAINT journeys_current_itinerary_fk
    FOREIGN KEY (id, current_version_number)
    REFERENCES itinerary_versions(journey_id, version_number)
    DEFERRABLE INITIALLY DEFERRED;

CREATE TABLE itinerary_items (
    id uuid PRIMARY KEY,
    itinerary_version_id uuid NOT NULL REFERENCES itinerary_versions(id) ON DELETE CASCADE,
    reservation_id uuid REFERENCES reservations(id) ON DELETE RESTRICT,
    provider_id uuid REFERENCES providers(id) ON DELETE RESTRICT,
    item_type varchar(32) NOT NULL,
    sequence_number integer NOT NULL,
    external_reservation_id varchar(128),
    title varchar(200) NOT NULL,
    status varchar(16) NOT NULL,
    starts_at timestamptz NOT NULL,
    ends_at timestamptz NOT NULL,
    start_time_zone varchar(64) NOT NULL,
    end_time_zone varchar(64) NOT NULL,
    origin_code varchar(64),
    destination_code varchar(64),
    source_amount_minor bigint,
    source_currency char(3),
    display_amount_minor bigint,
    display_currency char(3),
    fx_rate_value numeric(24,12),
    fx_source varchar(64),
    fx_effective_at timestamptz,
    snapshot jsonb NOT NULL,
    synthetic boolean NOT NULL DEFAULT true,
    source varchar(16) NOT NULL DEFAULT 'MOCK',
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT itinerary_items_item_type_check CHECK (item_type IN (
        'FERRY_OUTBOUND', 'ARRIVAL_BUFFER', 'TRANSPORT_PICKUP', 'HOSPITAL_APPOINTMENT',
        'ADDITIONAL_CARE', 'HOTEL_STAY', 'TRANSPORT_DROPOFF', 'DEPARTURE_BUFFER', 'FERRY_RETURN'
    )),
    CONSTRAINT itinerary_items_sequence_check CHECK (sequence_number >= 1),
    CONSTRAINT itinerary_items_title_length CHECK (char_length(btrim(title)) BETWEEN 1 AND 200),
    CONSTRAINT itinerary_items_status_check CHECK (status IN ('CONFIRMED', 'BUFFER', 'SUPERSEDED')),
    CONSTRAINT itinerary_items_time_window_check CHECK (ends_at > starts_at),
    CONSTRAINT itinerary_items_time_zone_length CHECK (
        char_length(start_time_zone) BETWEEN 3 AND 64 AND char_length(end_time_zone) BETWEEN 3 AND 64
    ),
    CONSTRAINT itinerary_items_booking_group_check CHECK (
        (provider_id IS NULL AND reservation_id IS NULL AND external_reservation_id IS NULL AND status = 'BUFFER')
        OR (provider_id IS NOT NULL AND reservation_id IS NOT NULL AND external_reservation_id IS NOT NULL)
    ),
    CONSTRAINT itinerary_items_money_group_check CHECK (
        (source_amount_minor IS NULL AND source_currency IS NULL AND display_amount_minor IS NULL
            AND display_currency IS NULL AND fx_rate_value IS NULL AND fx_source IS NULL AND fx_effective_at IS NULL)
        OR (source_amount_minor IS NOT NULL AND source_amount_minor >= 0
            AND source_currency IS NOT NULL AND source_currency ~ '^[A-Z]{3}$'
            AND display_amount_minor IS NOT NULL AND display_amount_minor >= 0
            AND display_currency IS NOT NULL AND display_currency ~ '^[A-Z]{3}$'
            AND fx_rate_value IS NOT NULL AND fx_rate_value > 0
            AND fx_source IS NOT NULL AND fx_effective_at IS NOT NULL)
    ),
    CONSTRAINT itinerary_items_snapshot_object CHECK (jsonb_typeof(snapshot) = 'object'),
    CONSTRAINT itinerary_items_synthetic_check CHECK (synthetic),
    CONSTRAINT itinerary_items_source_check CHECK (source = 'MOCK'),
    CONSTRAINT itinerary_items_version_sequence_key UNIQUE (itinerary_version_id, sequence_number)
);

CREATE INDEX itinerary_items_version_id_idx ON itinerary_items(itinerary_version_id, sequence_number);
CREATE INDEX itinerary_items_provider_reservation_idx
    ON itinerary_items(provider_id, external_reservation_id) WHERE provider_id IS NOT NULL;

CREATE FUNCTION protect_itinerary_items() RETURNS trigger AS $$
DECLARE
    old_parent_status varchar(16);
    new_parent_status varchar(16);
BEGIN
    IF TG_OP IN ('UPDATE', 'DELETE') THEN
        SELECT status INTO old_parent_status
        FROM itinerary_versions
        WHERE id = OLD.itinerary_version_id;
        IF old_parent_status IN ('ACTIVE', 'SUPERSEDED') THEN
            RAISE EXCEPTION 'items in an activated itinerary version are immutable';
        END IF;
    END IF;
    IF TG_OP IN ('INSERT', 'UPDATE') THEN
        SELECT status INTO new_parent_status
        FROM itinerary_versions
        WHERE id = NEW.itinerary_version_id;
        IF new_parent_status IN ('ACTIVE', 'SUPERSEDED') THEN
            RAISE EXCEPTION 'items in an activated itinerary version are immutable';
        END IF;
    END IF;
    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER itinerary_items_immutable
    BEFORE INSERT OR UPDATE OR DELETE ON itinerary_items
    FOR EACH ROW EXECUTE FUNCTION protect_itinerary_items();

CREATE FUNCTION protect_itinerary_versions() RETURNS trigger AS $$
BEGIN
    IF TG_OP = 'DELETE' AND OLD.status IN ('ACTIVE', 'SUPERSEDED') THEN
        RAISE EXCEPTION 'activated itinerary versions cannot be deleted';
    END IF;
    IF TG_OP = 'UPDATE' AND OLD.status = 'SUPERSEDED' THEN
        RAISE EXCEPTION 'superseded itinerary versions are immutable';
    END IF;
    IF TG_OP = 'UPDATE' AND OLD.status = 'ACTIVE' AND (
        NEW.status <> 'SUPERSEDED'
        OR NEW.id IS DISTINCT FROM OLD.id
        OR NEW.journey_id IS DISTINCT FROM OLD.journey_id
        OR NEW.version_number IS DISTINCT FROM OLD.version_number
        OR NEW.change_reason IS DISTINCT FROM OLD.change_reason
        OR NEW.source_disruption_id IS DISTINCT FROM OLD.source_disruption_id
        OR NEW.total_price_snapshot IS DISTINCT FROM OLD.total_price_snapshot
        OR NEW.activated_at IS DISTINCT FROM OLD.activated_at
        OR NEW.created_at IS DISTINCT FROM OLD.created_at
    ) THEN
        RAISE EXCEPTION 'active itinerary versions may only transition to SUPERSEDED';
    END IF;
    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER itinerary_versions_immutable
    BEFORE UPDATE OR DELETE ON itinerary_versions
    FOR EACH ROW EXECUTE FUNCTION protect_itinerary_versions();

CREATE TABLE provider_events (
    id uuid PRIMARY KEY,
    provider_id uuid NOT NULL REFERENCES providers(id) ON DELETE RESTRICT,
    journey_id uuid NOT NULL REFERENCES journeys(id) ON DELETE RESTRICT,
    external_event_id varchar(128) NOT NULL,
    request_fingerprint char(64) NOT NULL,
    event_type varchar(64) NOT NULL,
    occurred_at timestamptz NOT NULL,
    target_snapshot jsonb NOT NULL,
    actor_snapshot jsonb NOT NULL,
    event_payload jsonb NOT NULL,
    assessment_outcome varchar(32),
    synthetic boolean NOT NULL DEFAULT true,
    source varchar(16) NOT NULL DEFAULT 'MOCK',
    received_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT provider_events_external_event_length CHECK (char_length(external_event_id) BETWEEN 3 AND 128),
    CONSTRAINT provider_events_request_fingerprint_format CHECK (request_fingerprint ~ '^[0-9a-f]{64}$'),
    CONSTRAINT provider_events_event_type_check CHECK (event_type IN (
        'FERRY_DELAYED', 'FERRY_CANCELLED', 'HOSPITAL_APPOINTMENT_CHANGED',
        'HOSPITAL_ADDITIONAL_CARE_REQUESTED', 'HOSPITAL_CLINICAL_HOLD',
        'HOTEL_ROOM_UNAVAILABLE', 'HOTEL_RESERVATION_PROBLEM',
        'TRANSPORT_DELAYED', 'TRANSPORT_UNAVAILABLE'
    )),
    CONSTRAINT provider_events_target_object CHECK (jsonb_typeof(target_snapshot) = 'object'),
    CONSTRAINT provider_events_actor_object CHECK (jsonb_typeof(actor_snapshot) = 'object'),
    CONSTRAINT provider_events_payload_object CHECK (jsonb_typeof(event_payload) = 'object'),
    CONSTRAINT provider_events_assessment_outcome_check CHECK (
        assessment_outcome IS NULL OR assessment_outcome IN ('NO_ACTION', 'DISRUPTION_CREATED')
    ),
    CONSTRAINT provider_events_received_after_occurred CHECK (received_at >= occurred_at),
    CONSTRAINT provider_events_synthetic_check CHECK (synthetic),
    CONSTRAINT provider_events_source_check CHECK (source = 'MOCK'),
    CONSTRAINT provider_events_provider_external_key UNIQUE (provider_id, external_event_id),
    CONSTRAINT provider_events_id_journey_key UNIQUE (id, journey_id)
);

CREATE INDEX provider_events_journey_id_idx ON provider_events(journey_id, received_at DESC);

CREATE FUNCTION enforce_provider_event_type() RETURNS trigger AS $$
DECLARE
    actual_type varchar(16);
BEGIN
    SELECT provider_type INTO actual_type FROM providers WHERE id = NEW.provider_id;
    IF actual_type IS NULL OR split_part(NEW.event_type, '_', 1) <> actual_type THEN
        RAISE EXCEPTION 'provider event type is incompatible with provider type';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER provider_events_type_compatible
    BEFORE INSERT OR UPDATE OF provider_id, event_type ON provider_events
    FOR EACH ROW EXECUTE FUNCTION enforce_provider_event_type();

CREATE TABLE disruptions (
    id uuid PRIMARY KEY,
    provider_event_id uuid NOT NULL UNIQUE,
    journey_id uuid NOT NULL,
    analyzed_itinerary_version_id uuid NOT NULL,
    status varchar(32) NOT NULL,
    impact_summary jsonb NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    resolved_at timestamptz,
    CONSTRAINT disruptions_event_journey_fk FOREIGN KEY (provider_event_id, journey_id)
        REFERENCES provider_events(id, journey_id) ON DELETE RESTRICT,
    CONSTRAINT disruptions_itinerary_journey_fk FOREIGN KEY (analyzed_itinerary_version_id, journey_id)
        REFERENCES itinerary_versions(id, journey_id) ON DELETE RESTRICT,
    CONSTRAINT disruptions_status_check CHECK (status IN (
        'DETECTED', 'VALIDATING', 'ANALYZING', 'CLINICAL_HOLD', 'MANUAL_REVIEW',
        'REPLAN_READY', 'AWAITING_APPROVAL', 'APPLYING', 'RESOLVED', 'RECOVERY_FAILED'
    )),
    CONSTRAINT disruptions_impact_object CHECK (jsonb_typeof(impact_summary) = 'object'),
    CONSTRAINT disruptions_updated_after_created CHECK (updated_at >= created_at),
    CONSTRAINT disruptions_resolved_check CHECK (resolved_at IS NULL OR resolved_at >= created_at)
);

CREATE INDEX disruptions_journey_id_idx ON disruptions(journey_id, created_at DESC);

ALTER TABLE itinerary_versions
    ADD CONSTRAINT itinerary_versions_source_disruption_fk
    FOREIGN KEY (source_disruption_id)
    REFERENCES disruptions(id) ON DELETE RESTRICT
    DEFERRABLE INITIALLY DEFERRED;

CREATE TABLE recovery_options (
    id uuid PRIMARY KEY,
    disruption_id uuid NOT NULL REFERENCES disruptions(id) ON DELETE CASCADE,
    analysis_revision integer NOT NULL,
    rank integer NOT NULL,
    status varchar(16) NOT NULL DEFAULT 'PROPOSED',
    explanation jsonb NOT NULL DEFAULT '[]'::jsonb,
    price_delta_amount_minor bigint NOT NULL,
    price_delta_currency char(3) NOT NULL,
    price_delta_estimated boolean NOT NULL DEFAULT true,
    time_delta_minutes integer NOT NULL,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT recovery_options_analysis_revision_check CHECK (analysis_revision >= 1),
    CONSTRAINT recovery_options_rank_check CHECK (rank BETWEEN 1 AND 2),
    CONSTRAINT recovery_options_status_check CHECK (status IN ('PROPOSED', 'APPROVED', 'APPLYING', 'APPLIED', 'EXPIRED', 'FAILED')),
    CONSTRAINT recovery_options_explanation_array CHECK (jsonb_typeof(explanation) = 'array'),
    CONSTRAINT recovery_options_price_currency_format CHECK (price_delta_currency ~ '^[A-Z]{3}$'),
    CONSTRAINT recovery_options_price_estimated_check CHECK (price_delta_estimated),
    CONSTRAINT recovery_options_expiry_check CHECK (expires_at > created_at),
    CONSTRAINT recovery_options_updated_after_created CHECK (updated_at >= created_at),
    CONSTRAINT recovery_options_disruption_revision_rank_key UNIQUE (disruption_id, analysis_revision, rank)
);

CREATE UNIQUE INDEX recovery_options_one_selection_per_revision_key
    ON recovery_options(disruption_id, analysis_revision)
    WHERE status IN ('APPROVED', 'APPLYING', 'APPLIED');
CREATE INDEX recovery_options_disruption_id_idx ON recovery_options(disruption_id, analysis_revision DESC);

CREATE TABLE recovery_items (
    id uuid PRIMARY KEY,
    recovery_option_id uuid NOT NULL REFERENCES recovery_options(id) ON DELETE CASCADE,
    old_itinerary_item_id uuid REFERENCES itinerary_items(id) ON DELETE RESTRICT,
    change_type varchar(16) NOT NULL,
    sequence_number integer NOT NULL,
    replacement_offer_snapshot jsonb,
    item_delta jsonb NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT recovery_items_change_type_check CHECK (change_type IN ('ADDED', 'CHANGED', 'REMOVED', 'UNCHANGED')),
    CONSTRAINT recovery_items_sequence_check CHECK (sequence_number >= 1),
    CONSTRAINT recovery_items_old_item_check CHECK (
        (change_type = 'ADDED' AND old_itinerary_item_id IS NULL)
        OR (change_type <> 'ADDED' AND old_itinerary_item_id IS NOT NULL)
    ),
    CONSTRAINT recovery_items_replacement_check CHECK (
        (change_type = 'REMOVED' AND replacement_offer_snapshot IS NULL)
        OR (change_type <> 'REMOVED' AND replacement_offer_snapshot IS NOT NULL)
    ),
    CONSTRAINT recovery_items_replacement_object CHECK (
        replacement_offer_snapshot IS NULL OR jsonb_typeof(replacement_offer_snapshot) = 'object'
    ),
    CONSTRAINT recovery_items_delta_object CHECK (jsonb_typeof(item_delta) = 'object'),
    CONSTRAINT recovery_items_option_sequence_key UNIQUE (recovery_option_id, sequence_number)
);

CREATE INDEX recovery_items_recovery_option_id_idx ON recovery_items(recovery_option_id, sequence_number);

CREATE TABLE idempotency_records (
    id uuid PRIMARY KEY,
    auth_scope varchar(160) NOT NULL,
    operation varchar(128) NOT NULL,
    idempotency_key varchar(128) NOT NULL,
    request_fingerprint char(64) NOT NULL,
    response_status integer NOT NULL,
    response_body jsonb,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT idempotency_records_auth_scope_length CHECK (char_length(auth_scope) BETWEEN 1 AND 160),
    CONSTRAINT idempotency_records_operation_length CHECK (char_length(operation) BETWEEN 1 AND 128),
    CONSTRAINT idempotency_records_key_format CHECK (
        idempotency_key ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{7,127}$'
    ),
    CONSTRAINT idempotency_records_fingerprint_format CHECK (request_fingerprint ~ '^[0-9a-f]{64}$'),
    CONSTRAINT idempotency_records_response_status_check CHECK (response_status BETWEEN 100 AND 599),
    CONSTRAINT idempotency_records_response_body_object CHECK (
        response_body IS NULL OR jsonb_typeof(response_body) IN ('object', 'array')
    ),
    CONSTRAINT idempotency_records_expiry_check CHECK (expires_at > created_at),
    CONSTRAINT idempotency_records_scope_operation_key UNIQUE (auth_scope, operation, idempotency_key)
);

CREATE INDEX idempotency_records_expires_at_idx ON idempotency_records(expires_at);
