package scrape

import (
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils/timeutils"
)

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
	WeekHours timeutils.WeekHours
}
