package duffel

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestGetSeatmaps(t *testing.T) {
	defer gock.Off()
	a := assert.New(t)
	gock.New("https://api.duffel.com").
		Get("/air/seat_maps").
		MatchParam("offer_id", "off_00009htYpSCXrwaB9DnUm0").
		Reply(200).
		SetHeader("Ratelimit-Limit", "5").
		SetHeader("Ratelimit-Remaining", "5").
		SetHeader("Ratelimit-Reset", time.Now().Format(time.RFC1123)).
		SetHeader("Date", time.Now().Format(time.RFC1123)).
		File("fixtures/200-get-seatmap.json")

	ctx := context.TODO()

	client := New("duffel_test_123")
	iter := client.GetSeatmaps(ctx, "off_00009htYpSCXrwaB9DnUm0")

	iter.Next()
	data := iter.Current()
	err := iter.Err()

	a.NoError(err)
	a.NotNil(data)

	a.Equal("sea_00003hthlsHZ8W4LxXjkzo", data.ID)
	a.Equal("seg_00009htYpSCXrwaB9Dn456", data.SegmentID)
	a.Equal("sli_00009htYpSCXrwaB9Dn123", data.SliceID)
}
