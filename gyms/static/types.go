package static

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Gym struct {
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	ImageName  string      `json:"imageName"`
	ScrapeName string      `json:"scrapeName"`
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

func (t TimeString) parseTime() (struct {
	Hour   int
	Minute int
}, error,
) {
	parts := strings.Split(string(t), ":")
	emptyTime := struct {
		Hour   int
		Minute int
	}{Hour: 0, Minute: 0}

	if len(parts) != 2 {
		return emptyTime, fmt.Errorf("invalid time format: %s", t)
	}
	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return emptyTime, fmt.Errorf("invalid hours: %s", parts[0])
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return emptyTime, fmt.Errorf("invalid minutes: %s", parts[1])
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
	estLocation, err := time.LoadLocation("America/New_York")
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(date.Year(), date.Month(), date.Day(), parsed.Hour, parsed.Minute, 0, 0, estLocation), nil
}
