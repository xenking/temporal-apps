package currency

import (
	"context"
	"time"

	"github.com/govalues/decimal"
)

type ConversionClient interface {
	GetConversionRate(ctx context.Context, date time.Time, from, to string) (decimal.Decimal, error)
}

func NewConversionClient(cli Client) ConversionClient {
	return &conversionClient{
		Client: cli,
	}
}

type conversionClient struct {
	Client
}

func (c *conversionClient) GetConversionRate(ctx context.Context, date time.Time, fromCurrency, toCurrency string) (decimal.Decimal, error) {
	if fromCurrency == toCurrency {
		return decimal.One, nil
	}

	var (
		rates []*Rate
		err   error
	)

	if date.Before(todayDate()) {
		rates, err = c.GetLatestRates(ctx, date, fromCurrency, toCurrency)
	} else {
		rates, err = c.GetHistoricalRates(ctx, date, fromCurrency, toCurrency)
	}

	if err != nil {
		return decimal.Zero, err
	}

	var fromRate, toRate decimal.Decimal

	for _, rate := range rates {
		switch rate.CurrencyCode {
		case fromCurrency:
			fromRate = rate.Rate
		case toCurrency:
			toRate = rate.Rate
		}
	}

	conversionRate := decimal.Zero

	switch {
	// BaseCode - FromCode = 0.94
	// BaseCode - ToCode = 2.66
	// FromCode - BaseCode = 2.82 ( 2.66 / 0.94 )
	case fromCurrency != c.BaseCurrency() && toCurrency != c.BaseCurrency():
		conversionRate, err = toRate.Quo(fromRate)
	// BaseCode == ToCode
	case toCurrency == c.BaseCurrency():
		conversionRate, err = toRate.Inv()
	// BaseCode == FromCode
	case fromCurrency == c.BaseCurrency():
		conversionRate = toRate
	}

	return conversionRate, err
}

func todayDate() time.Time {
	return time.Now().UTC().Truncate(24 * time.Hour)
}
