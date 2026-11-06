package mbta

import (
	_ "embed"
	"encoding/json"
	"log"
)

//go:embed allowlist.json
var allowlistBytes []byte

// allowlistIDs are the only MBTA routes the server serves.
var allowlistIDs = loadAllowlistIDs()

func loadAllowlistIDs() []string {
	var ids []string

	err := json.Unmarshal(allowlistBytes, &ids)
	if err != nil {
		log.Fatalf("error unmarshalling mbta allowlist IDs: %v", err)
	}

	return ids
}

func isAllowlisted(routeID string) bool {
	for _, id := range allowlistIDs {
		if id == routeID {
			return true
		}
	}
	return false
}
