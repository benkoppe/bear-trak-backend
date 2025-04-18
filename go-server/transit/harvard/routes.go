package harvard

import (
	"fmt"
	"log"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/transit/shared"
	shared_gtfs "github.com/benkoppe/bear-trak-backend/go-server/transit/shared/gtfs"
	"github.com/benkoppe/bear-trak-backend/go-server/transit/shared/gtfs_rt"
	"github.com/benkoppe/bear-trak-backend/go-server/transit/shared/pasio"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
	"github.com/jamespfennell/gtfs"
)

type Caches struct {
	pasioCache    pasio.Cache
	staticCache   shared_gtfs.Cache
	realtimeCache gtfs_rt.Cache
}

func InitCaches(pasioBaseUrl, pasioSystemId, staticGtfsUrl, realtimeGtfsBaseUrl string) Caches {
	alerts, err := utils.ExtendUrl(realtimeGtfsBaseUrl, "serviceAlerts")
	if err != nil {
		log.Fatalf("failed to extend realtime GTFS alerts URL: %v", err)
	}
	tripUpdates, err := utils.ExtendUrl(realtimeGtfsBaseUrl, "tripUpdates")
	if err != nil {
		log.Fatalf("failed to extend realtime GTFS tripupdates URL: %v", err)
	}
	vehicles, err := utils.ExtendUrl(realtimeGtfsBaseUrl, "vehiclePositions")
	if err != nil {
		log.Fatalf("failed to extend realtime GTFS vehicle positions URL: %v", err)
	}

	harvardGtfsRealtime := gtfs_rt.RealtimeUrls{
		Alerts:           *alerts,
		TripUpdates:      *tripUpdates,
		VehiclePositions: *vehicles,
	}

	return Caches{
		pasioCache:    pasio.InitCache(pasioBaseUrl, pasioSystemId),
		staticCache:   shared_gtfs.InitCache(staticGtfsUrl),
		realtimeCache: gtfs_rt.InitCache(harvardGtfsRealtime),
	}
}

func GetRoutes(caches Caches) ([]api.BusRoute, error) {
	staticGtfs, err := caches.staticCache.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to load static data: %v", err)
	}

	pasioRoutes, err := caches.pasioCache.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to load pasio data: %v", err)
	}

	routes, err := getRoutes(pasioRoutes, *staticGtfs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse routes: %v", err)
	}

	realtimeGtfs, err := caches.realtimeCache.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to load realtime data: %v", err)
	}

	vehicles, err := getVehicles(*staticGtfs, *realtimeGtfs)
	if err != nil {
		return nil, fmt.Errorf("failed to load vehicles: %v", err)
	}

	routes = shared.AppendVehicles(routes, vehicles)

	return routes, nil
}

func getRoutes(pasioRoutes []pasio.Route, staticGtfs gtfs.Static) ([]api.BusRoute, error) {
	var routes []api.BusRoute

	for _, route := range pasioRoutes {
		var gtfsRoute *gtfs.Route
		for _, gtfsRouteOption := range staticGtfs.Routes {
			if gtfsRouteOption.Id == route.GroupId {
				gtfsRoute = &gtfsRouteOption
				break
			}
		}

		if gtfsRoute == nil {
			return nil, fmt.Errorf("failed to find GTFS route for route ID: %s", route.GroupId)
		}

		apiRoute := shared_gtfs.ConvertRoute(*gtfsRoute, staticGtfs)
		apiRoute.Name = route.Name

		routes = append(routes, apiRoute)
	}

	return routes, nil
}
