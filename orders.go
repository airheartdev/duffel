// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package duffel

import (
	"context"
	"net/url"
	"time"

	"github.com/bojanz/currency"
	"github.com/gorilla/schema"
)

const orderIDPrefix = "ord_"

type (
	ListOrdersSort string

	Order struct {
		ID               string           `json:"id"`
		LiveMode         bool             `json:"live_mode"`
		Metadata         Metadata         `json:"metadata"`
		RawBaseAmount    *string          `json:"base_amount,omitempty"`
		RawBaseCurrency  *string          `json:"base_currency,omitempty"`
		BookingReference string           `json:"booking_reference"`
		CancelledAt      *time.Time       `json:"cancelled_at,omitempty"`
		CreatedAt        time.Time        `json:"created_at"`
		Conditions       Conditions       `json:"conditions,omitempty"`
		Documents        []Document       `json:"documents,omitempty"`
		Owner            Airline          `json:"owner"`
		Passengers       []OrderPassenger `json:"passengers,omitempty"`
		PaymentStatus    PaymentStatus    `json:"payment_status"`
		Services         []Service        `json:"services,omitempty"`
		Slices           []Slice          `json:"slices,omitempty"`
		SyncedAt         time.Time        `json:"synced_at"`
		RawTaxAmount     *string          `json:"tax_amount,omitempty"`
		RawTaxCurrency   *string          `json:"tax_currency,omitempty"`
		RawTotalAmount   string           `json:"total_amount"`
		RawTotalCurrency string           `json:"total_currency"`
	}

	SliceConditions struct {
		ChangeBeforeDeparture *ChangeCondition `json:"change_before_departure,omitempty"`
	}

	Conditions struct {
		RefundBeforeDeparture *ChangeCondition `json:"refund_before_departure,omitempty"`
		ChangeBeforeDeparture *ChangeCondition `json:"change_before_departure,omitempty"`
	}

	ChangeCondition struct {
		Allowed            bool    `json:"allowed"`
		RawPenaltyAmount   *string `json:"penalty_amount,omitempty"`
		RawPenaltyCurrency *string `json:"penalty_currency,omitempty"`
	}

	Document struct {
		Type             string `json:"type"`
		UniqueIdentifier string `json:"unique_identifier"`
	}

	// NOTE: If you receive a 500 Internal Server Error when trying to create an order,
	// it may have still been created on the airline’s side.
	// Please contact Duffel support before trying the request again.

	OrderType string

	CreateOrderInput struct {
		Type OrderType `json:"type"`

		// Metadata contains a set of key-value pairs that you can attach to an object.
		// It can be useful for storing additional information about the object, in a
		// structured format. Duffel does not use this information.
		//
		// You should not store sensitive information in this field.
		Metadata Metadata `json:"metadata,omitempty"`

		// The personal details of the passengers, expanding on
		// the information initially provided when creating the offer request.
		Passengers []OrderPassenger `json:"passengers"`

		Payments []PaymentCreateInput `json:"payments,omitempty"`

		// The ids of the offers you want to book. You must specify an array containing exactly one selected offer.
		SelectedOffers []string `json:"selected_offers"`

		Services []ServiceCreateInput `json:"services,omitempty"`
	}

	// The services you want to book along with the first selected offer. This key should be omitted when the order’s type is hold, as we do not support services for hold orders yet.
	ServiceCreateInput struct {
		// The id of the service from the offer's available_services that you want to book
		ID string `json:"id"`

		// The quantity of the service to book. This will always be 1 for seat services.
		Quantity int `json:"quantity"`
	}

	Service struct {
		// Duffel's unique identifier for the booked service
		ID string `json:"id"`

		// The metadata varies by the type of service.
		// It includes further data about the service. For example, for
		// baggages, it may have data about size and weight restrictions.
		Metadata Metadata `json:"metadata"`

		// List of passenger ids the service applies to.
		// The service applies to all the passengers in this list.
		PassengerIDs []string `json:"passenger_ids"`

		// The quantity of the service that was booked
		Quantity int `json:"quantity"`

		// List of segment ids the service applies to. The service applies to all the segments in this list.
		SegmentIDs []string `json:"segment_ids"`

		// The total price of the service for all passengers and segments it applies to, accounting for quantity and including taxes
		RawTotalAmount   string `json:"total_amount,omitempty"`
		RawTotalCurrency string `json:"total_currency,omitempty"`

		// Possible values: "baggage" or "seat"
		Type string `json:"type"`
	}

	AvailableService struct {
		// Duffel's unique identifier for service
		ID               string                   `json:"id"`
		MaximumQuantity  int                      `json:"maximum_quantity"`
		Metadata         AvailableServiceMetadata `json:"metadata"`
		PassengerIDs     []string                 `json:"passenger_ids"`
		SegmentIDs       []string                 `json:"segment_ids"`
		RawTotalAmount   string                   `json:"total_amount"`
		RawTotalCurrency string                   `json:"total_currency"`

		// Possible values: "baggage"
		Type string `json:"type"`
	}

	AvailableServiceMetadata struct {
		MaximumDepthCM  int `json:"maximum_depth_cm,omitempty"`
		MaximumHeightCM int `json:"maximum_height_cm,omitempty"`
		MaximumLengthCM int `json:"maximum_length_cm,omitempty"`
		MaximumWeightKg int `json:"maximum_weight_kg,omitempty"`
		// Possible values: "checked", "carry_on"
		Type string `json:"type"`
	}

	// OrderUpdateParams is used as the input to UpdateOrder.
	// Only certain order fields are updateable.
	// Each field that can be updated is detailed in the `OrderUpdateParams` object.
	OrderUpdateParams struct {
		Metadata map[string]any
	}

	ListOrdersParams struct {
		// Filters orders by their booking reference.
		// The filter requires an exact match but is case insensitive.
		BookingReference string `url:"booking_reference,omitempty"`

		// Whether to filter orders that are awaiting payment or not.
		// If not specified, all orders regardless of their payment state will be returned.
		AwaitingPayment bool `url:"awaiting_payment,omitempty"`

		// By default, orders aren't returned in any specific order.
		// This parameter allows you to sort the list of orders by the payment_required_by date
		Sort ListOrdersSort `url:"sort,omitempty"`

		// Filters the returned orders by owner.id. Values must be valid airline.ids.
		// Check the Airline schema for details.
		OwnerIDs []string `url:"owner_id,omitempty"`

		// Filters the returned orders by origin. Values must be valid origin identifiers.
		// Check the Order schema for details.
		OriginIDs []string `url:"origin_id,omitempty"`

		// Filters the returned orders by destination. Values must be valid destination identifiers.
		// Check the Order schema for details.
		DestinationIDs []string `url:"destination_id,omitempty"`

		// Filters the returned orders by departure datetime.
		// Orders will be included if any of their segments matches the given criteria
		DepartingAt *TimeFilter `url:"departing_at,omitempty"`

		// Filters the returned orders by arrival datetime.
		// Orders will be included if any of their segments matches the given criteria.
		ArrivingAt *TimeFilter `url:"arriving_at,omitempty"`

		// Filters the returned orders by creation datetime.
		CreatedAt *TimeFilter `url:"created_at,omitempty"`

		// Orders will be included if any of their passengers matches any of the given names.
		// Matches are case insensitive, and include partial matches.
		PassengerNames []string `url:"passenger_name,omitempty"`
	}

	Metadata map[string]any

	TimeFilter struct {
		Before *time.Time `url:"before,omitempty"`
		After  *time.Time `url:"after,omitempty"`
	}

	OrderClient interface {
		// Get a single order by ID.
		GetOrder(ctx context.Context, id string) (*Order, error)

		// Update a single order by ID.
		UpdateOrder(ctx context.Context, id string, params OrderUpdateParams) (*Order, error)

		// List orders.
		ListOrders(ctx context.Context, params ...ListOrdersParams) *Iter[Order]

		// Create an order.
		CreateOrder(ctx context.Context, input CreateOrderInput) (*Order, error)
	}
)

