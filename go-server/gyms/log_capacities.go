package gyms

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/benkoppe/bear-trak-backend/go-server/db"
	"github.com/benkoppe/bear-trak-backend/go-server/gyms/external"
	"github.com/jackc/pgx/v5/pgtype"
)

func LogCapacities(ctx context.Context, externalUrl string, queries *db.Queries) error {
	gyms, err := external.FetchData(externalUrl)
	if err != nil {
		return fmt.Errorf("error fetching gyms: %v", err)
	}

	for _, gym := range gyms {
		latestCapacity, err := queries.GetLatestCapacity(ctx, int32(gym.LocationID))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// no rows found, should be logged
				logCapacity(ctx, queries, gym)
			} else {
				// other error, continue
				log.Printf("error fetching latest capacity: %v", err)
			}
			continue
		}

		if latestCapacity.LastUpdatedAt.Time.Equal(gym.LastUpdatedDateAndTime.ToTime()) {
			// don't log the same data unnecessarily
			continue
		}

		logCapacity(ctx, queries, gym)
	}

	return nil
}

func logCapacity(ctx context.Context, queries *db.Queries, gym external.Gym) error {
	newCapacity, err := queries.CreateGymCapacity(ctx, db.CreateGymCapacityParams{
		LocationID:    int32(gym.LocationID),
		Percentage:    int32(gym.GetPercentage()),
		LastUpdatedAt: pgtype.Timestamptz{Time: gym.LastUpdatedDateAndTime.ToTime(), Valid: true},
	})
	if err == nil {
		log.Printf("\t created new new capacity entry: %v", newCapacity)
	} else {
		log.Printf("\t error creating new capacity entry: %v", err)
	}
	return err
}
