package currency

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/govalues/decimal"
)

type Client interface {
	GetCurrencies(ctx context.Context) ([]Currency, error)
	GetHistoricalRates(ctx context.Context, date time.Time, currencies ...string) ([]*Rate, error)
	GetAllHistoricalRates(ctx context.Context, date time.Time) ([]*Rate, error)
	GetLatestRates(ctx context.Context, date time.Time, currencies ...string) ([]*Rate, error)
	GetAllLatestRates(ctx context.Context, date time.Time) ([]*Rate, error)
	BaseCurrency() string
}

type Middleware func(next Client) Client

type Rate struct {
	Rate         decimal.Decimal `json:"rate"`
	CurrencyCode string          `json:"currency"`
	Timestamp    int64           `json:"timestamp"`
}

// Currency represents a currency from OXR.
type Currency struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// client holds a connection to the OXR API.
type client struct {
	appID string
	options
}

// NewClient creates a new Client with the appropriate connection details and
// services used for communicating with the API.
func NewClient(appID string, opts ...Option) Client {
	return &client{
		appID:   appID,
		options: newOptions(opts...),
	}
}

func (c *client) BaseCurrency() string {
	return c.baseCurrency
}

// newRequest creates an authenticated API Request
func (c *client) newRequest(ctx context.Context, method, urlPath string, params url.Values) (*http.Request, error) {
	params.Set("app_id", c.appID)

	urlPath = fmt.Sprintf("%s/api/%s?%s", c.backendURL, urlPath, params.Encode())

	req, err := http.NewRequestWithContext(ctx, method, urlPath, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", c.userAgent)

	return req, nil
}

// fetch send an API newRequest and returns the API response. The API response is
// JSON decoded and stored in 'v', or returned as an error if an API (if found).
func (c *client) fetch(req *http.Request, v interface{}) (err error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	// Check for any errors that may have occurred.
	if status := resp.StatusCode; status >= http.StatusBadRequest {
		errorResponse := &ErrorResponse{Response: resp}

		if err = json.NewDecoder(resp.Body).Decode(errorResponse); err != nil {
			return err
		}

		return errorResponse
	}

	return json.NewDecoder(resp.Body).Decode(&v)
}

// An ErrorResponse reports the error caused by an API newRequest
type ErrorResponse struct {
	*http.Response
	ErrorCode   int64  `json:"status"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%d %v", r.Response.StatusCode, r.Description)
}
