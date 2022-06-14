// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package duffel

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestListAircraft(t *testing.T) {
	defer gock.Off()

	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Get("/air/aircraft").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-list-aircraft.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	iter := client.ListAircraft(ctx)

	iter.Next()
	aircraft := iter.Current()
	err := iter.Err()
	a.NoError(err)
	a.NotNil(aircraft)

	a.Equal("arc_00009UhD4ongolulWd91Ky", aircraft.ID)
	a.Equal("380", aircraft.IATACode)
	a.Equal("Airbus Industries A380", aircraft.Name)
}

func TestGetAircraftByID(t *testing.T) {
	defer gock.Off()

	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Get("/air/aircraft/arc_00009UhD4ongolulWd91Ky").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-get-aircraft.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	aircraft, err := client.GetAircraft(ctx, "arc_00009UhD4ongolulWd91Ky")
	a.NoError(err)
	a.NotNil(aircraft)
	a.Equal("arc_00009UhD4ongolulWd91Ky", aircraft.ID)
	a.Equal("380", aircraft.IATACode)
	a.Equal("Airbus Industries A380", aircraft.Name)
}
