package static

type Eatery struct {
	ID         int      `json:"id"`
	ScrapePath string   `json:"scrapePath"`
	ImageName  string   `json:"imageName"`
	Region     string   `json:"region"`
	Location   Location `json:"location"`
	Categories []string `json:"categories"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
