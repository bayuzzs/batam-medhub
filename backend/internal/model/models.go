package model

import (
	"encoding/json"
	"time"
)

type Patient struct {
	ID                string    `gorm:"column:id;type:uuid;primaryKey"`
	EmailNormalized   string    `gorm:"column:email_normalized;size:254;not null"`
	PasswordHash      string    `gorm:"column:password_hash;type:char(60);not null"`
	FullName          string    `gorm:"column:full_name;size:120;not null"`
	PreferredCurrency string    `gorm:"column:preferred_currency;type:char(3);not null"`
	Status            string    `gorm:"column:status;size:32;not null"`
	Synthetic         bool      `gorm:"column:synthetic;not null"`
	CreatedAt         time.Time `gorm:"column:created_at;not null"`
	UpdatedAt         time.Time `gorm:"column:updated_at;not null"`
}

type AuthSession struct {
	ID                  string     `gorm:"column:id;type:uuid;primaryKey"`
	PatientID           string     `gorm:"column:patient_id;type:uuid;not null"`
	RefreshTokenHash    string     `gorm:"column:refresh_token_hash;type:char(64);not null"`
	ReplacedBySessionID *string    `gorm:"column:replaced_by_session_id;type:uuid"`
	ExpiresAt           time.Time  `gorm:"column:expires_at;not null"`
	RevokedAt           *time.Time `gorm:"column:revoked_at"`
	CreatedAt           time.Time  `gorm:"column:created_at;not null"`
	LastUsedAt          *time.Time `gorm:"column:last_used_at"`
}

type Provider struct {
	ID           string    `gorm:"column:id;type:uuid;primaryKey"`
	ProviderType string    `gorm:"column:provider_type;size:16;not null"`
	Code         string    `gorm:"column:code;size:64;not null"`
	DisplayName  string    `gorm:"column:display_name;size:160;not null"`
	Status       string    `gorm:"column:status;size:32;not null"`
	Synthetic    bool      `gorm:"column:synthetic;not null"`
	Source       string    `gorm:"column:source;size:16;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;not null"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null"`
}

type ProviderCredential struct {
	ID            string     `gorm:"column:id;type:uuid;primaryKey"`
	ProviderID    string     `gorm:"column:provider_id;type:uuid;not null"`
	KeyPrefix     string     `gorm:"column:key_prefix;size:32;not null"`
	SecretHash    string     `gorm:"column:secret_hash;type:char(64);not null"`
	HashAlgorithm string     `gorm:"column:hash_algorithm;size:16;not null"`
	Status        string     `gorm:"column:status;size:32;not null"`
	ExpiresAt     *time.Time `gorm:"column:expires_at"`
	Synthetic     bool       `gorm:"column:synthetic;not null"`
	CreatedAt     time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;not null"`
}

type MedicalService struct {
	ID                     string    `gorm:"column:id;type:uuid;primaryKey"`
	Code                   string    `gorm:"column:code;size:64;not null"`
	Name                   string    `gorm:"column:name;size:160;not null"`
	Category               string    `gorm:"column:category;size:64;not null"`
	Description            *string   `gorm:"column:description;type:text"`
	DefaultDurationMinutes int       `gorm:"column:default_duration_minutes;not null"`
	Active                 bool      `gorm:"column:active;not null"`
	Synthetic              bool      `gorm:"column:synthetic;not null"`
	Source                 string    `gorm:"column:source;size:16;not null"`
	CreatedAt              time.Time `gorm:"column:created_at;not null"`
	UpdatedAt              time.Time `gorm:"column:updated_at;not null"`
}

type ProviderCapability struct {
	ID                string    `gorm:"column:id;type:uuid;primaryKey"`
	ProviderID        string    `gorm:"column:provider_id;type:uuid;not null"`
	MedicalServiceID  string    `gorm:"column:medical_service_id;type:uuid;not null"`
	ExternalServiceID string    `gorm:"column:external_service_id;size:128;not null"`
	Active            bool      `gorm:"column:active;not null"`
	Synthetic         bool      `gorm:"column:synthetic;not null"`
	Source            string    `gorm:"column:source;size:16;not null"`
	CreatedAt         time.Time `gorm:"column:created_at;not null"`
	UpdatedAt         time.Time `gorm:"column:updated_at;not null"`
}

