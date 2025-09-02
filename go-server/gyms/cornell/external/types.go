package external

import "github.com/benkoppe/bear-trak-backend/go-server/utils/timeutils"

type Gym struct {
	LocationID   int    `json:"LocationId"`
	LocationName string `json:"LocationName"`

	LastUpdatedDateAndTime timeutils.ESTTime `json:"LastUpdatedDateAndTime"`

	// used to calculate percentage
	TotalCapacity int `json:"TotalCapacity"`
	LastCount     int `json:"LastCount"`
}

func (g Gym) GetPercentage() int {
	return (g.LastCount * 100) / g.TotalCapacity
}
