package static

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed alerts.json
var alertBytes []byte

func GetAlerts() []Alert {
	var alerts []Alert

	err := json.Unmarshal(alertBytes, &alerts)
	if err != nil {
		fmt.Printf("error unmarshalling alerts: %v", err)
	}

	return alerts
}
