package hotel

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

type RoomType struct {
	ID               string      `gorm:"type:uuid;primaryKey" json:"id"`
	Code             string      `gorm:"type:varchar(64);uniqueIndex;not null" json:"code"`
	PropertyName     string      `gorm:"type:varchar(160);not null" json:"property_name"`
	Name             string      `gorm:"type:varchar(80);not null" json:"name"`
	MaxGuestsPerRoom int         `gorm:"not null" json:"max_guests_per_room"`
	Accessibility    StringArray `gorm:"type:jsonb;not null;default:'[]'" json:"accessibility"`
	Active           bool        `gorm:"not null;default:true" json:"active"`
	Synthetic        bool        `gorm:"not null;default:true" json:"synthetic"`
	Source           string      `gorm:"type:varchar(16);not null;default:'MOCK'" json:"source"`
	CreatedAt        time.Time   `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time   `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (RoomType) TableName() string {
	return "room_types"
}

type RoomInventoryDay struct {
	ID               string    `gorm:"type:uuid;primaryKey" json:"id"`
	RoomTypeID       string    `gorm:"type:uuid;not null;uniqueIndex:idx_room_inv_day" json:"room_type_id"`
	StayDate         time.Time `gorm:"type:date;not null;uniqueIndex:idx_room_inv_day" json:"stay_date"`
	RoomsTotal       int       `gorm:"not null" json:"rooms_total"`
	PriceAmountMinor int64     `gorm:"not null" json:"price_amount_minor"`
	PriceCurrency    string    `gorm:"type:char(3);not null" json:"price_currency"`
	Status           string    `gorm:"type:varchar(16);not null;default:'AVAILABLE'" json:"status"`
	CreatedAt        time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (RoomInventoryDay) TableName() string {
	return "room_inventory_days"
}

type HotelOffer struct {
	ID                  string             `gorm:"type:uuid;primaryKey" json:"id"`
	ExternalOfferID     string             `gorm:"type:varchar(128);uniqueIndex;not null" json:"external_offer_id"`
	ExternalInventoryID string             `gorm:"type:varchar(128);uniqueIndex;not null" json:"external_inventory_id"`
	RoomTypeID          string             `gorm:"type:uuid;not null" json:"room_type_id"`
	RoomType            *RoomType          `gorm:"foreignKey:RoomTypeID;references:ID" json:"room_type,omitempty"`
	CheckInDate         time.Time          `gorm:"type:date;not null" json:"check_in_date"`
	CheckOutDate        time.Time          `gorm:"type:date;not null" json:"check_out_date"`
	ServiceStartsAt     time.Time          `gorm:"type:timestamptz;not null" json:"service_starts_at"`
	ServiceEndsAt       time.Time          `gorm:"type:timestamptz;not null" json:"service_ends_at"`
	StartTimeZone       string             `gorm:"type:varchar(64);not null" json:"start_time_zone"`
	EndTimeZone         string             `gorm:"type:varchar(64);not null" json:"end_time_zone"`
	ValidUntil          time.Time          `gorm:"type:timestamptz;not null" json:"valid_until"`
	Status              string             `gorm:"type:varchar(16);not null;default:'AVAILABLE'" json:"status"`
	Synthetic           bool               `gorm:"not null;default:true" json:"synthetic"`
	Source              string             `gorm:"type:varchar(16);not null;default:'MOCK'" json:"source"`
	OfferNights         []HotelOfferNight  `gorm:"foreignKey:OfferPK;references:ID" json:"offer_nights,omitempty"`
	CreatedAt           time.Time          `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt           time.Time          `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (HotelOffer) TableName() string {
	return "hotel_offers"
}

type HotelOfferNight struct {
	OfferPK        string            `gorm:"type:uuid;primaryKey" json:"offer_pk"`
	InventoryDayID string            `gorm:"type:uuid;primaryKey" json:"inventory_day_id"`
	RoomTypeID     string            `gorm:"type:uuid;not null" json:"room_type_id"`
	InventoryDay   *RoomInventoryDay `gorm:"foreignKey:InventoryDayID;references:ID" json:"inventory_day,omitempty"`
}

func (HotelOfferNight) TableName() string {
	return "hotel_offer_nights"
}

type Hold struct {
	ID                    string      `gorm:"type:uuid;primaryKey" json:"id"`
	ExternalHoldID        string      `gorm:"type:varchar(128);uniqueIndex;not null" json:"external_hold_id"`
	ExternalReference     string      `gorm:"type:varchar(128);uniqueIndex;not null" json:"external_reference"`
	OfferPK               string      `gorm:"type:uuid;not null" json:"offer_pk"`
	OfferID               string      `gorm:"type:varchar(128);not null" json:"offer_id"`
	RoomTypeID            string      `gorm:"type:uuid;not null" json:"room_type_id"`
	RoomType              *RoomType   `gorm:"foreignKey:RoomTypeID;references:ID" json:"room_type,omitempty"`
	ClientReference       string      `gorm:"type:varchar(128);uniqueIndex;not null" json:"client_reference"`
	CheckInDate           time.Time   `gorm:"type:date;not null" json:"check_in_date"`
	CheckOutDate          time.Time   `gorm:"type:date;not null" json:"check_out_date"`
	TimeZone              string      `gorm:"type:varchar(64);not null" json:"time_zone"`
	RoomCount             int         `gorm:"not null" json:"room_count"`
	Accessibility         StringArray `gorm:"type:jsonb;not null;default:'[]'" json:"accessibility"`
	UnitPriceAmountMinor  int64       `gorm:"not null" json:"unit_price_amount_minor"`
	TotalPriceAmountMinor int64       `gorm:"not null" json:"total_price_amount_minor"`
	PriceCurrency         string      `gorm:"type:char(3);not null" json:"price_currency"`
	ServiceStartsAt       time.Time   `gorm:"type:timestamptz;not null" json:"service_starts_at"`
	ServiceEndsAt         time.Time   `gorm:"type:timestamptz;not null" json:"service_ends_at"`
	StartTimeZone         string      `gorm:"type:varchar(64);not null" json:"start_time_zone"`
	EndTimeZone           string      `gorm:"type:varchar(64);not null" json:"end_time_zone"`
	Status                string      `gorm:"type:varchar(16);not null;default:'HELD'" json:"status"`
	CreatedAt             time.Time   `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	ExpiresAt             time.Time   `gorm:"type:timestamptz;not null" json:"expires_at"`
	ConfirmedAt           *time.Time  `gorm:"type:timestamptz" json:"confirmed_at,omitempty"`
	ReleasedAt            *time.Time  `gorm:"type:timestamptz" json:"released_at,omitempty"`
	ExpiredAt             *time.Time  `gorm:"type:timestamptz" json:"expired_at,omitempty"`
	UpdatedAt             time.Time   `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Hold) TableName() string {
	return "holds"
}

type HoldNight struct {
	HoldID         string `gorm:"type:uuid;primaryKey" json:"hold_id"`
	OfferPK        string `gorm:"type:uuid;not null" json:"offer_pk"`
	InventoryDayID string `gorm:"type:uuid;primaryKey" json:"inventory_day_id"`
	RoomCount      int    `gorm:"not null" json:"room_count"`
}

func (HoldNight) TableName() string {
	return "hold_nights"
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
	RoomCount             int        `gorm:"not null" json:"room_count"`
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
	ClientScope        string     `gorm:"type:varchar(128);not null;uniqueIndex:idx_hotel_idemp_scope_key,priority:1" json:"client_scope"`
	Operation          string     `gorm:"type:varchar(128);not null;uniqueIndex:idx_hotel_idemp_scope_key,priority:2" json:"operation"`
	IdempotencyKey     string     `gorm:"type:varchar(128);not null;uniqueIndex:idx_hotel_idemp_scope_key,priority:3" json:"idempotency_key"`
	RequestFingerprint string     `gorm:"type:char(64);not null" json:"request_fingerprint"`
	ResponseStatus     int        `gorm:"not null" json:"response_status"`
	ResponseBody       []byte     `gorm:"type:jsonb;not null" json:"response_body"`
	CreatedAt          time.Time  `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	ExpiresAt          *time.Time `gorm:"type:timestamptz" json:"expires_at,omitempty"`
}

func (IdempotencyRecord) TableName() string {
	return "idempotency_records"
}
