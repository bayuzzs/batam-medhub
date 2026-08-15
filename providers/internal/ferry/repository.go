package ferry

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
	ErrNotFound            = errors.New("resource not found")
	ErrCapacityConflict    = errors.New("capacity conflict")
	ErrOfferChanged        = errors.New("offer changed")
	ErrOfferExpired        = errors.New("offer expired")
	ErrHoldExpired         = errors.New("hold expired")
	ErrInvalidState        = errors.New("invalid state transition")
	ErrIdempotencyConflict = errors.New("idempotency conflict")
)

type CapacityConflictDetails struct {
	Requested int
	Available int
}

type OfferExpiredDetails struct {
	ValidUntil time.Time
}

type HoldExpiredDetails struct {
	ExpiredAt time.Time
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

type SailingWithAvailableUnits struct {
	Sailing        Sailing
	AvailableUnits int
}

func (r *Repository) SearchAvailableSailings(
	ctx context.Context,
	originPortCode string,
	destinationPortCode string,
	windowStart time.Time,
	windowEnd time.Time,
	now time.Time,
	requestedUnits int,
) ([]SailingWithAvailableUnits, error) {
	var sailings []Sailing
	err := r.db.WithContext(ctx).
		Where("origin_port_code = ? AND destination_port_code = ? AND status = 'AVAILABLE' AND departs_at >= ? AND departs_at <= ? AND valid_until > ?",
			originPortCode, destinationPortCode, windowStart, windowEnd, now).
		Order("departs_at ASC").
		Find(&sailings).Error
	if err != nil {
		return nil, fmt.Errorf("query sailings: %w", err)
	}

	var results []SailingWithAvailableUnits
	for _, sailing := range sailings {
		available, err := r.calculateAvailableUnits(ctx, r.db, sailing.ID, sailing.SeatCapacity, now)
		if err != nil {
			return nil, fmt.Errorf("calculate available units: %w", err)
		}
		if available >= requestedUnits {
			results = append(results, SailingWithAvailableUnits{
				Sailing:        sailing,
				AvailableUnits: available,
			})
		}
	}

	return results, nil
}

func (r *Repository) calculateAvailableUnits(
	ctx context.Context,
	tx *gorm.DB,
	sailingID string,
	totalCapacity int,
	now time.Time,
) (int, error) {
	var activeHoldsSum int64
	err := tx.WithContext(ctx).
		Model(&Hold{}).
		Where("sailing_id = ? AND status = 'HELD' AND expires_at > ?", sailingID, now).
		Select("COALESCE(SUM(passenger_count), 0)").
		Scan(&activeHoldsSum).Error
	if err != nil {
		return 0, fmt.Errorf("sum active holds: %w", err)
	}

	var confirmedResSum int64
	err = tx.WithContext(ctx).
		Model(&Reservation{}).
		Where("sailing_id = ? AND status = 'CONFIRMED'", sailingID).
		Select("COALESCE(SUM(passenger_count), 0)").
		Scan(&confirmedResSum).Error
	if err != nil {
		return 0, fmt.Errorf("sum confirmed reservations: %w", err)
	}

	available := totalCapacity - int(activeHoldsSum) - int(confirmedResSum)
	if available < 0 {
		available = 0
	}
	return available, nil
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
	Hold                *Hold
	Replayed            bool
	ConflictDetails     *CapacityConflictDetails
	OfferExpiredDetails *OfferExpiredDetails
}

func (r *Repository) CreateHoldTx(ctx context.Context, params CreateHoldParams) (*CreateHoldResult, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer tx.Rollback()

	// PostgreSQL transaction advisory lock to serialize concurrent identical idempotency keys
	lockKey := fmt.Sprintf("%s:%s:%s", params.ClientScope, params.Operation, params.IdempotencyKey)
	if err := tx.Exec("SELECT pg_advisory_xact_lock(hashtext(?))", lockKey).Error; err != nil {
		return nil, fmt.Errorf("acquire advisory lock: %w", err)
	}

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

	// Lock the sailing row
	var sailing Sailing
	err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("offer_id = ?", params.OfferID).
		First(&sailing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if sailing.Status != "AVAILABLE" {
		return &CreateHoldResult{
			ConflictDetails: &CapacityConflictDetails{
				Requested: params.Units,
				Available: 0,
			},
		}, ErrCapacityConflict
	}

	if !sailing.ValidUntil.After(params.Now) {
		return &CreateHoldResult{
			OfferExpiredDetails: &OfferExpiredDetails{
				ValidUntil: sailing.ValidUntil,
			},
		}, ErrOfferExpired
	}

	if sailing.PriceAmountMinor != params.ExpectedUnitPrice || sailing.PriceCurrency != params.ExpectedCurrency {
		return nil, ErrOfferChanged
	}

	availableUnits, err := r.calculateAvailableUnits(ctx, tx, sailing.ID, sailing.SeatCapacity, params.Now)
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

	nowUTC := params.Now.UTC()
	expiresAt := nowUTC.Add(params.HoldDuration)

	holdID := NewUUID()
	externalHoldID := "ferry-hold-" + NewRandomHex(6)
	externalRef := "FERRY-HOLD-" + NewRandomHex(6)

	unitPrice := sailing.PriceAmountMinor
	totalPrice := unitPrice * int64(params.Units)

	hold := Hold{
		ID:                    holdID,
		ExternalHoldID:        externalHoldID,
		ExternalReference:     externalRef,
		SailingID:             sailing.ID,
		OfferID:               sailing.OfferID,
		ClientReference:       params.ClientReference,
		PassengerCount:        params.Units,
		UnitPriceAmountMinor:  unitPrice,
		TotalPriceAmountMinor: totalPrice,
		PriceCurrency:         sailing.PriceCurrency,
		ServiceStartsAt:       sailing.DepartsAt,
		ServiceEndsAt:         sailing.ArrivesAt,
		StartTimeZone:         sailing.DepartureTimeZone,
		EndTimeZone:           sailing.ArrivalTimeZone,
		Status:                "HELD",
		CreatedAt:             nowUTC,
		ExpiresAt:             expiresAt,
		UpdatedAt:             nowUTC,
	}

	if err := tx.Create(&hold).Error; err != nil {
		return nil, fmt.Errorf("create hold: %w", err)
	}

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
		CreatedAt:          nowUTC,
		UpdatedAt:          nowUTC,
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

	// Advisory lock to serialize concurrent confirmations
	lockKey := fmt.Sprintf("%s:%s:%s", params.ClientScope, params.Operation, params.IdempotencyKey)
	if err := tx.Exec("SELECT pg_advisory_xact_lock(hashtext(?))", lockKey).Error; err != nil {
		return nil, fmt.Errorf("acquire advisory lock: %w", err)
	}

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
	externalResID := "ferry-res-" + NewRandomHex(6)
	externalRef := "FERRY-RES-" + NewRandomHex(6)

	reservation := Reservation{
		ID:                    reservationID,
		ExternalReservationID: externalResID,
		ExternalReference:     externalRef,
		HoldID:                hold.ID,
		SailingID:             hold.SailingID,
		OfferID:               hold.OfferID,
		ClientReference:       hold.ClientReference,
		PassengerCount:        hold.PassengerCount,
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

	// Advisory lock
	lockKey := fmt.Sprintf("%s:%s:%s", params.ClientScope, params.Operation, params.IdempotencyKey)
	if err := tx.Exec("SELECT pg_advisory_xact_lock(hashtext(?))", lockKey).Error; err != nil {
		return nil, fmt.Errorf("acquire advisory lock: %w", err)
	}

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
		Preload("Sailing").
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

	// Advisory lock
	lockKey := fmt.Sprintf("%s:%s:%s", params.ClientScope, params.Operation, params.IdempotencyKey)
	if err := tx.Exec("SELECT pg_advisory_xact_lock(hashtext(?))", lockKey).Error; err != nil {
		return nil, fmt.Errorf("acquire advisory lock: %w", err)
	}

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

	// Also update originating hold status to RELEASED if not already
	var hold Hold
	if err := tx.Where("id = ?", reservation.HoldID).First(&hold).Error; err == nil {
		if hold.Status != "RELEASED" {
			hold.Status = "RELEASED"
			hold.ReleasedAt = &nowUTC
			hold.UpdatedAt = nowUTC
			_ = tx.Save(&hold).Error
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
