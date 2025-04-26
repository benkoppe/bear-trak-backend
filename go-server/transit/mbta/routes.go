package mbta

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/transit/shared"
	shared_gtfs "github.com/benkoppe/bear-trak-backend/go-server/transit/shared/gtfs"
	"github.com/benkoppe/bear-trak-backend/go-server/transit/shared/gtfs_rt"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
	"github.com/jamespfennell/gtfs"
)

type Caches struct {
	staticCache   shared_gtfs.Cache
	realtimeCache gtfs_rt.Cache
}

func InitCaches(staticGtfsUrl, realtimeGtfsBaseUrl string) Caches {
	alerts, err := utils.ExtendUrl(realtimeGtfsBaseUrl, "Alerts.pb")
	if err != nil {
		log.Fatalf("failed to extend realtime GTFS alerts URL: %v", err)
	}
	tripUpdates, err := utils.ExtendUrl(realtimeGtfsBaseUrl, "TripUpdates.pb")
	if err != nil {
		log.Fatalf("failed to extend realtime GTFS tripupdates URL: %v", err)
	}
	vehicles, err := utils.ExtendUrl(realtimeGtfsBaseUrl, "VehiclePositions.pb")
	if err != nil {
		log.Fatalf("failed to extend realtime GTFS vehicle positions URL: %v", err)
	}

	mbtaGtfsRealtime := gtfs_rt.RealtimeUrls{
		Alerts:           *alerts,
		TripUpdates:      *tripUpdates,
		VehiclePositions: *vehicles,
	}

	return Caches{
		staticCache:   shared_gtfs.InitCache(staticGtfsUrl),
		realtimeCache: gtfs_rt.InitCache(mbtaGtfsRealtime),
	}
}

func GetRoutes(caches Caches) ([]api.BusRoute, error) {
	staticGtfs, err := caches.staticCache.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to load static data: %v", err)
	}

	routes, err := getRoutes(*staticGtfs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse routes: %v", err)
	}

	realtimeGtfs, err := caches.realtimeCache.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to load realtime data: %v", err)
	}

	vehicles, err := getVehicles(*staticGtfs, *realtimeGtfs)
	if err != nil {
		return nil, fmt.Errorf("failed to load vehicles: %v", err)
	}

	routes = shared.AppendVehicles(routes, vehicles)

	return routes, nil
}

func getRoutes(staticGtfs gtfs.Static) ([]api.BusRoute, error) {
	var routes []api.BusRoute

	for _, route := range staticGtfs.Routes {
		apiRoute := convertRoute(route, staticGtfs)

		routes = append(routes, apiRoute)
	}

	resolveDuplicateCodes(routes)

	return routes, nil
}

func resolveDuplicateCodes(routes []api.BusRoute) {
	codeCounts := make(map[string]int)

	for i := range routes {
		originalCode := routes[i].Code

		count, exists := codeCounts[originalCode]

		if !exists {
			codeCounts[originalCode] = 1
		} else {
			newCode := originalCode + "-" + strconv.Itoa(count)
			routes[i].Code = newCode

			codeCounts[originalCode] = count + 1
		}
	}
}

func convertRoute(route gtfs.Route, staticGtfs gtfs.Static) api.BusRoute {
	apiRoute := shared_gtfs.ConvertRoute(route, staticGtfs)

	id, ok := apiRoute.ID.GetString()
	if !ok {
		return apiRoute
	}
	code := apiRoute.Code
	codeLen := len(code)

	if code == "Shuttle" {
		apiRoute.Code = firstThreeCaps(apiRoute.Name)
	} else if codeLen == 0 {
		// separate ID into name
		// take capital letters from ID for code
		apiRoute.Code = firstThreeCaps(id)
		apiRoute.Name = addSpacesAroundWordsAndDashes(id)
	} else if codeLen <= 7 {
	} else {
		// separate ID into name
		// take capital letters from ID for code
		apiRoute.Code = firstThreeCaps(id)
		apiRoute.Name = addSpacesAroundWordsAndDashes(id)
	}

	return apiRoute
}

// excludes "Shuttle" if the string starts with that
func firstThreeCaps(s string) string {
	var result strings.Builder
	count := 0

	processString := s

	if strings.HasPrefix(s, "Shuttle") {
		if len(s) > 7 {
			processString = s[7:]
		} else {
			return ""
		}
	}

	for _, r := range processString {
		if unicode.IsUpper(r) {
			result.WriteRune(r)
			count++
			if count == 3 {
				break
			}
		}
	}
	return result.String()
}

func addSpacesAroundWordsAndDashes(s string) string {
	var result strings.Builder
	inWord := false

	for i, r := range s {
		isUpper := unicode.IsUpper(r)
		isDash := r == '-'

		if isDash {
			if result.Len() > 0 && result.String()[result.Len()-1] != ' ' {
				result.WriteRune(' ')
			}
			result.WriteRune(r)
			result.WriteRune(' ')
			inWord = false
			continue
		}

		if isUpper {
			if !inWord && result.Len() > 0 && result.String()[result.Len()-1] != ' ' {
				result.WriteRune(' ')
			}
			result.WriteRune(r)
			inWord = true
		} else {
			result.WriteRune(r)
			if i+1 < len(s) {
				nextRune := rune(s[i+1])
				if unicode.IsUpper(nextRune) && !unicode.IsUpper(r) {
					result.WriteRune(' ')
				}
			}
			inWord = false
		}
	}
	return strings.TrimSpace(result.String())
}
