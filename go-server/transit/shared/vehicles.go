// Package shared includes all shared transit methods.
package shared

import "github.com/benkoppe/bear-trak-backend/go-server/api"

func AppendVehicles(routes []api.BusRoute, vehicles []api.Vehicle) []api.BusRoute {
	routeIdVehicles := make(map[interface{}]([]api.Vehicle))
	for _, vehicle := range vehicles {
		routeIdVehicles[vehicle.RouteId] = append(routeIdVehicles[vehicle.RouteId], vehicle)
	}

	for i := range routes {
		routes[i].Vehicles = routeIdVehicles[routes[i].ID]
	}

	return routes
}
