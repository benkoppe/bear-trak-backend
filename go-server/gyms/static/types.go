package static

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
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

type TimeString string

type Hours struct {
	Open  TimeString `json:"open"`
	Close TimeString `json:"close"`
}

type WeekHours struct {
	Monday    []Hours `json:"monday"`
	Tuesday   []Hours `json:"tuesday"`
	Wednesday []Hours `json:"wednesday"`
	Thursday  []Hours `json:"thursday"`
	Friday    []Hours `json:"friday"`
	Saturday  []Hours `json:"saturday"`
	Sunday    []Hours `json:"sunday"`
}

func (w WeekHours) GetHours(date time.Time) []Hours {
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

func (w *WeekHours) AddHours(day string, hours []Hours) error {
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

func (t TimeString) parseTime() (struct {
	Hour   int
	Minute int
}, error,
) {
	emptyTime := struct {
		Hour   int
		Minute int
	}{Hour: 0, Minute: 0}
	s := strings.ToLower(string(t))

	hasAM := strings.HasSuffix(s, "am")
	hasPM := strings.HasSuffix(s, "pm")

	// remove suffix
	if hasAM {
		s = strings.TrimSuffix(s, "am")
	}
	if hasPM {
		s = strings.TrimSuffix(s, "pm")
	}

	var hours, minutes int
	var err error

	if strings.Contains(s, ":") {
		parts := strings.Split(s, ":")
		if len(parts) != 2 {
			return emptyTime, fmt.Errorf("invalid time format: %s", t)
		}

		hours, err = strconv.Atoi(parts[0])
		if err != nil {
			return emptyTime, fmt.Errorf("invalid hours: %s", parts[0])
		}

		minutes, err = strconv.Atoi(parts[1])
		if err != nil {
			return emptyTime, fmt.Errorf("invalid minutes: %s", parts[1])
		}
	} else {
		hours, err = strconv.Atoi(s)
		if err != nil {
			return emptyTime, fmt.Errorf("invalid hours: %s", s)
		}
		minutes = 0
	}

	// validate hours and minutes
	if hours < 0 || hours > 12 && (hasAM || hasPM) || hours > 23 {
		return emptyTime, fmt.Errorf("invalid hours: %d", hours)
	}
	if minutes < 0 || minutes > 59 {
		return emptyTime, fmt.Errorf("invalid minutes: %d", minutes)
	}

	// adjust hours for PM
	if hasPM && hours < 12 {
		hours += 12
	}

	// adjust for 12am (midnight)
	if hasAM && hours == 12 {
		hours = 0
	}

	return struct {
		Hour   int
		Minute int
	}{Hour: hours, Minute: minutes}, nil
}

func (t TimeString) ToDate(date time.Time) (time.Time, error) {
	parsed, err := t.parseTime()
	if err != nil {
		return time.Time{}, err
	}
	est := utils.LoadEST()
	return time.Date(date.Year(), date.Month(), date.Day(), parsed.Hour, parsed.Minute, 0, 0, est), nil
}
