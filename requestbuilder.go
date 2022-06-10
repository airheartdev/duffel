// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package duffel

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"github.com/segmentio/encoding/json"
)

type RequestBuilder[Req any, Resp any] struct {
	client         *client[Req, Resp]
	method         string
	resourcePath   string
	requestOptions []RequestOption
	body           *Req
	params         url.Values
}

type RequestMiddleware func(r *http.Request) error
type ResponseMiddleware func(r *http.Response) error

type ParamEncoder[T any] interface {
	Encode(v url.Values) error
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

func WithEncodableParams[T any](params ...ParamEncoder[T]) RequestOption {
	return func(req *http.Request) error {
		q := req.URL.Query()

		for _, param := range params {
			err := param.Encode(q)
			if err != nil {
				return err
			}
		}

		req.URL.RawQuery = q.Encode()
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
				return err
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

// newRequestWithAPI returns a new fluent requst builder for the given API.
// The input types are the request payload and response payload.
// For Get requests, the request payload is used to type the URL params.
//
// The response payload is used to type the parsed response body, and is the type returned by the finalizers.
func newRequestWithAPI[ReqT any, ResponseT any](a *API) *RequestBuilder[ReqT, ResponseT] {
	return &RequestBuilder[ReqT, ResponseT]{
		client: newInternalClient[ReqT, ResponseT](a),
		method: http.MethodGet,
	}
}

// WithParam adds a single query param to the URL.
// These operations will be applied in defined order after the request is initialized.
func (r *RequestBuilder[Req, Resp]) WithParam(key, value string) *RequestBuilder[Req, Resp] {
	r.requestOptions = append(r.requestOptions, WithURLParam(key, value))
	return r
}

// WithParams sets the URL query params for the request.
// These operations will be applied in defined order after the request is initialized.
func (r *RequestBuilder[Req, Resp]) WithParams(obj ...ParamEncoder[Req]) *RequestBuilder[Req, Resp] {
	r.requestOptions = append(r.requestOptions, WithEncodableParams(obj...))
	return r
}

func (r *RequestBuilder[Req, Resp]) WithOptions(opts ...RequestOption) *RequestBuilder[Req, Resp] {
	r.requestOptions = append(r.requestOptions, opts...)
	return r
}

// Get sets the request method to GET and the request path to the given path. Global request options are applied.
func (r *RequestBuilder[Req, Resp]) Get(path string, opts ...RequestOption) *RequestBuilder[Req, Resp] {
	r.method = http.MethodGet
	r.resourcePath = path
	r.requestOptions = append(r.requestOptions, opts...)
	return r
}

// Getf is like Get but accepts a format string and args.
func (r *RequestBuilder[Req, Resp]) Getf(path string, a ...any) *RequestBuilder[Req, Resp] {
	return r.Get(fmt.Sprintf(path, a...))
}

// Post sets the request method to POST, the request path to the given path, and the request payload to body. Global request options are applied.
func (r *RequestBuilder[Req, Resp]) Post(path string, body *Req, opts ...RequestOption) *RequestBuilder[Req, Resp] {
	r.method = http.MethodPost
	r.resourcePath = path
	r.body = body
	r.requestOptions = append(r.requestOptions, opts...)
	return r
}

func (r *RequestBuilder[Req, Resp]) Postf(path string, a ...any) *RequestBuilder[Req, Resp] {
	r.method = http.MethodPost
	r.resourcePath = fmt.Sprintf(path, a...)
	return r
}

func (r *RequestBuilder[Req, Resp]) Body(body *Req) *RequestBuilder[Req, Resp] {
	r.body = body
	return r
}

// Patch sets the request method to PATCH, the request path to the given path, and the request payload to body. Global request options are applied.
func (r *RequestBuilder[Req, Resp]) Patch(path string, body *Req, opts ...RequestOption) *RequestBuilder[Req, Resp] {
	r.method = http.MethodPatch
	r.resourcePath = path
	r.requestOptions = append(r.requestOptions, opts...)
	r.body = body
	return r
}

// All finalizes the request and returns an iterator over the response.
func (r *RequestBuilder[Req, Resp]) All(ctx context.Context) *Iter[Resp] {
	return GetIter(func(lastMeta *ListMeta) (*List[Resp], error) {
		ctx, cancel := context.WithDeadline(ctx, time.Now().Add(90*time.Second))
		defer cancel()

		list := new(List[Resp])
		response, err := r.makeRequest(ctx, WithRequestPagination(lastMeta))
		if err != nil {
			return nil, errors.Wrap(err, "failed to make request")
		}

		container := new(ResponsePayload[[]*Resp])
		err = decodeResponse(response, &container)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode response")
		}

		list.SetListMeta(container.Meta)
		list.SetItems(container.Data)
		return list, nil
	})
}

// One finalizes the request and returns a single item response.
func (r *RequestBuilder[Req, Resp]) One(ctx context.Context) (*Resp, error) {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(90*time.Second))
	defer cancel()

	response, err := r.makeRequest(ctx)
	if err != nil {
		return nil, err
	}
	container := new(ResponsePayload[*Resp])
	err = decodeResponse(response, &container)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	return container.Data, nil
}

func (r *RequestBuilder[Req, Resp]) makeRequest(ctx context.Context, opts ...RequestOption) (*http.Response, error) {
	return r.client.Do(ctx, r.resourcePath, r.method, r.body, append(r.requestOptions, opts...)...)
}

func decodeResponse[T any](resp *http.Response, v T) error {
	reader, err := gzipResponseReader(resp)
	if err != nil {
		return err
	}
	defer reader.Close()

	return json.NewDecoder(reader).Decode(v)
}

// normalizeParams returns a slice of interfaces from the given params.
// This is only neeeded because Go doesn't allow slice conversion of slice spreads.
// See: https://github.com/golang/go/wiki/InterfaceSlice
func normalizeParams[T ParamEncoder[T]](params []T) []ParamEncoder[T] {
	p := make([]ParamEncoder[T], len(params))
	for i, param := range params {
		p[i] = ParamEncoder[T](param)
	}
	return p
}
