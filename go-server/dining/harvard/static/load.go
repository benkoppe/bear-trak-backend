// Package static loads static harvard dining content.
package static

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed eateries.json
var eateryBytes []byte

func GetEateries() []Eatery {
	var eateries []Eatery

	err := json.Unmarshal(eateryBytes, &eateries)
	if err != nil {
		fmt.Printf("error unmarshalling static eateries: %v\n", err)
	}

	return eateries
}