type FXRate struct {
	ID            string    `gorm:"column:id;type:uuid;primaryKey"`
	BaseCurrency  string    `gorm:"column:base_currency;type:char(3);not null"`
	QuoteCurrency string    `gorm:"column:quote_currency;type:char(3);not null"`
	Rate          string    `gorm:"column:rate;type:numeric(24,12);not null"`
	Source        string    `gorm:"column:source;size:64;not null"`
	EffectiveAt   time.Time `gorm:"column:effective_at;not null"`
	Estimated     bool      `gorm:"column:estimated;not null"`
	Synthetic     bool      `gorm:"column:synthetic;not null"`
	CreatedAt     time.Time `gorm:"column:created_at;not null"`
}

type TripRequest struct {
	ID                   string          `gorm:"column:id;type:uuid;primaryKey"`
	PatientID            string          `gorm:"column:patient_id;type:uuid;not null"`
	Status               string          `gorm:"column:status;size:32;not null"`
	MedicalServiceID     *string         `gorm:"column:medical_service_id;type:uuid"`
	RequestedServiceText *string         `gorm:"column:requested_service_text;size:300"`
	StructuredIntent     json.RawMessage `gorm:"column:structured_intent;type:jsonb"`
	ReferenceCurrency    string          `gorm:"column:reference_currency;type:char(3);not null"`
	PlanningRevision     int             `gorm:"column:planning_revision;not null"`
	SelectedPlanOptionID *string         `gorm:"column:selected_plan_option_id;type:uuid"`
	CreatedAt            time.Time       `gorm:"column:created_at;not null"`
	UpdatedAt            time.Time       `gorm:"column:updated_at;not null"`
}

type PlanOption struct {
	ID                 string          `gorm:"column:id;type:uuid;primaryKey"`
	TripRequestID      string          `gorm:"column:trip_request_id;type:uuid;not null"`
	PlanningRevision   int             `gorm:"column:planning_revision;not null"`
	Rank               int             `gorm:"column:rank;not null"`
	Status             string          `gorm:"column:status;size:16;not null"`
	Explanation        json.RawMessage `gorm:"column:explanation;type:jsonb;not null"`
	TotalPriceSnapshot json.RawMessage `gorm:"column:total_price_snapshot;type:jsonb;not null"`
	ExpiresAt          time.Time       `gorm:"column:expires_at;not null"`
	CreatedAt          time.Time       `gorm:"column:created_at;not null"`
	UpdatedAt          time.Time       `gorm:"column:updated_at;not null"`
}

type PlanItem struct {
	ID                 string          `gorm:"column:id;type:uuid;primaryKey"`
	PlanOptionID       string          `gorm:"column:plan_option_id;type:uuid;not null"`
	ProviderID         *string         `gorm:"column:provider_id;type:uuid"`
	ItemType           string          `gorm:"column:item_type;size:32;not null"`
	SequenceNumber     int             `gorm:"column:sequence_number;not null"`
	ExternalOfferID    *string         `gorm:"column:external_offer_id;size:128"`
	Title              string          `gorm:"column:title;size:200;not null"`
	StartsAt           time.Time       `gorm:"column:starts_at;not null"`
	EndsAt             time.Time       `gorm:"column:ends_at;not null"`
	StartTimeZone      string          `gorm:"column:start_time_zone;size:64;not null"`
	EndTimeZone        string          `gorm:"column:end_time_zone;size:64;not null"`
	OriginCode         *string         `gorm:"column:origin_code;size:64"`
	DestinationCode    *string         `gorm:"column:destination_code;size:64"`
	SourceAmountMinor  *int64          `gorm:"column:source_amount_minor"`
	SourceCurrency     *string         `gorm:"column:source_currency;type:char(3)"`
	DisplayAmountMinor *int64          `gorm:"column:display_amount_minor"`
	DisplayCurrency    *string         `gorm:"column:display_currency;type:char(3)"`
	FXRateID           *string         `gorm:"column:fx_rate_id;type:uuid"`
	OfferSnapshot      json.RawMessage `gorm:"column:offer_snapshot;type:jsonb;not null"`
	OfferExpiresAt     *time.Time      `gorm:"column:offer_expires_at"`
	OperationalNotes   json.RawMessage `gorm:"column:operational_notes;type:jsonb;not null"`
	Synthetic          bool            `gorm:"column:synthetic;not null"`
	Source             string          `gorm:"column:source;size:16;not null"`
	CreatedAt          time.Time       `gorm:"column:created_at;not null"`
}

type Journey struct {
	ID                   string    `gorm:"column:id;type:uuid;primaryKey"`
	TripRequestID        string    `gorm:"column:trip_request_id;type:uuid;not null"`
	PatientID            string    `gorm:"column:patient_id;type:uuid;not null"`
	Status               string    `gorm:"column:status;size:32;not null"`
	CurrentVersionNumber int       `gorm:"column:current_version_number;not null"`
	ActivatedAt          time.Time `gorm:"column:activated_at;not null"`
	CreatedAt            time.Time `gorm:"column:created_at;not null"`
	UpdatedAt            time.Time `gorm:"column:updated_at;not null"`
}

