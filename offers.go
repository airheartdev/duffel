package duffel

import (
	"context"
	"fmt"
	"net/http"
)

type (
	OfferClient interface {
		UpdateOfferPassenger(ctx context.Context, input *PassengerUpdateInput) (*Passenger, error)
		// GetOffer(ctx context.Context, id string) (*Offer, error)
		ListOffers(ctx context.Context, reqId string, options ...ListOffersParams) *Iter[Offer]
	}

	ListOffersSortParam string

	ListOffersParams struct {
		Sort           string `json:"sort"`
		MaxConnections int    `json:"max_connections"`
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
	c := newInternalClient[struct{}, Offer](a)
	return c.getIterator(ctx, http.MethodGet, "/air/offers", func(req *http.Request) {
		q := req.URL.Query()
		q.Add("offer_request_id", offerRequestId)
		if len(options) == 1 {
			if options[0].Sort != "" {
				q.Add("sort", string(options[0].Sort))
			}
			if options[0].MaxConnections != 0 {
				q.Add("max_connections", fmt.Sprintf("%d", options[0].MaxConnections))
			}
		}
		req.URL.RawQuery = q.Encode()
	})
}
