// Package convex loads additional cornell dining content from the Convex admin panel.
package convex

// QueryRequest is the body sent to the Convex HTTP query API.
type QueryRequest struct {
	Path   string `json:"path"`
	Args   any    `json:"args"`
	Format string `json:"format"`
}

// QueryResponse is the envelope returned by POST /api/query for the eateries:getAll function.
type QueryResponse struct {
	Status       string   `json:"status"`
	Value        []Eatery `json:"value"`
	ErrorMessage string   `json:"errorMessage,omitempty"`
	LogLines     []string `json:"logLines"`
}

// Eatery is the per-location record returned by eateries:getAll.
type Eatery struct {
	ID             string         `json:"_id"`
	Name           string         `json:"name"`
	NameShort      string         `json:"nameShort"`
	LocationCode   string         `json:"locationCode"`
	ImagePath      *string        `json:"imagePath"`
	Latitude       float64        `json:"latitude"`
	Longitude      float64        `json:"longitude"`
	Region         string         `json:"region"`
	PayMethods     []string       `json:"payMethods"`
	Categories     []string       `json:"categories"`
	Hours          []Hours        `json:"hours"`
	NextWeekEvents NextWeekEvents `json:"nextWeekEvents"`
	AllWeekMenu    []MenuCategory `json:"allWeekMenu"`
}

// Hours is a single open/close window expressed as RFC3339 strings.
type Hours struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// NextWeekEvents holds per-weekday event slices for the upcoming Mon–Sun week.
// A JSON null on any day means the location is closed that day.
type NextWeekEvents struct {
	Monday    []EateryEvent `json:"monday"`
	Tuesday   []EateryEvent `json:"tuesday"`
	Wednesday []EateryEvent `json:"wednesday"`
	Thursday  []EateryEvent `json:"thursday"`
	Friday    []EateryEvent `json:"friday"`
	Saturday  []EateryEvent `json:"saturday"`
	Sunday    []EateryEvent `json:"sunday"`
}

// EateryEvent is one service period within a day (e.g. Breakfast, Lunch).
type EateryEvent struct {
	Start          string         `json:"start"`
	End            string         `json:"end"`
	Name           string         `json:"name"`
	MenuCategories []MenuCategory `json:"menuCategories"`
}

// MenuCategory is a named group of menu items.
type MenuCategory struct {
	Name  string     `json:"name"`
	Items []MenuItem `json:"items"`
}

// MenuItem is a single dish or product.
type MenuItem struct {
	Name    string `json:"name"`
	Healthy bool   `json:"healthy"`
}
