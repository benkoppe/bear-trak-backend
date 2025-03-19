package static

import "github.com/benkoppe/bear-trak-backend/go-server/utils/time_utils"

type Library struct {
	ID              int                   `json:"id"`
	Name            string                `json:"name"`
	ImageName       string                `json:"imageName"`
	LID             *int                  `json:"lid"`
	LIDCardAccess   *int                  `json:"lid_card_access,omitempty"`
	ExternalMapNote *string               `json:"externalMapNote,omitempty"`
	Location        *Location             `json:"location,omitempty"`
	WeekHours       *time_utils.WeekHours `json:"weekHours,omitempty"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
