// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package duffel

import (
	"context"
	"testing"
	"time"

	"github.com/segmentio/encoding/json"
	"gopkg.in/h2non/gock.v1"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalOrderWithCurrency(t *testing.T) {
	a := assert.New(t)
	order := new(Order)

	err := json.Unmarshal([]byte(`{"base_amount": "120.00", "base_currency": "USD", "cancelled_at":"2020-04-11T15:48:11.642Z"}`), order)
	a.NoError(err)

	a.Equal("120.00 USD", order.BaseAmount().String())
	a.Equal("2020-04-11 15:48:11.642 +0000 UTC", order.CancelledAt.String())
}

func TestCreateOrder(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)

	// gock.Observe(gock.DumpRequest)

	gock.New("https://api.duffel.com").
		Post("/air/orders").
		File("fixtures/201-create-order-input.json").
		Reply(201).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/201-create-order.json")

	ctx := context.TODO()
	client := New("duffel_test_123")
	order, err := client.CreateOrder(ctx, CreateOrderInput{
		Type:     OrderTypeInstant,
		Metadata: Metadata{"seat_preference": "isle", "meal_preference": "NLML", "payment_intent_id": "pit_00009htYpSCXrwaB9DnUm2"},
		Services: []ServiceCreateInput{
			{
				ID:       "ase_00009hj8USM7Ncg31cB123",
				Quantity: 1,
			},
		},
		SelectedOffers: []string{"off_00009htyDGjIfajdNBZRlw"},
		Payments: []PaymentCreateInput{
			{
				Type:     "balance",
				Currency: "GBP",
				Amount:   "30.20",
			},
		},
		Passengers: []OrderPassenger{
			{
				Type:              PassengerTypeAdult,
				ID:                "pas_00009hj8USM7Ncg31cBCLL",
				Title:             PassengerTitleMrs,
				FamilyName:        "Earhart",
				GivenName:         "Amelia",
				BornOn:            Date(time.Date(1987, time.July, 24, 0, 0, 0, 0, time.UTC)),
				Gender:            GenderFemale,
				InfantPassengerID: "pas_00009hj8USM8Ncg32aTGHL",
				PhoneNumber:       "+442080160509",
				Email:             "amelia@duffel.com",
				IdentityDocuments: []IdentityDocument{
					{
						UniqueIdentifier:   "19KL56147",
						ExpiresOn:          Date(time.Date(2025, time.April, 25, 0, 0, 0, 0, time.UTC)),
						IssuingCountryCode: "GB",
						Type:               "passport",
					},
				},
			},
		},
	})
	a.NoError(err)
	a.Equal("RZPNX8", order.BookingReference)
	a.Equal("ord_00009hthhsUZ8W4LxQgkjo", order.ID)
	a.Equal("100.00 GBP", order.Conditions.ChangeBeforeDeparture.PenaltyAmount().String())
	a.Equal("100.00 GBP", order.Conditions.RefundBeforeDeparture.PenaltyAmount().String())
	a.Len(order.Slices, 1)
	a.Len(order.Slices[0].Segments, 1)
	a.Len(order.Slices[0].Segments[0].Passengers, 1)
	a.Equal("passenger_0", order.Slices[0].Segments[0].Passengers[0].ID)
	a.Equal(CabinClassEconomy, order.Slices[0].Segments[0].Passengers[0].CabinClass)
	a.Equal("Economy Basic", order.Slices[0].Segments[0].Passengers[0].CabinClassMarketingName)
	a.Equal("14B", order.Slices[0].Segments[0].Passengers[0].Seat.Designator)
	a.Equal("Exit row seat", order.Slices[0].Segments[0].Passengers[0].Seat.Name)
	a.Equal([]string{"Do not seat children in exit row seats", "Do not seat passengers with special needs in exit row seats"}, order.Slices[0].Segments[0].Passengers[0].Seat.Disclosures)
}

func TestListOrders(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)

	gock.New("https://api.duffel.com").
		Get("/air/orders").
		MatchParam("after", "g2wAAAACbQAAABBBZXJvbWlzdC1LaGFya2l2bQAAAB=").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-list-orders-page2.json")

	gock.New("https://api.duffel.com").
		Get("/air/orders").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-list-orders.json")

	ctx := context.TODO()
	client := New("duffel_test_123")
	iter := client.ListOrders(ctx, ListOrdersParams{
		BookingReference: "RZPNX8",
	})

	iter.Next()
	a.NoError(iter.Err())
	order1 := iter.Current()

	a.Equal("RZPNX8", order1.BookingReference)
	a.Equal("ord_00009hthhsUZ8W4LxQgkjo", order1.ID)

	a.Equal(&ListMeta{After: "g2wAAAACbQAAABBBZXJvbWlzdC1LaGFya2l2bQAAAB=", Limit: 50}, iter.Meta())

	iter.Next()
	a.NoError(iter.Err())
	order2 := iter.Current()

	a.Len(iter.List().GetItems(), 1, "iterator has 1 item on this page")
	a.Equal(&ListMeta{Limit: 50}, iter.Meta())

	a.Equal("ABC123", order2.BookingReference)
	a.Equal("ord_00009hthhsUZ8W4LxQgkjo", order2.ID)
}

func TestGetOrderByID(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)

	gock.New("https://api.duffel.com").
		Get("/air/orders/ord_00009hthhsUZ8W4LxQgkjo").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-get-order.json")

	ctx := context.TODO()
	client := New("duffel_test_123")
	order, err := client.GetOrder(ctx, "ord_00009hthhsUZ8W4LxQgkjo")
	a.NoError(err)

	a.Equal("RZPNX8", order.BookingReference)
	a.Equal("ord_00009hthhsUZ8W4LxQgkjo", order.ID)
}

func TestUpdateOrder(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)

	gock.New("https://api.duffel.com").
		Patch("/air/orders/ord_00009hthhsUZ8W4LxQgkjo").
		MatchType("json").
		JSON(Payload[OrderUpdateParams]{Data: OrderUpdateParams{
			Metadata: map[string]any{"seat_preference": "window"},
		}}).
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-update-order.json")

	ctx := context.TODO()
	client := New("duffel_test_123")
	order, err := client.UpdateOrder(ctx, "ord_00009hthhsUZ8W4LxQgkjo", OrderUpdateParams{
		Metadata: map[string]any{
			"seat_preference": "window",
		},
	})
	a.NoError(err)

	a.Equal("RZPNX8", order.BookingReference)
	a.Equal(Metadata{"seat_preference": "window"}, order.Metadata)
	a.Equal("ord_00009hthhsUZ8W4LxQgkjo", order.ID)
}
