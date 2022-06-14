// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

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
		Iter(ctx)
}

func (a *API) GetAircraft(ctx context.Context, id string) (*Aircraft, error) {
	return newRequestWithAPI[EmptyPayload, Aircraft](a).
		Getf("/air/aircraft/%s", id).
		Single(ctx)
}

var _ AircraftClient = (*API)(nil)
