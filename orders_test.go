package duffel

import (
	"testing"

	"github.com/segmentio/encoding/json"

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
