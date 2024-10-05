package currency

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// GetHistoricalRates will fetch all the latest rates for the base currency either from the store or the OXR api.
func (c *client) GetHistoricalRates(ctx context.Context, date time.Time, currencies ...string) ([]*Rate, error) {
	index := make(map[string]struct{})
	for _, currency := range currencies {
		index[currency] = struct{}{}
	}

	rates, err := c.GetAllHistoricalRates(ctx, date)
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

func (c *client) GetAllHistoricalRates(ctx context.Context, date time.Time) ([]*Rate, error) {
	const historicalAPIPath = "historical/%s.json"

	params := url.Values{
		"base": []string{c.baseCurrency},
	}

	request, err := c.newRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf(historicalAPIPath, date.Format("2006-01-02")),
		params,
	)
	if err != nil {
		return nil, err
	}

	// Make newRequest
	var resp *RateResponse

	err = c.fetch(request, &resp)
	if err != nil {
		return nil, err
	}

	var rates []*Rate

	for code, rate := range resp.Rates {
		rates = append(rates, &Rate{
			CurrencyCode: code,
			Rate:         rate,
			Timestamp:    resp.Timestamp,
		})
	}

	return rates, nil
}

// HistoricalRatesResponse holds our forex rates for a given base currency
type HistoricalRatesResponse struct {
	Rates     map[string]float64 `json:"rates"`
	Base      string             `json:"base"`
	Timestamp int64              `json:"timestamp"`
}
