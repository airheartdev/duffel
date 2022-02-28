package duffel

import (
	"context"
	"fmt"
	"strings"
)

const orderChangeRequestIDPrefix = "ocr_"

type (

	// OrderChangeRequest is the input to the OrderChange API.
	// To change an order, you'll need to create an order change request.
	// An order change request describes the slices of an existing paid order
	// that you want to remove and search criteria for new slices you want to add.
	OrderChangeRequest struct {
		ID                string                  `json:"id"`
		OrderID           string                  `json:"order_id"`
		Slices            []SliceChangesetRequest `json:"slices"`
		OrderChangeOffers []OrderChangeOffer      `json:"order_change_offers"`
		CreatedAt         string                  `json:"created_at"`
		UpdatedAt         string                  `json:"updated_at"`
		LiveMode          bool                    `json:"live_mode"`
	}

	OrderChangeOffer struct {
		ID                   string         `json:"id"`
		OrderChangeID        string         `json:"order_change_id"`
		UpdatedAt            string         `json:"updated_at"`
		Slices               SliceChangeset `json:"slices"`
		RefundTo             string         `json:"refund_to"`
		PenaltyTotalCurrency string         `json:"penalty_total_currency"`
		PenaltyTotalAmount   string         `json:"penalty_total_amount"`
		NewTotalCurrency     string         `json:"new_total_currency"`
		NewTotalAmount       string         `json:"new_total_amount"`
		LiveMode             bool           `json:"live_mode"`
		ExpiresAt            string         `json:"expires_at"`
		CreatedAt            string         `json:"created_at"`
		ChangeTotalCurrency  string         `json:"change_total_currency"`
		ChangeTotalAmount    string         `json:"change_total_amount"`
	}

	SliceChangeset struct {
		Add    []Slice `json:"add"`
		Remove []Slice `json:"remove"`
	}

	OrderChangeRequestParams struct {
		OrderID string                `json:"order_id"`
		Slices  SliceChangesetRequest `json:"slices"`
	}

	SliceAddRequest struct {
		OfferRequestSlice
		CabinClass CabinClass `json:"cabin_class"`
	}

	SliceRemoveRequest struct {
		SliceID string `json:"slice_id"`
	}

	SliceChangesetRequest struct {
		Add    []SliceAddRequest    `json:"add"`
		Remove []SliceRemoveRequest `json:"remove"`
	}

	OrderChangeClient interface {
		CreateOrderChangeRequest(ctx context.Context, orderID string, params OrderChangeRequestParams) (*OrderChangeRequest, error)
		GetOrderChangeRequest(ctx context.Context, id string) (*OrderChangeRequest, error)
	}
)

func (a *API) GetOrderChangeRequest(ctx context.Context, orderChangeRequestID string) (*OrderChangeRequest, error) {
	if orderChangeRequestID == "" {
		return nil, fmt.Errorf("orderChangeRequestID param is required")
	} else if !strings.HasPrefix(orderChangeRequestID, orderChangeRequestIDPrefix) {
		return nil, fmt.Errorf("orderChangeRequestID should begin with %s", orderChangeRequestIDPrefix)
	}

	return newRequestWithAPI[EmptyPayload, OrderChangeRequest](a).
		Get(fmt.Sprintf("/api/order_change_requests/%s", orderChangeRequestID)).
		One(ctx)
}
