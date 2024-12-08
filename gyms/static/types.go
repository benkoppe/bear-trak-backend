package static

type Gym struct {
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	ImageName  string      `json:"imageName"`
	ScrapeName string      `json:"scrapeName"`
	Location   Location    `json:"location"`
	Facilities []Facility  `json:"facilities"`
	Equipment  []Equipment `json:"equipment"`
	WeekHours  WeekHours   `json:"weekHours"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Facility struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type Equipment struct {
	Type  string   `json:"type"`
	Items []string `json:"items"`
}

type Hours struct {
	Open  string `json:"open"`
	Close string `json:"close"`
}

type WeekHours struct {
	Monday    []Hours `json:"monday"`
	Tuesday   []Hours `json:"tuesday"`
	Wednesday []Hours `json:"wednesday"`
	Thursday  []Hours `json:"thursday"`
	Friday    []Hours `json:"friday"`
	Saturday  []Hours `json:"saturday"`
	Sunday    []Hours `json:"sunday"`
}
