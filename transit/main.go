// package main
package transit

//
// import (
// 	"fmt"
//
// 	"github.com/amit7itz/goset"
// 	"github.com/benkoppe/bear-trak-backend/transit/external_gtfs"
// 	"github.com/jamespfennell/gtfs"
// 	"github.com/twpayne/go-polyline"
// )
//
// func main() {
// 	staticGtfs := external_gtfs.GetStaticGtfs("https://s3.amazonaws.com/tcat-gtfs/tcat-ny-us.zip")
// 	fmt.Println(staticGtfs)
//
// 	realtimeUrls := external_gtfs.RealtimeUrls{
// 		Alerts:           "https://realtimetcatbus.availtec.com/InfoPoint/GTFS-Realtime.ashx?&Type=Alert",
// 		VehiclePositions: "https://realtimetcatbus.availtec.com/InfoPoint/GTFS-Realtime.ashx?&Type=VehiclePosition",
// 		TripUpdates:      "https://realtimetcatbus.availtec.com/InfoPoint/GTFS-Realtime.ashx?&Type=TripUpdate",
// 	}
//
// 	realtimeData, _ := external_gtfs.GetRealtimeGtfs(realtimeUrls)
// 	fmt.Println(realtimeData)
//
// 	fmt.Println("\n")
//
// 	for _, route := range staticGtfs.Routes {
// 		fmt.Println(route.LongName)
// 		fmt.Println(route.Id)
//
// 		// directionTrips := getDirectionTrips(route, staticGtfs)
//
// 		// for _, trips := range directionTrips {
// 		// stops := getStops(trips)
// 		// fmt.Println(len(stops))
// 		// polylines := getPolylines(trips)
// 		// fmt.Println(polylines)
// 		// }
// 	}
// }
//
// func getDirectionTrips(route gtfs.Route, staticGtfs *gtfs.Static) map[gtfs.DirectionID][]gtfs.ScheduledTrip {
// 	var routeTrips []gtfs.ScheduledTrip
//
// 	for _, trip := range staticGtfs.Trips {
// 		if *trip.Route == route {
// 			routeTrips = append(routeTrips, trip)
// 		}
// 	}
//
// 	directionMappedTrips := make(map[gtfs.DirectionID][]gtfs.ScheduledTrip)
//
// 	for _, trip := range routeTrips {
// 		directionMappedTrips[trip.DirectionId] = append(directionMappedTrips[trip.DirectionId], trip)
// 	}
//
// 	return directionMappedTrips
// }
//
// func getStops(trips []gtfs.ScheduledTrip) []gtfs.Stop {
// 	stops := goset.NewSet[gtfs.Stop]()
//
// 	for _, trip := range trips {
// 		for _, stopTime := range trip.StopTimes {
// 			if stopTime.Stop != nil {
// 				stops.Add(*stopTime.Stop)
// 			}
// 		}
// 	}
//
// 	return stops.Items()
// }
//
// func getPolylines(trips []gtfs.ScheduledTrip) []string {
// 	var polylines []string
//
// 	for _, trip := range trips {
// 		shape := trip.Shape
// 		if shape == nil {
// 			continue
// 		}
//
// 		var coords [][]float64
//
// 		for _, point := range shape.Points {
// 			coords = append(coords, []float64{point.Latitude, point.Longitude})
// 		}
//
// 		if len(coords) < 2 {
// 			continue
// 		}
//
// 		line := string(polyline.EncodeCoords(coords))
// 		polylines = append(polylines, line)
// 	}
//
// 	return polylines
// }
