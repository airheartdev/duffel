// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package duffel

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestCreatePayment(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)

	gock.New("https://api.duffel.com").
		Post("/air/payments").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-create-payment.json")

	expected := &Payment{
		Type:      PaymentTypeBalance,
		Amount:    "30.20",
		ID:        "pay_00009hthhsUZ8W4LxQgkjo",
		Currency:  "GBP",
		LiveMode:  false,
		CreatedAt: DateTime(time.Date(2020, 0o4, 11, 15, 48, 11, 642000000, time.UTC)),
	}

	ctx := context.TODO()

	client := New("duffel_test_123")
	payment, err := client.CreatePayment(ctx, CreatePaymentRequest{
		OrderID: "ord_00003x8pVDGcS8y2AWCoWv",
		Payment: CreatePayment{
			Amount:   "30.20",
			Currency: "GBP",
			Type:     PaymentTypeBalance,
		},
	})

	a.NoError(err)
	a.Equal(expected, payment)
}
