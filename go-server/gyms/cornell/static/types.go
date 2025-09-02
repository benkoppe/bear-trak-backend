package static

import (
	"github.com/benkoppe/bear-trak-backend/go-server/utils/timeutils"
)

type Gym struct {
	ID         int                 `json:"id"`
	LocationID int                 `json:"locationId"`
	Name       string              `json:"name"`
	ScrapeName string              `json:"scrapeName"`
	ImageName  string              `json:"imageName"`
	Location   Location            `json:"location"`
	Facilities []Facility          `json:"facilities"`
	Equipment  []Equipment         `json:"equipment"`
	WeekHours  timeutils.WeekHours `json:"weekHours"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Facility struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type Equipment struct {
	Type  string   `json:"type"`
	Items []string `json:"items"`
}
