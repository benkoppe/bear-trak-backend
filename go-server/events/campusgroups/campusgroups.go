// Package campusgroups loads event data from campusgroups and converts to the api-expected format
package campusgroups

import (
	"fmt"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/events/campusgroups/geolocate"
	"github.com/benkoppe/bear-trak-backend/go-server/events/campusgroups/login"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

type Cache = *utils.Cache[[]api.Event]

func InitCache(baseURL string, loginParams login.LoginParams, locator *geolocate.GeoLocator) Cache {
	return utils.NewCache(
		"campusGroups",
		8*time.Hour,
		func() ([]api.Event, error) {
			return fetchAndConvert(baseURL, loginParams, locator)
		},
	)
}

func fetchAndConvert(baseURL string, loginParams login.LoginParams, locator *geolocate.GeoLocator) ([]api.Event, error) {
	events, err := fetchAllData(baseURL, loginParams, locator)
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	return convertAndSort(events)
}
