package duffel

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type (
	OfferClient interface {
		UpdateOfferPassenger(ctx context.Context, offerRequestID, passengerID string, input *PassengerUpdateInput) (*OfferRequestPassenger, error)
		ListOffers(ctx context.Context, reqId string, options ...ListOffersParams) *Iter[Offer]
		GetOffer(ctx context.Context, id string) (*Offer, error)
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
)

const (
	ListOffersSortTotalAmount   ListOffersSortParam = "total_amount"
	ListOffersSortTotalDuration ListOffersSortParam = "total_duration"
)

// UpdateOfferPassenger updates a single offer passenger.
func (a *API) UpdateOfferPassenger(ctx context.Context, offerRequestID, passengerID string, input *PassengerUpdateInput) (*Passenger, error) {
	client := newInternalClient[PassengerUpdateInput, Passenger](a)
	url := fmt.Sprintf("/air/offers/%s/passengers/%s", offerRequestID, passengerID)
	return client.makeRequestWithPayload(ctx, url, http.MethodPatch, input)
}

// ListOffers lists all the offers for an offer request. Returns an iterator.
func (a *API) ListOffers(ctx context.Context, offerRequestId string, options ...ListOffersParams) *Iter[Offer] {
	if offerRequestId == "" {
		return GetIter(func(meta *ListMeta) (*List[Offer], error) {
			return nil, fmt.Errorf("offerRequestId param is required")
		})
	} else if !strings.HasPrefix(offerRequestId, "orq_") {
		return GetIter(func(meta *ListMeta) (*List[Offer], error) {
			return nil, fmt.Errorf("offerRequestId should begin with orq_")
		})
	}

	c := newInternalClient[struct{}, Offer](a)
	return c.getIterator(ctx, http.MethodGet, "/air/offers",
		WithURLParam("offer_request_id", offerRequestId),
		WithURLParams(options...))
}

// GetOffer gets a single offer by ID.
func (a *API) GetOffer(ctx context.Context, id string) (*Offer, error) {
	c := newInternalClient[struct{}, Offer](a)
	return c.makeRequestWithPayload(ctx, "/air/offers/"+id, http.MethodGet, nil)
}
