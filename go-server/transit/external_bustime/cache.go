package external_bustime

import (
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

type Caches struct {
	RoutesCache   RoutesCache
	VehiclesCache VehiclesCache
}

func InitCaches(baseUrl string) Caches {
	routesCache := initRoutesCache(baseUrl)
	vehiclesCache := initVehiclesCache(baseUrl, routesCache)
	return Caches{
		RoutesCache:   routesCache,
		VehiclesCache: vehiclesCache,
	}
}

type RoutesCache = *utils.Cache[[]Route]

func initRoutesCache(baseUrl string) RoutesCache {
	return utils.NewCache(
		"busTimeRoutesCache",
		24*time.Hour,
		func() ([]Route, error) {
			return fetchRoutes(baseUrl)
		})
}

type VehiclesCache = *utils.Cache[[]Vehicle]

func initVehiclesCache(baseUrl string, routesCache RoutesCache) VehiclesCache {
	return utils.NewCache(
		"busTimeVehiclesCache",
		10*time.Second,
		func() ([]Vehicle, error) {
			routes, err := routesCache.Get()
			if err != nil {
				return nil, err
			}
			return chunkFetchVehicles(baseUrl, collectRouteIds(routes))
		})
}
