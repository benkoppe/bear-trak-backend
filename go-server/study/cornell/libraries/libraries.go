// Package libraries loads cornell library data.
package libraries

import (
	"fmt"
	"strings"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/schools/cornell/externalmap"
	"github.com/benkoppe/bear-trak-backend/go-server/study/cornell/libraries/external"
	"github.com/benkoppe/bear-trak-backend/go-server/study/cornell/libraries/static"
	"github.com/benkoppe/bear-trak-backend/go-server/study/shared/libcal"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

func Get(cache external.Cache, mapCache externalmap.Cache) ([]api.Library, error) {
	staticData := static.GetLibraries()

	if len(staticData) == 0 {
		return nil, fmt.Errorf("loaded empty static libraries")
	}

	externalData, err := cache.Get()
	if err != nil {
		return nil, fmt.Errorf("error loading external data: %v", err)
	}

	mapItems, err := mapCache.Get("Library")
	if err != nil {
		return nil, fmt.Errorf("error loading map data: %v", err)
	}

	var libraries []api.Library
	for _, staticLibrary := range staticData {
		library := convertStaticLibrary(staticLibrary)

		if staticLibrary.LID != nil {
			regularHours, err := getExternalHours(externalData, *staticLibrary.LID)
			if err != nil {
				fmt.Printf("regular hours not found for library %d: %v\n", staticLibrary.LID, err)
				continue
			}
			library.Hours = regularHours
		} else if staticLibrary.WeekHours != nil {
			library.Hours = staticLibrary.WeekHours.CreateFutureHours()
		} else {
			fmt.Printf("no hours for library %s\n", staticLibrary.Name)
			continue
		}

		if staticLibrary.LIDCardAccess != nil {
			cardHours, err := getExternalHours(externalData, *staticLibrary.LIDCardAccess)
			if err != nil {
				fmt.Printf("card access hours not found for library %d: %v\n", staticLibrary.LIDCardAccess, err)
				continue
			}
			library.CardAccessHours = cardHours
		}

		if staticLibrary.ExternalMapNote != nil {
			mapItem := utils.Find(mapItems, func(item externalmap.Item) bool {
				return strings.Contains(item.Notes, *staticLibrary.ExternalMapNote)
			})
			if mapItem == nil {
				fmt.Printf("no map item for library %d\n", staticLibrary.LID)
				continue
			}
			library.Latitude = mapItem.LatLng.Latitude
			library.Longitude = mapItem.LatLng.Longitude

		} else if staticLibrary.Location != nil {
			library.Latitude = staticLibrary.Location.Latitude
			library.Longitude = staticLibrary.Location.Longitude
		} else {
			fmt.Printf("no location for library %s\n", staticLibrary.Name)
			continue
		}

		libraries = append(libraries, library)
	}

	return libraries, nil
}

func convertStaticLibrary(static static.Library) api.Library {
	return api.Library{
		ID:               static.ID,
		Name:             static.Name,
		ImagePath:        utils.ImageNameToPath("study/cornell", static.ImageName),
		PrinterLocations: static.PrinterLocations,
	}
}

func getExternalHours(externalData []external.Library, lid int) ([]api.Hours, error) {
	externalLibrary := utils.Find(externalData, func(externalLibrary external.Library) bool {
		return externalLibrary.LID == lid
	})
	if externalLibrary == nil {
		return nil, fmt.Errorf("no external data for library %d", lid)
	}

	return libcal.ConvertToHours(externalLibrary.Weeks)
}
