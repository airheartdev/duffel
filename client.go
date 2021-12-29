package duffel

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

func (c *client) makeRequest(ctx context.Context, resourceName string, method string, body io.Reader, opts ...Option) (*http.Response, error) {
	u, err := c.buildRequestURL(resourceName)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("user-agent", c.options.UserAgent)
	req.Header.Add("Duffel-Version", c.options.Version)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.APIToken))

	resp, err := c.httpDoer.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 499 {
		return nil, fmt.Errorf("request failed with HTTP status: %s url=%s", resp.Status, req.URL.String())
	} else if resp.StatusCode > 399 {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			fallthrough
		case http.StatusNotFound:
			return nil, ErrNotFound
		default:
			return nil, fmt.Errorf("request failed with HTTP status: %s url=%s", resp.Status, req.URL.String())
		}
	}

	return resp, nil
}

func (c *client) buildRequestURL(resourceName string) (*url.URL, error) {
	u, err := url.Parse(path.Join(c.options.Host, resourceName))
	if err != nil {
		return nil, err
	}

	return u, nil
}
