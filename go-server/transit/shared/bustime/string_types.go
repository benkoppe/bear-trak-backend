package bustime

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils/time_utils"
)

// For some reason, some floats/ints are encoded as strings
type Float64String float64

func (f *Float64String) UnmarshalJSON(data []byte) error {
	// If the data is a JSON string (e.g. "1.23")
	if len(data) > 0 && data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		parsed, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		*f = Float64String(parsed)
		return nil
	}
	// Otherwise, assume it's a JSON number
	var num float64
	if err := json.Unmarshal(data, &num); err != nil {
		return err
	}
	*f = Float64String(num)
	return nil
}

type IntString int

func (i *IntString) UnmarshalJSON(data []byte) error {
	// Check if the JSON data is a quoted string (e.g., "123")
	if len(data) > 0 && data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		parsed, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*i = IntString(parsed)
		return nil
	}

	// Otherwise, assume it's a JSON number
	var num int
	if err := json.Unmarshal(data, &num); err != nil {
		return err
	}
	*i = IntString(num)
	return nil
}

// Defines a TransitTime
// this is encoded in the transit API as
// a string of the format "20250316 16:16" (YYYYMMDD HH:MM)
type TransitTime time.Time

func (tt *TransitTime) UnmarshalJSON(data []byte) error {
	// Remove the quotes from the JSON string
	s := string(data)
	if len(s) < 2 {
		return fmt.Errorf("invalid transit time format")
	}
	s = s[1 : len(s)-1]

	// Validate the format
	if len(s) != 14 || s[8] != ' ' {
		return fmt.Errorf("invalid transit time format: %s", s)
	}

	// Parse the date portion (YYYYMMDD)
	year, err := strconv.Atoi(s[0:4])
	if err != nil {
		return fmt.Errorf("invalid year in transit time: %s", s)
	}

	month, err := strconv.Atoi(s[4:6])
	if err != nil || month < 1 || month > 12 {
		return fmt.Errorf("invalid month in transit time: %s", s)
	}

	day, err := strconv.Atoi(s[6:8])
	if err != nil || day < 1 || day > 31 {
		return fmt.Errorf("invalid day in transit time: %s", s)
	}

	// Parse the time portion (HH:MM)
	if s[11] != ':' {
		return fmt.Errorf("invalid time format in transit time: %s", s)
	}

	hour, err := strconv.Atoi(s[9:11])
	if err != nil || hour < 0 || hour > 23 {
		return fmt.Errorf("invalid hour in transit time: %s", s)
	}

	minute, err := strconv.Atoi(s[12:14])
	if err != nil || minute < 0 || minute > 59 {
		return fmt.Errorf("invalid minute in transit time: %s", s)
	}

	// Create the time in local timezone
	est := time_utils.LoadEST()
	t := time.Date(year, time.Month(month), day, hour, minute, 0, 0, est)
	*tt = TransitTime(t)
	return nil
}

func (tt TransitTime) ToTime() time.Time {
	est := time_utils.LoadEST()
	return time.Time(tt).In(est)
}
