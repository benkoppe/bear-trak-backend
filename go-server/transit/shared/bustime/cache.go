package bustime

import (
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

type Caches struct {
	RoutesCache   RoutesCache
	VehiclesCache VehiclesCache
}

func InitCaches(baseURL, apiKey string) Caches {
	routesCache := initRoutesCache(baseURL, apiKey)
	vehiclesCache := initVehiclesCache(baseURL, apiKey, routesCache)
	return Caches{
		RoutesCache:   routesCache,
		VehiclesCache: vehiclesCache,
	}
}

type RoutesCache = *utils.Cache[[]Route]

func initRoutesCache(baseURL, apiKey string) RoutesCache {
	return utils.NewCache(
		"busTimeRoutesCache",
		24*time.Hour,
		func() ([]Route, error) {
			return fetchRoutes(baseURL, apiKey)
		})
}

type VehiclesCache = *utils.Cache[[]Vehicle]

func initVehiclesCache(baseURL, apiKey string, routesCache RoutesCache) VehiclesCache {
	return utils.NewCache(
		"busTimeVehiclesCache",
		20*time.Second,
		func() ([]Vehicle, error) {
			routes, err := routesCache.Get()
			if err != nil {
				return nil, err
			}
			return chunkFetchVehicles(baseURL, apiKey, collectRouteIds(routes))
		})
}
