package service

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"batam-medhub/internal/model"

	"gorm.io/gorm"
)

var (
	ErrUnsupportedCurrency = errors.New("unsupported currency")
	ErrFXRateNotFound      = errors.New("exchange rate not found")
	ErrInvalidAmount       = errors.New("amount must be non-negative")
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

// MoneyService performs currency conversions using static exchange rates and exact rational math.
type MoneyService struct {
	db *gorm.DB
}

// NewMoneyService constructs a MoneyService with the provided database handle.
func NewMoneyService(db *gorm.DB) *MoneyService {
	return &MoneyService{db: db}
}

// canonicalizeFXRate formats the exchange rate string to preserve at least 6 decimal places
// while trimming any superfluous trailing zeros beyond 6 decimal digits.
func canonicalizeFXRate(raw string) string {
	raw = strings.TrimSpace(raw)
	dotIdx := strings.IndexByte(raw, '.')
	if dotIdx == -1 {
		return raw + ".000000"
	}
	intPart := raw[:dotIdx]
	fracPart := raw[dotIdx+1:]
	if len(fracPart) < 6 {
		fracPart = fracPart + strings.Repeat("0", 6-len(fracPart))
		return intPart + "." + fracPart
	}
	trimmedFrac := strings.TrimRight(fracPart[6:], "0")
	return intPart + "." + fracPart[:6] + trimmedFrac
}

// Convert converts source money into target display currency using exact rational decimal arithmetic against seeded FX rows.
func (s *MoneyService) Convert(ctx context.Context, source Money, targetCurrency string) (*ConvertedMoney, error) {
	if source.AmountMinor < 0 {
		return nil, fmt.Errorf("%w: %d", ErrInvalidAmount, source.AmountMinor)
	}
	if targetCurrency != "SGD" && targetCurrency != "IDR" {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedCurrency, targetCurrency)
	}
	if source.Currency != "SGD" && source.Currency != "IDR" {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedCurrency, source.Currency)
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

	rateRat := new(big.Rat)
	if _, ok := rateRat.SetString(fx.Rate); !ok {
		return nil, fmt.Errorf("parse fx rate rational: %s", fx.Rate)
	}

	amountRat := new(big.Rat).SetInt64(source.AmountMinor)
	productRat := new(big.Rat).Mul(amountRat, rateRat)

	num := productRat.Num()
	denom := productRat.Denom()

	// Deterministic half-up rounding without floating-point math: (2*num + denom) / (2*denom)
	two := big.NewInt(2)
	scaledNum := new(big.Int).Mul(num, two)
	scaledDenom := new(big.Int).Mul(denom, two)

	if num.Sign() >= 0 {
		scaledNum.Add(scaledNum, denom)
	} else {
		scaledNum.Sub(scaledNum, denom)
	}

	roundedInt := new(big.Int).Quo(scaledNum, scaledDenom)
	if !roundedInt.IsInt64() {
		return nil, fmt.Errorf("converted amount minor overflow: %s", roundedInt.String())
	}

	return &ConvertedMoney{
		Source: source,
		Display: Money{
			AmountMinor: roundedInt.Int64(),
			Currency:    targetCurrency,
		},
		FXRate:        canonicalizeFXRate(fx.Rate),
		FXSource:      fx.Source,
		FXEffectiveAt: fx.EffectiveAt,
		Estimated:     true,
	}, nil
}
