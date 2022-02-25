package duffel

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type (
	AirportsClient interface {
		ListAirports(ctx context.Context, input *ListAirportsRequest) *Iter[Airport]
	}

	ListAirportsRequest struct {
		IATACountryCode string `json:"iata_country_code"`
		Limit           int    `json:"limit"`
	}

	ListAirportsResponse struct {
	}
)

func (a *API) ListAirports(ctx context.Context, input *ListAirportsRequest) *Iter[Airport] {
	c := newInternalClient[ListAirportsRequest, []*Airport](a)

	return GetIter(func(lastMeta *ListMeta) (*List[Airport], error) {
		list := new(List[Airport])
		response, err := c.makeIteratorRequest(ctx, "/air/airports", http.MethodGet, input, func(req *http.Request) {
			q := req.URL.Query()
			if lastMeta != nil && lastMeta.After != "" {
				q.Add("after", lastMeta.After)
			}
			if input != nil {
				q.Add("iata_country_code", input.IATACountryCode)
				if input.Limit > 0 {
					q.Add("limit", fmt.Sprintf("%d", input.Limit))
				}
			}
			req.URL.RawQuery = q.Encode()
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to list airports")
		}
		if response == nil {
			return nil, fmt.Errorf("internal: empty response")
		}
		list.SetListMeta(response.Meta)
		list.SetItems(response.Data)
		return list, nil
	})
}

var _ AirportsClient = (*API)(nil)
