package hotel

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

type OfferWithUnits struct {
	Offer          HotelOffer
	AvailableUnits int
	UnitPriceMinor int64
	PriceCurrency  string
}

func (r *Repository) SearchAvailableOffers(
	ctx context.Context,
	checkInDate time.Time,
	checkOutDate time.Time,
	now time.Time,
	roomCount int,
	guestCount int,
	accessibility []string,
) ([]OfferWithUnits, error) {
	var offers []HotelOffer
	err := r.db.WithContext(ctx).
		Preload("RoomType").
		Preload("OfferNights").
		Preload("OfferNights.InventoryDay").
		Where("check_in_date = ? AND check_out_date = ? AND status = 'AVAILABLE' AND valid_until > ?",
			checkInDate.Format("2006-01-02"), checkOutDate.Format("2006-01-02"), now).
		Find(&offers).Error
	if err != nil {
		return nil, fmt.Errorf("query hotel offers: %w", err)
	}

	var results []OfferWithUnits
	for _, offer := range offers {
		if offer.RoomType == nil {
			continue
		}
		if offer.RoomType.MaxGuestsPerRoom*roomCount < guestCount {
			continue
		}

		if len(accessibility) > 0 {
			if !hasAllAccessibility(offer.RoomType.Accessibility, accessibility) {
				continue
			}
		}

		minAvailable := -1
		var unitPriceTotal int64
		currency := ""

		for _, night := range offer.OfferNights {
			if night.InventoryDay == nil {
				continue
			}
			if currency == "" {
				currency = night.InventoryDay.PriceCurrency
			}
			unitPriceTotal += night.InventoryDay.PriceAmountMinor

			avail, err := r.calculateAvailableRoomsOnDay(ctx, r.db, night.InventoryDayID, night.InventoryDay.RoomsTotal, now)
			if err != nil {
				return nil, fmt.Errorf("calculate available rooms on day: %w", err)
			}
			if minAvailable == -1 || avail < minAvailable {
				minAvailable = avail
			}
		}

		if minAvailable < 0 {
			minAvailable = 0
		}

		if minAvailable >= roomCount {
			results = append(results, OfferWithUnits{
				Offer:          offer,
				AvailableUnits: minAvailable,
				UnitPriceMinor: unitPriceTotal,
				PriceCurrency:  currency,
			})
		}
	}

	return results, nil
}

func (r *Repository) calculateAvailableRoomsOnDay(
	ctx context.Context,
	tx *gorm.DB,
	inventoryDayID string,
	roomsTotal int,
	now time.Time,
) (int, error) {
	var activeHoldsSum int64
	err := tx.WithContext(ctx).
		Table("hold_nights hn").
		Joins("JOIN holds h ON h.id = hn.hold_id").
		Where("hn.inventory_day_id = ? AND h.status = 'HELD' AND h.expires_at > ?", inventoryDayID, now).
		Select("COALESCE(SUM(hn.room_count), 0)").
		Scan(&activeHoldsSum).Error
	if err != nil {
		return 0, fmt.Errorf("sum active hold nights: %w", err)
	}

	var confirmedResSum int64
	err = tx.WithContext(ctx).
		Table("hold_nights hn").
		Joins("JOIN holds h ON h.id = hn.hold_id").
		Where("hn.inventory_day_id = ? AND h.status = 'CONFIRMED'", inventoryDayID).
		Select("COALESCE(SUM(hn.room_count), 0)").
		Scan(&confirmedResSum).Error
	if err != nil {
		return 0, fmt.Errorf("sum confirmed reservation nights: %w", err)
	}

	available := roomsTotal - int(activeHoldsSum) - int(confirmedResSum)
	if available < 0 {
		available = 0
	}
	return available, nil
}

func hasAllAccessibility(available []string, requested []string) bool {
	availMap := make(map[string]bool, len(available))
	for _, a := range available {
		availMap[a] = true
	}
	for _, req := range requested {
		if !availMap[req] {
			return false
		}
	}
	return true
}

