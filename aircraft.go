package duffel

import (
	"context"
	"net/http"
)

type (
	AircraftClient interface {
		ListAircraft(ctx context.Context, params ...ListAirportsParams) *Iter[Aircraft]
		GetAircraft(ctx context.Context, id string) (*Aircraft, error)
	}
)

func (a *API) ListAircraft(ctx context.Context, params ...ListAirportsParams) *Iter[Aircraft] {
	c := newInternalClient[struct{}, Aircraft](a)
	return c.getIterator(ctx, http.MethodGet, "/air/aircraft", WithURLParams(params...))
}

func (a *API) GetAircraft(ctx context.Context, id string) (*Aircraft, error) {
	c := newInternalClient[struct{}, Aircraft](a)
	return c.makeRequestWithPayload(ctx, "/air/aircraft/"+id, http.MethodGet, nil)
}

var _ AircraftClient = (*API)(nil)
