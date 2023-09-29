// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package duffel

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/segmentio/encoding/json"
)

const RequestIDHeader = "x-request-id"

type Payload[T any] struct {
	Data T `json:"data"`
}

type ResponsePayload[T any] struct {
	Meta *ListMeta `json:"meta"`
	Data T         `json:"data"`
}

type RequestOption func(req *http.Request) error

type EmptyPayload struct{}

func buildRequestPayload[T any](data T) *Payload[T] {
	return &Payload[T]{
		Data: data,
	}
}

func encodePayload[T any](requestInput T) (io.ReadCloser, error) {
	payload := bytes.NewBuffer(nil)
	err := json.NewEncoder(payload).Encode(buildRequestPayload(requestInput))
	if err != nil {
		return nil, err
	}

	return io.NopCloser(payload), nil
}

func (c *client[R, T]) makeRequest(ctx context.Context, resourceName string, method string, body io.ReadCloser, opts ...RequestOption) (*http.Response, error) {
	if c.APIToken == "" {
		return nil, fmt.Errorf("duffel: missing API token")
	}

	u, err := c.buildRequestURL(resourceName)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if method != http.MethodGet {
		req.Body = body
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	if !c.options.Debug {
		req.Header.Add("Accept-Encoding", "gzip")
	}
	req.Header.Add("User-Agent", c.options.UserAgent)
	req.Header.Add("Duffel-Version", c.options.Version)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.APIToken))

	// Apply request options
	for _, o := range opts {
		if o != nil {
			err := o(req)
			if err != nil {
				return nil, err
			}
		}
	}

	if c.options.Debug {
		b, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			return nil, err
		}
		fmt.Printf("REQUEST:\n%s\n", string(b))
	}

	do := func(c *client[R, T], req *http.Request, reuse bool) (*http.Response, error) {
		if reuse && req.Body != nil {
			// In a way when we use retry functionality we have to copy
			// request and pass it to c.httpDoer.Do, but req.Clone() doesn't really deep clone Body
			// and we have to clone body manually as in httputil.DumpRequestOut
			//
			// Issue https://github.com/golang/go/issues/36095
			var b bytes.Buffer
			b.ReadFrom(req.Body)
			req.Body = ioutil.NopCloser(&b)

			cloneReq := req.Clone(ctx)
			cloneReq.Body = ioutil.NopCloser(bytes.NewReader(b.Bytes()))
			req = cloneReq
		}

		resp, err := c.httpDoer.Do(req)
		if err != nil {
			return nil, err
		}

		if c.options.Debug {
			b, err := httputil.DumpResponse(resp, true)
			if err != nil {
				return nil, err
			}
			fmt.Printf("RESPONSE:\n%s\n", string(b))
		}

		if resp.StatusCode >= 400 {
			err = decodeError(resp)
			return nil, err
		}
		return resp, nil
	}

	if c.retry == nil {
		// Do single request without using backoff retry mechanism
		return do(c, req, false)
	}

	for {
		resp, err := do(c, req, true)

		var isMatchedCond bool
		for _, cond := range c.options.Retry.Conditions {
			if ok := cond(resp, err); ok {
				isMatchedCond = true
				break
			}
		}
		if isMatchedCond {
			// Get next duration internval, sleep and make another request
			// till nextDuration != stopBackoff
			nextDuration := c.retry.next()
			if nextDuration == stopBackoff {
				c.retry.reset()
				return resp, err
			}
			time.Sleep(nextDuration)
			continue
		}

		// Break retries mechanism if conditions weren't matched
		return resp, err
	}
}

func (c *client[R, T]) buildRequestURL(resourceName string) (*url.URL, error) {
	u, err := url.Parse(c.options.Host)

	if err != nil {
		return nil, err
	}
	u.Path = resourceName

	return u, nil
}

func gzipResponseReader(response *http.Response) (io.ReadCloser, error) {
	var reader io.ReadCloser
	var err error
	if response.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
	} else {
		reader = response.Body
	}
	return reader, nil
}

func decodeError(response *http.Response) error {
	reader, err := gzipResponseReader(response)
	if err != nil {
		return err
	}

	contentType := response.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "text/html") {
		// Handle occasional HTML error pages at routing layer
		return &DuffelError{
			StatusCode: response.StatusCode,
			Retryable:  true,
			Errors: []Error{
				{
					Type:    ErrorType(InternalServerError),
					Title:   http.StatusText(response.StatusCode),
					Message: "An internal server error occurred. Please try again later.",
					Code:    InternalServerError,
				},
			},
		}
	}

	notRetryable := strings.HasPrefix(response.Request.URL.Path, "/air/orders") &&
		response.StatusCode == http.StatusInternalServerError

	derr := &DuffelError{
		StatusCode: response.StatusCode,
		Retryable:  !notRetryable,
	}
	err = json.NewDecoder(reader).Decode(derr)
	if err != nil {
		return err
	}
	return derr
}
