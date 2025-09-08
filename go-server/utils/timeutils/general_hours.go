package timeutils

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
)

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

func (w WeekHours) GetConvertedHours(date time.Time) []api.Hours {
	var convertedHours []api.Hours
	dayHours := w.GetHours(date)
	for _, hours := range dayHours {

		futureHour, err := hours.Convert(date)
		if err != nil {
			fmt.Printf("error converting hours: %v\n", err)
			continue
		}
		convertedHours = append(convertedHours, *futureHour)
	}
	return convertedHours
}

func (w WeekHours) CreateFutureHours() []api.Hours {
	est := LoadEST()
	now := time.Now().In(est)
	var futureHours []api.Hours

	for i := range [7]int{} {
		date := now.AddDate(0, 0, i)
		futureHours = append(futureHours, w.GetConvertedHours(date)...)
	}

	return futureHours
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

type Hours struct {
	Open  TimeString `json:"open"`
	Close TimeString `json:"close"`
}

func (h *Hours) Convert(date time.Time) (*api.Hours, error) {
	start, e1 := h.Open.ToDate(date)
	end, e2 := h.Close.ToDate(date)

	if e1 != nil {
		return nil, fmt.Errorf("error parsing hours: %v", e1)
	}
	if e2 != nil {
		return nil, fmt.Errorf("error parsing hours: %v", e2)
	}

	return &api.Hours{
		Start: start,
		End:   end,
	}, nil
}

type HoursParserOptions struct {
	RangeSeparators     []string // Separators between time ranges
	OpenCloseSeparators []string // Separators between open/close times
	ClosedKeywords      []string // Words indicating closed status
	Open24HoursKeywords []string // Words indicating 24-hour operation
}

func DefaultHoursParserOptions() *HoursParserOptions {
	return &HoursParserOptions{
		RangeSeparators:     []string{"/", ","},
		OpenCloseSeparators: []string{"-", "to", "–", "—"}, // includes en/em dashes
		ClosedKeywords:      []string{"Closed", "Open by Appointment"},
		Open24HoursKeywords: []string{"24 hours", "24/7", "All day", "Open 24 hours"},
	}
}

// ParseHours splits a cell value into one or more Hours.
// It handles values like "6am - 9pm" or "7am - 8:30am / 10am - 10:45pm".
func ParseHours(raw string, opts *HoursParserOptions) []Hours {
	if opts == nil {
		opts = DefaultHoursParserOptions()
	}

	var hrs []Hours
	rawLower := strings.ToLower(raw)

	for _, keyword := range opts.ClosedKeywords {
		if strings.Contains(rawLower, strings.ToLower(keyword)) {
			return hrs
		}
	}

	for _, keyword := range opts.Open24HoursKeywords {
		if strings.Contains(rawLower, strings.ToLower(keyword)) {
			return []Hours{{
				Open:  "00:00",
				Close: "23:59",
			}}
		}
	}

	// Build regex for range separators
	rangeSepPattern := strings.Join(opts.RangeSeparators, "|")
	rangeRegex := regexp.MustCompile("\\s*(" + rangeSepPattern + ")\\s*")

	ranges := rangeRegex.Split(raw, -1)
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}

		var openTime, closeTime string
		found := false
		for _, sep := range opts.OpenCloseSeparators {
			parts := strings.Split(r, sep)
			if len(parts) == 2 {
				openTime = strings.TrimSpace(parts[0])
				closeTime = strings.TrimSpace(parts[1])
				found = true
				break
			}
		}

		if !found {
			log.Printf("Couldn't format a general hours range: %s", r)
			continue
		}

		hrs = append(hrs, Hours{
			Open:  TimeString(openTime),
			Close: TimeString(closeTime),
		})
	}

	return hrs
}
