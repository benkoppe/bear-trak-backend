package transit

import (
	"fmt"
	"strconv"

	"github.com/amit7itz/goset"
	"github.com/benkoppe/bear-trak-backend/go-server/api"
	availtec "github.com/benkoppe/bear-trak-backend/go-server/transit/external_availtec"
	"github.com/benkoppe/bear-trak-backend/go-server/transit/external_gtfs"
	"github.com/jamespfennell/gtfs"
	"github.com/twpayne/go-polyline"
)

type Caches struct {
	availtecCache availtec.Cache
	staticCache   external_gtfs.Cache
}

func InitCaches(availtecUrl string, staticGtfsUrl string) Caches {
	return Caches{
		availtecCache: availtec.InitCache(availtecUrl),
		staticCache:   external_gtfs.InitCache(staticGtfsUrl),
	}
}

func GetRoutes(caches Caches) ([]api.BusRoute, error) {
	staticGtfs, err := caches.staticCache.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to load static data: %v", err)
	}

	availtecRoutes, err := caches.availtecCache.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to load availtec routes: %v", err)
	}

	routes, err := getRoutes(availtecRoutes, *staticGtfs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse routes: %v", err)
	}

	vehicles, err := getVehiclesFromRoutes(availtecRoutes)
	if err != nil {
		return nil, fmt.Errorf("failed to load vehicles: %v", err)
	}

	routeIdVehicles := make(map[int]([]api.Vehicle))
	for _, vehicle := range vehicles {
		routeIdVehicles[vehicle.RouteId] = append(routeIdVehicles[vehicle.RouteId], vehicle)
	}

	for i := range routes {
		routes[i].Vehicles = routeIdVehicles[routes[i].ID]
	}

	return routes, nil
}

func getRoutes(availtecRoutes []availtec.Route, staticGtfs gtfs.Static) ([]api.BusRoute, error) {
	var routes []api.BusRoute

	for _, route := range availtecRoutes {
		var gtfsRoute *gtfs.Route
		for _, gtfsRouteOption := range staticGtfs.Routes {
			if gtfsRouteOption.Id == strconv.Itoa(route.RouteId) {
				gtfsRoute = &gtfsRouteOption
				break
			}
		}

		if gtfsRoute == nil {
			return nil, fmt.Errorf("failed to find GTFS route for route ID: %v", route.RouteId)
		}

		tripsMap := getDirectionTrips(*gtfsRoute, staticGtfs)

		var polylines []string
		var directions []api.BusRouteDirection

		for directionId, trips := range tripsMap {
			polylines = append(polylines, getPolylines(trips)...)
			stops := getStops(trips)

			directions = append(directions, api.BusRouteDirection{
				Name:  convertStaticDirectionId(directionId),
				Stops: convertStaticStops(stops),
			})
		}

		routes = append(routes, api.BusRoute{
			ID:         route.RouteId,
			SortIdx:    route.SortOrder,
			Name:       route.GoogleDescription,
			Code:       route.RouteAbbreviation,
			Color:      route.Color,
			Directions: directions,
			Polylines:  polylines,
		})
	}

	return routes, nil
}

func convertStaticDirectionId(id gtfs.DirectionID) string {
	switch id {
	case gtfs.DirectionID_True:
		return "O"
	case gtfs.DirectionID_False:
		return "I"
	default:
		return "?"
	}
}

func convertStaticStops(stops []gtfs.Stop) []api.BusRouteDirectionStopsItem {
	var backendStops []api.BusRouteDirectionStopsItem

	for _, stop := range stops {
		backendStops = append(backendStops, api.BusRouteDirectionStopsItem{
			ID:        stop.Id,
			Name:      stop.Name,
			Longitude: *stop.Longitude,
			Latitude:  *stop.Latitude,
		})
	}

	return backendStops
}

func getDirectionTrips(route gtfs.Route, staticGtfs gtfs.Static) map[gtfs.DirectionID][]gtfs.ScheduledTrip {
	var routeTrips []gtfs.ScheduledTrip

	for _, trip := range staticGtfs.Trips {
		if *trip.Route == route {
			routeTrips = append(routeTrips, trip)
		}
	}

	directionMappedTrips := make(map[gtfs.DirectionID][]gtfs.ScheduledTrip)

	for _, trip := range routeTrips {
		directionMappedTrips[trip.DirectionId] = append(directionMappedTrips[trip.DirectionId], trip)
	}

	return directionMappedTrips
}

func getStops(trips []gtfs.ScheduledTrip) []gtfs.Stop {
	stops := goset.NewSet[gtfs.Stop]()

	for _, trip := range trips {
		for _, stopTime := range trip.StopTimes {
			if stopTime.Stop != nil {
				stops.Add(*stopTime.Stop)
			}
		}
	}

	return stops.Items()
}

func getPolylines(trips []gtfs.ScheduledTrip) []string {
	polylines := goset.NewSet[string]()

	for _, trip := range trips {
		shape := trip.Shape
		if shape == nil {
			continue
		}

		var coords [][]float64

		for _, point := range shape.Points {
			coords = append(coords, []float64{point.Latitude, point.Longitude})
		}

		if len(coords) < 2 {
			continue
		}

		line := string(polyline.EncodeCoords(coords))
		polylines.Add(line)
	}

	return polylines.Items()
}
