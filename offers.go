package duffel

import (
	"context"
	"net/http"
	"net/url"
)

type (
	OfferRequestClient interface {
		OfferRequest(ctx context.Context, requestInput *OfferRequestInput) (*OfferResponse, error)
	}
)

func (c *API) OfferRequest(ctx context.Context, requestInput *OfferRequestInput) (*OfferResponse, error) {
	client := newInternalClient[OfferRequestInput, OfferResponse](c)
	return client.makeRequestWithPayload(ctx, "/air/offer_requests", http.MethodPost, requestInput,
		func(_ *http.Request, u *url.URL) {
			if requestInput.ReturnOffers {
				q := u.Query()
				q.Add("return_offers", "true")
				u.RawQuery = q.Encode()
			}
		},
	)
}
