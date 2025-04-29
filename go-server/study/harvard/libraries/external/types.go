package external

import "github.com/benkoppe/bear-trak-backend/go-server/study/shared/libcal"

type Library struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Image       Image       `json:"image"`
	Coordinates Coordinates `json:"coordinates"`
	WeeksHours  WeeksHours  `json:"weeks_hours"`
}

type Image struct {
	Src string `json:"src"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type WeeksHours struct {
	Locations []LocationDetails `json:"locations"` // should always have exactly one item
}

type LocationDetails struct {
	LID   int                `json:"lid"`
	Title string             `json:"title"`
	Weeks []libcal.WeekHours `json:"weeks"`
}
