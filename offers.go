package duffel

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const offerIDPrefix = "off_"
const offerRequestIDPrefix = "orq_"

type (
	OfferClient interface {
		UpdateOfferPassenger(ctx context.Context, offerRequestID, passengerID string, input PassengerUpdateInput) (*OfferRequestPassenger, error)
		ListOffers(ctx context.Context, reqId string, options ...ListOffersParams) *Iter[Offer]
		GetOffer(ctx context.Context, id string, params ...GetOfferParams) (*Offer, error)
	}

	Offer struct {
		ID                                 string                  `json:"id"`
		LiveMode                           bool                    `json:"live_mode"`
		CreatedAt                          time.Time               `json:"created_at"`
		UpdatedAt                          time.Time               `json:"updated_at"`
		TotalEmissionsKg                   string                  `json:"total_emissions_kg"`
		TotalCurrency                      string                  `json:"total_currency"`
		TotalAmount                        string                  `json:"total_amount"`
		TaxCurrency                        string                  `json:"tax_currency"`
		TaxAmount                          string                  `json:"tax_amount"`
		Owner                              Airline                 `json:"owner"`
		Slices                             []Slice                 `json:"slices"`
		Passengers                         []OfferRequestPassenger `json:"passengers"`
		PassengerIdentityDocumentsRequired bool                    `json:"passenger_identity_documents_required"`
	}

	ListOffersSortParam string

	ListOffersParams struct {
		Sort           ListOffersSortParam `url:"sort,omitempty"`
		MaxConnections int                 `url:"max_connections,omitempty"`
	}

	GetOfferParams struct {
		ReturnAvailableServices bool
	}
)

const (
	ListOffersSortTotalAmount   ListOffersSortParam = "total_amount"
	ListOffersSortTotalDuration ListOffersSortParam = "total_duration"
)

// UpdateOfferPassenger updates a single offer passenger.
func (a *API) UpdateOfferPassenger(ctx context.Context, offerRequestID, passengerID string, input PassengerUpdateInput) (*OfferRequestPassenger, error) {
	url := fmt.Sprintf("/air/offers/%s/passengers/%s", offerRequestID, passengerID)
	return newRequestWithAPI[PassengerUpdateInput, OfferRequestPassenger](a).Patch(url, &input).One(ctx)
}

// ListOffers lists all the offers for an offer request. Returns an iterator.
func (a *API) ListOffers(ctx context.Context, offerRequestID string, options ...ListOffersParams) *Iter[Offer] {
	if offerRequestID == "" {
		return ErrIter[Offer](fmt.Errorf("offerRequestId param is required"))
	} else if !strings.HasPrefix(offerRequestID, offerRequestIDPrefix) {
		return ErrIter[Offer](fmt.Errorf("offerRequestId should begin with %s", offerRequestIDPrefix))
	}

	return newRequestWithAPI[ListOffersParams, Offer](a).Get("/air/offers").
		WithParam("offer_request_id", offerRequestID).
		WithParams(normalizeParams(options)...).
		All(ctx)
}

// GetOffer gets a single offer by ID.
func (a *API) GetOffer(ctx context.Context, offerID string, params ...GetOfferParams) (*Offer, error) {
	if !strings.HasPrefix(offerID, offerIDPrefix) {
		return nil, fmt.Errorf("offerID should begin with %s", offerIDPrefix)
	}

	return newRequestWithAPI[GetOfferParams, Offer](a).
		Get(fmt.Sprintf("/air/offers/%s", offerID)).
		WithParams(normalizeParams(params)...).
		One(ctx)
}

func (o ListOffersParams) Encode(q url.Values) error {
	if o.Sort != "" {
		q.Set("sort", string(o.Sort))
	}

	if o.MaxConnections != 0 {
		q.Set("max_connections", fmt.Sprintf("%d", o.MaxConnections))
	}

	return nil
}

func (o GetOfferParams) Encode(q url.Values) error {
	if o.ReturnAvailableServices {
		q.Set("return_available_services", "true")
	}
	return nil
}
