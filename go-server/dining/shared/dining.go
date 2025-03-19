package shared

import (
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/utils/time_utils"
)

func SelectNextWeekEvents(events []api.EateryEvent) api.EateryNextWeekEvents {
	est := time_utils.LoadEST()
	now := time.Now().In(est)
	weekStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekEnd := weekStart.AddDate(0, 0, 7)

	result := api.EateryNextWeekEvents{}

	for _, event := range events {
		if event.Start.After(weekStart) && event.Start.Before(weekEnd) {
			switch event.Start.Weekday() {
			case time.Monday:
				result.Monday = append(result.Monday, event)
			case time.Tuesday:
				result.Tuesday = append(result.Tuesday, event)
			case time.Wednesday:
				result.Wednesday = append(result.Wednesday, event)
			case time.Thursday:
				result.Thursday = append(result.Thursday, event)
			case time.Friday:
				result.Friday = append(result.Friday, event)
			case time.Saturday:
				result.Saturday = append(result.Saturday, event)
			case time.Sunday:
				result.Sunday = append(result.Sunday, event)
			}
		}
	}

	return result
}
