package duffel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

func newInternalClient[T any, R any](a *API) *client[T, R] {
	return &client[T, R]{
		httpDoer: a.httpDoer,
		options:  a.options,
		APIToken: a.APIToken,
		limiter:  rate.NewLimiter(rate.Every(1*time.Second), 5),
	}
}

// TODO: refactor this to use `makeIteratorRequest` instead of duplicating most of the logic
func (c *client[Req, Resp]) makeRequestWithPayload(ctx context.Context, resourceName string, method string, requestInput *Req, requestBuilders ...RequestOption) (*Resp, error) {
	resp, err := c.makeIteratorRequest(ctx, resourceName, method, requestInput, requestBuilders...)
	if err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *client[Req, Resp]) makeIteratorRequest(ctx context.Context, resourceName string, method string, requestInput *Req, requestBuilders ...RequestOption) (*ResponsePayload[Resp], error) {
	payload, err := encodePayload(requestInput)
	if err != nil {
		return nil, err
	}

	err = c.limiter.Wait(ctx) // This is a blocking call. Honors the rate limit
	if err != nil {
		return nil, err
	}

	resp, err := c.makeRequest(ctx, resourceName, method, payload, requestBuilders...)
	if err != nil {
		return nil, err
	}

	rateLimit, err := parseRateLimit(resp)
	if err != nil {
		return nil, err
	}

	c.rateLimit = rateLimit
	c.limiter.SetBurst(rateLimit.Limit)
	c.limiter.SetLimit(rate.Every(rateLimit.Period))

	if rateLimit.Remaining == 0 || resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("rate limit exceeded, reset in: %s, current limit: %d", rateLimit.Period.String(), rateLimit.Limit)
	}

	reader, err := decodeResponse(resp)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	container := new(ResponsePayload[Resp])
	err = json.NewDecoder(reader).Decode(&container)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	return container, nil
}
