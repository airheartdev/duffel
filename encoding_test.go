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
