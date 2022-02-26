package duffel

import (
	"context"
	"net/http"
)

type (
	OfferRequestClient interface {
		CreateOfferRequest(ctx context.Context, requestInput *OfferRequestInput) (*OfferRequest, error)
		GetOfferRequest(ctx context.Context, id string) (*OfferRequest, error)
		ListOfferRequests(ctx context.Context) *Iter[OfferRequest]
	}
)

func (a *API) CreateOfferRequest(ctx context.Context, requestInput *OfferRequestInput) (*OfferRequest, error) {
	client := newInternalClient[OfferRequestInput, OfferRequest](a)
	return client.makeRequestWithPayload(ctx, "/air/offer_requests", http.MethodPost, requestInput,
		func(req *http.Request) {
			if requestInput != nil && requestInput.ReturnOffers {
				q := req.URL.Query()
				q.Add("return_offers", "true")
				req.URL.RawQuery = q.Encode()
			}
		},
	)
}

func (a *API) GetOfferRequest(ctx context.Context, id string) (*OfferRequest, error) {
	client := newInternalClient[OfferRequestInput, OfferRequest](a)
	return client.makeRequestWithPayload(ctx, "/air/offer_requests/"+id, http.MethodGet, nil)
}

func (a *API) ListOfferRequests(ctx context.Context) *Iter[OfferRequest] {
	c := newInternalClient[ListAirportsRequest, OfferRequest](a)
	return c.getIterator(ctx, http.MethodGet, "/air/offer_requests", nil)
}
