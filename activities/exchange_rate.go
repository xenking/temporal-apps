package activities

import (
	"context"
	"time"

	"github.com/govalues/decimal"
)

type GetExchangeRateRequest struct {
	Date         time.Time
	CurrencyFrom string
	CurrencyTo   string
}

func (a *Activities) GetExchangeRate(ctx context.Context, in GetExchangeRateRequest) (decimal.Decimal, error) {
	rate, err := a.ConversionClient.GetConversionRate(ctx, in.Date, in.CurrencyFrom, in.CurrencyTo)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return rate, nil
}
