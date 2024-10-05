package activities

import (
	"github.com/xenking/temporal-apps/pkg/currency"
	"go.temporal.io/sdk/client"
)

type Activities struct {
	Client           client.Client
	ConversionClient currency.ConversionClient
}
