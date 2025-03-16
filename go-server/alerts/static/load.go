package static

import (
	"encoding/json"
	"fmt"
)

func GetAlerts(data []byte) []Alert {
	var alerts []Alert

	err := json.Unmarshal(data, &alerts)
	if err != nil {
		fmt.Printf("error unmarshalling alerts: %v", err)
	}

	return alerts
}
