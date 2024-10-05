package currency

import (
	"context"
	"net/http"
	"net/url"
)

func (c *client) GetCurrencies(ctx context.Context) ([]Currency, error) {
	const currenciesAPIPath = "currencies.json"

	req, err := c.newRequest(ctx,
		http.MethodGet,
		currenciesAPIPath,
		url.Values{},
	)
	if err != nil {
		return nil, err
	}

	resp := make(map[string]string)
	if err = c.fetch(req, &resp); err != nil {
		return nil, err
	}

	var currencies []Currency
	for code, name := range resp {
		currencies = append(currencies, Currency{
			Code: code,
			Name: name,
		})
	}

	return currencies, nil
}
