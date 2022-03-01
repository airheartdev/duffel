package duffel

import (
	"context"
	"net/url"
	"strconv"
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
		FamilyName               string                    `json:"family_name,omitempty"`
		GivenName                string                    `json:"given_name,omitempty"`
		Age                      int                       `json:"age,omitempty"`
		LoyaltyProgrammeAccounts []LoyaltyProgrammeAccount `json:"loyalty_programme_accounts,omitempty"`
		// deprecated
		Type PassengerType `json:"type,omitempty"`
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
	return newRequestWithAPI[OfferRequestInput, OfferRequest](a).
		Post("/air/offer_requests", &requestInput).
		WithParams(requestInput).
		One(ctx)
}

func (a *API) GetOfferRequest(ctx context.Context, id string) (*OfferRequest, error) {
	return newRequestWithAPI[EmptyPayload, OfferRequest](a).Getf("/air/offer_requests/%s", id).One(ctx)
}

func (a *API) ListOfferRequests(ctx context.Context) *Iter[OfferRequest] {
	return newRequestWithAPI[EmptyPayload, OfferRequest](a).Get("/air/offer_requests").All(ctx)
}

// Encode implements the ParamEncoder interface.
func (o OfferRequestInput) Encode(q url.Values) error {
	q.Set("return_offers", strconv.FormatBool(o.ReturnOffers))
	return nil
}
