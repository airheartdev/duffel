package duffel

import (
	"context"
	"net/http"
)

type (
	AirlinesClient interface {
		ListAirlines(ctx context.Context, params ...ListAirportsParams) *Iter[Airline]
		GetAirline(ctx context.Context, id string) (*Airline, error)
	}
)

func (a *API) ListAirlines(ctx context.Context, params ...ListAirportsParams) *Iter[Airline] {
	c := newInternalClient[struct{}, Airline](a)
	return c.getIterator(ctx, http.MethodGet, "/air/airlines", WithURLParams(params...))
}

func (a *API) GetAirline(ctx context.Context, id string) (*Airline, error) {
	c := newInternalClient[struct{}, Airline](a)
	return c.makeRequestWithPayload(ctx, "/air/airlines/"+id, http.MethodGet, nil)
}

var _ AirlinesClient = (*API)(nil)
