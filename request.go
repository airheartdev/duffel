package duffel

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

func newInternalClient[T any, R any](a *API) *client[T, R] {
	return &client[T, R]{
		httpDoer: a.httpDoer,
		options:  a.options,
		APIToken: a.APIToken,
	}
}

func (c *client[Req, Resp]) makeRequestWithPayload(ctx context.Context, resourceName string, method string, requestInput *Req, requestBuilders ...RequestOption) (*Resp, error) {
	payload, err := encodePayload(requestInput)
	if err != nil {
		return nil, err
	}

	resp, err := c.makeRequest(ctx, resourceName, method, payload)
	if err != nil {
		return nil, err
	}

	paylod := new(Resp)
	container := buildPayload(paylod)

	err = json.NewDecoder(resp.Body).Decode(&container)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	return paylod, nil
}
