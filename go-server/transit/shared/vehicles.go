// Package shared includes all shared transit methods.
package shared

import "github.com/benkoppe/bear-trak-backend/go-server/api"

func AppendVehicles(routes []api.BusRoute, vehicles []api.Vehicle) []api.BusRoute {
	routeIDVehicles := make(map[any]([]api.Vehicle))
	for _, vehicle := range vehicles {
		routeIDVehicles[vehicle.RouteId] = append(routeIDVehicles[vehicle.RouteId], vehicle)
	}

	for i := range routes {
		routes[i].Vehicles = routeIDVehicles[routes[i].ID]
	}

	return routes
}
