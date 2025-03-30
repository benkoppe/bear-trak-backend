package umich

import (
	"fmt"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/transit/shared/bustime"
)

func GetVehicles(caches Caches) ([]api.Vehicle, error) {
	bustimeVehicles, err := caches.bustimeCaches.VehiclesCache.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to load bustime vehicles: %w", err)
	}

	var vehicles []api.Vehicle
	for _, vehicle := range bustimeVehicles {
		vehicles = append(vehicles, convertVehicle(vehicle))
	}

	return vehicles, nil
}

func convertVehicle(vehicle bustime.Vehicle) api.Vehicle {
	return api.Vehicle{
		ID:            api.NewStringVehicleID(vehicle.Id),
		RouteId:       api.NewStringVehicleRouteId(vehicle.RouteId),
		Direction:     "",
		Heading:       int(vehicle.Heading),
		Speed:         float64(vehicle.Speed),
		Latitude:      float64(vehicle.Latitude),
		Longitude:     float64(vehicle.Longitude),
		DisplayStatus: vehicle.OccupancyStatus,
		Destination:   vehicle.Destination,
		LastUpdated:   vehicle.LastUpdated.ToTime(),
		LastStop:      api.NilString{Null: true},
	}
}
