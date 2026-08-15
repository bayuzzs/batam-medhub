package transport

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// StringArray is a helper for PostgreSQL JSONB string arrays.
type StringArray []string

func (a *StringArray) Scan(value any) error {
	if value == nil {
		*a = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		str, ok := value.(string)
		if !ok {
			return errors.New("cannot scan type into StringArray")
		}
		bytes = []byte(str)
	}
	return json.Unmarshal(bytes, a)
}

func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return "[]", nil
	}
	return json.Marshal(a)
}

type Vehicle struct {
	ID                string      `gorm:"type:uuid;primaryKey" json:"id"`
	VehicleCode       string      `gorm:"type:varchar(64);uniqueIndex;not null" json:"vehicle_code"`
	VehicleType       string      `gorm:"type:varchar(80);not null" json:"vehicle_type"`
	PassengerCapacity int         `gorm:"not null" json:"passenger_capacity"`
	Accessibility     StringArray `gorm:"type:jsonb;not null;default:'[]'" json:"accessibility"`
	Status            string      `gorm:"type:varchar(16);not null;default:'ACTIVE'" json:"status"`
	Synthetic         bool        `gorm:"not null;default:true" json:"synthetic"`
	Source            string      `gorm:"type:varchar(16);not null;default:'MOCK'" json:"source"`
	CreatedAt         time.Time   `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time   `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Vehicle) TableName() string {
	return "vehicles"
}

