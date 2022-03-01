package duffel

import (
	"context"
)

type (
	AircraftClient interface {
		ListAircraft(ctx context.Context) *Iter[Aircraft]
		GetAircraft(ctx context.Context, id string) (*Aircraft, error)
	}
)

func (a *API) ListAircraft(ctx context.Context) *Iter[Aircraft] {
	return newRequestWithAPI[ListAirportsParams, Aircraft](a).
		Get("/air/aircraft").
		All(ctx)
}

func (a *API) GetAircraft(ctx context.Context, id string) (*Aircraft, error) {
	return newRequestWithAPI[EmptyPayload, Aircraft](a).
		Getf("/air/aircraft/%s", id).
		One(ctx)
}

var _ AircraftClient = (*API)(nil)
