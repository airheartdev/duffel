package duffel

import (
	"context"
	"net/http"
	"time"
)

type (
	OfferRequestClient interface {
		CreateOfferRequest(ctx context.Context, requestInput OfferRequestInput) (*OfferRequest, error)
		GetOfferRequest(ctx context.Context, id string) (*OfferRequest, error)
		ListOfferRequests(ctx context.Context) *Iter[OfferRequest]
	}

	OfferRequestInput struct {
		Passengers   []OfferRequestPassenger `json:"passengers" url:"-"`
		Slices       []OfferRequestSlice     `json:"slices" url:"-"`
		CabinClass   CabinClass              `json:"cabin_class" url:"-"`
		ReturnOffers bool                    `json:"-" url:"return_offers"`
	}

	OfferRequestSlice struct {
		DepartureDate Date   `json:"departure_date"`
		Destination   string `json:"destination"`
		Origin        string `json:"origin"`
	}

	OfferRequestPassenger struct {
		FamilyName               string                    `json:"family_name"`
		GivenName                string                    `json:"given_name"`
		Age                      int                       `json:"age,omitempty"`
		LoyaltyProgrammeAccounts []LoyaltyProgrammeAccount `json:"loyalty_programme_accounts,omitempty"`
		// deprecated
		Type *PassengerType `json:"type,omitempty"`
	}

	// OfferRequest is the response from the OfferRequest endpoint, created using the OfferRequestInput.
	OfferRequest struct {
		ID         string                  `json:"id"`
		LiveMode   bool                    `json:"live_mode"`
		CreatedAt  time.Time               `json:"created_at"`
		Slices     []BaseSlice             `json:"slices"`
		Passengers []OfferRequestPassenger `json:"passengers"`
		CabinClass CabinClass              `json:"cabin_class"`
		Offers     []Offer                 `json:"offers"`
	}
)

func (a *API) CreateOfferRequest(ctx context.Context, requestInput OfferRequestInput) (*OfferRequest, error) {
	client := newInternalClient[OfferRequestInput, OfferRequest](a)
	return client.makeRequestWithPayload(ctx,
		"/air/offer_requests",
		http.MethodPost,
		&requestInput,
		WithURLParams(requestInput),
	)
}

func (a *API) GetOfferRequest(ctx context.Context, id string) (*OfferRequest, error) {
	c := newInternalClient[OfferRequestInput, OfferRequest](a)
	return c.makeRequestWithPayload(ctx, "/air/offer_requests/"+id, http.MethodGet, nil)
}

func (a *API) ListOfferRequests(ctx context.Context) *Iter[OfferRequest] {
	c := newInternalClient[ListAirportsParams, OfferRequest](a)
	return c.getIterator(ctx, http.MethodGet, "/air/offer_requests")
}
