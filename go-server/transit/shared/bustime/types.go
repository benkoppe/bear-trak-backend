package bustime

type bustimeResponse[T any] struct {
	Response T `json:"bustime-response"`
}

type routesResponse struct {
	Routes []Route `json:"routes"`
}

type Route struct {
	Id    string `json:"rt"`
	Name  string `json:"rtnm"`
	Color string `json:"rtclr"`
	Code  string `json:"rtdd"`
}

type vehiclesResponse struct {
	Vehicles []Vehicle `json:"vehicle"`
}

type Vehicle struct {
	Id              string        `json:"vid"`
	LastUpdated     TransitTime   `json:"tmstmp"`
	Latitude        Float64String `json:"lat"`
	Longitude       Float64String `json:"lon"`
	Heading         IntString     `json:"hdg"`
	RouteId         string        `json:"rt"`
	Destination     string        `json:"des"`
	Speed           int           `json:"spd"`
	TripId          string        `json:"origtatripno"`
	OccupancyStatus string        `json:"psgld"`
}

type Prediction struct{}

func collectRouteIds(routes []Route) []string {
	routeIds := make([]string, len(routes))
	for i, route := range routes {
		routeIds[i] = route.Id
	}
	return routeIds
}
