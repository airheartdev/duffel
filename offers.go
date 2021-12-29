package duffel

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type (
	OfferRequestClient interface {
		OfferRequest(ctx context.Context, requestInput *OfferRequestInput) (*OfferResponse, error)
	}
)

func (c *client) OfferRequest(ctx context.Context, requestInput *OfferRequestInput) (*OfferResponse, error) {
	payload := bytes.NewBuffer(nil)
	err := json.NewEncoder(payload).Encode(buildRequestPayload(requestInput))
	if err != nil {
		return nil, err
	}

	resp, err := c.makeRequest(ctx, "/air/offer_requests", http.MethodPost, payload, func(req *http.Request, u *url.URL) {
		if requestInput.ReturnOffers {
			q := u.Query()
			q.Add("return_offers", "true")
			u.RawQuery = q.Encode()
		}
	})

	if err != nil {
		return nil, err
	}

	// body := bytes.NewBuffer(nil)
	// _, err = body.ReadFrom(resp.Body)
	// if err != nil {
	// 	return nil, err
	// }

	offers := new(OfferResponse)
	respPayload := buildRequestPayload(offers)

	err = json.NewDecoder(resp.Body).Decode(&respPayload)
	if err != nil {
		return nil, err
	}

	return offers, nil
}
