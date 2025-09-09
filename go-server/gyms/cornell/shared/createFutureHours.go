// Package shared contains shared functions used for cornell gyms
package shared

import (
	"log"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/gyms/cornell/scrape"
	"github.com/benkoppe/bear-trak-backend/go-server/gyms/cornell/static"
	"github.com/benkoppe/bear-trak-backend/go-server/utils/timeutils"
)

func CreateFutureHours(static static.Gym, schedules []scrape.ParsedSchedule) []api.Hours {
	staticHours := static.WeekHours
	est := timeutils.LoadEST()
	now := time.Now().In(est)
	var futureHours []api.Hours

	for i := range [7]int{} {
		date := now.AddDate(0, 0, i)
		weekHours := staticHours
		overrideStatic := false

		// if a scraped schedule is found, override the static hours for this day
		schedule := scrape.DetermineRelevantSchedule(schedules, date)
		if schedule != nil {
			gymSchedule := scrape.GetGymSchedule(*schedule, static)
			if gymSchedule != nil {
				weekHours = gymSchedule.WeekHours
				overrideStatic = true
			}
		}

		if !overrideStatic {
			// log that static data was used for hours
			log.Printf("FALLBACK: using static hours for gym %s on %s", static.Name, date)
		}

		futureHours = append(futureHours, weekHours.GetConvertedHours(date)...)
	}

	return futureHours
}
