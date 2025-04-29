package libraries

import (
	"fmt"
	"log"
	"strconv"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/study/harvard/libraries/external"
	"github.com/benkoppe/bear-trak-backend/go-server/study/shared/libcal"
)

func Get(cache external.Cache) ([]api.Library, error) {
	externalData, err := cache.Get()
	if err != nil {
		return nil, fmt.Errorf("error loading external data: %v", err)
	}

	var libraries []api.Library
	for _, externalLibrary := range externalData {
		library, err := convertExternalLibrary(externalLibrary)
		if err != nil {
			log.Printf("error converting external library: %v", err)
			continue
		}

		libraries = append(libraries, *library)
	}

	return libraries, nil
}

func convertExternalLibrary(external external.Library) (*api.Library, error) {
	library := api.Library{
		Name:             external.Name,
		Latitude:         external.Coordinates.Latitude,
		Longitude:        external.Coordinates.Longitude,
		PrinterLocations: []string{},
	}

	id, err := strconv.Atoi(external.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert ID to int: %v", err)
	}
	library.ID = id

	if len(external.WeeksHours.Locations) != 1 {
		return nil, fmt.Errorf("expected 1 weeks_hours location, got=%d", len(external.WeeksHours.Locations))
	}
	details := external.WeeksHours.Locations[0]

	hours, err := libcal.ConvertToHours(details.Weeks)
	if err != nil {
		return nil, fmt.Errorf("failed to convert libcal hours: %v", err)
	}
	library.Hours = hours

	return &library, nil
}
