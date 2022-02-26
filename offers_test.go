package duffel

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestListOffers(t *testing.T) {
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

	a.Equal("228.60", data.TotalAmount)
	a.Len(data.Slices, 1)

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
