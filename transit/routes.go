package transit

import (
	"fmt"
	"strconv"

	"github.com/amit7itz/goset"
	backend "github.com/benkoppe/bear-trak-backend/backend"
	"github.com/benkoppe/bear-trak-backend/transit/external_gtfs"
	"github.com/jamespfennell/gtfs"
	"github.com/twpayne/go-polyline"
)

func GetRoutes(staticUrl string, realtimeUrls external_gtfs.RealtimeUrls) ([]backend.BusRoute, error) {
	staticGtfs := external_gtfs.GetStaticGtfs(staticUrl)

	realtimeGtfs, err := external_gtfs.GetRealtimeGtfs(realtimeUrls)
	if err != nil {
		return nil, fmt.Errorf("failed to load realtime gtfs data: %v", err)
	}

	routes, err := getRoutes(*staticGtfs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse routes: %v", err)
	}

	vehicles, err := getVehicles(*staticGtfs, *realtimeGtfs)
	if err != nil {
		return nil, fmt.Errorf("failed to load vehicles: %v", err)
	}

	routeIdVehicles := make(map[int]([]backend.Vehicle))
	for _, vehicle := range vehicles {
		routeIdVehicles[vehicle.RouteId] = append(routeIdVehicles[vehicle.RouteId], vehicle)
	}

	for i := range routes {
		routes[i].Vehicles = routeIdVehicles[routes[i].ID]
	}

	return routes, nil
}

func getRoutes(staticGtfs gtfs.Static) ([]backend.BusRoute, error) {
	var routes []backend.BusRoute

	for _, route := range staticGtfs.Routes {
		id, err := strconv.Atoi(route.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to parse route ID: %v", err)
		}

		routes = append(routes, backend.BusRoute{
			ID:         id,
			SortIdx:    int(*route.SortOrder),
			Name:       route.Description,
			Code:       route.ShortName,
			Color:      route.Color,
			Directions: deriveRouteDirections(route, staticGtfs),
		})
	}

	return routes, nil
}

func deriveRouteDirections(route gtfs.Route, staticGtfs gtfs.Static) []backend.BusRouteDirection {
	directionTrips := getDirectionTrips(route, staticGtfs)

	var directions []backend.BusRouteDirection

	for directionId, trips := range directionTrips {
		stops := getStops(trips)
		polylines := getPolylines(trips)
		directions = append(directions, backend.BusRouteDirection{
			ID:        convertStaticDirectionId(directionId),
			Polylines: polylines,
			Stops:     convertStaticStops(stops),
		})
	}

	return directions
}

func convertStaticDirectionId(id gtfs.DirectionID) backend.BusRouteDirectionID {
	switch id {
	case gtfs.DirectionID_True:
		return backend.BusRouteDirectionIDOutbound
	case gtfs.DirectionID_False:
		return backend.BusRouteDirectionIDInbound
	default:
		return backend.BusRouteDirectionIDUnspecified
	}
}

func convertStaticStops(stops []gtfs.Stop) []backend.BusRouteDirectionStopsItem {
	var backendStops []backend.BusRouteDirectionStopsItem

	for _, stop := range stops {
		backendStops = append(backendStops, backend.BusRouteDirectionStopsItem{
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
