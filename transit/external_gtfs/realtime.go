package external_gtfs

import (
	"fmt"
	"io"
	"net/http"

	"github.com/benkoppe/bear-trak-backend/utils"
	"github.com/jamespfennell/gtfs"
)

type RealtimeUrls struct {
	Alerts           string
	VehiclePositions string
	TripUpdates      string
}

func GetRealtimeGtfs(urls RealtimeUrls) (*gtfs.Realtime, error) {
	alerts, err := getSingleRealtimeGtfs(urls.Alerts)
	if err != nil {
		return nil, fmt.Errorf("error loading alerts: %v", err)
	}

	vehicles, err := getSingleRealtimeGtfs(urls.VehiclePositions)
	if err != nil {
		return nil, fmt.Errorf("error loading vehicles: %v", err)
	}

	trips, err := getSingleRealtimeGtfs(urls.TripUpdates)
	if err != nil {
		return nil, fmt.Errorf("error loading trips: %v", err)
	}

	vehicles.Trips = trips.Trips
	vehicles.Alerts = alerts.Alerts

	return vehicles, nil
}

func getSingleRealtimeGtfs(url string) (*gtfs.Realtime, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error loading realtime tcat data: %v", err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading realtime tcat data: %v", err)
	}

	est := utils.LoadEST()

	realtimeData, err := gtfs.ParseRealtime(b, &gtfs.ParseRealtimeOptions{Timezone: est})
	if err != nil {
		return nil, fmt.Errorf("error parsing realtime tcat data: %v", err)
	}

	return realtimeData, nil
}
