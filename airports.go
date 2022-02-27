package duffel

import (
	"context"
	"net/http"
)

type (
	AirportsClient interface {
		ListAirports(ctx context.Context, params ...ListAirportsParams) *Iter[Airport]
		GetAirport(ctx context.Context, id string) (*Airport, error)
	}

	ListAirportsParams struct {
		IATACountryCode string `url:"iata_country_code,omitempty"`
	}
)

func (a *API) ListAirports(ctx context.Context, params ...ListAirportsParams) *Iter[Airport] {
	c := newInternalClient[struct{}, Airport](a)
	return c.getIterator(ctx, http.MethodGet, "/air/airports", WithURLParams(params...))
}

func (a *API) GetAirport(ctx context.Context, id string) (*Airport, error) {
	c := newInternalClient[struct{}, Airport](a)
	return c.makeRequestWithPayload(ctx, "/air/airports/"+id, http.MethodGet, nil)
}

var _ AirportsClient = (*API)(nil)
