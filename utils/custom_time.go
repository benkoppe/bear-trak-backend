package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

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
	return time.Time(ut)
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
	return time.Time(mt)
}
