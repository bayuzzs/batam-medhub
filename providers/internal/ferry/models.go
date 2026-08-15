package ferry

import (
	"time"
)

// Sailing represents a scheduled ferry sailing in the database.
type Sailing struct {
	ID                  string    `gorm:"type:uuid;primaryKey" json:"id"`
	SailingCode         string    `gorm:"type:varchar(64);not null;uniqueIndex" json:"sailing_code"`
	ExternalSailingID   string    `gorm:"type:varchar(128);not null;uniqueIndex" json:"external_sailing_id"`
	OfferID             string    `gorm:"type:varchar(128);not null;uniqueIndex" json:"offer_id"`
	OperatorName        string    `gorm:"type:varchar(160);not null" json:"operator_name"`
	OriginPortCode      string    `gorm:"type:varchar(64);not null;index:sailings_search_idx,priority:1" json:"origin_port_code"`
	DestinationPortCode string    `gorm:"type:varchar(64);not null;index:sailings_search_idx,priority:2" json:"destination_port_code"`
	DepartsAt           time.Time `gorm:"type:timestamptz;not null;index:sailings_search_idx,priority:4" json:"departs_at"`
	ArrivesAt           time.Time `gorm:"type:timestamptz;not null" json:"arrives_at"`
	DepartureTimeZone   string    `gorm:"type:varchar(64);not null" json:"departure_time_zone"`
	ArrivalTimeZone     string    `gorm:"type:varchar(64);not null" json:"arrival_time_zone"`
	CheckInCutoffAt     time.Time `gorm:"type:timestamptz;not null" json:"check_in_cutoff_at"`
	SeatCapacity        int       `gorm:"type:integer;not null" json:"seat_capacity"`
	PriceAmountMinor    int64     `gorm:"type:bigint;not null" json:"price_amount_minor"`
	PriceCurrency       string    `gorm:"type:char(3);not null" json:"price_currency"`
	ValidUntil          time.Time `gorm:"type:timestamptz;not null;index:sailings_valid_until_idx" json:"valid_until"`
	Status              string    `gorm:"type:varchar(16);not null;default:'AVAILABLE';index:sailings_search_idx,priority:3" json:"status"`
	Synthetic           bool      `gorm:"type:boolean;not null;default:true" json:"synthetic"`
	Source              string    `gorm:"type:varchar(16);not null;default:'MOCK'" json:"source"`
	CreatedAt           time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt           time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Sailing) TableName() string {
	return "sailings"
}

