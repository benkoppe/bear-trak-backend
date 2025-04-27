package mbta

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed ignore-id-list.json
var ignoreBytes []byte

func getIgnoreIds() []string {
	var ignoreIds []string

	err := json.Unmarshal(ignoreBytes, &ignoreIds)
	if err != nil {
		fmt.Printf("error unmarshalling mbta ignore IDs: %v", err)
	}

	return ignoreIds
}