type CreateHoldParams struct {
	ClientScope        string
	Operation          string
	IdempotencyKey     string
	RequestFingerprint string
	OfferID            string
	RoomCount          int
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

	// Lock the hotel offer row
	var offer HotelOffer
	err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("RoomType").
		Preload("OfferNights").
		Preload("OfferNights.InventoryDay").
		Where("external_offer_id = ?", params.OfferID).
		First(&offer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if offer.Status != "AVAILABLE" {
		return &CreateHoldResult{
			ConflictDetails: &CapacityConflictDetails{
				Requested: params.RoomCount,
				Available: 0,
			},
		}, ErrCapacityConflict
	}

	if !offer.ValidUntil.After(params.Now) {
		return &CreateHoldResult{
			OfferExpiredDetails: &OfferExpiredDetails{
				ValidUntil: offer.ValidUntil,
			},
		}, ErrOfferExpired
	}

	// Lock each inventory day row and check prices & capacity
	var unitPriceTotal int64
	currency := ""
	minAvailable := -1

	for _, night := range offer.OfferNights {
		var invDay RoomInventoryDay
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", night.InventoryDayID).
			First(&invDay).Error
		if err != nil {
			return nil, fmt.Errorf("lock inventory day: %w", err)
		}

		if currency == "" {
			currency = invDay.PriceCurrency
		}
		unitPriceTotal += invDay.PriceAmountMinor

		avail, err := r.calculateAvailableRoomsOnDay(ctx, tx, invDay.ID, invDay.RoomsTotal, params.Now)
		if err != nil {
			return nil, fmt.Errorf("calculate available rooms on day: %w", err)
		}
		if minAvailable == -1 || avail < minAvailable {
			minAvailable = avail
		}
	}

	if unitPriceTotal != params.ExpectedUnitPrice || currency != params.ExpectedCurrency {
		return nil, ErrOfferChanged
	}

	if minAvailable < params.RoomCount {
		if minAvailable < 0 {
			minAvailable = 0
		}
		return &CreateHoldResult{
			ConflictDetails: &CapacityConflictDetails{
				Requested: params.RoomCount,
				Available: minAvailable,
			},
		}, ErrCapacityConflict
	}

	nowUTC := params.Now.UTC()
	expiresAt := nowUTC.Add(params.HoldDuration)

	holdID := NewUUID()
	externalHoldID := "hotel-hold-" + NewRandomHex(6)
	externalRef := "HOTEL-HOLD-" + NewRandomHex(6)
	totalPrice := unitPriceTotal * int64(params.RoomCount)

	accessibility := StringArray{}
	if offer.RoomType != nil {
		accessibility = offer.RoomType.Accessibility
	}

	hold := Hold{
		ID:                    holdID,
		ExternalHoldID:        externalHoldID,
		ExternalReference:     externalRef,
		OfferPK:               offer.ID,
		OfferID:               offer.ExternalOfferID,
		RoomTypeID:            offer.RoomTypeID,
		ClientReference:       params.ClientReference,
		CheckInDate:           offer.CheckInDate,
		CheckOutDate:          offer.CheckOutDate,
		TimeZone:              offer.StartTimeZone,
		RoomCount:             params.RoomCount,
		Accessibility:         accessibility,
		UnitPriceAmountMinor:  unitPriceTotal,
		TotalPriceAmountMinor: totalPrice,
		PriceCurrency:         currency,
		ServiceStartsAt:       offer.ServiceStartsAt,
		ServiceEndsAt:         offer.ServiceEndsAt,
		StartTimeZone:         offer.StartTimeZone,
		EndTimeZone:           offer.EndTimeZone,
		Status:                "HELD",
		CreatedAt:             nowUTC,
		ExpiresAt:             expiresAt,
		UpdatedAt:             nowUTC,
	}

	if err := tx.Create(&hold).Error; err != nil {
		return nil, fmt.Errorf("create hold: %w", err)
	}

	for _, night := range offer.OfferNights {
		holdNight := HoldNight{
			HoldID:         hold.ID,
			OfferPK:        offer.ID,
			InventoryDayID: night.InventoryDayID,
			RoomCount:      params.RoomCount,
		}
		if err := tx.Create(&holdNight).Error; err != nil {
			return nil, fmt.Errorf("create hold night: %w", err)
		}
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
	externalResID := "hotel-res-" + NewRandomHex(6)
	externalRef := "HOTEL-RES-" + NewRandomHex(6)

	reservation := Reservation{
		ID:                    reservationID,
		ExternalReservationID: externalResID,
		ExternalReference:     externalRef,
		HoldID:                hold.ID,
		OfferPK:               hold.OfferPK,
		OfferID:               hold.OfferID,
		ClientReference:       hold.ClientReference,
		RoomCount:             hold.RoomCount,
		TotalPriceAmountMinor: hold.TotalPriceAmountMinor,
		PriceCurrency:         hold.PriceCurrency,
		ServiceStartsAt:       hold.ServiceStartsAt,
		ServiceEndsAt:         hold.ServiceEndsAt,
		StartTimeZone:         hold.StartTimeZone,
		EndTimeZone:           hold.EndTimeZone,
		Status:                "CONFIRMED",
		ConfirmedAt:           nowUTC,
		ReleasedAt:            nil,
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
		Preload("Hold.RoomType").
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
