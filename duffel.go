package duffel

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

const userAgentString = "duffel-go/1.0"

type (
	Duffel interface {
		OfferRequest(ctx context.Context, requestInput *OfferRequestInput) (Offers, error)
	}

	OfferRequestInput struct {
		Passengers []Passenger `json:"passengers"`
	}

	Offer struct {
	}

	Passenger struct {
		FamilyName string        `json:"family_name"`
		GivenName  string        `json:"given_name"`
		Age        int           `json:"age"`
		Type       PassengerType `json:"type"`
	}

	PassengerType string

	Offers []Offer

	Option  func(*Options)
	Options struct {
		Version   string
		Host      string
		UserAgent string
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
)

func New(apiToken string, opts ...Option) Duffel {
	c := &client{
		httpDoer: http.DefaultClient,
	}
	options := &Options{
		Version:   "beta",
		UserAgent: userAgentString,
		Host:      "https://api.duffel.com",
	}
	for _, opt := range opts {
		opt(options)
	}
	c.options = options
	return c
}

// Assert that our interface matches
var (
	_ Duffel = (*client)(nil)
)

func (c *client) OfferRequest(ctx context.Context, requestInput *OfferRequestInput) (Offers, error) {
	payload := bytes.NewBuffer(nil)
	err := json.NewEncoder(payload).Encode(requestInput)
	if err != nil {
		return nil, err
	}

	resp, err := c.makeRequest(ctx, "air/offer_requests", http.MethodPost, payload)
	if err != nil {
		return nil, err
	}

	offers := make(Offers, 0)
	err = json.NewDecoder(resp.Body).Decode(&offers)
	if err != nil {
		return nil, err
	}

	return offers, nil
}