const (
	ListOrdersSortPaymentRequiredByAsc  ListOrdersSort = "payment_required_by"
	ListOrdersSortPaymentRequiredByDesc ListOrdersSort = "-payment_required_by"

	OrderTypeHold    OrderType = "hold"
	OrderTypeInstant OrderType = "instant"
)

// CreateOrder creates a new order.
func (a *API) CreateOrder(ctx context.Context, input CreateOrderInput) (*Order, error) {
	return newRequestWithAPI[CreateOrderInput, Order](a).Post("/air/orders", &input).Single(ctx)
}

func (a *API) UpdateOrder(ctx context.Context, id string, params OrderUpdateParams) (*Order, error) {
	return newRequestWithAPI[OrderUpdateParams, Order](a).Patch("/air/orders/"+id, &params).Single(ctx)
}

// CreateOrder creates a new order.
func (a *API) GetOrder(ctx context.Context, id string) (*Order, error) {
	return newRequestWithAPI[EmptyPayload, Order](a).Get("/air/orders/" + id).Single(ctx)
}

func (a *API) ListOrders(ctx context.Context, params ...ListOrdersParams) *Iter[Order] {
	return newRequestWithAPI[ListOrdersParams, Order](a).
		Get("/air/orders").
		WithParams(normalizeParams(params)...).
		Iter(ctx)
}

func (o *Order) BaseAmount() *currency.Amount {
	if o.RawBaseAmount != nil && o.RawBaseCurrency != nil {
		amount, err := currency.NewAmount(*o.RawBaseAmount, *o.RawBaseCurrency)
		if err != nil {
			return nil
		}
		return &amount
	}
	return nil
}

func (o *Order) TaxAmount() *currency.Amount {
	if o.RawTaxAmount != nil && o.RawTaxCurrency != nil {
		amount, err := currency.NewAmount(*o.RawTaxAmount, *o.RawTaxCurrency)
		if err != nil {
			return nil
		}
		return &amount
	}
	return nil
}

func (o *Order) TotalAmount() currency.Amount {
	amount, err := currency.NewAmount(o.RawTotalAmount, o.RawTotalCurrency)
	if err != nil {
		return currency.Amount{}
	}
	return amount
}

func (c *ChangeCondition) PenaltyAmount() *currency.Amount {
	if c.RawPenaltyAmount != nil && c.RawPenaltyCurrency != nil {
		amount, err := currency.NewAmount(*c.RawPenaltyAmount, *c.RawPenaltyCurrency)
		if err != nil {
			return nil
		}
		return &amount
	}

	return nil
}

func (s *Service) TotalAmount() currency.Amount {
	amount, err := currency.NewAmount(s.RawTotalAmount, s.RawTotalCurrency)
	if err != nil {
		return currency.Amount{}
	}
	return amount
}

func (o ListOrdersParams) Encode(q url.Values) error {
	enc := schema.NewEncoder()
	enc.SetAliasTag("url")
	return enc.Encode(o, q)
}
