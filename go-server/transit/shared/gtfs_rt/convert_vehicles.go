package gtfs_rt

import (
	"log"
	"sort"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
	"github.com/jamespfennell/gtfs"
)

func ConvertVehicle(vehicle gtfs.Vehicle, staticGtfs gtfs.Static, realtimeGtfs gtfs.Realtime) api.Vehicle {
	vehicleID := vehicle.GetID()
	tripID := vehicle.GetTrip().ID
	trip := utils.Find(staticGtfs.Trips, func(t gtfs.ScheduledTrip) bool {
		return t.ID == tripID.ID
	})

	apiVehicle := api.Vehicle{
		ID:            api.NewStringVehicleID(vehicleID.ID),
		LastUpdated:   *vehicle.Timestamp,
		DisplayStatus: "",
		LastStop:      api.NilString{Null: true},
	}

	if vehicle.Position.Bearing == nil {
		log.Printf("failed to find bearing for vehicle ID: %s", vehicleID)
	} else {
		apiVehicle.Heading = int(*vehicle.Position.Bearing)
	}
	if vehicle.Position.Latitude == nil {
		log.Printf("failed to find latitude for vehicle ID: %s", vehicleID)
	} else {
		apiVehicle.Latitude = float64(*vehicle.Position.Latitude)
	}
	if vehicle.Position.Longitude == nil {
		log.Printf("failed to find longitude for vehicle ID: %s", vehicleID)
	} else {
		apiVehicle.Longitude = float64(*vehicle.Position.Longitude)
	}
	if vehicle.Position.Speed == nil {
		apiVehicle.Speed = 0
	} else {
		apiVehicle.Speed = float64(*vehicle.Position.Speed)
	}

	if trip == nil {
		log.Printf("failed to find trip for vehicle ID: %s", vehicleID)
		if tripID.RouteID == "" {
			log.Printf("failed to find route ID for vehicle ID: %s", vehicleID)
		} else {
			apiVehicle.RouteId = api.NewStringVehicleRouteId(tripID.RouteID)
		}
	} else {
		apiVehicle.RouteId = api.NewStringVehicleRouteId(trip.Route.Id)
		apiVehicle.Direction = trip.DirectionId.String()
		apiVehicle.Destination = trip.Headsign
	}

	nextStop := getStop(vehicle, staticGtfs)
	if nextStop != nil {
		apiVehicle.Destination = nextStop.Name
		if trip != nil {
			lastStop := getStopBefore(*nextStop, *trip)
			if lastStop != nil {
				apiVehicle.LastStop = api.NewNilString(lastStop.Name)
			}
		}
	} else {
		log.Printf("failed to find next stop for vehicle ID: %s", vehicleID)
	}

	return apiVehicle
}

func getStop(vehicle gtfs.Vehicle, staticGtfs gtfs.Static) *gtfs.Stop {
	if vehicle.StopID == nil {
		return nil
	}

	for _, stop := range staticGtfs.Stops {
		if stop.Id == *vehicle.StopID {
			return &stop
		}
	}
	return nil
}

func getStopBefore(stop gtfs.Stop, trip gtfs.ScheduledTrip) *gtfs.Stop {
	sort.Slice(trip.StopTimes, func(i, j int) bool {
		return trip.StopTimes[i].StopSequence < trip.StopTimes[j].StopSequence
	})

	for i, stopTime := range trip.StopTimes[1:] {
		if stopTime.Stop.Id == stop.Id {
			return trip.StopTimes[i].Stop
		}
	}
	return nil
}
