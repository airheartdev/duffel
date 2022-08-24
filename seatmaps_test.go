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

func TestGetSeatmaps(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Get("/air/seat_maps").
		MatchParam("offer_id", "off_00009htYpSCXrwaB9DnUm0").
		Reply(200).
		SetHeader(RequestIDHeader, "FvxRwfnMtKgc0EwCCoXE").
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-get-seatmap.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	seats, err := client.GetSeatmap(ctx, "off_00009htYpSCXrwaB9DnUm0")
	a.NoError(err)
	reqID, _ := client.LastRequestID()

	a.NotNil(seats)
	a.Equal("FvxRwfnMtKgc0EwCCoXE", reqID)

	seat := seats[0]
	a.Equal("sea_00003hthlsHZ8W4LxXjkzo", seat.ID)
	a.Equal("seg_00009htYpSCXrwaB9Dn456", seat.SegmentID)
	a.Equal("sli_00009htYpSCXrwaB9Dn123", seat.SliceID)

	serviceAmount := seat.Cabins[0].Rows[0].Sections[0].Elements[0].AvailableServices[0].TotalAmount().String()
	a.Equal("30.00 GBP", serviceAmount)
}
