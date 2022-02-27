package iso8601

import (
	"testing"
	"time"
)

func TestDurationParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
	}{
		{input: "P2WT1M", expected: (time.Hour * 24 * 14) + time.Minute},
		{input: "P1WT1S", expected: (time.Hour * 24 * 7) + (1 * time.Second)},
		{input: "P1DT1H", expected: testDuration("25h")},
		{input: "PT12H58M", expected: testDuration("12h58m")},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			d, err := ParseDuration(test.input)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if d != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, d)
			}
		})
	}
}

func testDuration(str string) time.Duration {
	d, err := time.ParseDuration(str)
	if err != nil {
		panic(err)
	}
	return d
}
