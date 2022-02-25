package duffel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Payload[T any] struct {
	Data T `json:"data"`
}

type RequestOption func(req *http.Request, u *url.URL)

func buildPayload[T any](data T) *Payload[T] {
	return &Payload[T]{
		Data: data,
	}
}

func encodePayload[T any](requestInput T) (*bytes.Buffer, error) {
	payload := bytes.NewBuffer(nil)
	err := json.NewEncoder(payload).Encode(buildPayload(requestInput))
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (c *client[R, T]) makeRequest(ctx context.Context, resourceName string, method string, body io.Reader, opts ...RequestOption) (*http.Response, error) {
	if c.APIToken == "" {
		return nil, fmt.Errorf("duffel: missing API token")
	}

	u, err := c.buildRequestURL(resourceName)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", c.options.UserAgent)
	req.Header.Add("Duffel-Version", c.options.Version)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.APIToken))
	for _, o := range opts {
		o(req, u)
	}

	resp, err := c.httpDoer.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 499 {
		return nil, fmt.Errorf("request failed with HTTP status: %s url=%s", resp.Status, req.URL.String())
	} else if resp.StatusCode > 399 {
		derr := &DuffelError{}
		err := json.NewDecoder(resp.Body).Decode(derr)
		if err != nil {
			return nil, err
		}
		return nil, derr
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
