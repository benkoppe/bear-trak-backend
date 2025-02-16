package external

import (
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

func FetchData(url string) ([]Gym, error) {
	gyms, err := utils.DoGetRequest[[]Gym](url)
	if gyms == nil {
		return []Gym{}, err
	}
	return *gyms, err
}
