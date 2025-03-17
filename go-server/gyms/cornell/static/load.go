package static

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed gyms.json
var gymBytes []byte

func GetGyms() []Gym {
	var gyms []Gym

	err := json.Unmarshal(gymBytes, &gyms)
	if err != nil {
		fmt.Printf("error unmarshalling static gyms: %v", err)
	}

	return gyms
}
