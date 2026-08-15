package hospital

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrNotFound             = errors.New("not found")
	ErrCapacityConflict     = errors.New("capacity conflict")
	ErrOfferChanged         = errors.New("offer changed")
	ErrOfferExpired         = errors.New("offer expired")
	ErrHoldExpired          = errors.New("hold expired")
	ErrInvalidState         = errors.New("invalid state")
	ErrIdempotencyConflict  = errors.New("idempotency conflict")
)

type CapacityConflictDetails struct {
	Requested int
	Available int
}

type HoldExpiredDetails struct {
	ExpiredAt time.Time
}

type OfferExpiredDetails struct {
	ValidUntil time.Time
}

type SlotWithAvailableUnits struct {
	Slot           AppointmentSlot
	AvailableUnits int
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) DB() *gorm.DB {
	return r.db
}

func (r *Repository) FindMedicalServiceByCode(ctx context.Context, code string) (*MedicalService, error) {
	var service MedicalService
	err := r.db.WithContext(ctx).Where("code = ? AND active = true", code).First(&service).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &service, nil
}

func (r *Repository) SearchAvailableSlots(
	ctx context.Context,
	medicalServiceID string,
	windowStartsAt time.Time,
	windowEndsAt time.Time,
	now time.Time,
	requiredAccessibility []string,
	patientCount int,
) ([]SlotWithAvailableUnits, error) {
	var slots []AppointmentSlot
	err := r.db.WithContext(ctx).
		Where("medical_service_id = ? AND status = 'AVAILABLE' AND starts_at >= ? AND ends_at <= ? AND valid_until > ?",
			medicalServiceID, windowStartsAt, windowEndsAt, now).
		Order("starts_at ASC, price_amount_minor ASC").
		Find(&slots).Error
	if err != nil {
		return nil, err
	}

	results := make([]SlotWithAvailableUnits, 0, len(slots))
	for _, slot := range slots {
		if !satisfiesAccessibility(slot.Accessibility, requiredAccessibility) {
			continue
		}

		availableUnits, err := r.calculateAvailableUnits(ctx, r.db, slot.ID, slot.CapacityTotal, now)
		if err != nil {
			return nil, err
		}

		if availableUnits >= patientCount {
			results = append(results, SlotWithAvailableUnits{
				Slot:           slot,
				AvailableUnits: availableUnits,
			})
		}
	}

	return results, nil
}

func satisfiesAccessibility(slotAccessibility []string, required []string) bool {
	if len(required) == 0 {
		return true
	}
	existing := make(map[string]bool, len(slotAccessibility))
	for _, a := range slotAccessibility {
		existing[a] = true
	}
	for _, req := range required {
		if !existing[req] {
			return false
		}
	}
	return true
}

func (r *Repository) calculateAvailableUnits(ctx context.Context, tx *gorm.DB, slotID string, totalCapacity int, now time.Time) (int, error) {
	var activeHolds int
	err := tx.WithContext(ctx).
		Model(&Hold{}).
		Select("COALESCE(SUM(patient_count), 0)").
		Where("appointment_slot_id = ? AND status = 'HELD' AND expires_at > ?", slotID, now).
		Scan(&activeHolds).Error
	if err != nil {
		return 0, err
	}

	var confirmedReservations int
	err = tx.WithContext(ctx).
		Model(&Reservation{}).
		Select("COALESCE(SUM(patient_count), 0)").
		Where("appointment_slot_id = ? AND status = 'CONFIRMED'", slotID).
		Scan(&confirmedReservations).Error
	if err != nil {
		return 0, err
	}

	available := totalCapacity - (activeHolds + confirmedReservations)
	if available < 0 {
		available = 0
	}
	return available, nil
}

