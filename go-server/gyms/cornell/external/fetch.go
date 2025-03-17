package external

import (
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

type Cache = *utils.Cache[[]Gym]

func InitCache(url string) Cache {
	return utils.NewCache(
		"gymExternal",
		1*time.Minute,
		func() ([]Gym, error) {
			return FetchData(url)
		})
}

func FetchData(url string) ([]Gym, error) {
	gyms, err := utils.DoGetRequest[[]Gym](url)
	if gyms == nil {
		return []Gym{}, err
	}
	return *gyms, err
}
