// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package duffel

import (
	"context"
)

type (
	PaymentType string

	Payment struct {
		Amount    string      `json:"amount"`
		CreatedAt DateTime    `json:"created_at"`
		Currency  string      `json:"currency"`
		ID        string      `json:"id"`
		LiveMode  bool        `json:"live_mode"`
		Type      PaymentType `json:"type"`
	}

	CreatePaymentRequest struct {
		OrderID string        `json:"order_id"`
		Payment CreatePayment `json:"payment"`
	}

	CreatePayment struct {
		Amount   string      `json:"amount"`
		Currency string      `json:"currency"`
		Type     PaymentType `json:"type"`
	}

	OrderPaymentClient interface {
		CreatePayment(ctx context.Context, req CreatePaymentRequest) (*Payment, error)
	}
)

const (
	PaymentTypeBalance = PaymentType("balance")
	PaymentTypeCash    = PaymentType("arc_bsp_cash")
)

func (a *API) CreatePayment(ctx context.Context, req CreatePaymentRequest) (*Payment, error) {
	return newRequestWithAPI[CreatePaymentRequest,Payment](a).Post("/air/payments", &req).One(ctx)
}

var _ OrderPaymentClient = (*API)(nil)
