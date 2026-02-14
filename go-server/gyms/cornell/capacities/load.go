// Package capacities loads capacities content from the db.
package capacities

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/db"
	"github.com/benkoppe/bear-trak-backend/go-server/gyms/cornell/scrape"
	"github.com/benkoppe/bear-trak-backend/go-server/gyms/cornell/shared"
	"github.com/benkoppe/bear-trak-backend/go-server/gyms/cornell/static"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
	"github.com/benkoppe/bear-trak-backend/go-server/utils/timeutils"
	"github.com/jackc/pgx/v5/pgtype"
)

type Cache = *utils.Cache[[]api.GymCapacityData]

func InitCache(queries *db.Queries, hoursCache scrape.Cache) Cache {
	return utils.NewCache(
		"gymCapacities",
		1*time.Minute,
		func() ([]api.GymCapacityData, error) {
			return LoadData(queries, hoursCache)
		})
}

func LoadData(queries *db.Queries, hoursCache scrape.Cache) ([]api.GymCapacityData, error) {
	scrapedSchedules, err := hoursCache.Get()
	if err != nil {
		fmt.Printf("error fetching scraped schedules: %v\n", err)
	}
	staticData := static.GetGyms()

	est := timeutils.LoadEST()
	now := time.Now().In(est)
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	capacities, err := queries.GetGymCapacitiesBetween(ctx, db.GetGymCapacitiesBetweenParams{
		LastUpdatedAt:   pgtype.Timestamptz{Time: dayStart, Valid: true},
		LastUpdatedAt_2: pgtype.Timestamptz{Time: dayEnd, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching capacities from db: %w", err)
	}

	// group by location id
	byLocation := make(map[int32][]db.GymCapacity)
	for _, capacity := range capacities {
		byLocation[capacity.LocationID] = append(byLocation[capacity.LocationID], capacity)
	}

	var result []api.GymCapacityData

	for internalLocationID, capacities := range byLocation {
		locationStaticData := utils.Find(staticData, func(data static.Gym) bool {
			return data.ID == int(internalLocationID)
		})
		if locationStaticData == nil {
			fmt.Printf("couldn't find static data for internal location ID %d\n", internalLocationID)
			continue
		}

		points := make([]api.GymCapacityDataPoint, 0, len(capacities))
		for _, entry := range capacities {
			points = append(points, convertDB(entry))
		}

		// perform time-based smoothing
		points = utils.SmoothTime(points,
			20*time.Minute,
			func(p api.GymCapacityDataPoint) float64 { return float64(p.Count) },
			func(p api.GymCapacityDataPoint, v float64) api.GymCapacityDataPoint {
				p.Count = int(math.Round(v))
				return p
			},
			func(p api.GymCapacityDataPoint) time.Time { return p.LastUpdated },
		)

		hours := shared.CreateFutureHours(*locationStaticData, scrapedSchedules)
		firstOpen, lastClose := timeutils.FirstOpenAndLastClose(hours, now)
		filteredPoints := shared.FilterByTimeRange(points, func(p api.GymCapacityDataPoint) time.Time { return p.LastUpdated }, firstOpen, lastClose)

		result = append(result, api.GymCapacityData{
			LocationId: int(internalLocationID),
			Points:     filteredPoints,
		})

	}

	return result, nil
}

func convertDB(entry db.GymCapacity) api.GymCapacityDataPoint {
	return api.GymCapacityDataPoint{
		LastUpdated: entry.LastUpdatedAt.Time,
		Count:       int(entry.Count),
	}
}
