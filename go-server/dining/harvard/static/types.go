package static

import (
	"fmt"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/utils/timeutils"
)

type Eatery struct {
	ID          int              `json:"id"`
	Name        string           `json:"name"`
	APINumber   *string          `json:"apiNumber,omitempty"`
	ImageName   string           `json:"imageName"`
	Region      string           `json:"region"`
	Location    Location         `json:"location"`
	Categories  []string         `json:"categories"`
	AllWeekMenu *[]MenuCategory  `json:"allWeekMenu,omitempty"`
	WeekHours   HarvardWeekHours `json:"weekHours"`
}

type MenuCategory struct {
	Category string     `json:"category"`
	Items    []MenuItem `json:"items"`
}

type MenuItem struct {
	Item    string `json:"item"`
	Healthy bool   `json:"healthy"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type HarvardHours struct {
	MealNumber *int                 `json:"mealNumber,omitempty"`
	Open       timeutils.TimeString `json:"open"`
	Close      timeutils.TimeString `json:"close"`
}

type HarvardWeekHours struct {
	Monday    []HarvardHours `json:"monday"`
	Tuesday   []HarvardHours `json:"tuesday"`
	Wednesday []HarvardHours `json:"wednesday"`
	Thursday  []HarvardHours `json:"thursday"`
	Friday    []HarvardHours `json:"friday"`
	Saturday  []HarvardHours `json:"saturday"`
	Sunday    []HarvardHours `json:"sunday"`
}

func (w HarvardWeekHours) GetHours(date time.Time) []HarvardHours {
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

func (w HarvardWeekHours) GetConvertedHours(date time.Time) []api.Hours {
	var convertedHours []api.Hours
	dayHours := w.GetHours(date)
	for _, hours := range dayHours {

		futureHour, err := hours.Convert(date)
		if err != nil {
			fmt.Printf("error converting hours: %v", err)
			continue
		}
		convertedHours = append(convertedHours, *futureHour)
	}
	return convertedHours
}

func (w HarvardWeekHours) CreateFutureHours() []api.Hours {
	est := timeutils.LoadEST()
	now := time.Now().In(est)
	var futureHours []api.Hours

	for i := range [7]int{} {
		date := now.AddDate(0, 0, i)
		futureHours = append(futureHours, w.GetConvertedHours(date)...)
	}

	return futureHours
}

func (h *HarvardHours) Convert(date time.Time) (*api.Hours, error) {
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
