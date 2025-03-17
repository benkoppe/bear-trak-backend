package umich

import (
	"fmt"
	"log"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/transit/shared"
	"github.com/benkoppe/bear-trak-backend/go-server/transit/shared/bustime"
	shared_gtfs "github.com/benkoppe/bear-trak-backend/go-server/transit/shared/gtfs"
	"github.com/jamespfennell/gtfs"
)

type Caches struct {
	bustimeCaches bustime.Caches
	staticCache   shared_gtfs.Cache
}

func InitCaches(bustimeUrl, bustimeApiKey, staticGtfsUrl string) Caches {
	return Caches{
		staticCache:   shared_gtfs.InitCache(staticGtfsUrl),
		bustimeCaches: bustime.InitCaches(bustimeUrl, bustimeApiKey),
	}
}

func GetRoutes(caches Caches) ([]api.BusRoute, error) {
	staticGtfs, err := caches.staticCache.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to load static data: %v", err)
	}

	bustimeRoutes, err := caches.bustimeCaches.RoutesCache.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to load bustime routes: %v", err)
	}

	routes, err := getRoutes(bustimeRoutes, *staticGtfs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse routes: %v", err)
	}

	vehicles, err := GetVehicles(caches)
	if err != nil {
		log.Printf("failed to load vehicles: %v", err)
	} else {
		routes = shared.AppendVehicles(routes, vehicles)
	}

	return routes, nil
}

func getRoutes(bustimeRoutes []bustime.Route, staticGtfs gtfs.Static) ([]api.BusRoute, error) {
	var routes []api.BusRoute

	for _, route := range bustimeRoutes {
		var gtfsRoute *gtfs.Route
		for _, gtfsRouteOption := range staticGtfs.Routes {
			if gtfsRouteOption.Id == route.Id {
				gtfsRoute = &gtfsRouteOption
				break
			}
		}

		if gtfsRoute == nil {
			return nil, fmt.Errorf("failed to find GTFS route for route ID: %v", route.Id)
		}

		routes = append(routes, shared_gtfs.ConvertRoute(*gtfsRoute, staticGtfs))
	}

	return routes, nil
}
