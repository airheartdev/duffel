package duffel

import (
	"context"
	"net/http"
)

type (
	AirportsClient interface {
		ListAirports(ctx context.Context, params ...ListAirportsParams) *Iter[Airport]
	}

	ListAirportsParams struct {
		IATACountryCode string `url:"iata_country_code,omitempty"`
	}
)

func (a *API) ListAirports(ctx context.Context, params ...ListAirportsParams) *Iter[Airport] {
	c := newInternalClient[struct{}, Airport](a)
	return c.getIterator(ctx, http.MethodGet, "/air/airports", WithURLParams(params...))
}

var _ AirportsClient = (*API)(nil)