type TransportOffer struct {
	ID                     string             `gorm:"type:uuid;primaryKey" json:"id"`
	ExternalOfferID        string             `gorm:"type:varchar(128);uniqueIndex;not null" json:"external_offer_id"`
	ExternalAvailabilityID string             `gorm:"type:varchar(128);uniqueIndex;not null" json:"external_availability_id"`
	VehicleType            string             `gorm:"type:varchar(80);not null" json:"vehicle_type"`
	PickupLocationCode     string             `gorm:"type:varchar(64);not null" json:"pickup_location_code"`
	DropoffLocationCode    string             `gorm:"type:varchar(64);not null" json:"dropoff_location_code"`
	PassengerCapacity      int                `gorm:"not null" json:"passenger_capacity"`
	Accessibility          StringArray        `gorm:"type:jsonb;not null;default:'[]'" json:"accessibility"`
	ServiceStartsAt        time.Time          `gorm:"type:timestamptz;not null" json:"service_starts_at"`
	ServiceEndsAt          time.Time          `gorm:"type:timestamptz;not null" json:"service_ends_at"`
	StartTimeZone          string             `gorm:"type:varchar(64);not null" json:"start_time_zone"`
	EndTimeZone            string             `gorm:"type:varchar(64);not null" json:"end_time_zone"`
	PriceAmountMinor       int64              `gorm:"not null" json:"price_amount_minor"`
	PriceCurrency          string             `gorm:"type:char(3);not null" json:"price_currency"`
	ValidUntil             time.Time          `gorm:"type:timestamptz;not null" json:"valid_until"`
	Status                 string             `gorm:"type:varchar(16);not null;default:'AVAILABLE'" json:"status"`
	Synthetic              bool               `gorm:"not null;default:true" json:"synthetic"`
	Source                 string             `gorm:"type:varchar(16);not null;default:'MOCK'" json:"source"`
	AvailabilitySlots      []AvailabilitySlot `gorm:"foreignKey:OfferPK;references:ID" json:"availability_slots,omitempty"`
	CreatedAt              time.Time          `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt              time.Time          `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (TransportOffer) TableName() string {
	return "transport_offers"
}

type AvailabilitySlot struct {
	ID        string          `gorm:"type:uuid;primaryKey" json:"id"`
	SlotCode  string          `gorm:"type:varchar(128);uniqueIndex;not null" json:"slot_code"`
	OfferPK   string          `gorm:"type:uuid;not null;uniqueIndex:idx_offer_vehicle" json:"offer_pk"`
	VehicleID string          `gorm:"type:uuid;not null;uniqueIndex:idx_offer_vehicle" json:"vehicle_id"`
	Vehicle   *Vehicle        `gorm:"foreignKey:VehicleID;references:ID" json:"vehicle,omitempty"`
	StartsAt  time.Time       `gorm:"type:timestamptz;not null" json:"starts_at"`
	EndsAt    time.Time       `gorm:"type:timestamptz;not null" json:"ends_at"`
	TimeZone  string          `gorm:"type:varchar(64);not null" json:"time_zone"`
	Status    string          `gorm:"type:varchar(16);not null;default:'AVAILABLE'" json:"status"`
	CreatedAt time.Time       `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time       `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (AvailabilitySlot) TableName() string {
	return "availability_slots"
}

type Hold struct {
	ID                      string      `gorm:"type:uuid;primaryKey" json:"id"`
	ExternalHoldID          string      `gorm:"type:varchar(128);uniqueIndex;not null" json:"external_hold_id"`
	ExternalReference       string      `gorm:"type:varchar(128);uniqueIndex;not null" json:"external_reference"`
	OfferPK                 string      `gorm:"type:uuid;not null" json:"offer_pk"`
	OfferID                 string      `gorm:"type:varchar(128);not null" json:"offer_id"`
	ClientReference         string      `gorm:"type:varchar(128);uniqueIndex;not null" json:"client_reference"`
	Units                   int         `gorm:"not null" json:"units"`
	PassengerCount          int         `gorm:"not null" json:"passenger_count"`
	PickupLocationCode      string      `gorm:"type:varchar(64);not null" json:"pickup_location_code"`
	DropoffLocationCode     string      `gorm:"type:varchar(64);not null" json:"dropoff_location_code"`
	RequestedPickupStartsAt time.Time   `gorm:"type:timestamptz;not null" json:"requested_pickup_starts_at"`
	RequestedPickupEndsAt   time.Time   `gorm:"type:timestamptz;not null" json:"requested_pickup_ends_at"`
	RequestedStartTimeZone  string      `gorm:"type:varchar(64);not null" json:"requested_start_time_zone"`
	RequestedEndTimeZone    string      `gorm:"type:varchar(64);not null" json:"requested_end_time_zone"`
	Accessibility           StringArray `gorm:"type:jsonb;not null;default:'[]'" json:"accessibility"`
	UnitPriceAmountMinor    int64       `gorm:"not null" json:"unit_price_amount_minor"`
	TotalPriceAmountMinor   int64       `gorm:"not null" json:"total_price_amount_minor"`
	PriceCurrency           string      `gorm:"type:char(3);not null" json:"price_currency"`
	ServiceStartsAt         time.Time   `gorm:"type:timestamptz;not null" json:"service_starts_at"`
	ServiceEndsAt           time.Time   `gorm:"type:timestamptz;not null" json:"service_ends_at"`
	StartTimeZone           string      `gorm:"type:varchar(64);not null" json:"start_time_zone"`
	EndTimeZone             string      `gorm:"type:varchar(64);not null" json:"end_time_zone"`
	Status                  string      `gorm:"type:varchar(16);not null;default:'HELD'" json:"status"`
	CreatedAt               time.Time   `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	ExpiresAt               time.Time   `gorm:"type:timestamptz;not null" json:"expires_at"`
	ConfirmedAt             *time.Time  `gorm:"type:timestamptz" json:"confirmed_at,omitempty"`
	ReleasedAt              *time.Time  `gorm:"type:timestamptz" json:"released_at,omitempty"`
	ExpiredAt               *time.Time  `gorm:"type:timestamptz" json:"expired_at,omitempty"`
	UpdatedAt               time.Time   `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Hold) TableName() string {
	return "holds"
}

type HoldAssignment struct {
	HoldID             string `gorm:"type:uuid;primaryKey" json:"hold_id"`
	OfferPK            string `gorm:"type:uuid;not null" json:"offer_pk"`
	AvailabilitySlotID string `gorm:"type:uuid;primaryKey" json:"availability_slot_id"`
}

func (HoldAssignment) TableName() string {
	return "hold_assignments"
}

type Reservation struct {
	ID                    string     `gorm:"type:uuid;primaryKey" json:"id"`
	ExternalReservationID string     `gorm:"type:varchar(128);uniqueIndex;not null" json:"external_reservation_id"`
	ExternalReference     string     `gorm:"type:varchar(128);uniqueIndex;not null" json:"external_reference"`
	HoldID                string     `gorm:"type:uuid;uniqueIndex;not null" json:"hold_id"`
	Hold                  *Hold      `gorm:"foreignKey:HoldID;references:ID" json:"hold,omitempty"`
	OfferPK               string     `gorm:"type:uuid;not null" json:"offer_pk"`
	OfferID               string     `gorm:"type:varchar(128);not null" json:"offer_id"`
	ClientReference       string     `gorm:"type:varchar(128);uniqueIndex;not null" json:"client_reference"`
	Units                 int        `gorm:"not null" json:"units"`
	PassengerCount        int        `gorm:"not null" json:"passenger_count"`
	TotalPriceAmountMinor int64      `gorm:"not null" json:"total_price_amount_minor"`
	PriceCurrency         string     `gorm:"type:char(3);not null" json:"price_currency"`
	ServiceStartsAt       time.Time  `gorm:"type:timestamptz;not null" json:"service_starts_at"`
	ServiceEndsAt         time.Time  `gorm:"type:timestamptz;not null" json:"service_ends_at"`
	StartTimeZone         string     `gorm:"type:varchar(64);not null" json:"start_time_zone"`
	EndTimeZone           string     `gorm:"type:varchar(64);not null" json:"end_time_zone"`
	Status                string     `gorm:"type:varchar(16);not null;default:'CONFIRMED'" json:"status"`
	ConfirmedAt           time.Time  `gorm:"type:timestamptz;not null" json:"confirmed_at"`
	ReleasedAt            *time.Time `gorm:"type:timestamptz" json:"released_at,omitempty"`
	UpdatedAt             time.Time  `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Reservation) TableName() string {
	return "reservations"
}

type IdempotencyRecord struct {
	ID                 string     `gorm:"type:uuid;primaryKey" json:"id"`
	ClientScope        string     `gorm:"type:varchar(128);not null;uniqueIndex:idx_transport_idemp_scope_key,priority:1" json:"client_scope"`
	Operation          string     `gorm:"type:varchar(128);not null;uniqueIndex:idx_transport_idemp_scope_key,priority:2" json:"operation"`
	IdempotencyKey     string     `gorm:"type:varchar(128);not null;uniqueIndex:idx_transport_idemp_scope_key,priority:3" json:"idempotency_key"`
	RequestFingerprint string     `gorm:"type:char(64);not null" json:"request_fingerprint"`
	ResponseStatus     int        `gorm:"not null" json:"response_status"`
	ResponseBody       []byte     `gorm:"type:jsonb;not null" json:"response_body"`
	CreatedAt          time.Time  `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	ExpiresAt          *time.Time `gorm:"type:timestamptz" json:"expires_at,omitempty"`
}

func (IdempotencyRecord) TableName() string {
	return "idempotency_records"
}
