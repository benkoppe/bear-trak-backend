// Package gtfs handles gtfs data.
package gtfs

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
	"github.com/jamespfennell/gtfs"
)

type Cache = *utils.Cache[*gtfs.Static]

func InitCache(source string) Cache {
	return InitFilteredCache(source, nil)
}

// InitFilteredCache keeps only trips on routeIDs. Empty routeIDs keeps all.
func InitFilteredCache(source string, routeIDs []string) Cache {
	return utils.NewCache(
		"transitExternalGtfs",
		utils.NoExpiration,
		func() (*gtfs.Static, error) {
			return loadData(source, routeIDs)
		},
	)
}

func loadData(source string, routeIDs []string) (*gtfs.Static, error) {
	tcatGtfsData, err := loadTcatGtfs(source)
	if err != nil {
		return nil, fmt.Errorf("error loading tcat data: %v", err)
	}

	staticData, err := gtfs.ParseStatic(tcatGtfsData, gtfs.ParseStaticOptions{})
	if err != nil {
		return nil, fmt.Errorf("error parsing tcat data: %v", err)
	}

	if len(routeIDs) > 0 {
		keepRoutes(staticData, routeIDs)
	}

	return staticData, nil
}

// keepRoutes drops trips not on routeIDs to save memory on large feeds.
// Kept trips still reference their own shapes.
func keepRoutes(staticData *gtfs.Static, routeIDs []string) {
	var trips []gtfs.ScheduledTrip
	for _, trip := range staticData.Trips {
		if trip.Route != nil && slices.Contains(routeIDs, trip.Route.Id) {
			trips = append(trips, trip)
		}
	}

	// new slice so dropped trips can be freed
	staticData.Trips = trips
	staticData.Shapes = nil
}

func loadTcatGtfs(source string) ([]byte, error) {
	var reader io.ReadCloser
	var err error

	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		// Load from HTTP
		resp, err := http.Get(source)
		if err != nil {
			return nil, fmt.Errorf("error fetching gtfs ZIP: %v", err)
		}
		reader = resp.Body
	} else {
		// Load from file
		reader, err = os.Open(source)
		if err != nil {
			return nil, fmt.Errorf("error opening gtfs ZIP file: %v", err)
		}
	}

	defer func() { _ = reader.Close() }()

	originalGtfsData, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading gtfs ZIP: %v", err)
	}

	return originalGtfsData, nil
}
