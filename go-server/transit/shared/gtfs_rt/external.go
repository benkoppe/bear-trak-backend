// Package gtfs_rt handles gtfs-rt data.
package gtfs_rt

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
	"github.com/benkoppe/bear-trak-backend/go-server/utils/time_utils"
	"github.com/jamespfennell/gtfs"
)

type RealtimeUrls struct {
	Alerts           string
	VehiclePositions string
	TripUpdates      string
}

type Cache = *utils.Cache[*gtfs.Realtime]

func InitCache(urls RealtimeUrls) Cache {
	return utils.NewCache(
		"transitExternalRealtimeGtfs",
		time.Second*3,
		func() (*gtfs.Realtime, error) {
			return getRealtimeGtfs(urls)
		},
	)
}

func getRealtimeGtfs(urls RealtimeUrls) (*gtfs.Realtime, error) {
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

	estLocation := time_utils.LoadEST()
	realtimeData, err := gtfs.ParseRealtime(b, &gtfs.ParseRealtimeOptions{Timezone: estLocation})
	if err != nil {
		return nil, fmt.Errorf("error parsing realtime tcat data: %v", err)
	}

	return realtimeData, nil
}
