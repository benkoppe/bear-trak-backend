// Package libcal includes libcla study content types.
package libcal

import (
	"fmt"
	"sort"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/utils/time_utils"
)

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

func GetAllDays(weeks []WeekHours) []Day {
	var days []Day

	for _, week := range weeks {
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

func ConvertToHours(weeks []WeekHours) ([]api.Hours, error) {
	est := time_utils.LoadEST()
	now := time.Now().In(est)
	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0, est,
	)
	weekAhead := today.AddDate(0, 0, 7)

	var hours []api.Hours
	for _, day := range GetAllDays(weeks) {
		if !day.Date.ToTime().Before(today) && day.Date.ToTime().Before(weekAhead) {
			if day.Times.Status == "24hours" {
				hours = append(hours, api.Hours{
					Start: day.Date.ToTime(),
					End:   day.Date.ToTime().AddDate(0, 0, 1),
				})
				continue
			}
			hours = append(hours, convertHours(day.Date.ToTime(), day.Times.Hours)...)
		}
	}

	return hours, nil
}

func convertHours(date time.Time, libcalHours []Hours) []api.Hours {
	var hours []api.Hours
	for _, libcalHour := range libcalHours {
		start, e1 := libcalHour.From.ToDate(date)
		end, e2 := libcalHour.To.ToDate(date)

		if e1 != nil {
			fmt.Printf("error parsing hours: %v", e1)
			continue
		}
		if e2 != nil {
			fmt.Printf("error parsing hours: %v", e2)
			continue
		}

		// if end is before the start (ie 12am was parsed as the wrong day, add a day)
		if end.Before(start) {
			end = end.AddDate(0, 0, 1)
		}

		hours = append(hours, api.Hours{
			Start: start,
			End:   end,
		})
	}

	return hours
}
