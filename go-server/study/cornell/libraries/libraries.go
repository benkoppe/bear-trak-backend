package libraries

import (
	"fmt"
	"strings"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/schools/cornell/external_map"
	"github.com/benkoppe/bear-trak-backend/go-server/study/cornell/libraries/external"
	"github.com/benkoppe/bear-trak-backend/go-server/study/cornell/libraries/static"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

func Get(cache external.Cache, mapCache external_map.Cache) ([]api.Library, error) {
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
			mapItem := utils.Find(mapItems, func(item external_map.Item) bool {
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
		ID:        static.ID,
		Name:      static.Name,
		ImagePath: utils.ImageNameToPath("study", static.ImageName),
	}
}

func getExternalHours(externalData []external.Library, lid int) ([]api.Hours, error) {
	externalLibrary := utils.Find(externalData, func(externalLibrary external.Library) bool {
		return externalLibrary.LID == lid
	})
	if externalLibrary == nil {
		return nil, fmt.Errorf("no external data for library %d", lid)
	}

	today := time.Now().Truncate(24 * time.Hour)
	weekAhead := today.AddDate(0, 0, 7)

	var hours []api.Hours
	for _, day := range externalLibrary.GetAllDays() {
		if day.Date.ToTime().After(today) && day.Date.ToTime().Before(weekAhead) {
			if day.Times.Status == "24hours" {
				hours = append(hours, api.Hours{
					Start: day.Date.ToTime(),
					End:   day.Date.ToTime().AddDate(0, 0, 1),
				})
				continue
			}
			hours = append(hours, convertExternalHours(day.Date.ToTime(), day.Times.Hours)...)
		}
	}

	return hours, nil
}

func convertExternalHours(date time.Time, externalHours []external.Hours) []api.Hours {
	var hours []api.Hours
	for _, externalHour := range externalHours {
		start, e1 := externalHour.From.ToDate(date)
		end, e2 := externalHour.To.ToDate(date)

		if e1 != nil {
			fmt.Printf("error parsing hours: %v", e1)
			continue
		}
		if e2 != nil {
			fmt.Printf("error parsing hours: %v", e2)
			continue
		}

		hours = append(hours, api.Hours{
			Start: start,
			End:   end,
		})
	}

	return hours
}
