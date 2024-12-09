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

		staticTrip := matchingTrip(tripId.ID, staticGtfs)
		if staticTrip == nil {
			return nil, fmt.Errorf("failed to find a matching trip given ID: %s", tripId.ID)
		}

		nextStopTime := nextStopTime(vehicle, *staticTrip)
		if nextStopTime == nil {
			return nil, fmt.Errorf("failed to find a next stop time for vehicle ID: %d", id)
		}
		nextStop := nextStopTime.Stop.Name

		lastStopTime := stopTimeBefore(*staticTrip, *nextStopTime)
		var lastStop backend.NilString
		if lastStopTime == nil {
			lastStop = backend.NilString{Null: true}
		} else {
			lastStop = backend.NewNilString(lastStopTime.Stop.Name)
		}

		vehicles = append(vehicles, backend.Vehicle{
			ID:            id,
			RouteId:       routeId,
			Name:          vehicleId.Label,
			DirectionId:   convertStaticDirectionId(tripId.DirectionID),
			Heading:       int(*vehicle.Position.Bearing),
			Longitude:     float64(*vehicle.Position.Longitude),
			Latitude:      float64(*vehicle.Position.Latitude),
			NextStop:      nextStop,
			LastStop:      lastStop,
			DisplayStatus: vehicle.CurrentStatus.String(),
			LastUpdated:   *vehicle.Timestamp,
		})
	}

	return vehicles, nil
}

func matchingTrip(tripId string, staticGtfs gtfs.Static) *gtfs.ScheduledTrip {
	for _, trip := range staticGtfs.Trips {
		if trip.ID == tripId {
			return &trip
		}
	}
	return nil
}

func nextStopTime(vehicle gtfs.Vehicle, staticTrip gtfs.ScheduledTrip) *gtfs.ScheduledStopTime {
	stopId := vehicle.StopID

	for _, stopTime := range staticTrip.StopTimes {
		if stopTime.Stop.Id == *stopId {
			return &stopTime
		}
	}

	return nil
}

func stopTimeBefore(staticTrip gtfs.ScheduledTrip, stopTime gtfs.ScheduledStopTime) *gtfs.ScheduledStopTime {
	targetSequence := stopTime.StopSequence - 1

	for _, stopTime := range staticTrip.StopTimes {
		if stopTime.StopSequence == targetSequence {
			return &stopTime
		}
	}

	return nil
}

func nillableStop(stop *gtfs.Stop) backend.NilString {
	if stop == nil {
		return backend.NilString{Null: true}
	} else {
		return backend.NewNilString(stop.Name)
	}
}

// func getLastStop(vehicle gtfs.Vehicle, staticGtfs gtfs.Static) string {
//   for _, trip := range staticGtfs.Trips {
//     if trip.ID == tripId {
//       trip.
//
//     }
//   }
// }
