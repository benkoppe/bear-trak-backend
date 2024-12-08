package utils

import (
	"encoding/json"
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
