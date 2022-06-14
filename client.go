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
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/segmentio/encoding/json"
)

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

	if resp.StatusCode > 399 {
		err = decodeError(resp)
		return nil, err
	}

	return resp, nil
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
