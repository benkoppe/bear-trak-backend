package external_map

type overlayResponse struct {
	Overlays []itemCategory `json:"overlays"`
}

type itemCategory struct {
	DOM_ID string `json:"DOM_ID"`
	Items  []Item `json:"items"`
}

type Item struct {
	Notes   string `json:"NOTES"`
	Name    string `json:"NAME"`
	Address string `json:"ADDRESS"`
	LatLng  LatLng `json:"LatLng"`
}

type LatLng struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}
