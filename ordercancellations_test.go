package duffel

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestCancelOrder(t *testing.T) {
	defer gock.Off()
	// gock.Observe(gock.DumpRequest)

	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Post("/air/order_cancellations").
		Reply(201).
		JSON(`{"data":{"order_id":"ord_00009hthhsUZ8W4LxQgkjo"}}`).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/201-create-order-cancellation.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	data, err := client.CreateOrderCancellation(ctx, "ord_00009hthhsUZ8W4LxQgkjo")
	a.NoError(err)
	a.NotNil(data)
	a.Equal("90.80 GBP", data.RefundAmount().String())
}

func TestConfirmCancelOrder(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Post("/air/order_cancellations/ore_00009qzZWzjDipIkqpaUAj/actions/confirm").
		BodyString("").
		Reply(201).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/201-create-order-cancellation.json") // Uses the same payload as the create endpoint.

	ctx := context.TODO()

	client := New("duffel_test_123")
	data, err := client.ConfirmOrderCancellation(ctx, "ore_00009qzZWzjDipIkqpaUAj")
	a.NoError(err)
	a.NotNil(data)
	a.Equal("ore_00009qzZWzjDipIkqpaUAj", data.ID)
	a.Equal("90.80 GBP", data.RefundAmount().String())
}

func TestGetOrderCancellation(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Get("/air/order_cancellations/ore_00009qzZWzjDipIkqpaUAj").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-get-order-cancellation.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	data, err := client.GetOrderCancellation(ctx, "ore_00009qzZWzjDipIkqpaUAj")
	a.NoError(err)
	a.NotNil(data)
	a.Equal("90.80 GBP", data.RefundAmount().String())
}
