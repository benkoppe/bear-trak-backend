package scrape

import (
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/gyms/static"
)

type tableData struct {
	Caption string
	Headers []string
	Rows    []rowData
}

type rowData struct {
	Columns []string
}

type ParsedSchedule struct {
	Caption      captionData
	GymSchedules []GymSchedule
}

type captionData struct {
	Title     string
	StartDate *time.Time
	EndDate   *time.Time
}

type GymSchedule struct {
	GymName   string
	WeekHours static.WeekHours
}
