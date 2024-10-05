package currency

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/govalues/decimal"
)

func (c *client) GetLatestRates(ctx context.Context, date time.Time, currencies ...string) ([]*Rate, error) {
	index := make(map[string]struct{})
	for _, currency := range currencies {
		index[currency] = struct{}{}
	}

	rates, err := c.GetAllLatestRates(ctx, date)
	if err != nil {
		return nil, err
	}

	var result []*Rate
	for _, rate := range rates {
		if _, ok := index[rate.CurrencyCode]; ok {
			result = append(result, rate)
		}
	}

	return result, nil
}

// GetAllLatestRates will fetch all rates for the base currency for the given time.Time.
func (c *client) GetAllLatestRates(ctx context.Context, date time.Time) ([]*Rate, error) {
	const latestAPIPath = "latest.json"

	params := url.Values{
		"base": []string{c.baseCurrency},
	}

	request, err := c.newRequest(
		ctx,
		http.MethodGet,
		latestAPIPath,
		params,
	)
	if err != nil {
		return nil, err
	}

	response := &RateResponse{}

	err = c.fetch(request, response)
	if err != nil {
		return nil, err
	}

	var rates []*Rate
	for currency, rate := range response.Rates {
		rates = append(rates, &Rate{
			CurrencyCode: currency,
			Rate:         rate,
			Timestamp:    response.Timestamp,
		})
	}

	return rates, nil
}

// RateResponse holds our forex rates for a given base currency
type RateResponse struct {
	Rates        map[string]decimal.Decimal `json:"rates"`
	BaseCurrency string                     `json:"base"`
	Timestamp    int64                      `json:"timestamp"`
}
