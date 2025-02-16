package external

import (
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

func FetchData(url string) (*EateryAPIResponse, error) {
	return utils.DoGetRequest[EateryAPIResponse](url)
}
