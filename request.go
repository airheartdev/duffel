package duffel

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

func newInternalClient[Req any, Resp any](a *API) *client[Req, Resp] {
	return &client[Req, Resp]{
		httpDoer: a.httpDoer,
		options:  a.options,
		APIToken: a.APIToken,
		limiter:  rate.NewLimiter(rate.Every(1*time.Second), 5),
	}
}

func (c *client[Req, Resp]) Do(ctx context.Context, resourceName string, method string, body *Req, opts ...RequestOption) (*http.Response, error) {
	payload, err := encodePayload(body)
	if err != nil {
		return nil, err
	}

	err = c.limiter.Wait(ctx) // This is a blocking call. Honors the rate limit
	if err != nil {
		return nil, err
	}

	resp, err := c.makeRequest(ctx, resourceName, method, payload, opts...)
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
	return resp, nil
}