type Reservation struct {
	ID                    string          `gorm:"column:id;type:uuid;primaryKey"`
	TripRequestID         string          `gorm:"column:trip_request_id;type:uuid;not null"`
	JourneyID             *string         `gorm:"column:journey_id;type:uuid"`
	PlanItemID            *string         `gorm:"column:plan_item_id;type:uuid"`
	ProviderID            string          `gorm:"column:provider_id;type:uuid;not null"`
	Status                string          `gorm:"column:status;size:16;not null"`
	ExternalOfferID       string          `gorm:"column:external_offer_id;size:128;not null"`
	ExternalHoldID        *string         `gorm:"column:external_hold_id;size:128"`
	ExternalReservationID *string         `gorm:"column:external_reservation_id;size:128"`
	HoldExpiresAt         *time.Time      `gorm:"column:hold_expires_at"`
	ProviderSnapshot      json.RawMessage `gorm:"column:provider_snapshot;type:jsonb;not null"`
	CleanupRequired       bool            `gorm:"column:cleanup_required;not null"`
	CreatedAt             time.Time       `gorm:"column:created_at;not null"`
	UpdatedAt             time.Time       `gorm:"column:updated_at;not null"`
}

type ItineraryVersion struct {
	ID                 string          `gorm:"column:id;type:uuid;primaryKey"`
	JourneyID          string          `gorm:"column:journey_id;type:uuid;not null"`
	VersionNumber      int             `gorm:"column:version_number;not null"`
	Status             string          `gorm:"column:status;size:16;not null"`
	ChangeReason       string          `gorm:"column:change_reason;size:500;not null"`
	SourceDisruptionID *string         `gorm:"column:source_disruption_id;type:uuid"`
	TotalPriceSnapshot json.RawMessage `gorm:"column:total_price_snapshot;type:jsonb;not null"`
	ActivatedAt        *time.Time      `gorm:"column:activated_at"`
	CreatedAt          time.Time       `gorm:"column:created_at;not null"`
}

type ItineraryItem struct {
	ID                    string          `gorm:"column:id;type:uuid;primaryKey"`
	ItineraryVersionID    string          `gorm:"column:itinerary_version_id;type:uuid;not null"`
	ReservationID         *string         `gorm:"column:reservation_id;type:uuid"`
	ProviderID            *string         `gorm:"column:provider_id;type:uuid"`
	ItemType              string          `gorm:"column:item_type;size:32;not null"`
	SequenceNumber        int             `gorm:"column:sequence_number;not null"`
	ExternalReservationID *string         `gorm:"column:external_reservation_id;size:128"`
	Title                 string          `gorm:"column:title;size:200;not null"`
	Status                string          `gorm:"column:status;size:16;not null"`
	StartsAt              time.Time       `gorm:"column:starts_at;not null"`
	EndsAt                time.Time       `gorm:"column:ends_at;not null"`
	StartTimeZone         string          `gorm:"column:start_time_zone;size:64;not null"`
	EndTimeZone           string          `gorm:"column:end_time_zone;size:64;not null"`
	OriginCode            *string         `gorm:"column:origin_code;size:64"`
	DestinationCode       *string         `gorm:"column:destination_code;size:64"`
	SourceAmountMinor     *int64          `gorm:"column:source_amount_minor"`
	SourceCurrency        *string         `gorm:"column:source_currency;type:char(3)"`
	DisplayAmountMinor    *int64          `gorm:"column:display_amount_minor"`
	DisplayCurrency       *string         `gorm:"column:display_currency;type:char(3)"`
	FXRateValue           *string         `gorm:"column:fx_rate_value;type:numeric(24,12)"`
	FXSource              *string         `gorm:"column:fx_source;size:64"`
	FXEffectiveAt         *time.Time      `gorm:"column:fx_effective_at"`
	Snapshot              json.RawMessage `gorm:"column:snapshot;type:jsonb;not null"`
	Synthetic             bool            `gorm:"column:synthetic;not null"`
	Source                string          `gorm:"column:source;size:16;not null"`
	CreatedAt             time.Time       `gorm:"column:created_at;not null"`
}

