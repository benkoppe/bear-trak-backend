package transit

import (
	"fmt"
	"strconv"

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

		id, err := strconv.Atoi(vehicleId.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse vehicle ID as integer: %v", err)
		}

		routeId, err := strconv.Atoi(tripId.RouteID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse route ID as integer: %v", err)
		}

		nextStop := getStop(vehicle, staticGtfs)
		if nextStop == nil {
			return nil, fmt.Errorf("failed to find a next stop for vehicle ID: %d", id)
		}

		lastStop := getStopBefore(*nextStop, tripId.ID, staticGtfs)
		var lastStopName backend.NilString
		if lastStop == nil {
			lastStopName = backend.NilString{Null: true}
		} else {
			lastStopName = backend.NewNilString(lastStop.Name)
		}

		vehicles = append(vehicles, backend.Vehicle{
			ID:            id,
			RouteId:       routeId,
			Name:          vehicleId.Label,
			DirectionId:   convertStaticDirectionId(tripId.DirectionID),
			Heading:       int(*vehicle.Position.Bearing),
			Longitude:     float64(*vehicle.Position.Longitude),
			Latitude:      float64(*vehicle.Position.Latitude),
			NextStop:      nextStop.Name,
			LastStop:      lastStopName,
			DisplayStatus: vehicle.CurrentStatus.String(),
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

func getStopBefore(stop gtfs.Stop, tripId string, staticGtfs gtfs.Static) *gtfs.Stop {
	trip := getTrip(tripId, staticGtfs)

	if trip == nil {
		return nil
	}

	var sequence int = -1

	for _, stopTime := range trip.StopTimes {
		if *stopTime.Stop == stop {
			sequence = stopTime.StopSequence - 1
			break
		}
	}

	if sequence == -1 {
		return nil
	}

	for _, stopTime := range trip.StopTimes {
		if stopTime.StopSequence == sequence {
			return stopTime.Stop
		}
	}
	return nil
}

func getTrip(tripId string, staticGtfs gtfs.Static) *gtfs.ScheduledTrip {
	for _, trip := range staticGtfs.Trips {
		if trip.ID == tripId {
			return &trip
		}
	}
	return nil
}
