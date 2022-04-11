package duffel

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/segmentio/encoding/json"
	"github.com/thetreep/duffel/iso8601"
)

type (
	Date     time.Time
	DateTime time.Time
	Duration time.Duration
	Distance float64
)

const DateFormat = "2006-01-02"

func (t Date) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(DateFormat))
	return []byte(stamp), nil
}

func (t Date) String() string {
	return time.Time(t).Format(DateFormat)
}

// UnmarshalJSON implements the json.Unmarshaler from date string to time.Time
func (t *Date) UnmarshalJSON(b []byte) error {
	str, err := parseJSONBytesToString(b)
	if err != nil {
		if err == ErrNullValue {
			return nil
		}

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

func (t DateTime) String() string {
	return time.Time(t).Format(time.RFC3339)
}

var timeFormats = []string{
	"2006-01-02T15:04:05",
	"2006-01-02",
	time.RFC3339,
	time.RFC3339Nano,
}

// UnmarshalJSON implements the json.Unmarshaler from date string to time.Time
func (t *DateTime) UnmarshalJSON(b []byte) error {
	str, err := parseJSONBytesToString(b)
	if err != nil {
		if err == ErrNullValue {
			return nil
		}
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

func (t Duration) String() string {
	return iso8601.FormatDuration(time.Duration(t))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (d *Duration) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid duration value: %v", v)
	}

	dur, err := iso8601.ParseDuration(str)
	if err != nil {
		return err
	}

	*d = Duration(dur)

	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (d Duration) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(d.String()))
}

// UnmarshalJSON implements the json.Unmarshaler from date string to time.Time
func (t *Duration) UnmarshalJSON(b []byte) error {
	f, err := parseJSONBytesToString(b)
	if err != nil {
		if err == ErrNullValue {
			return nil
		}
		return err
	}

	d, err := iso8601.ParseDuration(f)
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
	f, err := parseJSONBytesToString(b)
	if err != nil {
		if err == ErrNullValue {
			return nil
		}
		return err
	}

	d, err := strconv.ParseFloat(f, 16)
	if err != nil {
		return err
	}

	*t = Distance(d)

	return nil
}

var ErrNullValue = fmt.Errorf("null value")

func parseJSONBytesToString(b []byte) (string, error) {
	b = json.Unescape(b)
	if len(b) == 0 {
		return "", ErrNullValue
	}

	return string(b), nil
}
