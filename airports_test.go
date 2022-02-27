package duffel

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestListAirports(t *testing.T) {
	defer gock.Off()

	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Get("/air/airports").
		MatchParam("iata_country_code", "GB").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-list-airports.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	iter := client.ListAirports(ctx, ListAirportsParams{
		IATACountryCode: "GB",
	})

	iter.Next()
	airport := iter.Current()
	err := iter.Err()
	a.NoError(err)
	a.NotNil(airport)

	a.Equal("arp_lhr_gb", airport.ID)
}

func TestGetAirportByID(t *testing.T) {
	defer gock.Off()

	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Get("/air/airports/arp_lhr_gb").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-get-airport.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	airport, err := client.GetAirport(ctx, "arp_lhr_gb")
	a.NoError(err)
	a.NotNil(airport)
	a.Equal("arp_lhr_gb", airport.ID)
}
