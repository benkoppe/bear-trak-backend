package external

import "github.com/benkoppe/bear-trak-backend/go-server/utils/time_utils"

type Gym struct {
	LocationID   int    `json:"LocationId"`
	LocationName string `json:"LocationName"`

	LastUpdatedDateAndTime time_utils.ESTTime `json:"LastUpdatedDateAndTime"`

	// used to calculate percentage
	TotalCapacity int `json:"TotalCapacity"`
	LastCount     int `json:"LastCount"`
}

func (g Gym) GetPercentage() int {
	return (g.LastCount * 100) / g.TotalCapacity
}
