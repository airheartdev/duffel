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

func TestCreateOffersRequest(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Post("/air/offer_requests").
		MatchParam("return_offers", "true").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-get-offer-request.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	data, err := client.CreateOfferRequest(ctx, OfferRequestInput{
		Passengers: []OfferRequestPassenger{
			{
				FamilyName: "Earhardt",
				GivenName:  "Amelia",
				Type:       PassengerTypeAdult,
			},
			{
				Age: 14,
			},
		},
		CabinClass:   CabinClassEconomy,
		ReturnOffers: true,
		Slices: []OfferRequestSlice{
			{
				DepartureDate: Date(time.Now().AddDate(0, 0, 7)),
				Origin:        "JFK",
				Destination:   "AUS",
			},
		},
	})
	a.NoError(err)
	a.NotNil(data)

	a.Equal("1390.66 GBP", data.Offers[0].TotalAmount().String())
	a.Equal("116.08 GBP", data.Offers[0].TaxAmount().String())
	a.Len(data.Slices, 1)
	a.Equal("2021-12-30", data.Slices[0].DepartureDate.String())
	a.Equal("arp_jfk_us", data.Slices[0].Origin.ID)
	a.Equal("airport", data.Slices[0].OriginType)
	a.Equal("2021-12-30", data.Slices[0].DepartureDate.String())
}

func TestGetOfferRequest(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Get("/air/offer_requests/orq_0000AEtEexyvXbB0OhB5jk").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-get-offer-request.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	data, err := client.GetOfferRequest(ctx, "orq_0000AEtEexyvXbB0OhB5jk")
	a.NoError(err)
	a.NotNil(data)
	a.Equal("1390.66 GBP", data.Offers[0].TotalAmount().String())
	a.Equal("116.08 GBP", data.Offers[0].TaxAmount().String())
	a.Equal(false, data.Offers[0].LiveMode)
	a.Equal("137", data.Offers[0].TotalEmissionsKg)
	a.Equal(false, data.Offers[0].PassengerIdentityDocumentsRequired)
	a.Equal("airport", data.Offers[0].Slices[0].DestinationType)
	a.Equal(false, data.Offers[0].Slices[0].Changeable)
	a.Equal("Refundable Main Cabin", data.Offers[0].Slices[0].FareBrandName)

	// Assert that departing at is in the correct timezone
	dep, err := data.Offers[0].Slices[0].Segments[0].DepartingAt()
	a.NoError(err)
	est, _ := time.LoadLocation("America/New_York")
	a.True(dep.Equal(time.Date(2021, time.December, 30, 8, 55, 0, 0, est)), "Departure time should be in EST")

	arr, err := data.Offers[0].Slices[0].Segments[0].ArrivingAt()
	a.NoError(err)
	a.True(arr.Equal(time.Date(2021, time.December, 30, 12, 0o7, 0, 0, est)), "Arrival time should be in EST")
}

func TestListOfferRequests(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Get("/air/offer_requests").
		// MatchParam("return_offers", "true").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-list-offer-requests.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	iter := client.ListOfferRequests(ctx)

	iter.Next()
	data := iter.Current()
	err := iter.Err()

	a.NoError(err)
	a.NotNil(data)
	a.Equal("arp_jfk_us", data.Slices[0].Origin.ID)
	a.Equal("cit_aus_us", data.Slices[0].Destination.ID)
}
