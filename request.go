package duffel

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

func newInternalClient[T any, R any](a *API) *client[T, R] {
	return &client[T, R]{
		httpDoer: a.httpDoer,
		options:  a.options,
		APIToken: a.APIToken,
		rl:       rate.NewLimiter(rate.Every(1*time.Second), 5),
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

	err = c.rl.Wait(ctx) // This is a blocking call. Honors the rate limit
	if err != nil {
		return nil, err
	}

	resp, err := c.makeRequest(ctx, resourceName, method, payload, requestBuilders...)
	if err != nil {
		return nil, err
	}

	var reader io.ReadCloser
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer reader.Close()
	} else {
		reader = resp.Body
	}

	c.limit, err = strconv.Atoi(resp.Header.Get("Ratelimit-Limit"))
	if err != nil {
		return nil, err
	}
	c.limitRemaining, err = strconv.Atoi(resp.Header.Get("Ratelimit-Remaining"))
	if err != nil {
		return nil, err
	}
	c.limitReset, err = time.Parse(time.RFC1123, resp.Header.Get("Ratelimit-Reset"))
	if err != nil {
		return nil, err
	}
	timestamp, err := time.Parse(time.RFC1123, resp.Header.Get("Date"))
	if err != nil {
		return nil, err
	}

	period := c.limitReset.Sub(timestamp)
	c.rl.SetBurst(c.limit)
	c.rl.SetLimit(rate.Every(period))

	if c.limitRemaining == 0 || resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("rate limit exceeded, reset in: %s, current limit: %d", period.String(), c.limit)
	}

	container := new(ResponsePayload[Resp])
	err = json.NewDecoder(reader).Decode(&container)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	return container, nil
}
