// Package cornell loads cornell transit content.
package cornell

import (
	"fmt"
	"strconv"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/transit/shared"
	"github.com/benkoppe/bear-trak-backend/go-server/transit/shared/availtec"
	shared_gtfs "github.com/benkoppe/bear-trak-backend/go-server/transit/shared/gtfs"
	"github.com/jamespfennell/gtfs"
)

type Caches struct {
	availtecCache availtec.Cache
	staticCache   shared_gtfs.Cache
}

func InitCaches(availtecURL string, staticGtfsURL string) Caches {
	return Caches{
		availtecCache: availtec.InitCache(availtecURL),
		staticCache:   shared_gtfs.InitCache(staticGtfsURL),
	}
}

func GetRoutes(caches Caches) ([]api.BusRoute, error) {
	staticGtfs, err := caches.staticCache.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to load static data: %v", err)
	}

	availtecRoutes, err := caches.availtecCache.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to load availtec routes: %v", err)
	}

	routes, err := getRoutes(availtecRoutes, *staticGtfs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse routes: %v", err)
	}

	vehicles, err := getVehiclesFromRoutes(availtecRoutes)
	if err != nil {
		return nil, fmt.Errorf("failed to load vehicles: %v", err)
	}

	routes = shared.AppendVehicles(routes, vehicles)

	return routes, nil
}

func getRoutes(availtecRoutes []availtec.Route, staticGtfs gtfs.Static) ([]api.BusRoute, error) {
	var routes []api.BusRoute

	for _, route := range availtecRoutes {
		var gtfsRoute *gtfs.Route
		for _, gtfsRouteOption := range staticGtfs.Routes {
			if gtfsRouteOption.Id == strconv.Itoa(route.RouteID) {
				gtfsRoute = &gtfsRouteOption
				break
			}
		}

		if gtfsRoute == nil {
			return nil, fmt.Errorf("failed to find GTFS route for route ID: %v", route.RouteID)
		}

		apiRoute := shared_gtfs.ConvertRoute(*gtfsRoute, staticGtfs)
		apiRoute.ID = api.NewIntBusRouteID(route.RouteID)

		routes = append(routes, apiRoute)
	}

	return routes, nil
}
