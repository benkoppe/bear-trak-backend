package external

import (
	"sort"

	"github.com/benkoppe/bear-trak-backend/go-server/utils/time_utils"
)

type librariesResponse struct {
	Locations []Library `json:"locations"`
}

type Library struct {
	LID      int         `json:"lid"`
	Name     string      `json:"name"`
	Category string      `json:"category"`
	ParentID int         `json:"parent_lid,omitempty"`
	Weeks    []WeekHours `json:"weeks"`
}

type WeekHours struct {
	Sunday    Day `json:"Sunday"`
	Monday    Day `json:"Monday"`
	Tuesday   Day `json:"Tuesday"`
	Wednesday Day `json:"Wednesday"`
	Thursday  Day `json:"Thursday"`
	Friday    Day `json:"Friday"`
	Saturday  Day `json:"Saturday"`
}

type Day struct {
	Date     time_utils.DateTimeString `json:"date"`
	Times    Times                     `json:"times"`
	Rendered string                    `json:"rendered"`
}

type Times struct {
	Status        string  `json:"status"`
	Hours         []Hours `json:"hours,omitempty"`
	CurrentlyOpen bool    `json:"currently_open"`
}

type Hours struct {
	From time_utils.TimeString `json:"from"`
	To   time_utils.TimeString `json:"to"`
}

func (lib *Library) GetAllDays() []Day {
	var days []Day

	for _, week := range lib.Weeks {
		days = append(days,
			week.Sunday,
			week.Monday,
			week.Tuesday,
			week.Wednesday,
			week.Thursday,
			week.Friday,
			week.Saturday,
		)
	}

	// Sort by date
	sort.Slice(days, func(i, j int) bool {
		return days[i].Date.ToTime().Before(days[j].Date.ToTime())
	})

	return days
}
