package external

import "github.com/benkoppe/bear-trak-backend/utils"

type Gym struct {
	LocationID   int    `json:"LocationId"`
	LocationName string `json:"LocationName"`

	LastUpdatedDateAndTime utils.ESTTime `json:"LastUpdatedDateAndTime"`

	// used to calculate percentage
	TotalCapacity int `json:"TotalCapacity"`
	LastCount     int `json:"LastCount"`
}

func (g Gym) GetPercentage() int {
	return (g.LastCount * 100) / g.TotalCapacity
}
