package transit

import (
	"fmt"

	backend "github.com/benkoppe/bear-trak-backend/backend"
	availtec "github.com/benkoppe/bear-trak-backend/transit/external_availtec"
)

func GetVehicles(availtecUrl string) ([]backend.Vehicle, error) {
	availtecRoutes, err := availtec.FetchRoutes(availtecUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to load availtec routes: %v", err)
	}

	return getVehiclesFromRoutes(availtecRoutes)
}

func getVehiclesFromRoutes(availtecRoutes []availtec.Route) ([]backend.Vehicle, error) {
	var availtecVehicles []availtec.Vehicle

	for _, route := range availtecRoutes {
		availtecVehicles = append(availtecVehicles, route.Vehicles...)
	}

	return getVehicles(availtecVehicles)
}

func getVehicles(availtecVehicles []availtec.Vehicle) ([]backend.Vehicle, error) {
	var vehicles []backend.Vehicle

	for _, vehicle := range availtecVehicles {
		vehicles = append(vehicles, backend.Vehicle{
			ID:            vehicle.VehicleId,
			RouteId:       vehicle.RouteId,
			Direction:     vehicle.Direction,
			Heading:       vehicle.Heading,
			Latitude:      vehicle.Latitude,
			Longitude:     vehicle.Longitude,
			DisplayStatus: vehicle.OccupancyStatusReportLabel,
			Destination:   vehicle.Destination,
			LastStop:      vehicle.LastStop,
			LastUpdated:   vehicle.LastUpdated.ToTime(),
		})
	}

	return vehicles, nil
}