// Hold represents a temporary capacity reservation on a ferry sailing.
type Hold struct {
	ID                    string     `gorm:"type:uuid;primaryKey" json:"id"`
	ExternalHoldID        string     `gorm:"type:varchar(128);not null;uniqueIndex" json:"external_hold_id"`
	ExternalReference     string     `gorm:"type:varchar(128);not null;uniqueIndex" json:"external_reference"`
	SailingID             string     `gorm:"type:uuid;not null;index:holds_capacity_idx,priority:1" json:"sailing_id"`
	Sailing               *Sailing   `gorm:"foreignKey:SailingID;references:ID" json:"sailing,omitempty"`
	OfferID               string     `gorm:"type:varchar(128);not null" json:"offer_id"`
	ClientReference       string     `gorm:"type:varchar(128);not null;uniqueIndex" json:"client_reference"`
	PassengerCount        int        `gorm:"type:integer;not null" json:"passenger_count"`
	UnitPriceAmountMinor  int64      `gorm:"type:bigint;not null" json:"unit_price_amount_minor"`
	TotalPriceAmountMinor int64      `gorm:"type:bigint;not null" json:"total_price_amount_minor"`
	PriceCurrency         string     `gorm:"type:char(3);not null" json:"price_currency"`
	ServiceStartsAt       time.Time  `gorm:"type:timestamptz;not null" json:"service_starts_at"`
	ServiceEndsAt         time.Time  `gorm:"type:timestamptz;not null" json:"service_ends_at"`
	StartTimeZone         string     `gorm:"type:varchar(64);not null" json:"start_time_zone"`
	EndTimeZone           string     `gorm:"type:varchar(64);not null" json:"end_time_zone"`
	Status                string     `gorm:"type:varchar(16);not null;index:holds_capacity_idx,priority:2" json:"status"`
	CreatedAt             time.Time  `gorm:"type:timestamptz;not null" json:"created_at"`
	ExpiresAt             time.Time  `gorm:"type:timestamptz;not null;index:holds_capacity_idx,priority:3" json:"expires_at"`
	ConfirmedAt           *time.Time `gorm:"type:timestamptz" json:"confirmed_at,omitempty"`
	ReleasedAt            *time.Time `gorm:"type:timestamptz" json:"released_at,omitempty"`
	ExpiredAt             *time.Time `gorm:"type:timestamptz" json:"expired_at,omitempty"`
	UpdatedAt             time.Time  `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Hold) TableName() string {
	return "holds"
}

// Reservation represents a confirmed booking on a ferry sailing.
type Reservation struct {
	ID                    string     `gorm:"type:uuid;primaryKey" json:"id"`
	ExternalReservationID string     `gorm:"type:varchar(128);not null;uniqueIndex" json:"external_reservation_id"`
	ExternalReference     string     `gorm:"type:varchar(128);not null;uniqueIndex" json:"external_reference"`
	HoldID                string     `gorm:"type:uuid;not null;uniqueIndex" json:"hold_id"`
	Hold                  *Hold      `gorm:"foreignKey:HoldID;references:ID" json:"hold,omitempty"`
	SailingID             string     `gorm:"type:uuid;not null;index:reservations_sailing_status_idx,priority:1" json:"sailing_id"`
	Sailing               *Sailing   `gorm:"foreignKey:SailingID;references:ID" json:"sailing,omitempty"`
	OfferID               string     `gorm:"type:varchar(128);not null" json:"offer_id"`
	ClientReference       string     `gorm:"type:varchar(128);not null;index:reservations_client_reference_idx" json:"client_reference"`
	PassengerCount        int        `gorm:"type:integer;not null" json:"passenger_count"`
	TotalPriceAmountMinor int64      `gorm:"type:bigint;not null" json:"total_price_amount_minor"`
	PriceCurrency         string     `gorm:"type:char(3);not null" json:"price_currency"`
	ServiceStartsAt       time.Time  `gorm:"type:timestamptz;not null" json:"service_starts_at"`
	ServiceEndsAt         time.Time  `gorm:"type:timestamptz;not null" json:"service_ends_at"`
	StartTimeZone         string     `gorm:"type:varchar(64);not null" json:"start_time_zone"`
	EndTimeZone           string     `gorm:"type:varchar(64);not null" json:"end_time_zone"`
	Status                string     `gorm:"type:varchar(16);not null;index:reservations_sailing_status_idx,priority:2" json:"status"`
	ConfirmedAt           time.Time  `gorm:"type:timestamptz;not null" json:"confirmed_at"`
	ReleasedAt            *time.Time `gorm:"type:timestamptz" json:"released_at,omitempty"`
	CreatedAt             time.Time  `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt             time.Time  `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Reservation) TableName() string {
	return "reservations"
}

// IdempotencyRecord stores request hashes and cached response bodies.
type IdempotencyRecord struct {
	ID                 string    `gorm:"type:uuid;primaryKey" json:"id"`
	ClientScope        string    `gorm:"type:varchar(128);not null;uniqueIndex:idx_ferry_idemp_scope_key,priority:1" json:"client_scope"`
	Operation          string    `gorm:"type:varchar(160);not null;uniqueIndex:idx_ferry_idemp_scope_key,priority:2" json:"operation"`
	IdempotencyKey     string    `gorm:"type:varchar(128);not null;uniqueIndex:idx_ferry_idemp_scope_key,priority:3" json:"idempotency_key"`
	RequestFingerprint string    `gorm:"type:char(64);not null" json:"request_fingerprint"`
	ResponseStatus     int       `gorm:"type:integer;not null" json:"response_status"`
	ResponseBody       []byte    `gorm:"type:jsonb;not null" json:"response_body"`
	CreatedAt          time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (IdempotencyRecord) TableName() string {
	return "idempotency_records"
}
