package duffel

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestPlacesSuggestions(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)

	gock.New("https://api.duffel.com").
		Get("/places/suggestions").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-get-places-suggestion.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	places, err := client.PlaceSuggestions(ctx, "Lond")
	a.NoError(err)
	a.NotNil(places)

	a.Equal("Heathrow", places[0].Name)
	a.Equal("EGLL", places[0].ICAOCode)
	a.Equal("London", places[0].City.Name)
	a.Equal("London", places[0].CityName)
	a.Equal("Heathrow", places[0].Airports[0].Name)
}