type ProviderEvent struct {
	ID                 string          `gorm:"column:id;type:uuid;primaryKey"`
	ProviderID         string          `gorm:"column:provider_id;type:uuid;not null"`
	JourneyID          string          `gorm:"column:journey_id;type:uuid;not null"`
	ExternalEventID    string          `gorm:"column:external_event_id;size:128;not null"`
	RequestFingerprint string          `gorm:"column:request_fingerprint;type:char(64);not null"`
	EventType          string          `gorm:"column:event_type;size:64;not null"`
	OccurredAt         time.Time       `gorm:"column:occurred_at;not null"`
	TargetSnapshot     json.RawMessage `gorm:"column:target_snapshot;type:jsonb;not null"`
	ActorSnapshot      json.RawMessage `gorm:"column:actor_snapshot;type:jsonb;not null"`
	EventPayload       json.RawMessage `gorm:"column:event_payload;type:jsonb;not null"`
	AssessmentOutcome  *string         `gorm:"column:assessment_outcome;size:32"`
	Synthetic          bool            `gorm:"column:synthetic;not null"`
	Source             string          `gorm:"column:source;size:16;not null"`
	ReceivedAt         time.Time       `gorm:"column:received_at;not null"`
}

type Disruption struct {
	ID                         string          `gorm:"column:id;type:uuid;primaryKey"`
	ProviderEventID            string          `gorm:"column:provider_event_id;type:uuid;not null"`
	JourneyID                  string          `gorm:"column:journey_id;type:uuid;not null"`
	AnalyzedItineraryVersionID string          `gorm:"column:analyzed_itinerary_version_id;type:uuid;not null"`
	Status                     string          `gorm:"column:status;size:32;not null"`
	ImpactSummary              json.RawMessage `gorm:"column:impact_summary;type:jsonb;not null"`
	CreatedAt                  time.Time       `gorm:"column:created_at;not null"`
	UpdatedAt                  time.Time       `gorm:"column:updated_at;not null"`
	ResolvedAt                 *time.Time      `gorm:"column:resolved_at"`
}

type RecoveryOption struct {
	ID                    string          `gorm:"column:id;type:uuid;primaryKey"`
	DisruptionID          string          `gorm:"column:disruption_id;type:uuid;not null"`
	AnalysisRevision      int             `gorm:"column:analysis_revision;not null"`
	Rank                  int             `gorm:"column:rank;not null"`
	Status                string          `gorm:"column:status;size:16;not null"`
	Explanation           json.RawMessage `gorm:"column:explanation;type:jsonb;not null"`
	PriceDeltaAmountMinor int64           `gorm:"column:price_delta_amount_minor;not null"`
	PriceDeltaCurrency    string          `gorm:"column:price_delta_currency;type:char(3);not null"`
	PriceDeltaEstimated   bool            `gorm:"column:price_delta_estimated;not null"`
	TimeDeltaMinutes      int             `gorm:"column:time_delta_minutes;not null"`
	ExpiresAt             time.Time       `gorm:"column:expires_at;not null"`
	CreatedAt             time.Time       `gorm:"column:created_at;not null"`
	UpdatedAt             time.Time       `gorm:"column:updated_at;not null"`
}

type RecoveryItem struct {
	ID                       string          `gorm:"column:id;type:uuid;primaryKey"`
	RecoveryOptionID         string          `gorm:"column:recovery_option_id;type:uuid;not null"`
	OldItineraryItemID       *string         `gorm:"column:old_itinerary_item_id;type:uuid"`
	ChangeType               string          `gorm:"column:change_type;size:16;not null"`
	SequenceNumber           int             `gorm:"column:sequence_number;not null"`
	ReplacementOfferSnapshot json.RawMessage `gorm:"column:replacement_offer_snapshot;type:jsonb"`
	ItemDelta                json.RawMessage `gorm:"column:item_delta;type:jsonb;not null"`
	CreatedAt                time.Time       `gorm:"column:created_at;not null"`
}

type IdempotencyRecord struct {
	ID                 string          `gorm:"column:id;type:uuid;primaryKey"`
	AuthScope          string          `gorm:"column:auth_scope;size:160;not null"`
	Operation          string          `gorm:"column:operation;size:128;not null"`
	IdempotencyKey     string          `gorm:"column:idempotency_key;size:128;not null"`
	RequestFingerprint string          `gorm:"column:request_fingerprint;type:char(64);not null"`
	ResponseStatus     int             `gorm:"column:response_status;not null"`
	ResponseBody       json.RawMessage `gorm:"column:response_body;type:jsonb"`
	ExpiresAt          time.Time       `gorm:"column:expires_at;not null"`
	CreatedAt          time.Time       `gorm:"column:created_at;not null"`
}
