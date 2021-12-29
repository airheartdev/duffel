package duffel

import (
	"fmt"
	"time"
)

type Date time.Time

func (t Date) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02"))
	return []byte(stamp), nil
}

// UnmarshalJSON implements the json.Unmarshaler from date string to time.Time
func (t *Date) UnmarshalJSON(b []byte) error {
	stamp, err := time.Parse("\"2006-01-02\"", string(b))
	if err != nil {
		return err
	}
	*t = Date(stamp)
	return nil
}
