// Package availtec loads transit data from availtec.
package availtec

import (
	"fmt"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

type Cache = *utils.Cache[[]Route]

func InitCache(baseURL string) Cache {
	return utils.NewCache(
		"transitExternalAvailtec",
		time.Second*3,
		func() ([]Route, error) {
			return fetchRoutes(baseURL)
		})
}

func fetchRoutes(baseURL string) ([]Route, error) {
	fullURL, err := utils.ExtendURL(baseURL, "Routes/GetVisibleRoutes")
	if fullURL == nil {
		return []Route{}, fmt.Errorf("failed to extend url: %w", err)
	}

	routes, err := utils.DoGetRequest[[]Route](*fullURL, nil)
	if routes == nil {
		return []Route{}, err
	}

	return *routes, nil
}
