// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package duffel

import (
	"context"
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
		Iter(ctx)
}

func (a *API) GetAirport(ctx context.Context, id string) (*Airport, error) {
	return newRequestWithAPI[EmptyPayload, Airport](a).
		Getf("/air/airports/%s", id).
		Single(ctx)
}

func (p ListAirportsParams) Encode(q url.Values) error {
	if p.IATACountryCode != "" {
		q.Set("iata_country_code", p.IATACountryCode)
	}
	return nil
}

var _ AirportsClient = (*API)(nil)
