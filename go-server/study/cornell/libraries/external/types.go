package external

import "github.com/benkoppe/bear-trak-backend/go-server/study/shared/libcal"

type librariesResponse struct {
	Locations []Library `json:"locations"`
}

type Library struct {
	LID      int    `json:"lid"`
	Name     string `json:"name"`
	Category string `json:"category"`
	ParentID int    `json:"parent_lid,omitempty"`

	Weeks []libcal.WeekHours `json:"weeks"`
}
