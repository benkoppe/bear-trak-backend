package external

import (
	"log"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

type Cache = *utils.Cache[*EateryAPIResponse]

func InitCache(url string) Cache {
	return utils.NewCache(
		1*time.Minute,
		func() (*EateryAPIResponse, error) {
			return fetchData(url)
		})
}

func fetchData(url string) (*EateryAPIResponse, error) {
	log.Println("Fetching data from URL:", url)
	return utils.DoGetRequest[EateryAPIResponse](url)
}
