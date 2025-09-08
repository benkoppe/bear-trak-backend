// Package capacities loads capacities content from the db.
package capacities

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/db"
	"github.com/benkoppe/bear-trak-backend/go-server/gyms/cornell/external"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
	"github.com/benkoppe/bear-trak-backend/go-server/utils/timeutils"
	"github.com/jackc/pgx/v5/pgtype"
)

type Cache = *utils.Cache[[]api.GymCapacityData]

func InitCache(queries *db.Queries, externalCache external.Cache) Cache {
	return utils.NewCache(
		"gymCapacities",
		1*time.Minute,
		func() ([]api.GymCapacityData, error) {
			return LoadData(queries, externalCache)
		})
}

func LoadData(queries *db.Queries, externalCache external.Cache) ([]api.GymCapacityData, error) {
	externalData, err := externalCache.Get()
	if err != nil {
		return nil, fmt.Errorf("error fetching external data: %w", err)
	}

	est := timeutils.LoadEST()
	now := time.Now().In(est)
	dayStart := now.Truncate(24 * time.Hour)
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

	for locationID, capacities := range byLocation {
		locationExternalData := utils.Find(externalData, func(data external.Gym) bool {
			return data.LocationID == int(locationID)
		})

		if locationExternalData == nil {
			fmt.Printf("couldn't find external data for location ID %d\n", locationID)
			continue
		}

		points := make([]api.GymCapacityDataPoint, 0, len(capacities))
		for _, entry := range capacities {
			points = append(points, convertDB(entry, *locationExternalData))
		}

		result = append(result, api.GymCapacityData{
			LocationId: int(locationID),
			Points:     points,
		})

		return result, nil
	}

	return nil, nil
}

func convertDB(entry db.GymCapacity, externalData external.Gym) api.GymCapacityDataPoint {
	return api.GymCapacityDataPoint{
		LastUpdated: entry.LastUpdatedAt.Time,
		Count:       int(math.Round(float64(entry.Percentage) * float64(externalData.TotalCapacity) / 100.0)),
	}
}
