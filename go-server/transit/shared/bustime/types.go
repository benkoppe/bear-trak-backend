package bustime

type Route struct {
	ID    string `json:"rt"`
	Name  string `json:"rtnm"`
	Color string `json:"rtclr"`
	Code  string `json:"rtdd"`
}

type Vehicle struct {
	ID              string        `json:"vid"`
	LastUpdated     TransitTime   `json:"tmstmp"`
	Latitude        Float64String `json:"lat"`
	Longitude       Float64String `json:"long"`
	Heading         IntString     `json:"hdg"`
	RouteID         string        `json:"rt"`
	Destination     string        `json:"des"`
	Speed           int           `json:"spd"`
	TripID          string        `json:"origtatripno"`
	OccupancyStatus string        `json:"psgld"`
}

type Prediction struct{}

func collectRouteIds(routes []Route) []string {
	routeIds := make([]string, len(routes))
	for i, route := range routes {
		routeIds[i] = route.ID
	}
	return routeIds
}
