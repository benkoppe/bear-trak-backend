package external

import (
	"github.com/benkoppe/bear-trak-backend/go-server/utils/timeutils"
)

type EateryAPIResponse struct {
	Status string                `json:"status"`
	Data   EateryAPIResponseData `json:"data"`
}

type EateryAPIResponseData struct {
	Eateries []Eatery `json:"eateries"`
}

type Eatery struct {
	ID             int             `json:"id"`
	Slug           string          `json:"slug"`
	Name           string          `json:"name"`
	NameShort      string          `json:"nameshort"`
	About          string          `json:"about"`
	AboutShort     string          `json:"aboutshort"`
	CornellDining  bool            `json:"cornellDining"`
	OnlineOrdering bool            `json:"onlineOrdering"`
	OnlineOrderUrl *string         `json:"onlineOrderUrl"`
	ContactPhone   *string         `json:"contactPhone"`
	ContactEmail   *string         `json:"contactEmail"`
	CampusArea     CampusArea      `json:"campusArea"`
	Latitude       float64         `json:"latitude"`
	Longitude      float64         `json:"longitude"`
	Location       string          `json:"location"`
	OperatingHours []OperatingHour `json:"operatingHours"`
	EateryTypes    []EateryType    `json:"eateryTypes"`
	DiningCuisines []DiningCuisine `json:"diningCuisines"`
	PayMethods     []PayMethod     `json:"payMethods"`
}

type CampusArea struct {
	Descr      string `json:"descr"`
	Descrshort string `json:"descrshort"`
}

type OperatingHour struct {
	Date   string  `json:"date"`
	Status string  `json:"status"`
	Events []Event `json:"events"`
}

type Event struct {
	Descr          string              `json:"descr"`
	StartTimestamp timeutils.UnixTime `json:"startTimestamp"`
	EndTimestamp   timeutils.UnixTime `json:"endTimestamp"`
	Start          string              `json:"start"`
	End            string              `json:"end"`
	CalSummary     string              `json:"calSummary"`
	Menu           []MenuCategory      `json:"menu"`
}

type MenuCategory struct {
	Category string     `json:"category"`
	SortIdx  int        `json:"sortIdx"`
	Items    []MenuItem `json:"items"`
}

type MenuItem struct {
	Item    string `json:"item"`
	Healthy bool   `json:"healthy"`
	SortIdx int    `json:"sortIdx"`
}

type PayMethod struct {
	Descr      string `json:"descr"`
	DescrShort string `json:"descrshort"`
}

type EateryType struct {
	Descr      string `json:"descr"`
	DescSshort string `json:"descrshort"`
}

type DiningCuisine struct {
	Name      string  `json:"name"`
	NameShort string  `json:"nameshort"`
	Descr     *string `json:"descr"`
}
