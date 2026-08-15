package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"batam-medhub/internal/model"

	"gorm.io/gorm"
)

var (
	ErrUnsupportedCurrency = errors.New("unsupported currency")
	ErrFXRateNotFound      = errors.New("exchange rate not found")
)

// Money represents an amount in minor units (e.g. cents) with an ISO currency code.
type Money struct {
	AmountMinor int64  `json:"amount_minor"`
	Currency    string `json:"currency"`
}

// ConvertedMoney represents money converted to a display currency with conversion metadata.
type ConvertedMoney struct {
	Source        Money     `json:"source"`
	Display       Money     `json:"display"`
	FXRate        string    `json:"fx_rate"`
	FXSource      string    `json:"fx_source"`
	FXEffectiveAt time.Time `json:"fx_effective_at"`
	Estimated     bool      `json:"estimated"`
}

// MoneyService performs currency conversions using static exchange rates.
type MoneyService struct {
	db *gorm.DB
}

// NewMoneyService constructs a MoneyService with the provided database handle.
func NewMoneyService(db *gorm.DB) *MoneyService {
	return &MoneyService{db: db}
}

// Convert converts source money into target display currency using stored static FX rates.
func (s *MoneyService) Convert(ctx context.Context, source Money, targetCurrency string) (*ConvertedMoney, error) {
	if targetCurrency != "SGD" && targetCurrency != "IDR" {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedCurrency, targetCurrency)
	}

	if source.Currency == targetCurrency {
		now := time.Now().UTC()
		return &ConvertedMoney{
			Source:        source,
			Display:       source,
			FXRate:        "1.000000",
			FXSource:      "IDENTITY",
			FXEffectiveAt: now,
			Estimated:     true,
		}, nil
	}

	var fx model.FXRate
	err := s.db.WithContext(ctx).
		Where("base_currency = ? AND quote_currency = ?", source.Currency, targetCurrency).
		Order("effective_at DESC").
		First(&fx).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %s to %s", ErrFXRateNotFound, source.Currency, targetCurrency)
		}
		return nil, fmt.Errorf("lookup fx rate: %w", err)
	}

	rateFloat, err := strconv.ParseFloat(fx.Rate, 64)
	if err != nil {
		return nil, fmt.Errorf("parse fx rate: %w", err)
	}

	displayAmountMinor := int64(math.Round(float64(source.AmountMinor) * rateFloat))

	return &ConvertedMoney{
		Source: source,
		Display: Money{
			AmountMinor: displayAmountMinor,
			Currency:    targetCurrency,
		},
		FXRate:        fx.Rate,
		FXSource:      fx.Source,
		FXEffectiveAt: fx.EffectiveAt,
		Estimated:     true,
	}, nil
}
