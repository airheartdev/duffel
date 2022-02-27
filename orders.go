package duffel

import (
	"context"
	"net/http"
	"time"
)

const orderIDPrefix = "ord_"

type (
	ListOrdersSort string

	Order struct{}

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
		Metadata map[string]any `json:"metadata,omitempty"`

		// The personal details of the passengers, expanding on
		// the information initially provided when creating the offer request.
		Passengers []PassengerCreateInput `json:"passengers"`

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
		ListOrders(ctx context.Context, params ...ListOrdersParams) Iter[*Order]

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
	c := newInternalClient[CreateOrderInput, Order](a)
	return c.makeRequestWithPayload(ctx, "/air/orders", http.MethodPost, &input)
}
