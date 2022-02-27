package duffel

import (
	"context"
	"testing"
	"time"

	"github.com/segmentio/encoding/json"
	"gopkg.in/h2non/gock.v1"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalOrderWithCurrency(t *testing.T) {
	a := assert.New(t)
	order := new(Order)

	err := json.Unmarshal([]byte(`{"base_amount": "120.00", "base_currency": "USD", "cancelled_at":"2020-04-11T15:48:11.642Z"}`), order)
	a.NoError(err)

	a.Equal("120.00 USD", order.BaseAmount().String())
	a.Equal("2020-04-11 15:48:11.642 +0000 UTC", order.CancelledAt.String())
}

func TestCreateOrder(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)

	gock.New("https://api.duffel.com").
		Post("/air/orders").
		Reply(201).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/201-create-order.json")

	ctx := context.TODO()
	client := New("duffel_test_123")
	order, err := client.CreateOrder(ctx, CreateOrderInput{})
	a.NoError(err)
	a.Equal("RZPNX8", order.BookingReference)
	a.Equal("ord_00009hthhsUZ8W4LxQgkjo", order.ID)
}

func TestListOrders(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)

	gock.New("https://api.duffel.com").
		Get("/air/orders").
		MatchParam("after", "g2wAAAACbQAAABBBZXJvbWlzdC1LaGFya2l2bQAAAB=").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-list-orders-page2.json")

	gock.New("https://api.duffel.com").
		Get("/air/orders").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-list-orders.json")

	ctx := context.TODO()
	client := New("duffel_test_123")
	iter := client.ListOrders(ctx, ListOrdersParams{
		BookingReference: "RZPNX8",
	})

	iter.Next()
	a.NoError(iter.Err())
	order1 := iter.Current()

	a.Equal("RZPNX8", order1.BookingReference)
	a.Equal("ord_00009hthhsUZ8W4LxQgkjo", order1.ID)

	iter.Next()
	a.NoError(iter.Err())
	order2 := iter.Current()

	a.Equal("ABC123", order2.BookingReference)
	a.Equal("ord_00009hthhsUZ8W4LxQgkjo", order2.ID)
}

func TestGetOrderByID(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)

	gock.New("https://api.duffel.com").
		Get("/air/orders/ord_00009hthhsUZ8W4LxQgkjo").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-get-order.json")

	ctx := context.TODO()
	client := New("duffel_test_123")
	order, err := client.GetOrder(ctx, "ord_00009hthhsUZ8W4LxQgkjo")
	a.NoError(err)

	a.Equal("RZPNX8", order.BookingReference)
	a.Equal("ord_00009hthhsUZ8W4LxQgkjo", order.ID)
}

func TestUpdateOrder(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)

	gock.New("https://api.duffel.com").
		Patch("/air/orders/ord_00009hthhsUZ8W4LxQgkjo").
		MatchType("json").
		JSON(Payload[OrderUpdateParams]{Data: OrderUpdateParams{
			Metadata: map[string]any{"seat_preference": "window"},
		}}).
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-update-order.json")

	ctx := context.TODO()
	client := New("duffel_test_123")
	order, err := client.UpdateOrder(ctx, "ord_00009hthhsUZ8W4LxQgkjo", OrderUpdateParams{
		Metadata: map[string]any{
			"seat_preference": "window",
		},
	})
	a.NoError(err)

	a.Equal("RZPNX8", order.BookingReference)
	a.Equal(Metadata{"seat_preference": "window"}, order.Metadata)
	a.Equal("ord_00009hthhsUZ8W4LxQgkjo", order.ID)
}
