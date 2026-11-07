package harvard

import (
	"slices"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
)

// longwoodRouteCodes maps the Longwood routes Harvard shows (the M2s) to their codes.
var longwoodRouteCodes = map[string]string{
	"8526": "M2",
	"8527": "M2C",
}

func filterLongwoodRoutes(routes []api.BusRoute) []api.BusRoute {
	var m2Routes []api.BusRoute
	for _, route := range routes {
		if code, ok := longwoodRouteCodes[route.ID.String]; ok {
			route.Code = code
			m2Routes = append(m2Routes, route)
		}
	}
	return m2Routes
}

func filterLongwoodVehicles(vehicles []api.Vehicle) []api.Vehicle {
	return slices.DeleteFunc(vehicles, func(v api.Vehicle) bool {
		_, ok := longwoodRouteCodes[v.RouteId.String]
		return !ok
	})
}
