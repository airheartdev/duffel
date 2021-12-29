package duffel

import (
	"net/http"
	"time"
)

const userAgentString = "duffel-go/1.0"
const defaultHost = "https://api.duffel.com/"

type (
	Duffel interface {
		OfferRequestClient
	}

	OfferRequestInput struct {
		Passengers   []Passenger         `json:"passengers"`
		Slices       []OfferRequestSlice `json:"slices"`
		CabinClass   CabinClass          `json:"cabin_class"`
		ReturnOffers bool                `json:"-"`
	}

	OfferRequestSlice struct {
		DepartureDate Date   `json:"departure_date"`
		Destination   string `json:"destination"`
		Origin        string `json:"origin"`
	}

	OfferResponse struct {
		Offers     []Offer     `json:"offers"`
		Slices     []Slice     `json:"slices"`
		Passengers []Passenger `json:"passengers"`
	}

	Offer struct {
		Slices     []Slice     `json:"slices"`
		Passengers []Passenger `json:"passengers"`
	}

	Slice struct {
		OriginType      string    `json:"origin_type"`
		Origin          Location  `json:"origin"`
		DestinationType string    `json:"destination_type"`
		Destination     Location  `json:"destination"`
		DepartureDate   Date      `json:"departure_date,omitempty"`
		CreatedAt       time.Time `json:"created_at,omitempty"`
		Segments        []Segment `json:"segments,omitempty"`
	}

	Segment struct {
	}

	Location struct {
		ID              string     `json:"id"`
		Type            string     `json:"type"`
		TimeZone        string     `json:"time_zone"`
		Name            string     `json:"name"`
		Longitude       *float64   `json:"longitude,omitempty"`
		Latitude        *float64   `json:"latitude,omitempty"`
		ICAOCode        *string    `json:"icao_code,omitempty"`
		IATACountryCode *string    `json:"iata_country_code,omitempty"`
		IATACode        *string    `json:"iata_code,omitempty"`
		IATACityCode    *string    `json:"iata_city_code,omitempty"`
		CityName        *string    `json:"city_name,omitempty"`
		City            *Location  `json:"city,omitempty"`
		Airports        []Location `json:"airports,omitempty"`
	}

	Passenger struct {
		ID                       string                    `json:"id,omitempty"`
		FamilyName               string                    `json:"family_name"`
		GivenName                string                    `json:"given_name"`
		Age                      *int                      `json:"age,omitempty"`
		Type                     *PassengerType            `json:"type,omitempty"`
		LoyaltyProgrammeAccounts []LoyaltyProgrammeAccount `json:"loyalty_programme_accounts,omitempty"`
	}

	LoyaltyProgrammeAccount struct {
		IATACode string `json:"iata_code"`
	}

	PassengerType string

	CabinClass string

	Offers []Offer

	Option  func(*Options)
	Options struct {
		Version   string
		Host      string
		UserAgent string
		HttpDoer  *http.Client
	}

	client struct {
		httpDoer *http.Client
		APIToken string
		options  *Options
	}
)

const (
	PassengerTypeAdult             PassengerType = "adult"
	PassengerTypeChild             PassengerType = "child"
	PassengerTypeInfantWithoutSeat PassengerType = "infant_without_seat"

	CabinClassEconomy  CabinClass = "economy"
	CabinClassPremium  CabinClass = "premium"
	CabinClassBusiness CabinClass = "business"
	CabinClassFirst    CabinClass = "first"
)

func New(apiToken string, opts ...Option) Duffel {
	options := &Options{
		Version:   "beta",
		UserAgent: userAgentString,
		Host:      defaultHost,
		HttpDoer:  http.DefaultClient,
	}
	for _, opt := range opts {
		opt(options)
	}
	c := &client{
		httpDoer: options.HttpDoer,
		APIToken: apiToken,
		options:  options,
	}
	return c
}

// Assert that our interface matches
var (
	_ Duffel = (*client)(nil)
)
