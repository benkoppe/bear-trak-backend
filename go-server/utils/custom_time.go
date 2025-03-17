package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// lots of definitions for custom date types
// these are mostly useful for parsing JSON

// Defines a UnixTime
// this is encoded in the dining API as
// an integer number of seconds since 1970.

type UnixTime time.Time

func (ut *UnixTime) UnmarshalJSON(data []byte) error {
	var seconds int64
	if err := json.Unmarshal(data, &seconds); err != nil {
		return err
	}

	*ut = UnixTime(time.Unix(seconds, 0))
	return nil
}

func (ut UnixTime) ToTime() time.Time {
	est := LoadEST()
	return time.Time(ut).In(est)
}

// Defines a MicrosoftTime
// this is encoded in the transit API as
// a string of the format "/Date(1737414014000-0500)/"

type MicrosoftTime time.Time

func (mt *MicrosoftTime) UnmarshalJSON(data []byte) error {
	// b is something like: "/Date(1737412318000-0500)/" in quotes
	s := strings.Trim(string(data), `"`)
	// Replace any escaped `\/` with `/`
	s = strings.ReplaceAll(s, `\/`, `/`)
	if !strings.HasPrefix(s, "/Date(") || !strings.HasSuffix(s, ")/") {
		return errors.New("invalid date format: " + s)
	}

	// Strip off /Date( and )/
	s = s[6 : len(s)-2] // leaves "1737412318000-0500" or possibly just milliseconds

	// Extract the milliseconds part and any time-zone offset if present
	var msPart, offsetPart string
	// Check for + or - sign that indicates an offset
	plusIdx := strings.Index(s, "+")
	minusIdx := strings.Index(s, "-")

	if plusIdx > 0 {
		msPart, offsetPart = s[:plusIdx], s[plusIdx:]
	} else if minusIdx > 0 {
		msPart, offsetPart = s[:minusIdx], s[minusIdx:]
	} else {
		// No offset
		msPart = s
	}

	msVal, err := strconv.ParseInt(msPart, 10, 64)
	if err != nil {
		return err
	}

	// By default, parse as UTC
	loc := time.UTC

	// If there’s an offset, convert it to a fixed location
	if offsetPart != "" {
		if len(offsetPart) != 5 && len(offsetPart) != 6 {
			// e.g. "-0500" (5) or "-0500)" with some trailing?
			return fmt.Errorf("unexpected offset format: %q", offsetPart)
		}
		sign := offsetPart[0] // '+' or '-'
		hh, _ := strconv.Atoi(offsetPart[1:3])
		mm, _ := strconv.Atoi(offsetPart[3:5])

		offsetSeconds := hh*3600 + mm*60
		if sign == '-' {
			offsetSeconds = -offsetSeconds
		}
		loc = time.FixedZone(offsetPart, offsetSeconds)
	}

	*mt = MicrosoftTime(time.UnixMilli(msVal).In(loc))
	return nil
}

func (mt MicrosoftTime) ToTime() time.Time {
	est := LoadEST()
	return time.Time(mt).In(est)
}

// Defines an EST time
// this is encoded in the gyms capacities API as
// a typical time string, but missing the EST timezone,
// as in: "2025-01-25T23:39:21.53" for 01/25/2025 11:39 PM

type ESTTime time.Time

func (et *ESTTime) UnmarshalJSON(data []byte) error {
	// Remove the quotes from the JSON string
	s := string(data)
	if len(s) < 2 {
		return fmt.Errorf("invalid time format")
	}
	s = s[1 : len(s)-1] // Trim the surrounding quotes

	// Split the time string into main part and fractional seconds
	parts := strings.Split(s, ".")
	if len(parts) == 2 {
		fractional := parts[1]
		// Pad the fractional part to ensure three digits
		if len(fractional) < 3 {
			fractional = fractional + strings.Repeat("0", 3-len(fractional))
		} else if len(fractional) > 3 {
			fractional = fractional[:3]
		}
		s = parts[0] + "." + fractional
	} else {
		// If no fractional part, add ".000"
		s = s + ".000"
	}

	layout := "2006-01-02T15:04:05.000" // layout matches 2025-01-25T23:39:21.53
	est := LoadEST()

	t, err := time.ParseInLocation(layout, s, est)
	if err != nil {
		return err
	}

	*et = ESTTime(t)
	return nil
}

func (et ESTTime) ToTime() time.Time {
	est := LoadEST()
	return time.Time(et).In(est)
}
