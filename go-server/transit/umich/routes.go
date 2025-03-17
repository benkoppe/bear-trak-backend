package umich

import (
	"fmt"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/transit/external_bustime"
	"github.com/benkoppe/bear-trak-backend/go-server/transit/external_gtfs"
)

type Caches struct {
	bustimeCaches external_bustime.Caches
	staticCache   external_gtfs.Cache
}

func InitCaches(bustimeUrl, staticGtfsUrl string) Caches {
	return Caches{
		staticCache:   external_gtfs.InitCache(staticGtfsUrl),
		bustimeCaches: external_bustime.InitCaches(bustimeUrl),
	}
}

func GetRoutes(caches Caches) ([]api.BusRoute, error) {
	staticGtfs, err := caches.staticCache.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to load static data: %v", err)
	}

	_, err = caches.bustimeCaches.RoutesCache.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to load bustime routes: %v", err)
	}

	for _, route := range staticGtfs.Routes {
		fmt.Printf("route: %v \n", route.Id)
	}

	return nil, fmt.Errorf("not implemented yet")
}
