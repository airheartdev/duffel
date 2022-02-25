package duffel

import (
	"context"
	"net/http"
)

type (
	OfferRequestClient interface {
		CreateOfferRequest(ctx context.Context, requestInput *OfferRequestInput) (*OfferResponse, error)
		GetOfferRequest(ctx context.Context, id string) (*OfferResponse, error)
	}
)

func (c *API) CreateOfferRequest(ctx context.Context, requestInput *OfferRequestInput) (*OfferResponse, error) {
	client := newInternalClient[OfferRequestInput, OfferResponse](c)
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

func (c *API) GetOfferRequest(ctx context.Context, id string) (*OfferResponse, error) {
	client := newInternalClient[OfferRequestInput, OfferResponse](c)
	return client.makeRequestWithPayload(ctx, "/air/offer_requests/"+id, http.MethodGet, nil)
}
