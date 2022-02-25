package duffel

import (
	"context"
	"net/http"
)

type (
	OfferRequestClient interface {
		OfferRequest(ctx context.Context, requestInput *OfferRequestInput) (*OfferResponse, error)
	}
)

func (c *API) OfferRequest(ctx context.Context, requestInput *OfferRequestInput) (*OfferResponse, error) {
	client := newInternalClient[OfferRequestInput, OfferResponse](c)
	return client.makeRequestWithPayload(ctx, "/air/offer_requests", http.MethodPost, requestInput,
		func(req *http.Request) {
			if requestInput.ReturnOffers {
				q := req.URL.Query()
				q.Add("return_offers", "true")
				req.URL.RawQuery = q.Encode()
			}
		},
	)
}
