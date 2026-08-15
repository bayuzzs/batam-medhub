package adapter

import "encoding/json"

// Provider category constants.
const (
	ProviderTypeHospital  = "HOSPITAL"
	ProviderTypeFerry     = "FERRY"
	ProviderTypeHotel     = "HOTEL"
	ProviderTypeTransport = "TRANSPORT"
)

// TimeWindow represents an operational UTC time window with start/end timezones.
type TimeWindow struct {
	StartsAt      string `json:"starts_at"`
	EndsAt        string `json:"ends_at"`
	StartTimeZone string `json:"start_time_zone"`
	EndTimeZone   string `json:"end_time_zone"`
}

// Money represents amount in integer minor units and currency code.
type Money struct {
	AmountMinor int64  `json:"amount_minor"`
	Currency    string `json:"currency"`
}

// HealthResponse represents provider service health.
type HealthResponse struct {
	Status         string `json:"status"`
	ProviderID     string `json:"provider_id"`
	ProviderType   string `json:"provider_type"`
	DatabaseStatus string `json:"database_status"`
	CheckedAt      string `json:"checked_at"`
}

// SearchRequest represents provider-owned availability search.
type SearchRequest struct {
	ProviderType string `json:"provider_type"`
	Criteria     any    `json:"criteria"`
}

// HospitalSearchCriteria represents criteria for querying appointment slots.
type HospitalSearchCriteria struct {
	ServiceCode       string     `json:"service_code"`
	PatientCount      int        `json:"patient_count"`
	AppointmentWindow TimeWindow `json:"appointment_window"`
	Accessibility     []string   `json:"accessibility"`
}

// FerrySearchCriteria represents criteria for querying ferry sailings.
type FerrySearchCriteria struct {
	OriginPortCode      string     `json:"origin_port_code"`
	DestinationPortCode string     `json:"destination_port_code"`
	PassengerCount      int        `json:"passenger_count"`
	DepartureWindow     TimeWindow `json:"departure_window"`
}

// HotelSearchCriteria represents criteria for querying hotel room inventory.
type HotelSearchCriteria struct {
	CheckInDate   string   `json:"check_in_date"`
	CheckOutDate  string   `json:"check_out_date"`
	LocalTimezone string   `json:"local_timezone"`
	RoomCount     int      `json:"room_count"`
	GuestCount    int      `json:"guest_count"`
	Accessibility []string `json:"accessibility"`
}

// TransportSearchCriteria represents criteria for querying ground transport offers.
type TransportSearchCriteria struct {
	PickupLocationCode  string     `json:"pickup_location_code"`
	DropoffLocationCode string     `json:"dropoff_location_code"`
	PassengerCount      int        `json:"passenger_count"`
	PickupWindow        TimeWindow `json:"pickup_window"`
	Accessibility       []string   `json:"accessibility"`
}

// Offer represents a provider offer item returned by search.
type Offer struct {
	OfferID        string          `json:"offer_id"`
	ProviderID     string          `json:"provider_id"`
	ProviderType   string          `json:"provider_type"`
	Status         string          `json:"status"`
	ServiceWindow  TimeWindow      `json:"service_window"`
	AvailableUnits int             `json:"available_units"`
	UnitPrice      Money           `json:"unit_price"`
	ValidUntil     string          `json:"valid_until"`
	Synthetic      bool            `json:"synthetic"`
	Source         string          `json:"source"`
	Details        json.RawMessage `json:"details"`
}

// HospitalOfferDetails represents provider-specific hospital offer details.
type HospitalOfferDetails struct {
	ProviderType      string   `json:"provider_type"`
	ServiceCode       string   `json:"service_code"`
	AppointmentSlotID string   `json:"appointment_slot_id"`
	FacilityName      string   `json:"facility_name"`
	DurationMinutes   int      `json:"duration_minutes"`
	Languages         []string `json:"languages"`
	Accessibility     []string `json:"accessibility"`
}

// FerryOfferDetails represents provider-specific ferry offer details.
type FerryOfferDetails struct {
	ProviderType        string `json:"provider_type"`
	SailingID           string `json:"sailing_id"`
	OperatorName        string `json:"operator_name"`
	OriginPortCode      string `json:"origin_port_code"`
	DestinationPortCode string `json:"destination_port_code"`
	CheckInCutoffAt     string `json:"check_in_cutoff_at"`
}

