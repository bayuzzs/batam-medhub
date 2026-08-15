package hospital

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

type MedicalService struct {
	ID                     string    `gorm:"type:uuid;primaryKey" json:"id"`
	Code                   string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"code"`
	Name                   string    `gorm:"type:varchar(160);not null" json:"name"`
	DefaultDurationMinutes int       `gorm:"not null" json:"default_duration_minutes"`
	PatientRequestable     bool      `gorm:"not null;default:true" json:"patient_requestable"`
	Active                 bool      `gorm:"not null;default:true" json:"active"`
	Synthetic              bool      `gorm:"not null;default:true" json:"synthetic"`
	Source                 string    `gorm:"type:varchar(16);not null;default:'MOCK'" json:"source"`
	CreatedAt              time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt              time.Time `gorm:"not null" json:"updated_at"`
}

func (MedicalService) TableName() string {
	return "medical_services"
}

type AppointmentSlot struct {
	ID                string          `gorm:"type:uuid;primaryKey" json:"id"`
	MedicalServiceID  string          `gorm:"type:uuid;not null" json:"medical_service_id"`
	ExternalSlotID    string          `gorm:"type:varchar(128);uniqueIndex;not null" json:"external_slot_id"`
	OfferID           string          `gorm:"type:varchar(128);uniqueIndex;not null" json:"offer_id"`
	FacilityName      string          `gorm:"type:varchar(160);not null" json:"facility_name"`
	StartsAt          time.Time       `gorm:"not null" json:"starts_at"`
	EndsAt            time.Time       `gorm:"not null" json:"ends_at"`
	StartTimeZone     string          `gorm:"type:varchar(64);not null" json:"start_time_zone"`
	EndTimeZone       string          `gorm:"type:varchar(64);not null" json:"end_time_zone"`
	DurationMinutes   int             `gorm:"not null" json:"duration_minutes"`
	CapacityTotal     int             `gorm:"not null" json:"capacity_total"`
	PriceAmountMinor  int64           `gorm:"not null" json:"price_amount_minor"`
	PriceCurrency     string          `gorm:"type:char(3);not null" json:"price_currency"`
	ValidUntil        time.Time       `gorm:"not null" json:"valid_until"`
	Languages         StringArray     `gorm:"type:jsonb;not null;default:'[]'" json:"languages"`
	Accessibility     StringArray     `gorm:"type:jsonb;not null;default:'[]'" json:"accessibility"`
	Status            string          `gorm:"type:varchar(16);not null;default:'AVAILABLE'" json:"status"`
	Synthetic         bool            `gorm:"not null;default:true" json:"synthetic"`
	Source            string          `gorm:"type:varchar(16);not null;default:'MOCK'" json:"source"`
	CreatedAt         time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"not null" json:"updated_at"`
	MedicalService    *MedicalService `gorm:"foreignKey:MedicalServiceID" json:"medical_service,omitempty"`
}

func (AppointmentSlot) TableName() string {
	return "appointment_slots"
}

type Hold struct {
	ID                    string     `gorm:"type:uuid;primaryKey" json:"id"`
	ExternalHoldID        string     `gorm:"type:varchar(128);uniqueIndex;not null" json:"external_hold_id"`
	ExternalReference     string     `gorm:"type:varchar(128);uniqueIndex;not null" json:"external_reference"`
	AppointmentSlotID     string     `gorm:"type:uuid;not null" json:"appointment_slot_id"`
	OfferID               string     `gorm:"type:varchar(128);not null" json:"offer_id"`
	ClientReference       string     `gorm:"type:varchar(128);uniqueIndex;not null" json:"client_reference"`
	PatientCount          int        `gorm:"not null" json:"patient_count"`
	UnitPriceAmountMinor  int64      `gorm:"not null" json:"unit_price_amount_minor"`
	TotalPriceAmountMinor int64      `gorm:"not null" json:"total_price_amount_minor"`
	PriceCurrency         string     `gorm:"type:char(3);not null" json:"price_currency"`
	ServiceStartsAt       time.Time  `gorm:"not null" json:"service_starts_at"`
	ServiceEndsAt         time.Time  `gorm:"not null" json:"service_ends_at"`
	StartTimeZone         string     `gorm:"type:varchar(64);not null" json:"start_time_zone"`
	EndTimeZone           string     `gorm:"type:varchar(64);not null" json:"end_time_zone"`
	Status                string     `gorm:"type:varchar(16);not null" json:"status"`
	CreatedAt             time.Time  `gorm:"not null" json:"created_at"`
	ExpiresAt             time.Time  `gorm:"not null" json:"expires_at"`
	ConfirmedAt           *time.Time `json:"confirmed_at"`
	ReleasedAt            *time.Time `json:"released_at"`
	ExpiredAt             *time.Time `json:"expired_at"`
	UpdatedAt             time.Time  `gorm:"not null" json:"updated_at"`
}

func (Hold) TableName() string {
	return "holds"
}

type Reservation struct {
	ID                    string     `gorm:"type:uuid;primaryKey" json:"id"`
	ExternalReservationID string     `gorm:"type:varchar(128);uniqueIndex;not null" json:"external_reservation_id"`
	ExternalReference     string     `gorm:"type:varchar(128);uniqueIndex;not null" json:"external_reference"`
	HoldID                string     `gorm:"type:uuid;uniqueIndex;not null" json:"hold_id"`
	AppointmentSlotID     string     `gorm:"type:uuid;not null" json:"appointment_slot_id"`
	OfferID               string     `gorm:"type:varchar(128);not null" json:"offer_id"`
	ClientReference       string     `gorm:"type:varchar(128);not null" json:"client_reference"`
	PatientCount          int        `gorm:"not null" json:"patient_count"`
	TotalPriceAmountMinor int64      `gorm:"not null" json:"total_price_amount_minor"`
	PriceCurrency         string     `gorm:"type:char(3);not null" json:"price_currency"`
	ServiceStartsAt       time.Time  `gorm:"not null" json:"service_starts_at"`
	ServiceEndsAt         time.Time  `gorm:"not null" json:"service_ends_at"`
	StartTimeZone         string     `gorm:"type:varchar(64);not null" json:"start_time_zone"`
	EndTimeZone           string     `gorm:"type:varchar(64);not null" json:"end_time_zone"`
	Status                string     `gorm:"type:varchar(16);not null" json:"status"`
	ConfirmedAt           time.Time  `gorm:"not null" json:"confirmed_at"`
	ReleasedAt            *time.Time `json:"released_at"`
	CreatedAt             time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt             time.Time  `gorm:"not null" json:"updated_at"`
	Hold                  *Hold      `gorm:"foreignKey:HoldID" json:"hold,omitempty"`
}

func (Reservation) TableName() string {
	return "reservations"
}

type IdempotencyRecord struct {
	ID                 string          `gorm:"type:uuid;primaryKey" json:"id"`
	ClientScope        string          `gorm:"type:varchar(128);not null" json:"client_scope"`
	Operation          string          `gorm:"type:varchar(160);not null" json:"operation"`
	IdempotencyKey     string          `gorm:"type:varchar(128);not null" json:"idempotency_key"`
	RequestFingerprint string          `gorm:"type:char(64);not null" json:"request_fingerprint"`
	ResponseStatus     int             `gorm:"not null" json:"response_status"`
	ResponseBody       json.RawMessage `gorm:"type:jsonb;not null" json:"response_body"`
	CreatedAt          time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt          time.Time       `gorm:"not null" json:"updated_at"`
}

func (IdempotencyRecord) TableName() string {
	return "idempotency_records"
}
