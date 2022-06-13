package duffel

import "context"

type (
	PlacesClient interface {
		PlaceSuggestions(ctx context.Context, query string) ([]*Place, error)
	}

	Place struct {
		ID              string     `json:"id"`
		Airports        []*Airport `json:"airports"`
		City            *City      `json:"city"`
		CityName        string     `json:"city_name"`
		CountryName     string     `json:"country_name"`
		IATACityCode    string     `json:"iata_city_code"`
		IATACode        string     `json:"iata_code"`
		IATACountryCode string     `json:"iata_country_code"`
		ICAOCode        string     `json:"icao_code"`
		Latitude        float64    `json:"latitude"`
		Longitude       float64    `json:"longitude"`
		Name            string     `json:"name"`
		TimeZone        string     `json:"time_zone"`
		Type            PlaceType  `json:"type"`
	}

	PlaceType string
)

const PlaceTypeAirport = "airport"
const PlaceTypeCity = "city"

func (a *API) PlaceSuggestions(ctx context.Context, query string) ([]*Place, error) {
	return newRequestWithAPI[EmptyPayload, Place](a).
		Get("/places/suggestions").WithParam("query", query).
		Slice(ctx)
}
