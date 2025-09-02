// Package gyms includes all general gym methods.
package gyms

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/db"
	"github.com/jackc/pgx/v5/pgtype"
)

func LogCapacities(ctx context.Context, handler api.Handler, queries *db.Queries) error {
	gyms, err := handler.GetV1Gyms(ctx)
	if err != nil {
		return fmt.Errorf("error fetching gyms: %v", err)
	}

	for _, gym := range gyms {
		capacity, capacitySet := gym.Capacity.Get()
		if !capacitySet {
			continue
		}

		latestCapacity, err := queries.GetLatestCapacity(ctx, int32(gym.ID))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// no rows found, should be logged
				logCapacity(ctx, queries, gym, capacity)
			} else {
				// other error, continue
				log.Printf("error fetching latest capacity: %v", err)
			}
			continue
		}

		if latestCapacity.LastUpdatedAt.Time.Equal(capacity.LastUpdated) {
			// don't log the same data unnecessarily
			continue
		}

		logCapacity(ctx, queries, gym, capacity)
	}

	return nil
}

func logCapacity(ctx context.Context, queries *db.Queries, gym api.Gym, capacity api.GymCapacity) error {
	percentage, percentageSet := capacity.Percentage.Get()
	if !percentageSet {
		return nil
	}

	newCapacity, err := queries.CreateGymCapacity(ctx, db.CreateGymCapacityParams{
		LocationID:    int32(gym.ID),
		Percentage:    int32(percentage),
		LastUpdatedAt: pgtype.Timestamptz{Time: capacity.LastUpdated, Valid: true},
	})
	if err == nil {
		log.Printf("\t created new new capacity entry: %v", newCapacity)
	} else {
		log.Printf("\t error creating new capacity entry: %v", err)
	}
	return err
}