// HotelOfferDetails represents provider-specific hotel offer details.
type HotelOfferDetails struct {
	ProviderType      string   `json:"provider_type"`
	RoomInventoryID   string   `json:"room_inventory_id"`
	PropertyName      string   `json:"property_name"`
	RoomType          string   `json:"room_type"`
	CheckInDate       string   `json:"check_in_date"`
	CheckOutDate      string   `json:"check_out_date"`
	MaxGuestsPerRoom  int      `json:"max_guests_per_room"`
	Accessibility     []string `json:"accessibility"`
}

// TransportOfferDetails represents provider-specific transport offer details.
type TransportOfferDetails struct {
	ProviderType        string   `json:"provider_type"`
	AvailabilityID      string   `json:"availability_id"`
	VehicleType         string   `json:"vehicle_type"`
	PickupLocationCode  string   `json:"pickup_location_code"`
	DropoffLocationCode string   `json:"dropoff_location_code"`
	PassengerCapacity   int      `json:"passenger_capacity"`
	Accessibility       []string `json:"accessibility"`
}

// SearchResponse represents the response envelope of a search query.
type SearchResponse struct {
	ProviderID   string  `json:"provider_id"`
	ProviderType string  `json:"provider_type"`
	Offers       []Offer `json:"offers"`
}

// TransportBookingRequirements specifies passenger and route requirements for transport holds.
type TransportBookingRequirements struct {
	PassengerCount      int        `json:"passenger_count"`
	PickupLocationCode  string     `json:"pickup_location_code"`
	DropoffLocationCode string     `json:"dropoff_location_code"`
	PickupWindow        TimeWindow `json:"pickup_window"`
	Accessibility       []string   `json:"accessibility"`
}

// CreateHoldRequest represents the hold creation payload.
type CreateHoldRequest struct {
	ProviderID          string                        `json:"provider_id"`
	ProviderType        string                        `json:"provider_type"`
	OfferID             string                        `json:"offer_id"`
	Units               int                           `json:"units"`
	ExpectedUnitPrice   Money                         `json:"expected_unit_price"`
	ClientReference     string                        `json:"client_reference"`
	BookingRequirements *TransportBookingRequirements `json:"booking_requirements,omitempty"`
}

// Hold represents a live capacity reservation held by a provider.
type Hold struct {
	HoldID            string     `json:"hold_id"`
	ExternalReference string     `json:"external_reference"`
	ProviderID        string     `json:"provider_id"`
	ProviderType      string     `json:"provider_type"`
	OfferID           string     `json:"offer_id"`
	ClientReference   string     `json:"client_reference"`
	Status            string     `json:"status"`
	Units             int        `json:"units"`
	UnitPrice         Money      `json:"unit_price"`
	TotalPrice        Money      `json:"total_price"`
	ServiceWindow     TimeWindow `json:"service_window"`
	CreatedAt         string     `json:"created_at"`
	ExpiresAt         string     `json:"expires_at"`
}

// Reservation represents a confirmed provider booking.
type Reservation struct {
	ReservationID     string     `json:"reservation_id"`
	ExternalReference string     `json:"external_reference"`
	HoldID            string     `json:"hold_id"`
	ProviderID        string     `json:"provider_id"`
	ProviderType      string     `json:"provider_type"`
	OfferID           string     `json:"offer_id"`
	ClientReference   string     `json:"client_reference"`
	Status            string     `json:"status"`
	Units             int        `json:"units"`
	TotalPrice        Money      `json:"total_price"`
	ServiceWindow     TimeWindow `json:"service_window"`
	ConfirmedAt       string     `json:"confirmed_at"`
	ReleasedAt        *string    `json:"released_at"`
}

// ReleaseResult represents the outcome of releasing a hold or reservation.
type ReleaseResult struct {
	ResourceType      string `json:"resource_type"`
	ResourceID        string `json:"resource_id"`
	ExternalReference string `json:"external_reference"`
	ProviderID        string `json:"provider_id"`
	ProviderType      string `json:"provider_type"`
	Status            string `json:"status"`
	ReleasedAt        string `json:"released_at"`
}

// ErrorDetail represents detail items in provider error envelopes.
type ErrorDetail struct {
	Field  *string `json:"field,omitempty"`
	Reason string  `json:"reason"`
}

// ErrorBody represents the inner error body of a provider error envelope.
type ErrorBody struct {
	Code      string        `json:"code"`
	Message   string        `json:"message"`
	Retryable bool          `json:"retryable"`
	RequestID string        `json:"request_id"`
	Details   []ErrorDetail `json:"details"`
}

// ErrorEnvelope represents the standard provider error response.
type ErrorEnvelope struct {
	Error ErrorBody `json:"error"`
}
