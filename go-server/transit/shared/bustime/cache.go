package bustime

import (
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

type Caches struct {
	RoutesCache   RoutesCache
	VehiclesCache VehiclesCache
}

func InitCaches(baseUrl, apiKey string) Caches {
	routesCache := initRoutesCache(baseUrl, apiKey)
	vehiclesCache := initVehiclesCache(baseUrl, apiKey, routesCache)
	return Caches{
		RoutesCache:   routesCache,
		VehiclesCache: vehiclesCache,
	}
}

type RoutesCache = *utils.Cache[[]Route]

func initRoutesCache(baseUrl, apiKey string) RoutesCache {
	return utils.NewCache(
		"busTimeRoutesCache",
		24*time.Hour,
		func() ([]Route, error) {
			return fetchRoutes(baseUrl, apiKey)
		})
}

type VehiclesCache = *utils.Cache[[]Vehicle]

func initVehiclesCache(baseUrl, apiKey string, routesCache RoutesCache) VehiclesCache {
	return utils.NewCache(
		"busTimeVehiclesCache",
		20*time.Second,
		func() ([]Vehicle, error) {
			routes, err := routesCache.Get()
			if err != nil {
				return nil, err
			}
			return chunkFetchVehicles(baseUrl, apiKey, collectRouteIds(routes))
		})
}
