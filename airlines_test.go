package duffel

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestListAirlines(t *testing.T) {
	defer gock.Off()

	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Get("/air/airlines").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-list-airlines.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	iter := client.ListAirlines(ctx)

	iter.Next()
	airline := iter.Current()
	err := iter.Err()
	a.NoError(err)
	a.NotNil(airline)

	a.Equal("aln_00001876aqC8c5umZmrRds", airline.ID)
	a.Equal("BA", airline.IATACode)
	a.Equal("British Airways", airline.Name)
}

func TestGetAirlineByID(t *testing.T) {
	defer gock.Off()

	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Get("/air/airlines/aln_00001876aqC8c5umZmrRds").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-get-airline.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	airline, err := client.GetAirline(ctx, "aln_00001876aqC8c5umZmrRds")
	a.NoError(err)
	a.NotNil(airline)
	a.Equal("aln_00001876aqC8c5umZmrRds", airline.ID)
}
