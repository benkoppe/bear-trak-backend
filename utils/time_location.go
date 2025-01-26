package utils

import (
	"fmt"
	"time"
)

func LoadEST() *time.Location {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(fmt.Sprintf("failed to load EST time zone: %v", err))
	}
	return loc
}
