package external_availtec

import (
	"fmt"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

type Cache = *utils.Cache[[]Route]

func InitCache(baseUrl string) Cache {
	return utils.NewCache(
		"transitExternalAvailtec",
		time.Second*3,
		func() ([]Route, error) {
			return fetchRoutes(baseUrl)
		})
}

func fetchRoutes(baseUrl string) ([]Route, error) {
	fullUrl, err := utils.ExtendUrl(baseUrl, "Routes/GetVisibleRoutes")
	if fullUrl == nil {
		return []Route{}, fmt.Errorf("failed to extend url: %w", err)
	}

	routes, err := utils.DoGetRequest[[]Route](*fullUrl)
	if routes == nil {
		return []Route{}, err
	}

	return *routes, nil
}
