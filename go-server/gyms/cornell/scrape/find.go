package scrape

import (
	"log"
	"strings"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/gyms/cornell/static"
)

// functions used elsewhere in the package for working with fetched schedules

func DetermineRelevantSchedule(schedules []ParsedSchedule, date time.Time) *ParsedSchedule {
	// search for exact matches first
	for _, s := range schedules {
		if s.Caption.StartDate != nil && s.Caption.EndDate != nil {
			if date.After(*s.Caption.StartDate) && date.Before(*s.Caption.EndDate) {
				return &s
			}
		}
	}

	// search for half-matches next
	for _, s := range schedules {
		if s.Caption.StartDate != nil && s.Caption.EndDate == nil {
			if s.Caption.StartDate.Before(date) {
				return &s
			}
		} else if s.Caption.StartDate == nil && s.Caption.EndDate != nil {
			if s.Caption.EndDate.After(date) {
				return &s
			}
		}
	}

	// finally search for general matches
	for _, s := range schedules {
		if strings.Contains(s.Caption.Title, "Regular Semester Hours") {
			return &s
		}
	}

	log.Printf("No relevant schedule found for date %v", date)
	return nil
}

func GetGymSchedule(schedule ParsedSchedule, staticGym static.Gym) *GymSchedule {
	for _, s := range schedule.GymSchedules {
		if s.GymName == staticGym.ScrapeName {
			return &s
		}
	}
	return nil
}
