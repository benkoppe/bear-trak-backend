package transit

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	backend "github.com/benkoppe/bear-trak-backend/backend"
	"github.com/benkoppe/bear-trak-backend/transit/external_gtfs"
	"github.com/jamespfennell/gtfs"
)

func GetVehicles(staticUrl string, realtimeUrls external_gtfs.RealtimeUrls) ([]backend.Vehicle, error) {
	staticGtfs := external_gtfs.GetStaticGtfs(staticUrl)

	realtimeGtfs, err := external_gtfs.GetRealtimeGtfs(realtimeUrls)
	if err != nil {
		return nil, fmt.Errorf("failed to load realtime gtfs data: %v", err)
	}

	return getVehicles(*staticGtfs, *realtimeGtfs)
}

func getVehicles(staticGtfs gtfs.Static, realtimeGtfs gtfs.Realtime) ([]backend.Vehicle, error) {
	var vehicles []backend.Vehicle

	for _, vehicle := range realtimeGtfs.Vehicles {
		vehicleId := vehicle.GetID()
		tripId := vehicle.GetTrip().ID
		trip := getTrip(tripId, staticGtfs)

		id, err := strconv.Atoi(vehicleId.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse vehicle ID as integer: %v", err)
		}

		routeId, err := strconv.Atoi(tripId.RouteID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse route ID as integer: %v", err)
		}

		lastStop := getStop(vehicle, staticGtfs)
		if lastStop == nil {
			return nil, fmt.Errorf("failed to find a last stop for vehicle ID: %d", id)
		}

		directionId := tripId.DirectionID
		destination := backend.NilString{Null: true}
		nextStopName := backend.NilString{Null: true}

		if trip != nil {
			directionId = trip.DirectionId
			destination = backend.NewNilString(trip.Headsign)

			nextStop := getStopAfter(*lastStop, *trip)
			if nextStop != nil {
				nextStopName = backend.NewNilString(nextStop.Name)
			}
		}

		vehicles = append(vehicles, backend.Vehicle{
			ID:            id,
			RouteId:       routeId,
			Name:          vehicleId.Label,
			DirectionId:   convertStaticDirectionId(directionId),
			Heading:       int(*vehicle.Position.Bearing),
			Longitude:     float64(*vehicle.Position.Longitude),
			Latitude:      float64(*vehicle.Position.Latitude),
			Destination:   destination,
			LastStop:      lastStop.Name,
			NextStop:      nextStopName,
			DisplayStatus: capitalizeWords(vehicle.OccupancyStatus.String()),
			LastUpdated:   *vehicle.Timestamp,
		})
	}

	return vehicles, nil
}

func getStop(vehicle gtfs.Vehicle, staticGtfs gtfs.Static) *gtfs.Stop {
	for _, stop := range staticGtfs.Stops {
		if stop.Id == *vehicle.StopID {
			return &stop
		}
	}
	return nil
}

func getStopAfter(stop gtfs.Stop, trip gtfs.ScheduledTrip) *gtfs.Stop {
	sort.Slice(trip.StopTimes, func(i, j int) bool {
		return trip.StopTimes[i].StopSequence < trip.StopTimes[j].StopSequence
	})

	for i, stopTime := range trip.StopTimes {
		if *stopTime.Stop == stop && (i-1) >= 0 {
			return trip.StopTimes[i+1].Stop
		}
	}

	return nil
}

func getTrip(tripId gtfs.TripID, staticGtfs gtfs.Static) *gtfs.ScheduledTrip {
	for _, trip := range staticGtfs.Trips {
		if trip.ID == tripId.ID {
			return &trip
		}
	}
	return nil
}

func capitalizeWords(input string) string {
	words := strings.Fields(input)
	for i, word := range words {
		words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
	}
	return strings.Join(words, " ")
}
