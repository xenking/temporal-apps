package middleware

import (
	"context"
	"time"

	"github.com/xenking/temporal-apps/pkg/currency"
)

type CacheProvider interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
	Delete(key string) error
}

type cachingMiddleware struct {
	cache CacheProvider
	next  currency.Client
}

func NewCaching(cache CacheProvider) currency.Middleware {
	return func(c currency.Client) currency.Client {
		return &cachingMiddleware{cache, c}
	}
}

func (mw *cachingMiddleware) GetHistoricalRates(ctx context.Context, date time.Time, currencies ...string) ([]*currency.Rate, error) {
	index := make(map[string]struct{})
	for _, code := range currencies {
		index[code] = struct{}{}
	}

	// do not call GetHistoricalRates here because we will use
	rates, err := mw.GetAllHistoricalRates(ctx, date)
	if err != nil {
		return nil, err
	}

	filtered := make([]*currency.Rate, 0, len(rates))
	for _, rate := range rates {
		if _, ok := index[rate.CurrencyCode]; ok {
			filtered = append(filtered, rate)
		}
	}

	return filtered, nil
}

func (mw *cachingMiddleware) GetAllHistoricalRates(ctx context.Context, date time.Time) ([]*currency.Rate, error) {
	key := "all_historical_rates_" + date.Format("2006-01-02")
	if rates, err := mw.cache.Get(key); err == nil {
		return rates.([]*currency.Rate), nil
	}
	rates, err := mw.next.GetAllHistoricalRates(ctx, date)
	if err != nil {
		return nil, err
	}
	if err := mw.cache.Set(key, rates); err != nil {
		return rates, err
	}
	return rates, nil
}

func (mw *cachingMiddleware) GetLatestRates(ctx context.Context, date time.Time, currencies ...string) ([]*currency.Rate, error) {
	index := make(map[string]struct{})
	for _, code := range currencies {
		index[code] = struct{}{}
	}

	rates, err := mw.GetAllLatestRates(ctx, date)
	if err != nil {
		return nil, err
	}

	filtered := make([]*currency.Rate, 0, len(rates))
	for _, rate := range rates {
		if _, ok := index[rate.CurrencyCode]; ok {
			filtered = append(filtered, rate)
		}
	}

	return filtered, nil
}

func (mw *cachingMiddleware) GetAllLatestRates(ctx context.Context, date time.Time) ([]*currency.Rate, error) {
	key := "all_latest_rates_" + date.Format("2006-01-02")
	if rates, err := mw.cache.Get(key); err == nil {
		return rates.([]*currency.Rate), nil
	}
	rates, err := mw.next.GetAllLatestRates(ctx, date)
	if err != nil {
		return nil, err
	}
	if err := mw.cache.Set(key, rates); err != nil {
		return nil, err
	}
	return rates, nil
}

func (mw *cachingMiddleware) GetCurrencies(ctx context.Context) ([]currency.Currency, error) {
	return mw.next.GetCurrencies(ctx)
}

func (mw *cachingMiddleware) BaseCurrency() string {
	return mw.next.BaseCurrency()
}
