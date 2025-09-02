// Package external loads external harvard library data.
package external

import (
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

type Cache = *utils.Cache[[]Library]

func InitCache(url string) Cache {
	return utils.NewCache(
		"harvardLibrariesExternal",
		1*time.Hour,
		func() ([]Library, error) {
			return fetchData(url)
		},
	)
}

func fetchData(url string) ([]Library, error) {
	libraries, err := utils.DoGetRequest[[]Library](url, nil)
	if libraries == nil {
		return []Library{}, err
	}
	return *libraries, err
}
