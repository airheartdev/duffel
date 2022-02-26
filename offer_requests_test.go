package duffel

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestCreateOffersRequest(t *testing.T) {
	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Post("/air/offer_requests").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-offer-response.json")

	ctx := context.TODO()

	adult := PassengerTypeAdult
	age := 30

	client := New("duffel_test_123")
	data, err := client.CreateOfferRequest(ctx, OfferRequestInput{
		Passengers: []Passenger{
			{
				ID:         "1",
				FamilyName: "Earhardt",
				GivenName:  "Amelia",
				Age:        &age,
				Type:       &adult,
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

	a.Equal("1390.66", data.Offers[0].TotalAmount)
	a.Equal("GBP", data.Offers[0].TotalCurrency)
	a.Equal("GBP", data.Offers[0].TaxCurrency)
	a.Equal("116.08", data.Offers[0].TaxAmount)
}

func TestGetOfferRequest(t *testing.T) {
	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Get("/air/offer_requests/orq_0000AEtEexyvXbB0OhB5jk").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-offer-response.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	data, err := client.GetOfferRequest(ctx, "orq_0000AEtEexyvXbB0OhB5jk")
	a.NoError(err)
	a.NotNil(data)
	a.Equal("1390.66", data.Offers[0].TotalAmount)
	a.Equal("GBP", data.Offers[0].TotalCurrency)
	a.Equal("GBP", data.Offers[0].TaxCurrency)
	a.Equal("116.08", data.Offers[0].TaxAmount)
}

func TestListOfferRequests(t *testing.T) {
	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Get("/air/offer_requests").
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
