package time_utils

import (
	"fmt"
	"time"
)

// resillient datetime parser (good for both dates and times)
// will try multiple layouts in order before failing
// requires a list of templates
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

// wrapper around `ParseDateTime` with a list of common formats
func ParseCommonDateTime(s string) (time.Time, error) {
	commonLayouts := []string{
		time.RFC3339,                  // "2006-01-02T15:04:05Z07:00"
		time.RFC1123Z,                 // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC1123,                  // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC822,                   // "02 Jan 06 15:04 MST"
		time.RFC822Z,                  // "02 Jan 06 15:04 -0700"
		"2006-01-02 15:04:05",         // Common datetime format
		"2006-01-02",                  // ISO date format
		"01/02/2006",                  // US date format (MM/DD/YYYY)
		"2 Jan 2006",                  // Day Month Year without leading zeros
		"January 2, 2006",             // Long month name format
		"Mon, Jan 2, 2006",            // Abbreviated weekday with date
		"Mon Jan 2 15:04:05 MST 2006", // Full date with weekday
	}

	t, err := ParseDateTime(s, commonLayouts)
	if err != nil {
		return time.Time{}, err
	}

	// Check if time components are zero, if so set to midnight EST
	if t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 {
		est := LoadEST()
		// Create new time at midnight EST while preserving the date
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, est)
	}

	return t, nil
}

// parses common datetimes, including formats without a year
// sets the year to the current year if it isn't provided
func ParseCommonDateTimeYearOptional(s string) (time.Time, error) {
	// Layouts that include a year.
	layouts := []string{
		time.RFC3339,
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822,
		time.RFC822Z,
		"2006-01-02 15:04:05",
		"2006-01-02",
		"01/02/2006",
		"2 Jan 2006",
		"January 2, 2006",
		"Mon, Jan 2, 2006",
		"Mon Jan 2 15:04:05 MST 2006",
		"Jan 2",
		"January 2",
		"Mon, Jan 2",
		"Mon Jan 2",
		"Jan 2 Monday",
	}

	// Next, try layouts without a year.
	t, err := ParseDateTime(s, layouts)
	if err != nil {
		return time.Time{}, err
	}

	// If the parsed time has a zero year, update it with the current year.
	if t.Year() == 0 {
		now := time.Now()
		t = time.Date(now.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	}

	// Check if time components are zero, if so set to midnight EST
	if t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 {
		est := LoadEST()
		// Create new time at midnight EST while preserving the date
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, est)
	}

	return t, nil
}
