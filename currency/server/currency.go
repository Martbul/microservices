// ! gRPC is code generating and you must implement the (code-generated) interfaces
package server

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/martbul/currency/data"
	protos "github.com/martbul/currency/protos/currency"
)

// Currency is a gRPC server it implements the methods defined by the CurrencyServer interface
type Currency struct {
	protos.UnimplementedCurrencyServer // Embed the UnimplementedCurrencyServer struct
	rates                              *data.ExchangeRates
	log                                hclog.Logger
}

// NewCurrency creates a new Currency server
func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	return &Currency{rates: r, log: l}
}

// GetRate implements the CurrencyServer GetRate method and returns the currency exchange rate
// for the two given currencies.
func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle request for GetRate", "base", rr.GetBase(), "dest", rr.GetDestination())

	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil,err
	}
	return &protos.RateResponse{Rate: rate}, nil
}
