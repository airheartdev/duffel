package duffel

import "time"

func (f *Flight) DepartingAt() (time.Time, error) {
	loc, err := time.LoadLocation(f.Origin.TimeZone)
	if err != nil {
		return time.Time{}, err
	}

	t, err := time.ParseInLocation("2006-01-02T15:04:05", f.RawDepartingAt, loc)
	if err != nil {
		return time.Time{}, err
	}

	return t, err
}

func (f *Flight) ArrivingAt() (time.Time, error) {
	loc, err := time.LoadLocation(f.Destination.TimeZone)
	if err != nil {
		return time.Time{}, err
	}

	t, err := time.ParseInLocation("2006-01-02T15:04:05", f.RawArrivingAt, loc)
	if err != nil {
		return time.Time{}, err
	}

	return t, err
}
