// Package gyms includes all general gym methods.
package gyms

import (
	"context"
	"crypto/subtle"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/db"
	"github.com/jackc/pgx/v5/pgtype"
)

func NewCapacitiesExportHandler(queries *db.Queries, token string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		requestToken := r.Header.Get("X-Internal-Token")
		if subtle.ConstantTimeCompare([]byte(requestToken), []byte(token)) != 1 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		query := r.URL.Query()
		start, err := parseRFC3339QueryParam(query.Get("start"), "start")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		end, err := parseRFC3339QueryParam(query.Get("end"), "end")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if start != nil && end != nil && start.After(*end) {
			http.Error(w, "start must be before or equal to end", http.StatusBadRequest)
			return
		}

		locationIDFilter, hasLocationIDFilter, err := parseOptionalLocationID(query.Get("locationId"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 120*time.Second)
		defer cancel()

		rows, err := loadCapacityRows(ctx, queries, start, end)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to load capacity rows: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")

		writer := csv.NewWriter(w)
		defer writer.Flush()

		if err := writer.Write([]string{"id", "location_id", "last_updated_at", "percentage", "count", "total_capacity"}); err != nil {
			http.Error(w, fmt.Sprintf("failed to write csv header: %v", err), http.StatusInternalServerError)
			return
		}

		for _, row := range rows {
			if hasLocationIDFilter && row.LocationID != locationIDFilter {
				continue
			}

			record := []string{
				strconv.FormatInt(int64(row.ID), 10),
				strconv.FormatInt(int64(row.LocationID), 10),
				row.LastUpdatedAt.Time.Format(time.RFC3339Nano),
				strconv.FormatInt(int64(row.Percentage), 10),
				strconv.FormatInt(int64(row.Count), 10),
				strconv.FormatInt(int64(row.TotalCapacity), 10),
			}
			if err := writer.Write(record); err != nil {
				http.Error(w, fmt.Sprintf("failed to write csv row: %v", err), http.StatusInternalServerError)
				return
			}
		}

		if err := writer.Error(); err != nil {
			http.Error(w, fmt.Sprintf("failed to flush csv: %v", err), http.StatusInternalServerError)
			return
		}
	})
}

func parseRFC3339QueryParam(value string, name string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, fmt.Errorf("invalid %s: must be RFC3339", name)
	}

	return &parsed, nil
}

func parseOptionalLocationID(value string) (int32, bool, error) {
	if value == "" {
		return 0, false, nil
	}

	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, false, fmt.Errorf("invalid locationId: must be an integer")
	}

	return int32(parsed), true, nil
}

func loadCapacityRows(ctx context.Context, queries *db.Queries, start *time.Time, end *time.Time) ([]db.GymCapacity, error) {
	if start != nil && end != nil {
		return queries.GetGymCapacitiesBetween(ctx, db.GetGymCapacitiesBetweenParams{
			LastUpdatedAt:   pgtype.Timestamptz{Time: *start, Valid: true},
			LastUpdatedAt_2: pgtype.Timestamptz{Time: *end, Valid: true},
		})
	}

	if start != nil {
		return queries.GetGymCapacitiesFrom(ctx, pgtype.Timestamptz{Time: *start, Valid: true})
	}

	if end != nil {
		return queries.GetGymCapacitiesTo(ctx, pgtype.Timestamptz{Time: *end, Valid: true})
	}

	return queries.GetGymCapacitiesAll(ctx)
}
