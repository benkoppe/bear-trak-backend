package gyms

import "github.com/benkoppe/bear-trak-backend/go-server/api"

func GetCapacityPoints(caches Caches) ([]api.GymCapacityData, error) {
	capacitiesData, err := caches.capacitiesCache.Get()

	return capacitiesData, err
}
