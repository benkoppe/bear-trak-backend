package utils

import (
	"fmt"
	"time"
)

// resillient datetime parser (good for both dates and times)
// will try multiple layouts in order before failing
func ParseDateTime(s string, layouts []string) (time.Time, error) {
	var err error
	for _, layout := range layouts {
		if t, err2 := time.Parse(layout, s); err2 == nil {
			return t, nil
		} else {
			err = err2
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse datetime %q: %v", s, err)
}

type TimeString string

func (t TimeString) parseTime() (struct {
	Hour   int
	Minute int
}, error,
) {
	emptyTime := struct {
		Hour   int
		Minute int
	}{Hour: 0, Minute: 0}
	time, err2 := ParseDateTime(string(t), []string{
		"3:04pm",
		"3:04 pm",
		"3 pm",
		"3pm",
		"15:04",
	})
	if err2 != nil {
		return emptyTime, fmt.Errorf("invalid time: %d", err2)
	}

	return struct {
		Hour   int
		Minute int
	}{Hour: time.Hour(), Minute: time.Minute()}, nil
}

func (t TimeString) ToDate(date time.Time) (time.Time, error) {
	parsed, err := t.parseTime()
	if err != nil {
		return time.Time{}, err
	}
	est := LoadEST()
	return time.Date(date.Year(), date.Month(), date.Day(), parsed.Hour, parsed.Minute, 0, 0, est), nil
}
