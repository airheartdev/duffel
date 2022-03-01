package duffel

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestCreateOrderChangeRequest(t *testing.T) {
	defer gock.Off()
	gock.Observe(gock.DumpRequest)

	gock.New("https://api.duffel.com").
		Post("/air/order_change_requests").
		JSON(`{
			"data": {
				"slices": {
					"remove": [
						{
							"slice_id": "sli_00009htYpSCXrwaB9Dn123"
						}
					],
					"add": [
						{
							"origin": "LHR",
							"destination": "JFK",
							"departure_date": "2020-04-24",
							"cabin_class": "economy"
						}
					]
				},
				"order_id": "ord_0000A3bQ8FJIQoEfuC07n6"
			}
		}`).
		Reply(201).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/201-create-order-change-request.json")

	a := assert.New(t)

	ctx := context.TODO()
	client := New("duffel_test_123")
	order, err := client.CreateOrderChangeRequest(ctx, OrderChangeRequestParams{
		OrderID: "ord_0000A3bQ8FJIQoEfuC07n6",
		Slices: SliceChange{
			Remove: []SliceRemove{
				{SliceID: "sli_00009htYpSCXrwaB9Dn123"},
			},
			Add: []SliceAdd{
				{
					Origin:        "LHR",
					Destination:   "JFK",
					DepartureDate: Date(time.Date(2020, time.April, 24, 0, 0, 0, 0, time.UTC)),
					CabinClass:    CabinClassEconomy,
				},
			},
		},
	})
	a.NoError(err)
	a.Equal("ocr_0000A3bQP9RLVfNUcdpLpw", order.ID)
	a.Equal("ord_0000A3bQ8FJIQoEfuC07n6", order.OrderID)
	a.Equal(false, order.LiveMode)
	a.Len(order.OrderChangeOffers, 1)
}

func TestGetOrderChangeRequest(t *testing.T) {
	defer gock.Off()
	gock.Observe(gock.DumpRequest)

	gock.New("https://api.duffel.com").
		Get("/air/order_change_requests/ocr_0000A3bQP9RLVfNUcdpLpw").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-get-order-change-request.json")

	a := assert.New(t)

	ctx := context.TODO()
	client := New("duffel_test_123")

	data, err := client.GetOrderChangeRequest(ctx, "ocr_0000A3bQP9RLVfNUcdpLpw")
	a.NoError(err)
	a.Equal("ocr_0000A3bQP9RLVfNUcdpLpw", data.ID)
	a.Equal("ord_0000A3bQ8FJIQoEfuC07n6", data.OrderID)
	a.Equal("35.50 GBP", data.OrderChangeOffers[0].NewTotalAmount().String())
	a.Equal("10.50 GBP", data.OrderChangeOffers[0].PenaltyTotalAmount().String())
	a.Equal("90.80 GBP", data.OrderChangeOffers[0].ChangeTotalAmount().String())
}

func TestCreatePendingOrderChange(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.duffel.com").
		Post("/air/order_changes").
		Reply(201).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/201-create-pending-order-change.json")

	a := assert.New(t)

	ctx := context.TODO()
	client := New("duffel_test_123")

	data, err := client.CreatePendingOrderChange(ctx, "oco_0000A3vUda8dKRtUSQPSXw")
	a.NoError(err)
	a.Equal("ocr_0000A3tQSmKyqOrcySrGbo", data.ID)
	a.Equal("ord_0000A3tQcCRZ9R8OY0QlxA", data.OrderID)
}
