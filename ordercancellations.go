// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package duffel

import (
	"context"
	"fmt"
	"strings"

	"github.com/bojanz/currency"
)

const orderCancellationIDPrefix = "ore_"

type (
	// OrderCancellationRequest is response from the OrderCancellation API.
	//
	// Once you've created a pending order cancellation, you'll know
	// the `refund_amount` you're due to get back.
	//
	// To actually cancel the order, you'll need to confirm the cancellation.
	// The booking with the airline will be cancelled, and the `refund_amount` will be
	// returned to the original payment method (i.e. your Duffel balance).
	// You'll then need to refund your customer (e.g. back to their credit/debit card).
	OrderCancellation struct {
		ID                string        `json:"id"`
		OrderID           string        `json:"order_id"`
		RefundTo          PaymentMethod `json:"refund_to"`
		RawRefundCurrency string        `json:"refund_currency"`
		RawRefundAmount   string        `json:"refund_amount"`
		ExpiresAt         string        `json:"expires_at"`
		CreatedAt         string        `json:"created_at"`
		ConfirmedAt       string        `json:"confirmed_at"`
		LiveMode          bool          `json:"live_mode"`
	}

	OrderCancellationRequest struct {
		OrderID string `json:"order_id"`
	}

	// OrderCancellationClient
	OrderCancellationClient interface {
		CreateOrderCancellation(ctx context.Context, orderID string) (*OrderCancellation, error)
		ConfirmOrderCancellation(ctx context.Context, orderCancellationID string) (*OrderCancellation, error)
		GetOrderCancellation(ctx context.Context, orderCancellationID string) (*OrderCancellation, error)
	}
)

func (a *API) CreateOrderCancellation(ctx context.Context, orderID string) (*OrderCancellation, error) {
	return newRequestWithAPI[OrderCancellationRequest, OrderCancellation](a).
		Post("/air/order_cancellations", &OrderCancellationRequest{
			OrderID: orderID,
		}).
		Single(ctx)
}

func (a *API) ConfirmOrderCancellation(ctx context.Context, orderCancellationID string) (*OrderCancellation, error) {
	if !strings.HasPrefix(orderCancellationID, orderCancellationIDPrefix) {
		return nil, fmt.Errorf("orderCancellationID should have prefix %s, got %s", orderCancellationIDPrefix, orderCancellationID[:4])
	}

	return newRequestWithAPI[EmptyPayload, OrderCancellation](a).
		Post(fmt.Sprintf("/air/order_cancellations/%s/actions/confirm", orderCancellationID), nil).
		Single(ctx)
}

func (a *API) GetOrderCancellation(ctx context.Context, orderCancellationID string) (*OrderCancellation, error) {
	if !strings.HasPrefix(orderCancellationID, orderCancellationIDPrefix) {
		return nil, fmt.Errorf("orderCancellationID should have prefix %s, got %s", orderCancellationIDPrefix, orderCancellationID[:4])
	}

	return newRequestWithAPI[EmptyPayload, OrderCancellation](a).
		Getf("/air/order_cancellations/%s", orderCancellationID).
		Single(ctx)
}

func (o *OrderCancellation) RefundAmount() currency.Amount {
	amount, err := currency.NewAmount(o.RawRefundAmount, o.RawRefundCurrency)
	if err != nil {
		return currency.Amount{}
	}
	return amount
}

var _ OrderCancellationClient = (*API)(nil)
