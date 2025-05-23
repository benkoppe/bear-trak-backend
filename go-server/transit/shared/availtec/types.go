package availtec

import "github.com/benkoppe/bear-trak-backend/go-server/utils/time_utils"

type Route struct {
	Color             string    `json:"Color"`
	GoogleDescription string    `json:"GoogleDescription"`
	IsVisible         bool      `json:"IsVisible"`
	RouteAbbreviation string    `json:"RouteAbbreviation"`
	RouteId           int       `json:"RouteId"`
	SortOrder         int       `json:"SortOrder"`
	Vehicles          []Vehicle `json:"Vehicles"`
	// Messages                 []string  `json:"Messages"`
	DetourActiveMessageCount int `json:"DetourActiveMessageCount"`
	// Stops                    *string   `json:"Stops"`
	// RouteStops               *string   `json:"RouteStops"`
}

type Vehicle struct {
	Destination                string                   `json:"Destination"`
	Deviation                  int                      `json:"Deviation"`
	Direction                  string                   `json:"Direction"`
	DirectionLong              string                   `json:"DirectionLong"`
	DisplayStatus              string                   `json:"DisplayStatus"`
	StopId                     int                      `json:"StopId"`
	Heading                    int                      `json:"Heading"`
	LastStop                   string                   `json:"LastStop"`
	LastUpdated                time_utils.MicrosoftTime `json:"LastUpdated"`
	Latitude                   float64                  `json:"Latitude"`
	Longitude                  float64                  `json:"Longitude"`
	RouteId                    int                      `json:"RouteId"`
	Speed                      int                      `json:"Speed"`
	TripId                     int                      `json:"TripId"`
	VehicleId                  int                      `json:"VehicleId"`
	SeatingCapacity            *int                     `json:"SeatingCapacity"`
	TotalCapacity              *int                     `json:"TotalCapacity"`
	OccupancyStatusReportLabel string                   `json:"OccupancyStatusReportLabel"`
}
