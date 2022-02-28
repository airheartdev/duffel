package duffel

import (
	"context"
	"fmt"
	"net/url"
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
	return newRequestWithAPI[ListAirportsParams, Airport](a).
		Get("/air/airports").
		WithParams(normalizeParams(params)...).
		All(ctx)
}

func (a *API) GetAirport(ctx context.Context, id string) (*Airport, error) {
	return newRequestWithAPI[EmptyPayload, Airport](a).Get(fmt.Sprintf("/air/airports/%s", id)).One(ctx)
}

func (p ListAirportsParams) Encode(q url.Values) error {
	if p.IATACountryCode != "" {
		q.Set("iata_country_code", p.IATACountryCode)
	}
	return nil
}

var _ AirportsClient = (*API)(nil)
