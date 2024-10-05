package currency

import (
	"net/http"
	"time"
)

type options struct {
	httpClient   *http.Client
	backendURL   string
	baseCurrency string
	userAgent    string
}

const backendURL = "https://openexchangerates.org"

type Option func(*options)

func newOptions(opts ...Option) options {
	o := &options{
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		backendURL:   backendURL,
		baseCurrency: "USD",
		userAgent:    "currency-api-client/1.0",
	}
	for _, opt := range opts {
		opt(o)
	}
	return *o
}

func WithHTTPClient(c *http.Client) Option {
	return func(o *options) {
		o.httpClient = c
	}
}

func WithBackendURL(url string) Option {
	return func(o *options) {
		o.backendURL = url
	}
}

func WithBaseCurrency(bc string) Option {
	return func(o *options) {
		o.baseCurrency = bc
	}
}

func WithUserAgent(ua string) Option {
	return func(o *options) {
		o.userAgent = ua
	}
}
