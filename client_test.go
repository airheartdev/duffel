// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package duffel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestClientError(t *testing.T) {
	ctx := context.TODO()
	a := assert.New(t)
	gock.New("https://api.duffel.com/air/offer_requests").
		MatchParam("return_offers", "false").
		Reply(400).
		File("fixtures/400-bad-request.json")
	defer gock.Off()

	client := New("duffel_test_123")
	data, err := client.CreateOfferRequest(ctx, OfferRequestInput{
		ReturnOffers: false,
	})
	a.Error(err)
	a.Nil(data)

	a.Equal("duffel: The airline responded with an unexpected error, please contact support", err.Error())

	derr := err.(*DuffelError)
	a.True(derr.IsType(AirlineError))
	a.True(derr.IsCode(AirlineUnknown))
	a.True(IsErrorType(err, AirlineError))
	a.True(IsErrorCode(err, AirlineUnknown))
	a.True(ErrIsRetryable(err))

	reqId, ok := RequestIDFromError(err)
	a.True(ok)
	a.Equal("FZW0H3HdJwKk5HMAAKxB", reqId)
}