func (r *Repository) GetIdempotencyRecord(ctx context.Context, clientScope, operation, key string) (*IdempotencyRecord, error) {
	var record IdempotencyRecord
	err := r.db.WithContext(ctx).
		Where("client_scope = ? AND operation = ? AND idempotency_key = ?", clientScope, operation, key).
		First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

type CreateHoldParams struct {
	ClientScope        string
	Operation          string
	IdempotencyKey     string
	RequestFingerprint string
	OfferID            string
	Units              int
	ExpectedUnitPrice  int64
	ExpectedCurrency   string
	ClientReference    string
	HoldDuration       time.Duration
	Now                time.Time
}

type CreateHoldResult struct {
	Hold               *Hold
	Replayed           bool
	ConflictDetails    *CapacityConflictDetails
	HoldExpiredDetails *HoldExpiredDetails
	OfferExpiredDetails *OfferExpiredDetails
}

func (r *Repository) CreateHoldTx(ctx context.Context, params CreateHoldParams) (*CreateHoldResult, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer tx.Rollback()

	// Check idempotency first inside transaction
	var existingRecord IdempotencyRecord
	err := tx.Where("client_scope = ? AND operation = ? AND idempotency_key = ?",
		params.ClientScope, params.Operation, params.IdempotencyKey).
		First(&existingRecord).Error
	if err == nil {
		if existingRecord.RequestFingerprint != params.RequestFingerprint {
			return nil, ErrIdempotencyConflict
		}
		var hold Hold
		if err := json.Unmarshal(existingRecord.ResponseBody, &hold); err != nil {
			return nil, fmt.Errorf("unmarshal replayed hold: %w", err)
		}
		return &CreateHoldResult{Hold: &hold, Replayed: true}, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Lock the appointment slot row
	var slot AppointmentSlot
	err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("offer_id = ?", params.OfferID).
		First(&slot).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if slot.Status != "AVAILABLE" {
		return &CreateHoldResult{
			ConflictDetails: &CapacityConflictDetails{
				Requested: params.Units,
				Available: 0,
			},
		}, ErrCapacityConflict
	}

	if !slot.ValidUntil.After(params.Now) {
		return &CreateHoldResult{
			OfferExpiredDetails: &OfferExpiredDetails{
				ValidUntil: slot.ValidUntil,
			},
		}, ErrOfferExpired
	}

	if slot.PriceAmountMinor != params.ExpectedUnitPrice || slot.PriceCurrency != params.ExpectedCurrency {
		return nil, ErrOfferChanged
	}

	availableUnits, err := r.calculateAvailableUnits(ctx, tx, slot.ID, slot.CapacityTotal, params.Now)
	if err != nil {
		return nil, err
	}

	if availableUnits < params.Units {
		return &CreateHoldResult{
			ConflictDetails: &CapacityConflictDetails{
				Requested: params.Units,
				Available: availableUnits,
			},
		}, ErrCapacityConflict
	}

	holdID := NewUUID()
	externalHoldID := "hosp-hold-" + NewRandomHex(6)
	externalRef := "HOSP-HOLD-" + NewRandomHex(6)
	createdAt := params.Now.UTC()
	expiresAt := createdAt.Add(params.HoldDuration)

	hold := Hold{
		ID:                    holdID,
		ExternalHoldID:        externalHoldID,
		ExternalReference:     externalRef,
		AppointmentSlotID:     slot.ID,
		OfferID:               slot.OfferID,
		ClientReference:       params.ClientReference,
		PatientCount:          params.Units,
		UnitPriceAmountMinor:  slot.PriceAmountMinor,
		TotalPriceAmountMinor: slot.PriceAmountMinor * int64(params.Units),
		PriceCurrency:         slot.PriceCurrency,
		ServiceStartsAt:       slot.StartsAt,
		ServiceEndsAt:         slot.EndsAt,
		StartTimeZone:         slot.StartTimeZone,
		EndTimeZone:           slot.EndTimeZone,
		Status:                "HELD",
		CreatedAt:             createdAt,
		ExpiresAt:             expiresAt,
		UpdatedAt:             createdAt,
	}

	if err := tx.Create(&hold).Error; err != nil {
		return nil, fmt.Errorf("create hold: %w", err)
	}

	// Prepare idempotency record
	responseBytes, err := json.Marshal(hold)
	if err != nil {
		return nil, fmt.Errorf("marshal hold response: %w", err)
	}

	idempRecord := IdempotencyRecord{
		ID:                 NewUUID(),
		ClientScope:        params.ClientScope,
		Operation:          params.Operation,
		IdempotencyKey:     params.IdempotencyKey,
		RequestFingerprint: params.RequestFingerprint,
		ResponseStatus:     201,
		ResponseBody:       responseBytes,
		CreatedAt:          createdAt,
		UpdatedAt:          createdAt,
	}
	if err := tx.Create(&idempRecord).Error; err != nil {
		return nil, fmt.Errorf("create idempotency record: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("commit hold transaction: %w", err)
	}

	return &CreateHoldResult{Hold: &hold, Replayed: false}, nil
}

type ConfirmHoldParams struct {
	ClientScope        string
	Operation          string
	IdempotencyKey     string
	RequestFingerprint string
	ExternalHoldID     string
	Now                time.Time
}

type ConfirmHoldResult struct {
	Reservation        *Reservation
	Replayed           bool
	HoldExpiredDetails *HoldExpiredDetails
}

func (r *Repository) ConfirmHoldTx(ctx context.Context, params ConfirmHoldParams) (*ConfirmHoldResult, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer tx.Rollback()

	// Check idempotency first inside transaction
	var existingRecord IdempotencyRecord
	err := tx.Where("client_scope = ? AND operation = ? AND idempotency_key = ?",
		params.ClientScope, params.Operation, params.IdempotencyKey).
		First(&existingRecord).Error
	if err == nil {
		if existingRecord.RequestFingerprint != params.RequestFingerprint {
			return nil, ErrIdempotencyConflict
		}
		var res Reservation
		if err := json.Unmarshal(existingRecord.ResponseBody, &res); err != nil {
			return nil, fmt.Errorf("unmarshal replayed reservation: %w", err)
		}
		return &ConfirmHoldResult{Reservation: &res, Replayed: true}, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Lock the hold row
	var hold Hold
	err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("external_hold_id = ?", params.ExternalHoldID).
		First(&hold).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if hold.Status == "EXPIRED" {
		return &ConfirmHoldResult{
			HoldExpiredDetails: &HoldExpiredDetails{ExpiredAt: hold.ExpiresAt},
		}, ErrHoldExpired
	}
	if hold.Status == "RELEASED" || hold.Status == "CONFIRMED" {
		return nil, ErrInvalidState
	}

	// If now >= ExpiresAt, mark EXPIRED lazily
	if !params.Now.Before(hold.ExpiresAt) {
		nowUTC := params.Now.UTC()
		hold.Status = "EXPIRED"
		hold.ExpiredAt = &nowUTC
		hold.UpdatedAt = nowUTC
		if err := tx.Save(&hold).Error; err != nil {
			return nil, err
		}
		if err := tx.Commit().Error; err != nil {
			return nil, err
		}
		return &ConfirmHoldResult{
			HoldExpiredDetails: &HoldExpiredDetails{ExpiredAt: hold.ExpiresAt},
		}, ErrHoldExpired
	}

	nowUTC := params.Now.UTC()
	hold.Status = "CONFIRMED"
	hold.ConfirmedAt = &nowUTC
	hold.UpdatedAt = nowUTC
	if err := tx.Save(&hold).Error; err != nil {
		return nil, fmt.Errorf("update hold status: %w", err)
	}

	reservationID := NewUUID()
	externalResID := "hospital-res-" + NewRandomHex(6)
	externalRef := "HOSP-RES-" + NewRandomHex(6)

	reservation := Reservation{
		ID:                    reservationID,
		ExternalReservationID: externalResID,
		ExternalReference:     externalRef,
		HoldID:                hold.ID,
		AppointmentSlotID:     hold.AppointmentSlotID,
		OfferID:               hold.OfferID,
		ClientReference:       hold.ClientReference,
		PatientCount:          hold.PatientCount,
		TotalPriceAmountMinor: hold.TotalPriceAmountMinor,
		PriceCurrency:         hold.PriceCurrency,
		ServiceStartsAt:       hold.ServiceStartsAt,
		ServiceEndsAt:         hold.ServiceEndsAt,
		StartTimeZone:         hold.StartTimeZone,
		EndTimeZone:           hold.EndTimeZone,
		Status:                "CONFIRMED",
		ConfirmedAt:           nowUTC,
		ReleasedAt:            nil,
		CreatedAt:             nowUTC,
		UpdatedAt:             nowUTC,
	}

	if err := tx.Create(&reservation).Error; err != nil {
		return nil, fmt.Errorf("create reservation: %w", err)
	}

	responseBytes, err := json.Marshal(reservation)
	if err != nil {
		return nil, fmt.Errorf("marshal reservation response: %w", err)
	}

	idempRecord := IdempotencyRecord{
		ID:                 NewUUID(),
		ClientScope:        params.ClientScope,
		Operation:          params.Operation,
		IdempotencyKey:     params.IdempotencyKey,
		RequestFingerprint: params.RequestFingerprint,
		ResponseStatus:     201,
		ResponseBody:       responseBytes,
		CreatedAt:          nowUTC,
		UpdatedAt:          nowUTC,
	}
	if err := tx.Create(&idempRecord).Error; err != nil {
		return nil, fmt.Errorf("create idempotency record: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("commit confirm transaction: %w", err)
	}

	return &ConfirmHoldResult{Reservation: &reservation, Replayed: false}, nil
}

type ReleaseHoldParams struct {
	ClientScope        string
	Operation          string
	IdempotencyKey     string
	RequestFingerprint string
	ExternalHoldID     string
	Now                time.Time
}

type ReleaseHoldResult struct {
	Hold       *Hold
	Replayed   bool
	ReleasedAt time.Time
}

func (r *Repository) ReleaseHoldTx(ctx context.Context, params ReleaseHoldParams) (*ReleaseHoldResult, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer tx.Rollback()

	// Check idempotency first inside transaction
	var existingRecord IdempotencyRecord
	err := tx.Where("client_scope = ? AND operation = ? AND idempotency_key = ?",
		params.ClientScope, params.Operation, params.IdempotencyKey).
		First(&existingRecord).Error
	if err == nil {
		if existingRecord.RequestFingerprint != params.RequestFingerprint {
			return nil, ErrIdempotencyConflict
		}
		var hold Hold
		if err := json.Unmarshal(existingRecord.ResponseBody, &hold); err != nil {
			return nil, fmt.Errorf("unmarshal replayed release hold: %w", err)
		}
		releasedAt := params.Now
		if hold.ReleasedAt != nil {
			releasedAt = *hold.ReleasedAt
		}
		return &ReleaseHoldResult{Hold: &hold, Replayed: true, ReleasedAt: releasedAt}, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Lock the hold row
	var hold Hold
	err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("external_hold_id = ?", params.ExternalHoldID).
		First(&hold).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if hold.Status == "RELEASED" {
		releasedAt := hold.UpdatedAt
		if hold.ReleasedAt != nil {
			releasedAt = *hold.ReleasedAt
		}
		return &ReleaseHoldResult{Hold: &hold, Replayed: false, ReleasedAt: releasedAt}, nil
	}
	if hold.Status == "CONFIRMED" || hold.Status == "EXPIRED" {
		return nil, ErrInvalidState
	}

	nowUTC := params.Now.UTC()
	hold.Status = "RELEASED"
	hold.ReleasedAt = &nowUTC
	hold.UpdatedAt = nowUTC
	if err := tx.Save(&hold).Error; err != nil {
		return nil, fmt.Errorf("update hold released: %w", err)
	}

	responseBytes, err := json.Marshal(hold)
	if err != nil {
		return nil, fmt.Errorf("marshal release hold response: %w", err)
	}

	idempRecord := IdempotencyRecord{
		ID:                 NewUUID(),
		ClientScope:        params.ClientScope,
		Operation:          params.Operation,
		IdempotencyKey:     params.IdempotencyKey,
		RequestFingerprint: params.RequestFingerprint,
		ResponseStatus:     200,
		ResponseBody:       responseBytes,
		CreatedAt:          nowUTC,
		UpdatedAt:          nowUTC,
	}
	if err := tx.Create(&idempRecord).Error; err != nil {
		return nil, fmt.Errorf("create idempotency record: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("commit release hold transaction: %w", err)
	}

	return &ReleaseHoldResult{Hold: &hold, Replayed: false, ReleasedAt: nowUTC}, nil
}

func (r *Repository) GetReservation(ctx context.Context, externalReservationID string) (*Reservation, error) {
	var res Reservation
	err := r.db.WithContext(ctx).
		Preload("Hold").
		Where("external_reservation_id = ?", externalReservationID).
		First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &res, nil
}

type ReleaseReservationParams struct {
	ClientScope           string
	Operation             string
	IdempotencyKey        string
	RequestFingerprint    string
	ExternalReservationID string
	Now                   time.Time
}

type ReleaseReservationResult struct {
	Reservation *Reservation
	Replayed    bool
	ReleasedAt  time.Time
}

func (r *Repository) ReleaseReservationTx(ctx context.Context, params ReleaseReservationParams) (*ReleaseReservationResult, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer tx.Rollback()

	// Check idempotency first inside transaction
	var existingRecord IdempotencyRecord
	err := tx.Where("client_scope = ? AND operation = ? AND idempotency_key = ?",
		params.ClientScope, params.Operation, params.IdempotencyKey).
		First(&existingRecord).Error
	if err == nil {
		if existingRecord.RequestFingerprint != params.RequestFingerprint {
			return nil, ErrIdempotencyConflict
		}
		var res Reservation
		if err := json.Unmarshal(existingRecord.ResponseBody, &res); err != nil {
			return nil, fmt.Errorf("unmarshal replayed release reservation: %w", err)
		}
		releasedAt := params.Now
		if res.ReleasedAt != nil {
			releasedAt = *res.ReleasedAt
		}
		return &ReleaseReservationResult{Reservation: &res, Replayed: true, ReleasedAt: releasedAt}, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Lock the reservation row
	var reservation Reservation
	err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("external_reservation_id = ?", params.ExternalReservationID).
		First(&reservation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if reservation.Status == "RELEASED" {
		releasedAt := reservation.UpdatedAt
		if reservation.ReleasedAt != nil {
			releasedAt = *reservation.ReleasedAt
		}
		return &ReleaseReservationResult{Reservation: &reservation, Replayed: false, ReleasedAt: releasedAt}, nil
	}

	nowUTC := params.Now.UTC()
	reservation.Status = "RELEASED"
	reservation.ReleasedAt = &nowUTC
	reservation.UpdatedAt = nowUTC
	if err := tx.Save(&reservation).Error; err != nil {
		return nil, fmt.Errorf("update reservation released: %w", err)
	}

	// Also lock and update originating hold
	var hold Hold
	err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", reservation.HoldID).
		First(&hold).Error
	if err == nil {
		hold.Status = "RELEASED"
		hold.ReleasedAt = &nowUTC
		hold.UpdatedAt = nowUTC
		if err := tx.Save(&hold).Error; err != nil {
			return nil, fmt.Errorf("update originating hold released: %w", err)
		}
	}

	responseBytes, err := json.Marshal(reservation)
	if err != nil {
		return nil, fmt.Errorf("marshal release reservation response: %w", err)
	}

	idempRecord := IdempotencyRecord{
		ID:                 NewUUID(),
		ClientScope:        params.ClientScope,
		Operation:          params.Operation,
		IdempotencyKey:     params.IdempotencyKey,
		RequestFingerprint: params.RequestFingerprint,
		ResponseStatus:     200,
		ResponseBody:       responseBytes,
		CreatedAt:          nowUTC,
		UpdatedAt:          nowUTC,
	}
	if err := tx.Create(&idempRecord).Error; err != nil {
		return nil, fmt.Errorf("create idempotency record: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("commit release reservation transaction: %w", err)
	}

	return &ReleaseReservationResult{Reservation: &reservation, Replayed: false, ReleasedAt: nowUTC}, nil
}

// NewUUID generates a valid RFC 4122 v4 UUID.
func NewUUID() string {
	var uuid [16]byte
	_, _ = rand.Read(uuid[:])
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant RFC 4122
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])
}

// NewRandomHex generates a random lowercase hex string of specified byte length.
func NewRandomHex(n int) string {
	bytes := make([]byte, n)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
