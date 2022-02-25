package duffel

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_MarshallingDate(t *testing.T) {
	a := assert.New(t)

	testDate := Date(time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC))
	type test struct {
		DepartureDate Date `json:"departure_date"`
	}

	testStruct := test{
		DepartureDate: testDate,
	}

	payload, err := json.Marshal(testStruct)
	a.NoError(err)
	a.Equal(`{"departure_date":"2016-01-01"}`, string(payload))

	var unmarshalled test
	err = json.Unmarshal(payload, &unmarshalled)
	a.NoError(err)
	a.Equal(testDate, unmarshalled.DepartureDate)
}

func TestDateTime(t *testing.T) {
	tz, _ := time.LoadLocation("Asia/Bangkok")
	tests := []struct {
		Input    string
		Expected time.Time
	}{
		{Input: "{\"date_time\": \"2022-02-22T12:01:00Z\"}", Expected: time.Date(2022, 2, 22, 12, 1, 0, 0, time.UTC)},
		{Input: "{\"date_time\": \"2022-02-22T12:01:00+07:00\"}", Expected: time.Date(2022, 2, 22, 12, 1, 0, 0, tz)},
		{Input: "{\"date_time\": \"2022-02-22T12:01:00\"}", Expected: time.Date(2022, 2, 22, 12, 1, 0, 0, time.UTC)},
	}

	type container struct {
		DateTime DateTime `json:"date_time"`
	}

	for _, test := range tests {
		d := new(container)
		err := json.Unmarshal([]byte(test.Input), d)
		if err != nil {
			t.Fatal(err)
		}

		actual := time.Time(d.DateTime)

		if actual.Unix() != test.Expected.Unix() {
			t.Errorf("%s: expected %s, got %s", test.Input, test.Expected, actual)
		}
	}
}
