package duffel

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type (
	OfferRequestClient interface {
		CreateOfferRequest(ctx context.Context, requestInput *OfferRequestInput) (*OfferResponse, error)
		GetOfferRequest(ctx context.Context, id string) (*OfferResponse, error)
		ListOfferRequests(ctx context.Context) *Iter[OfferResponse]
	}
)

func (a *API) CreateOfferRequest(ctx context.Context, requestInput *OfferRequestInput) (*OfferResponse, error) {
	client := newInternalClient[OfferRequestInput, OfferResponse](a)
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

func (a *API) GetOfferRequest(ctx context.Context, id string) (*OfferResponse, error) {
	client := newInternalClient[OfferRequestInput, OfferResponse](a)
	return client.makeRequestWithPayload(ctx, "/air/offer_requests/"+id, http.MethodGet, nil)
}

func (a *API) ListOfferRequests(ctx context.Context) *Iter[OfferResponse] {
	c := newInternalClient[ListAirportsRequest, []*OfferResponse](a)

	return GetIter(func(lastMeta *ListMeta) (*List[OfferResponse], error) {
		list := new(List[OfferResponse])
		response, err := c.makeIteratorRequest(ctx, "/air/offer_requests", http.MethodGet, nil, func(req *http.Request) {
			q := req.URL.Query()
			if lastMeta != nil && lastMeta.After != "" {
				q.Add("after", lastMeta.After)
			}
			req.URL.RawQuery = q.Encode()
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to list airports")
		}
		if response == nil {
			return nil, fmt.Errorf("internal: empty response")
		}
		list.SetListMeta(response.Meta)
		list.SetItems(response.Data)
		return list, nil
	})
}
