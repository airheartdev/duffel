package duffel

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestOffersRequest(t *testing.T) {
	a := assert.New(t)
	gock.New("https://api.duffel.com/air/offer_requests").
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
	// client := &internalClient[OfferRequestInput, OfferResponse]{client: c}
	data, err := client.CreateOfferRequest(ctx, &OfferRequestInput{
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
