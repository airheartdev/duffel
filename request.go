package duffel

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/segmentio/encoding/json"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
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

func (c *client[R, W]) Debug() *client[R, W] {
	c.debug = true
	return c
}

// Get makes a GET request to the specified resource.
func (c *client[Req, Resp]) Get(ctx context.Context, resourcePath string, opts ...RequestOption) (*http.Response, error) {
	err := c.limiter.Wait(ctx) // This is a blocking call. Honors the rate limit
	if err != nil {
		return nil, err
	}

	resp, err := c.makeRequest(ctx, resourcePath, http.MethodGet, nil, opts...)
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

func (c *client[Req, Resp]) Post(ctx context.Context, resourcePath string, requestInput *Req, opts ...RequestOption) (*Resp, error) {
	resp, err := c.makeIteratorRequest(ctx, resourcePath, http.MethodPost, requestInput, opts...)
	if err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *client[Req, Resp]) makeRequestWithPayload(ctx context.Context, resourceName string, method string, requestInput *Req, requestBuilders ...RequestOption) (*Resp, error) {
	resp, err := c.makeIteratorRequest(ctx, resourceName, method, requestInput, requestBuilders...)
	if err != nil {
		return nil, err
	}

	return &resp.Data, nil
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

func (c *client[Req, Resp]) makeIteratorRequest(ctx context.Context, resourceName string, method string, requestInput *Req, requestBuilders ...RequestOption) (*ResponsePayload[Resp], error) {
	resp, err := c.Do(ctx, resourceName, method, requestInput, requestBuilders...)
	if err != nil {
		return nil, err
	}

	reader, err := gzipResponseReader(resp)
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

func (c *client[Req, Resp]) makeListRequest(ctx context.Context, resourceName string, method string, requestInput *Req, requestBuilders ...RequestOption) (*ResponsePayload[[]*Resp], error) {
	resp, err := c.Do(ctx, resourceName, method, requestInput, requestBuilders...)
	if err != nil {
		return nil, err
	}

	reader, err := gzipResponseReader(resp)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	container := new(ResponsePayload[[]*Resp])
	err = json.NewDecoder(reader).Decode(&container)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	return container, nil
}

func (c *client[Req, Resp]) getIterator(ctx context.Context, method string, resourceName string, requestBuilders ...RequestOption) *Iter[Resp] {
	return GetIter(func(lastMeta *ListMeta) (*List[Resp], error) {
		list := new(List[Resp])
		response, err := c.makeListRequest(ctx, resourceName, method, nil, append(requestBuilders, WithRequestPagination(lastMeta))...)
		if err != nil {
			return nil, err
		}
		if response == nil {
			return nil, fmt.Errorf("internal: empty response")
		}
		list.SetListMeta(response.Meta)
		list.SetItems(response.Data)
		return list, nil
	})
}

func WithRequestPagination(meta *ListMeta) RequestOption {
	return func(req *http.Request) error {
		q := req.URL.Query()

		if meta != nil {
			enc := schema.NewEncoder()
			enc.SetAliasTag("url")
			err := enc.Encode(meta, q)
			if err != nil {
				return (err)
			}
			req.URL.RawQuery = q.Encode()
		}
		return nil
	}
}

func WithURLParams[T any](params ...T) RequestOption {
	return func(req *http.Request) error {
		q := req.URL.Query()
		enc := schema.NewEncoder()
		enc.SetAliasTag("url")

		for _, param := range params {
			err := enc.Encode(param, q)
			if err != nil {
				return (err)
			}
		}
		req.URL.RawQuery = q.Encode()
		return nil
	}
}

func WithURLParam(key, value string) RequestOption {
	return func(req *http.Request) error {
		q := req.URL.Query()
		q.Add(key, value)
		req.URL.RawQuery = q.Encode()
		return nil
	}
}
