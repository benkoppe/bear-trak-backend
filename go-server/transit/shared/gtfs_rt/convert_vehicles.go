package gtfs_rt

import (
	"log"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
	"github.com/jamespfennell/gtfs"
)

func ConvertVehicle(vehicle gtfs.Vehicle, staticGtfs gtfs.Static, realtimeGtfs gtfs.Realtime) api.Vehicle {
	vehicleId := vehicle.GetID()
	tripId := vehicle.GetTrip().ID
	trip := utils.Find(staticGtfs.Trips, func(t gtfs.ScheduledTrip) bool {
		return t.ID == tripId.ID
	})

	if trip == nil {
		log.Printf("failed to find trip for vehicle ID: %s", vehicleId)
	}

	bearing := *vehicle.Position.Bearing
	latitude := *vehicle.Position.Latitude
	longitude := *vehicle.Position.Longitude
	var speed float32 = 0
	if vehicle.Position.Speed != nil {
		speed = *vehicle.Position.Speed
	}

	lastStop := api.NilString{Null: true}
	directionId := trip.DirectionId.String()
	destination := trip.Headsign

	return api.Vehicle{
		ID:            api.NewStringVehicleID(vehicleId.ID),
		RouteId:       api.NewStringVehicleRouteId(trip.Route.Id),
		Direction:     directionId,
		Heading:       int(bearing),
		Longitude:     float64(longitude),
		Latitude:      float64(latitude),
		Speed:         float64(speed),
		Destination:   destination,
		DisplayStatus: "",
		LastStop:      lastStop,
		LastUpdated:   *vehicle.Timestamp,
	}
}
