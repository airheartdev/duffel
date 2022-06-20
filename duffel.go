// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

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
		OrderClient
		OrderChangeClient
		OrderCancellationClient
		OrderPaymentClient
		SeatmapClient
		AirportsClient
		AirlinesClient
		AircraftClient
		PlacesClient
	}

	Gender string

	BaseSlice struct {
		OriginType      string    `json:"origin_type"`
		Origin          Location  `json:"origin"`
		DestinationType string    `json:"destination_type"`
		Destination     Location  `json:"destination"`
		DepartureDate   Date      `json:"departure_date,omitempty"`
		CreatedAt       time.Time `json:"created_at,omitempty"`
	}

	// TODO: We probably need an OfferRequestSlice and an OrderSlice since not all fields apply to both.
	Slice struct {
		*BaseSlice
		ID string `json:"id"`

		// Whether this slice can be changed. This can only be true for paid orders.
		Changeable bool `json:"changeable,omitempty"`

		// The conditions associated with this slice, describing the kinds of modifications you can make and any penalties that will apply to those modifications.
		Conditions    Conditions `json:"conditions,omitempty"`
		Duration      Duration   `json:"duration,omitempty"`
		Segments      []Flight   `json:"segments,omitempty"`
		FareBrandName string     `json:"fare_brand_name,omitempty"`
	}

	Flight struct {
		ID                           string             `json:"id"`
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
		Aircraft                     Aircraft           `json:"aircraft"`
	}

	SegmentPassenger struct {
		ID                      string     `json:"passenger_id"`
		FareBasisCode           string     `json:"fare_basis_code"`
		CabinClassMarketingName string     `json:"cabin_class_marketing_name"`
		CabinClass              CabinClass `json:"cabin_class"`
		Baggages                []Baggage  `json:"baggages"`
		Seat                    Seat       `json:"seat"`
	}

	Seat struct {
		Name        string   `json:"name,omitempty"`
		Disclosures []string `json:"disclosures,omitempty"`
		Designator  string   `json:"designator,omitempty"`
	}

	Aircraft struct {
		IATACode string `json:"iata_code"`
		ID       string `json:"id"`
		Name     string `json:"name"`
	}

	Airline struct {
		Name          string `json:"name"`
		IATACode      string `json:"iata_code"`
		ID            string `json:"id"`
		LogoSymbolURL string `json:"logo_symbol_url"`
		LogoLockupURL string `json:"logo_lockup_url"`
	}

	Airport struct {
		ID              string  `json:"id" `
		Name            string  `json:"name" `
		City            City    `json:"city,omitempty" `
		CityName        string  `json:"city_name" `
		IATACode        string  `json:"iata_code" `
		IATACountryCode string  `json:"iata_country_code" `
		ICAOCode        string  `json:"icao_code" `
		Latitude        float32 `json:"latitude" `
		Longitude       float32 `json:"longitude" `
		TimeZone        string  `json:"time_zone" `
	}

	Baggage struct {
		Quantity int    `json:"quantity"`
		Type     string `json:"type"`
	}

	Location struct {
		ID              string    `json:"id"`
		Type            string    `json:"type"`
		Name            string    `json:"name"`
		TimeZone        string    `json:"time_zone,omitempty"`
		Longitude       *float64  `json:"longitude,omitempty"`
		Latitude        *float64  `json:"latitude,omitempty"`
		ICAOCode        string    `json:"icao_code,omitempty"`
		IATACountryCode *string   `json:"iata_country_code,omitempty"`
		IATACode        string    `json:"iata_code,omitempty"`
		IATACityCode    *string   `json:"iata_city_code,omitempty"`
		CityName        *string   `json:"city_name,omitempty"`
		City            *City     `json:"city,omitempty"`
		Airports        []Airport `json:"airports,omitempty"`
	}

	City struct {
		ID              string  `json:"id,omitempty" csv:"city_id"`
		Name            string  `json:"name" csv:"city_name"`
		IATACountryCode *string `json:"iata_country_code,omitempty" csv:"city_iata_country_code"`
		IATACode        string  `json:"iata_code,omitempty" csv:"city_iata_code"`
	}

	OrderPassenger struct {
		// ID is id of the passenger, returned when the offer request was created
		ID string `json:"id"`
		// Title is passengers' title. Possible values: "mr", "ms", "mrs", or "miss"
		Title PassengerTitle `json:"title"`
		// FamilyName is the family name of the passenger.
		FamilyName string `json:"family_name"`
		// GivenName is the passenger's given name.
		GivenName string `json:"given_name"`
		// BornOn is the passengers DoB according to their travel documents.
		BornOn Date `json:"born_on"`
		// Email is the passengers email address.
		Email string `json:"email"`
		// Gender is the passengers gender.
		Gender Gender `json:"gender"`
		// The passenger's identity documents. You may only provide one identity document per passenger.
		IdentityDocuments []IdentityDocument `json:"identity_documents,omitempty"`
		// The `id` of the infant associated with this passenger
		InfantPassengerID string `json:"infant_passenger_id,omitempty"`
		// The Loyalty Programme Accounts for this passenger
		LoyaltyProgrammeAccounts []LoyaltyProgrammeAccount `json:"loyalty_programme_accounts,omitempty"`
		// (Required) The passenger's phone number in E.164 (international) format
		PhoneNumber string `json:"phone_number"`

		// Type is the type of passenger. This field is deprecated.
		// @Deprecated
		// Possible values: "adult", "child", or "infant_without_seat"
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
		Type PaymentMethod `json:"type"`
	}

	// The payment status for an order.
	PaymentStatus struct {
		AwaitingPayment         bool       `json:"awaiting_payment"`
		PaymentRequiredBy       *time.Time `json:"payment_required_by,omitempty"`
		PriceGuaranteeExpiresAt *time.Time `json:"price_guarantee_expires_at,omitempty"`
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

	PaymentMethod string

	// Offers is slice of offers that implement the sort.Sort interface
	// By default, offers are sorted cheapest first.
	Offers []Offer

	Option  func(*Options)
	Options struct {
		Version   string
		Host      string
		UserAgent string
		HttpDoer  *http.Client
		Debug     bool
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
	CabinClassPremium  CabinClass = "premium_economy"
	CabinClassBusiness CabinClass = "business"
	CabinClassFirst    CabinClass = "first"

	GenderMale   Gender = "m"
	GenderFemale Gender = "f"

	PassengerTitleMr   PassengerTitle = "mr"
	PassengerTitleMs   PassengerTitle = "ms"
	PassengerTitleMrs  PassengerTitle = "mrs"
	PassengerTitleMiss PassengerTitle = "miss"

	PaymentMethodBalance  PaymentMethod = "balance"
	ARCBSPCash            PaymentMethod = "arc_bsp_cash"
	Card                  PaymentMethod = "card"
	Voucher               PaymentMethod = "voucher"
	AwaitingPayment       PaymentMethod = "awaiting_payment"
	OriginalFormOfPayment PaymentMethod = "original_form_of_payment"
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

func (p PaymentMethod) String() string {
	return string(p)
}

func (p CabinClass) String() string {
	return string(p)
}

func (p Gender) String() string {
	return string(p)
}

func (p PassengerTitle) String() string {
	return string(p)
}
