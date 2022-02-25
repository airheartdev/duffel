package duffel

import (
	"net/http"
	"strconv"
	"time"
)

type (
	RateLimit struct {
		Limit     int
		Remaining int
		ResetAt   time.Time
		Period    time.Duration
	}
)

func parseRateLimit(resp *http.Response) (*RateLimit, error) {
	rl := &RateLimit{}
	var err error

	rl.Limit, err = strconv.Atoi(resp.Header.Get("Ratelimit-Limit"))
	if err != nil {
		return nil, err
	}

	rl.Remaining, err = strconv.Atoi(resp.Header.Get("Ratelimit-Remaining"))
	if err != nil {
		return nil, err
	}

	resetAt, err := time.Parse(time.RFC1123, resp.Header.Get("Ratelimit-Reset"))
	if err != nil {
		return nil, err
	}

	date, err := time.Parse(time.RFC1123, resp.Header.Get("Date"))
	if err != nil {
		return nil, err
	}

	rl.ResetAt = resetAt
	rl.Period = resetAt.Sub(date)

	return rl, nil
}
