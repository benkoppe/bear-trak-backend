// Package external loads external cornell library data.
package external

import (
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

type Cache = *utils.Cache[[]Library]

func InitCache(url string) Cache {
	return utils.NewCache(
		"libraryExternal",
		1*time.Hour,
		func() ([]Library, error) {
			return FetchData(url)
		})
}

func FetchData(url string) ([]Library, error) {
	library, err := utils.DoGetRequest[librariesResponse](url, nil)
	if library == nil {
		return []Library{}, err
	}
	return library.Locations, err
}
