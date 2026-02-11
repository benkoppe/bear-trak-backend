// Package external loads external umich dining content.
package external

import "time"

type Address struct {
	City       string `json:"city"`
	PostalCode string `json:"postalCode"`
	Street1    string `json:"street1"`
	State      string `json:"state"`
	Street2    string `json:"street2"`
}

type Contact struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type Location struct {
	OfficialBuildingID    int     `json:"officialbuildingid"`
	Image                 string  `json:"image"`
	Address               Address `json:"address"`
	BuildingPreferredName string  `json:"buildingpreferredname"`
	Lng                   float64 `json:"lng"`
	DisplayName           string  `json:"displayName"`
	Restricted            bool    `json:"restricted"`
	Campus                string  `json:"campus"`
	Contact               Contact `json:"contact"`
	Name                  string  `json:"name"`
	Type                  string  `json:"type"`
	Lat                   float64 `json:"lat"`
}

type MealHoursResponse struct {
	Meal  []MealItem  `json:"meal"`
	Hours []EventHour `json:"hours"`
}

type MealItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Hours       string `json:"hours"`
	HasMenu     bool   `json:"hasMenu"`
}

type EventHour struct {
	EventDayStart    string    `json:"event_day_start"`
	EventTimeEnd     time.Time `json:"event_time_end"`
	EventDescription any       `json:"event_description"`
	EventMapLink     string    `json:"event_maplink"`
	EventTimeStart   time.Time `json:"event_time_start"`
	EventDayEnd      string    `json:"event_day_end"`
	EventTitle       string    `json:"event_title"`
}

type MenuResponse struct {
	Menu Menu `json:"menu"`
}

type Menu struct {
	Name     string         `json:"name"`
	Category []MenuCategory `json:"category"`
}

type MenuCategory struct {
	Name     string     `json:"name"`
	MenuItem []MenuItem `json:"menuItem"`
}

type MenuItem struct {
	Name      string   `json:"name"`
	Attribute []string `json:"attribute"`
	Allergens []string `json:"allergens"`
}

type LocationDiningData struct {
	Location Location          `json:"location"`
	Days     []LocationDayData `json:"days"`
}

type LocationDayData struct {
	Date  string        `json:"date"`
	Meals []DayMealData `json:"meals"`
}

type DayMealData struct {
	MealName string      `json:"mealName"`
	Meal     MealItem    `json:"meal"`
	Hours    []EventHour `json:"hours"`
	Menu     Menu        `json:"menu"`
}
