package time_utils

import (
	"log"
	"regexp"
	"strings"
)

type Hours struct {
	Open  TimeString `json:"open"`
	Close TimeString `json:"close"`
}

type HoursParserOptions struct {
	RangeSeparators     []string // Separators between time ranges
	OpenCloseSeparators []string // Separators between open/close times
	ClosedKeywords      []string // Words indicating closed status
	Open24HoursKeywords []string // Words indicating 24-hour operation
}

func DefaultHoursParserOptions() *HoursParserOptions {
	return &HoursParserOptions{
		RangeSeparators:     []string{"/", ","},
		OpenCloseSeparators: []string{"-", "to", "–", "—"}, // includes en/em dashes
		ClosedKeywords:      []string{"Closed"},
		Open24HoursKeywords: []string{"24 hours", "24/7", "All day", "Open 24 hours"},
	}
}

// ParseHours splits a cell value into one or more Hours.
// It handles values like "6am - 9pm" or "7am - 8:30am / 10am - 10:45pm".
func ParseHours(raw string, opts *HoursParserOptions) []Hours {
	if opts == nil {
		opts = DefaultHoursParserOptions()
	}

	var hrs []Hours
	rawLower := strings.ToLower(raw)

	for _, keyword := range opts.ClosedKeywords {
		if strings.Contains(rawLower, strings.ToLower(keyword)) {
			return hrs
		}
	}

	for _, keyword := range opts.Open24HoursKeywords {
		if strings.Contains(rawLower, strings.ToLower(keyword)) {
			return []Hours{{
				Open:  "00:00",
				Close: "23:59",
			}}
		}
	}

	// Build regex for range separators
	rangeSepPattern := strings.Join(opts.RangeSeparators, "|")
	rangeRegex := regexp.MustCompile("\\s*(" + rangeSepPattern + ")\\s*")

	ranges := rangeRegex.Split(raw, -1)
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}

		var openTime, closeTime string
		found := false
		for _, sep := range opts.OpenCloseSeparators {
			parts := strings.Split(r, sep)
			if len(parts) == 2 {
				openTime = strings.TrimSpace(parts[0])
				closeTime = strings.TrimSpace(parts[1])
				found = true
				break
			}
		}

		if !found {
			log.Printf("Couldn't format a general hours range: %s", r)
			continue
		}

		hrs = append(hrs, Hours{
			Open:  TimeString(openTime),
			Close: TimeString(closeTime),
		})
	}

	return hrs
}
