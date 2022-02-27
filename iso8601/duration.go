package iso8601

import (
	"fmt"
	"time"

	"github.com/rickb777/date/period"
)

func ParseDuration(value string) (time.Duration, error) {
	p, err := period.Parse(value, true)
	if err != nil {
		return 0, err
	}

	d, _ := p.Duration()
	return d, nil
}

func FormatDuration(duration time.Duration) string {
	// we're not doing negative durations
	if duration.Seconds() <= 0 {
		return "PT0S"
	}

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) - (hours * 60)
	seconds := int(duration.Seconds()) - (hours*3600 + minutes*60)

	// we're not doing Y,M,W
	s := "PT"
	if hours > 0 {
		s = fmt.Sprintf("%s%dH", s, hours)
	}
	if minutes > 0 {
		s = fmt.Sprintf("%s%dM", s, minutes)
	}
	if seconds > 0 {
		s = fmt.Sprintf("%s%dS", s, seconds)
	}

	return s
}
