DROP TABLE IF EXISTS idempotency_records;
DROP TABLE IF EXISTS recovery_items;
DROP TABLE IF EXISTS recovery_options;

ALTER TABLE IF EXISTS itinerary_versions
    DROP CONSTRAINT IF EXISTS itinerary_versions_source_disruption_fk;
DROP TABLE IF EXISTS disruptions;
DROP TABLE IF EXISTS provider_events;
DROP FUNCTION IF EXISTS enforce_provider_event_type();

DROP TABLE IF EXISTS itinerary_items;
DROP FUNCTION IF EXISTS protect_itinerary_items();

ALTER TABLE IF EXISTS journeys
    DROP CONSTRAINT IF EXISTS journeys_current_itinerary_fk;
DROP TABLE IF EXISTS itinerary_versions;
DROP FUNCTION IF EXISTS protect_itinerary_versions();

DROP TABLE IF EXISTS reservations;
DROP TABLE IF EXISTS journeys;
DROP TABLE IF EXISTS plan_items;

ALTER TABLE IF EXISTS trip_requests
    DROP CONSTRAINT IF EXISTS trip_requests_selected_plan_option_fk;
DROP TABLE IF EXISTS plan_options;
DROP TABLE IF EXISTS trip_requests;
DROP TABLE IF EXISTS fx_rates;

DROP TABLE IF EXISTS provider_capabilities;
DROP FUNCTION IF EXISTS enforce_hospital_provider_capability();
DROP TABLE IF EXISTS medical_services;
DROP TABLE IF EXISTS provider_credentials;
DROP TABLE IF EXISTS providers;
DROP TABLE IF EXISTS auth_sessions;
DROP TABLE IF EXISTS patients;
