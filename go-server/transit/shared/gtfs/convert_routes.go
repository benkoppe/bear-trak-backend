package gtfs

import (
	"github.com/amit7itz/goset"
	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/jamespfennell/gtfs"
	"github.com/twpayne/go-polyline"
)

// ConvertRoute converts route into API routes
func ConvertRoute(route gtfs.Route, staticGtfs gtfs.Static) api.BusRoute {
	tripsMap := getDirectionTrips(route, staticGtfs)

	var polylines []string
	var directions []api.BusRouteDirection

	for directionID, trips := range tripsMap {
		polylines = append(polylines, getPolylines(trips)...)
		stops := getStops(trips)

		directions = append(directions, api.BusRouteDirection{
			Name:  ConvertDirectionID(directionID),
			Stops: convertStops(stops),
		})
	}

	apiRoute := api.BusRoute{
		ID:         api.NewStringBusRouteID(route.Id),
		Name:       route.Description,
		Code:       route.ShortName,
		Color:      route.Color,
		Directions: directions,
		Polylines:  polylines,
	}

	if route.SortOrder != nil {
		apiRoute.SortIdx = int(*route.SortOrder)
	}

	return apiRoute
}

func ConvertDirectionID(id gtfs.DirectionID) string {
	switch id {
	case gtfs.DirectionID_True:
		return "O"
	case gtfs.DirectionID_False:
		return "I"
	default:
		return "?"
	}
}

func convertStops(stops []gtfs.Stop) []api.BusRouteDirectionStopsItem {
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

// maps directions to their corresponding trips
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
