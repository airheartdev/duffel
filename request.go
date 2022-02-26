package duffel

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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

// TODO: refactor this to use `makeIteratorRequest` instead of duplicating most of the logic
func (c *client[Req, Resp]) makeRequestWithPayload(ctx context.Context, resourceName string, method string, requestInput *Req, requestBuilders ...RequestOption) (*Resp, error) {
	resp, err := c.makeIteratorRequest(ctx, resourceName, method, requestInput, requestBuilders...)
	if err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *client[Req, Resp]) do(ctx context.Context, resourceName string, method string, requestInput *Req, requestBuilders ...RequestOption) (*http.Response, error) {
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
	return resp, nil
}

func (c *client[Req, Resp]) makeIteratorRequest(ctx context.Context, resourceName string, method string, requestInput *Req, requestBuilders ...RequestOption) (*ResponsePayload[Resp], error) {
	resp, err := c.do(ctx, resourceName, method, requestInput, requestBuilders...)
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
	resp, err := c.do(ctx, resourceName, method, requestInput, requestBuilders...)
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

func WithRequestPagination(meta *ListMeta) func(req *http.Request) {
	return func(req *http.Request) {
		q := req.URL.Query()

		if meta != nil {
			enc := schema.NewEncoder()
			enc.SetAliasTag("url")
			err := enc.Encode(meta, q)
			if err != nil {
				panic(err)
			}
			log.Println(q.Encode())
			req.URL.RawQuery = q.Encode()
		}
	}
}

func WithURLParams[T any](params ...T) func(req *http.Request) {
	return func(req *http.Request) {
		q := req.URL.Query()
		enc := schema.NewEncoder()
		enc.SetAliasTag("url")

		for _, param := range params {
			err := enc.Encode(param, q)
			if err != nil {
				panic(err)
			}
		}
		req.URL.RawQuery = q.Encode()
	}
}

func WithURLParam(key, value string) func(req *http.Request) {
	return func(req *http.Request) {
		q := req.URL.Query()
		q.Add(key, value)
		req.URL.RawQuery = q.Encode()
	}
}
