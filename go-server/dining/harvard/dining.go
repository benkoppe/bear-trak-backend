package harvard

import (
	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/dining/harvard/external"
)

type Cache = external.Caches

func InitCache(baseUrl, apiKey string) Cache {
	return external.InitCaches(baseUrl, apiKey)
}

func Get(
	externalCache Cache,
) ([]api.Eatery, error) {
	return []api.Eatery{}, nil
}
