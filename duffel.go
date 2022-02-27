package duffel

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

const userAgentString = "duffel-go/1.0"
const defaultHost = "https://api.duffel.com/"

type (
	Duffel interface {
		OfferRequestClient
		OfferClient
		AirportsClient
	}

	Gender string

	OfferRequestInput struct {
		Passengers   []Passenger         `json:"passengers" url:"-"`
		Slices       []OfferRequestSlice `json:"slices" url:"-"`
		CabinClass   CabinClass          `json:"cabin_class" url:"-"`
		ReturnOffers bool                `json:"-" url:"return_offers,defualt=true"`
	}

	OfferRequestSlice struct {
		DepartureDate Date   `json:"departure_date"`
		Destination   string `json:"destination"`
		Origin        string `json:"origin"`
	}

	OfferRequest struct {
		ID         string      `json:"id"`
		LiveMode   bool        `json:"live_mode"`
		Offers     []Offer     `json:"offers"`
		Slices     []Slice     `json:"slices"`
		Passengers []Passenger `json:"passengers"`
		CreatedAt  time.Time   `json:"created_at"`
		CabinClass CabinClass  `json:"cabin_class"`
	}

	Offer struct {
		Slices           []Slice     `json:"slices"`
		Passengers       []Passenger `json:"passengers"`
		UpdatedAt        time.Time   `json:"updated_at"`
		TotalEmissionsKg string      `json:"total_emissions_kg"`
		TotalCurrency    string      `json:"total_currency"`
		TotalAmount      string      `json:"total_amount"`
		TaxCurrency      string      `json:"tax_currency"`
		TaxAmount        string      `json:"tax_amount"`
	}

	Slice struct {
		OriginType      string    `json:"origin_type"`
		Origin          Location  `json:"origin"`
		DestinationType string    `json:"destination_type"`
		Destination     Location  `json:"destination"`
		DepartureDate   Date      `json:"departure_date,omitempty"`
		CreatedAt       time.Time `json:"created_at,omitempty"`
		Segments        []Flight  `json:"segments,omitempty"`
	}

	Flight struct {
		Passengers                   []SegmentPassenger `json:"passengers"`
		Origin                       Location           `json:"origin"`
		OriginTerminal               string             `json:"origin_terminal"`
		OperatingCarrierFlightNumber string             `json:"operating_carrier_flight_number"`
		OperatingCarrier             Airline            `json:"operating_carrier"`
		MarketingCarrierFlightNumber string             `json:"marketing_carrier_flight_number"`
		MarketingCarrier             Airline            `json:"marketing_carrier"`
		Duration                     Duration           `json:"duration"`
		Distance                     Distance           `json:"distance,omitempty"`
		DestinationTerminal          string             `json:"destination_terminal"`
		Destination                  Location           `json:"destination"`
		DepartingAt                  DateTime           `json:"departing_at"`
		ArrivingAt                   DateTime           `json:"arriving_at"`
		Aircraft                     Equipment          `json:"aircraft"`
	}

	SegmentPassenger struct {
		ID                      string     `json:"id"`
		FareBasisCode           string     `json:"fare_basis_code"`
		CabinClassMarketingName string     `json:"cabin_class_marketing_name"`
		CabinClass              CabinClass `json:"cabin_class"`
		Baggages                []Baggage  `json:"baggages"`
	}

	Equipment struct {
		IATACode string `json:"iata_code"`
		ID       string `json:"id"`
		Name     string `json:"name"`
	}

	Airline struct {
		Name     string `json:"name"`
		IATACode string `json:"iata_code"`
		ID       string `json:"id"`
	}

	Airport struct {
		ID              string   `json:"id"`
		Name            string   `json:"name"`
		City            Location `json:"city"`
		CityName        string   `json:"city_name"`
		IATACode        string   `json:"iata_code"`
		IATACountryCode string   `json:"iata_country_code"`
		ICAOCode        string   `json:"icao_code"`
		Latitude        float32  `json:"latitude"`
		Longitude       float32  `json:"longitude"`
		TimeZone        string   `json:"time_zone"`
	}

	Baggage struct {
		Quantity int    `json:"quantity"`
		Type     string `json:"type"`
	}

	Location struct {
		ID              string     `json:"id"`
		Type            string     `json:"type"`
		Name            string     `json:"name"`
		TimeZone        string     `json:"time_zone,omitempty"`
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
		ID         string `json:"id,omitempty"`
		FamilyName string `json:"family_name"`
		GivenName  string `json:"given_name"`
		Age        *int   `json:"age,omitempty"`
		// deprecated
		Type                     *PassengerType            `json:"type,omitempty"`
		LoyaltyProgrammeAccounts []LoyaltyProgrammeAccount `json:"loyalty_programme_accounts,omitempty"`
	}

	PassengerCreateInput struct {
		ID         string         `json:"id"`
		Title      PassengerTitle `json:"title"`
		FamilyName string         `json:"family_name"`
		GivenName  string         `json:"given_name"`
		BornOn     Date           `json:"born_on"`
		Email      string         `json:"email"`
		Gender     Gender         `json:"gender"`
		// The passenger's identity documents. You may only provide one identity document per passenger.
		IdentityDocuments []IdentityDocument `json:"identity_documents,omitempty"`

		InfantPassengerID string `json:"infant_passenger_id,omitempty"`

		// (Required) The passenger's phone number in E.164 (international) format
		PhoneNumber string `json:"phone_number"`

		// deprecated
		Type PassengerType `json:"type"`
	}

	PassengerUpdateInput struct {
		FamilyName               string                    `json:"family_name"`
		GivenName                string                    `json:"given_name"`
		LoyaltyProgrammeAccounts []LoyaltyProgrammeAccount `json:"loyalty_programme_accounts,omitempty"`
	}

	// The payment details to use to pay for the order.
	// This key should be omitted when the order’s type is hold.
	PaymentCreateInput struct {
		// The amount of the payment. This should be the same as the total_amount of the offer specified in selected_offers, plus the total_amount of all the services specified in services.
		Amount string `json:"amount"`
		// The currency of the amount, as an ISO 4217 currency code. This should be the same as the total_currency of the offer specified in selected_offers.
		Currency string `json:"currency"`

		// Possible values: "arc_bsp_cash" or "balance"
		Type string `json:"type"`
	}

	LoyaltyProgrammeAccount struct {
		AirlineIATACode string `json:"airline_iata_code"`
		AccountNumber   string `json:"account_number"`
	}

	IdentityDocument struct {
		// The unique identifier of the identity document.
		// We currently only support passport so this would be the passport number.
		UniqueIdentifier string `json:"unique_identifier"`

		// The date on which the identity document expires
		ExpiresOn Date `json:"expires_on"`

		// The ISO 3166-1 alpha-2 code of the country that issued this identity document
		IssuingCountryCode string `json:"issuing_country_code"`

		Type string `json:"type"`
	}

	PassengerType string

	PassengerTitle string

	CabinClass string

	Offers []Offer

	Option  func(*Options)
	Options struct {
		Version   string
		Host      string
		UserAgent string
		HttpDoer  *http.Client
	}

	client[Req any, Resp any] struct {
		httpDoer  *http.Client
		APIToken  string
		options   *Options
		limiter   *rate.Limiter
		rateLimit *RateLimit
	}

	API struct {
		httpDoer *http.Client
		APIToken string
		options  *Options
	}
)

const (
	// deprecated
	PassengerTypeAdult PassengerType = "adult"
	// deprecated
	PassengerTypeChild PassengerType = "child"
	// deprecated
	PassengerTypeInfantWithoutSeat PassengerType = "infant_without_seat"

	CabinClassEconomy  CabinClass = "economy"
	CabinClassPremium  CabinClass = "premium"
	CabinClassBusiness CabinClass = "business"
	CabinClassFirst    CabinClass = "first"

	GenderMale   Gender = "m"
	GenderFemale Gender = "f"

	PassengerTitleMr   PassengerTitle = "mr"
	PassengerTitleMs   PassengerTitle = "ms"
	PassengerTitleMrs  PassengerTitle = "mrs"
	PassengerTitleMiss PassengerTitle = "miss"
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

	return &API{
		httpDoer: options.HttpDoer,
		APIToken: apiToken,
		options:  options,
	}
}

// Assert that our interface matches
var (
	_ Duffel = (*API)(nil)
)

func (p PassengerType) String() string {
	return string(p)
}
