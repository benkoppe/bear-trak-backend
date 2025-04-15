package static

type Eatery struct {
	ID                 int             `json:"id"`
	HouseDinnerSubject *string         `json:"houseDinnerSubject,omitempty"`
	AllWeekMenu        *[]MenuCategory `json:"allWeekMenu,omitempty"`
}

type MenuCategory struct {
	Category string     `json:"category"`
	Items    []MenuItem `json:"items"`
}

type MenuItem struct {
	Item    string `json:"item"`
	Healthy bool   `json:"healthy"`
}
