package duffel

import (
	"context"
	"net/http"
)

type (
	OfferRequestClient interface {
		CreateOfferRequest(ctx context.Context, requestInput OfferRequestInput) (*OfferRequest, error)
		GetOfferRequest(ctx context.Context, id string) (*OfferRequest, error)
		ListOfferRequests(ctx context.Context) *Iter[OfferRequest]
	}
)

func (a *API) CreateOfferRequest(ctx context.Context, requestInput OfferRequestInput) (*OfferRequest, error) {
	client := newInternalClient[OfferRequestInput, OfferRequest](a)
	return client.makeRequestWithPayload(ctx,
		"/air/offer_requests",
		http.MethodPost,
		&requestInput,
		WithURLParams(requestInput),
	)
}

func (a *API) GetOfferRequest(ctx context.Context, id string) (*OfferRequest, error) {
	c := newInternalClient[OfferRequestInput, OfferRequest](a)
	return c.makeRequestWithPayload(ctx, "/air/offer_requests/"+id, http.MethodGet, nil)
}

func (a *API) ListOfferRequests(ctx context.Context) *Iter[OfferRequest] {
	c := newInternalClient[ListAirportsParams, OfferRequest](a)
	return c.getIterator(ctx, http.MethodGet, "/air/offer_requests")
}
