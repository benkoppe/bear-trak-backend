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
		fmt.Printf("error unmarshalling static eateries: %v", err)
	}

	return eateries
}
