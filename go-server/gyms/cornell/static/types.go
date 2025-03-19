package static

import (
	"fmt"
	"strings"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils/time_utils"
)

type Gym struct {
	ID         int         `json:"id"`
	LocationID int         `json:"locationId"`
	Name       string      `json:"name"`
	ScrapeName string      `json:"scrapeName"`
	ImageName  string      `json:"imageName"`
	Location   Location    `json:"location"`
	Facilities []Facility  `json:"facilities"`
	Equipment  []Equipment `json:"equipment"`
	WeekHours  WeekHours   `json:"weekHours"`
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

type WeekHours struct {
	Monday    []time_utils.Hours `json:"monday"`
	Tuesday   []time_utils.Hours `json:"tuesday"`
	Wednesday []time_utils.Hours `json:"wednesday"`
	Thursday  []time_utils.Hours `json:"thursday"`
	Friday    []time_utils.Hours `json:"friday"`
	Saturday  []time_utils.Hours `json:"saturday"`
	Sunday    []time_utils.Hours `json:"sunday"`
}

func (w WeekHours) GetHours(date time.Time) []time_utils.Hours {
	switch date.Weekday() {
	case time.Monday:
		return w.Monday
	case time.Tuesday:
		return w.Tuesday
	case time.Wednesday:
		return w.Wednesday
	case time.Thursday:
		return w.Thursday
	case time.Friday:
		return w.Friday
	case time.Saturday:
		return w.Saturday
	case time.Sunday:
		return w.Sunday
	default:
		return nil
	}
}

func (w *WeekHours) AddHours(day string, hours []time_utils.Hours) error {
	switch strings.ToLower(day) {
	case "monday":
		w.Monday = append(w.Monday, hours...)
	case "tuesday":
		w.Tuesday = append(w.Tuesday, hours...)
	case "wednesday":
		w.Wednesday = append(w.Wednesday, hours...)
	case "thursday":
		w.Thursday = append(w.Thursday, hours...)
	case "friday":
		w.Friday = append(w.Friday, hours...)
	case "saturday":
		w.Saturday = append(w.Saturday, hours...)
	case "sunday":
		w.Sunday = append(w.Sunday, hours...)
	default:
		return fmt.Errorf("invalid day: %s", day)
	}
	return nil
}

func (w WeekHours) IsOpen(date time.Time) bool {
	dayHours := w.GetHours(date)

	for _, hours := range dayHours {
		open, err := hours.Open.ToDate(date)
		if err != nil {
			continue
		}
		close, err := hours.Close.ToDate(date)
		if err != nil {
			continue
		}

		if open.Before(date) && date.Before(close) {
			return true
		}
	}

	return false
}
