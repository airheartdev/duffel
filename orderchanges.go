// Order change flow:
// 1. Get an existing order by ID using client.GetOrder(...)
// 2. Create a new order change request using client.CreateOrderChangeRequest(...)
// 3. Get the order change offer using client.CreatePendingOrderChange(...)
package duffel

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bojanz/currency"
)

const (
	orderChangeRequestIDPrefix = "ocr_"
	orderChangeOfferIDPrefix   = "oco_"
	orderChangeIDPrefix        = "oce_"
)

type (

	// OrderChangeRequest is the input to the OrderChange API.
	// To change an order, you'll need to create an order change request.
	// An order change request describes the slices of an existing paid order
	// that you want to remove and search criteria for new slices you want to add.
	OrderChangeRequest struct {
		ID                string             `json:"id"`
		OrderID           string             `json:"order_id"`
		Slices            SliceChange        `json:"slices"`
		OrderChangeOffers []OrderChangeOffer `json:"order_change_offers"`
		CreatedAt         string             `json:"created_at"`
		UpdatedAt         string             `json:"updated_at"`
		LiveMode          bool               `json:"live_mode"`
	}

	OrderChangeOffer struct {
		ID string `json:"id"`
		// OrderChangeID is the ID for an order change if one has already been created from this order change offer
		OrderChangeID           string         `json:"order_change_id"`
		Slices                  SliceChangeset `json:"slices"`
		RefundTo                PaymentMethod  `json:"refund_to"`
		RawPenaltyTotalCurrency string         `json:"penalty_total_currency"`
		RawPenaltyTotalAmount   string         `json:"penalty_total_amount"`
		RawNewTotalCurrency     string         `json:"new_total_currency"`
		RawNewTotalAmount       string         `json:"new_total_amount"`
		RawChangeTotalCurrency  string         `json:"change_total_currency"`
		RawChangeTotalAmount    string         `json:"change_total_amount"`
		ExpiresAt               time.Time      `json:"expires_at"`
		CreatedAt               time.Time      `json:"created_at"`
		UpdatedAt               time.Time      `json:"updated_at"`
		LiveMode                bool           `json:"live_mode"`
	}

	OrderChange struct {
		ID                      string         `json:"id"`
		OrderID                 string         `json:"order_id"`
		Slices                  SliceChangeset `json:"slices"`
		RefundTo                PaymentMethod  `json:"refund_to"`
		RawPenaltyTotalCurrency string         `json:"penalty_total_currency"`
		RawPenaltyTotalAmount   string         `json:"penalty_total_amount"`
		RawNewTotalCurrency     string         `json:"new_total_currency"`
		RawNewTotalAmount       string         `json:"new_total_amount"`
		RawChangeTotalCurrency  string         `json:"change_total_currency"`
		RawChangeTotalAmount    string         `json:"change_total_amount"`
		ExpiresAt               string         `json:"expires_at"`
		CreatedAt               string         `json:"created_at"`
		UpdatedAt               string         `json:"updated_at"`
		LiveMode                bool           `json:"live_mode"`
	}

	SliceChangeset struct {
		Add    []Slice `json:"add"`
		Remove []Slice `json:"remove"`
	}

	OrderChangeRequestParams struct {
		OrderID string      `json:"order_id"`
		Slices  SliceChange `json:"slices,omitempty"`
	}

	SliceAdd struct {
		DepartureDate Date       `json:"departure_date"`
		Destination   string     `json:"destination"`
		Origin        string     `json:"origin"`
		CabinClass    CabinClass `json:"cabin_class"`
	}

	SliceRemove struct {
		SliceID string `json:"slice_id"`
	}

	SliceChange struct {
		Add    []SliceAdd    `json:"add,omitempty"`
		Remove []SliceRemove `json:"remove,omitempty"`
	}

	OrderChangeClient interface {
		CreateOrderChangeRequest(ctx context.Context, params OrderChangeRequestParams) (*OrderChangeRequest, error)
		GetOrderChangeRequest(ctx context.Context, id string) (*OrderChangeRequest, error)
		CreatePendingOrderChange(ctx context.Context, orderChangeRequestID string) (*OrderChange, error)
		ConfirmOrderChange(ctx context.Context, id string, payment PaymentCreateInput) (*OrderChange, error)
	}
)

func (a *API) CreateOrderChangeRequest(ctx context.Context, params OrderChangeRequestParams) (*OrderChangeRequest, error) {
	return newRequestWithAPI[OrderChangeRequestParams, OrderChangeRequest](a).
		Post("/air/order_change_requests", &params).
		One(ctx)
}

func (a *API) GetOrderChangeRequest(ctx context.Context, orderChangeRequestID string) (*OrderChangeRequest, error) {
	if err := validateID(orderChangeRequestID, orderChangeRequestIDPrefix); err != nil {
		return nil, err
	}
	return newRequestWithAPI[EmptyPayload, OrderChangeRequest](a).
		Getf("/air/order_change_requests/%s", orderChangeRequestID).
		One(ctx)
}

func (a *API) CreatePendingOrderChange(ctx context.Context, offerID string) (*OrderChange, error) {
	if err := validateID(offerID, orderChangeOfferIDPrefix); err != nil {
		return nil, err
	}
	return newRequestWithAPI[map[string]string, OrderChange](a).
		Postf("/air/order_changes").
		Body(&map[string]string{"selected_order_change_offer": offerID}).
		One(ctx)
}

func (a *API) ConfirmOrderChange(ctx context.Context, orderChangeRequestID string, payment PaymentCreateInput) (*OrderChange, error) {
	if err := validateID(orderChangeRequestID, orderChangeRequestIDPrefix); err != nil {
		return nil, err
	}
	return newRequestWithAPI[PaymentCreateInput, OrderChange](a).
		Postf("/air/order_changes/%s/actions/confirm", orderChangeRequestID).
		Body(&payment).
		One(ctx)
}

var _ OrderChangeClient = (*API)(nil)

func validateID(id, prefix string) error {
	if id == "" {
		return fmt.Errorf("id param is required")
	} else if !strings.HasPrefix(id, prefix) {
		return fmt.Errorf("id should begin with %s", prefix)
	}
	return nil
}

func (o *OrderChangeOffer) ChangeTotalAmount() currency.Amount {
	amount, err := currency.NewAmount(o.RawChangeTotalAmount, o.RawChangeTotalCurrency)
	if err != nil {
		return currency.Amount{}
	}
	return amount
}

func (o *OrderChangeOffer) NewTotalAmount() currency.Amount {
	amount, err := currency.NewAmount(o.RawNewTotalAmount, o.RawNewTotalCurrency)
	if err != nil {
		return currency.Amount{}
	}
	return amount
}

// PenaltyTotalAmount returns the penalty imposed by the airline for making this change.
func (o *OrderChangeOffer) PenaltyTotalAmount() currency.Amount {
	amount, err := currency.NewAmount(o.RawPenaltyTotalAmount, o.RawPenaltyTotalCurrency)
	if err != nil {
		return currency.Amount{}
	}
	return amount
}
