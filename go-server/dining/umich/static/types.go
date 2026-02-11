package static

type Eatery struct {
	ID                  int      `json:"id"`
	OfficialBuildingID  int      `json:"officialBuildingId"`
	LocationDisplayName string   `json:"locationDisplayName"`
	ImageName           string   `json:"imageName"`
	Region              string   `json:"region"`
	Categories          []string `json:"categories"`
	PayMethods          []string `json:"payMethods"`
}
