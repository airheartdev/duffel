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

func TestListOffers(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Get("/air/offers").
		MatchParam("offer_request_id", "orq_00009htyDGjIfajdNBZRlw").
		MatchParam("sort", "total_amount").
		MatchParam("max_connections", "1").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-offers-orq_0000AGqEDX9VCvWmHLBywi.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	iter := client.ListOffers(ctx, "orq_00009htyDGjIfajdNBZRlw", ListOffersParams{
		Sort:           ListOffersSortTotalAmount,
		MaxConnections: 1,
	})

	iter.Next()
	data := iter.Current()
	err := iter.Err()

	a.NoError(err)
	a.NotNil(data)

	a.Equal("228.60 USD", data.TotalAmount().String())
	a.Len(data.Slices, 1)
}

func TestGetOfferByID(t *testing.T) {
	defer gock.Off()

	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Get("/air/offers/off_00009htYpSCXrwaB9DnUm0").
		MatchParam("return_available_services", "true").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-offers-off_00009htYpSCXrwaB9DnUm0.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	data, err := client.GetOffer(ctx, "off_00009htYpSCXrwaB9DnUm0", GetOfferParams{
		ReturnAvailableServices: true,
	})
	a.NoError(err)
	a.NotNil(data)
	a.Equal("45.00 GBP", data.TotalAmount().String())
	a.Len(data.Slices, 1)
	a.False(data.PaymentRequirements.RequiresInstantPayment)
	a.Equal(time.Date(2020, 1, 17, 10, 42, 14, 0, time.UTC).Unix(), time.Time(*data.PaymentRequirements.PriceGuaranteeExpiresAt).Unix())
	a.Equal(time.Date(2020, 1, 17, 10, 42, 14, 0, time.UTC).Unix(), time.Time(*data.PaymentRequirements.PaymentRequiredBy).Unix())
}

func TestUpdateOffserPassenger(t *testing.T) {
	defer gock.Off()
	// gock.Observe(gock.DumpRequest)

	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Patch("/air/offers/orq_123/passengers/pas_00009hj8USM7Ncg31cBCL").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-update-offer-passenger.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	data, err := client.UpdateOfferPassenger(ctx, "orq_123", "pas_00009hj8USM7Ncg31cBCL", PassengerUpdateInput{
		FamilyName: "Earhardt",
		GivenName:  "Amelia",
		LoyaltyProgrammeAccounts: []LoyaltyProgrammeAccount{
			{
				AirlineIATACode: "AA",
				AccountNumber:   "AA1234567",
			},
		},
	})

	a.NoError(err)
	a.NotNil(data)

	a.Equal("pas_00009hj8USM7Ncg31cBCL", data.ID)
	a.Equal("Earhardt", data.FamilyName)
	a.Equal("Amelia", data.GivenName)
	a.Equal("adult", data.Type.String())
	a.Len(data.LoyaltyProgrammeAccounts, 1)
}

func TestListOffers_InavlidID(t *testing.T) {
	a := assert.New(t)
	ctx := context.TODO()

	client := New("duffel_test_123")
	iter := client.ListOffers(ctx, "fake-id", ListOffersParams{
		Sort:           ListOffersSortTotalAmount,
		MaxConnections: 1,
	})

	iter.Next()
	data := iter.Current()
	err := iter.Err()

	a.EqualError(err, "offerRequestId should begin with orq_")
	a.Nil(data)
}
