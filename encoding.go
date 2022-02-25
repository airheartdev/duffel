package duffel

import (
	"fmt"
	"strconv"
	"time"

	"github.com/airheartdev/duffel/iso8601"
	"github.com/pkg/errors"
)

type (
	Date     time.Time
	DateTime time.Time
	Duration time.Duration
	Distance float64
)

func (t Date) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02"))
	return []byte(stamp), nil
}

// UnmarshalJSON implements the json.Unmarshaler from date string to time.Time
func (t *Date) UnmarshalJSON(b []byte) error {
	str := string(b)
	if str == "null" {
		return nil
	}

	str, err := strconv.Unquote(str)
	if err != nil {
		return err
	}

	stamp, err := time.Parse("2006-01-02", str)
	if err != nil {
		return err
	}
	*t = Date(stamp)
	return nil
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(time.RFC3339Nano))
	return []byte(stamp), nil
}

var timeFormats = []string{
	"2006-01-02T15:04:05",
	"2006-01-02",
	time.RFC3339,
	time.RFC3339Nano,
}

// UnmarshalJSON implements the json.Unmarshaler from date string to time.Time
func (t *DateTime) UnmarshalJSON(b []byte) error {
	str := string(b)
	if str == "null" {
		return nil
	}

	str, err := strconv.Unquote(str)
	if err != nil {
		return err
	}

	stamp := time.Time{}
	for _, format := range timeFormats {
		stamp, err = time.Parse(format, str)
		if err != nil {
			continue
		}
		break
	}
	if stamp.IsZero() {
		return fmt.Errorf("failed to parse timestamp: '%s'", str)
	}

	*t = DateTime(stamp)
	return nil
}

func (t Duration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", iso8601.FormatDuration(time.Duration(t)))), nil
}

// UnmarshalJSON implements the json.Unmarshaler from date string to time.Time
func (t *Duration) UnmarshalJSON(b []byte) error {
	d, err := iso8601.ParseDuration(string(b))
	if err != nil {
		return err
	}

	*t = Duration(d)
	return nil
}

func (t Distance) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(fmt.Sprintf("%f", t))), nil
}

// UnmarshalJSON implements the json.Unmarshaler from date string to time.Time
func (t *Distance) UnmarshalJSON(b []byte) error {
	str := string(b)
	if str == "null" {
		return nil
	}

	f, err := strconv.Unquote(str)
	if err != nil {
		return errors.Wrap(err, "unquote distance")
	}

	d, err := strconv.ParseFloat(f, 16)
	if err != nil {
		return err
	}

	*t = Distance(d)

	return nil
}
