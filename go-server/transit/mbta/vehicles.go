package mbta

import (
	"fmt"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/transit/shared/gtfs_rt"
	"github.com/jamespfennell/gtfs"
)

func GetVehicles(caches Caches) ([]api.Vehicle, error) {
	staticGtfs, err := caches.staticCache.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to load static data: %v", err)
	}

	realtimeGtfs, err := caches.realtimeCache.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to load realtime data: %v", err)
	}

	vehicles, err := getVehicles(*staticGtfs, *realtimeGtfs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse vehicles: %v", err)
	}

	return vehicles, nil
}

func getVehicles(staticGtfs gtfs.Static, realtimeGtfs gtfs.Realtime) ([]api.Vehicle, error) {
	var vehicles []api.Vehicle

	for _, vehicle := range realtimeGtfs.Vehicles {
		vehicle := gtfs_rt.ConvertVehicle(vehicle, staticGtfs, realtimeGtfs)
		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}
