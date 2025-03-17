package external

import (
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

type Cache = *utils.Cache[*EateryAPIResponse]

func InitCache(url string) Cache {
	return utils.NewCache(
		"diningExternal",
		1*time.Minute,
		func() (*EateryAPIResponse, error) {
			return fetchData(url)
		})
}

func fetchData(url string) (*EateryAPIResponse, error) {
	return utils.DoGetRequest[EateryAPIResponse](url)
}
