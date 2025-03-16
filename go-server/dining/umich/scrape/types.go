package scrape

import (
	"fmt"
	"time"
)

type Eatery struct {
	Name           string
	Address        string
	Phone          string
	Email          string
	PaymentMethods []string
	Hours          []Hours // maps e.g. "Breakfast"
	Menus          []Menu
}

type Hours struct {
	Name      string // e.g. "Breakfast", "Lunch", "Dinner"
	StartTime time.Time
	EndTime   time.Time
}

type Menu struct {
	Name       string // e.g. "Breakfast", "Lunch", "Dinner"
	Categories []MenuCategory
}

type MenuCategory struct {
	Title string
	Items []MenuItem
}

type MenuItem struct {
	Name          string
	Traits        []string // e.g., "Vegetarian", "Gluten Free"
	NutrientLevel string   // e.g., "Nutrient Dense High"
	CarbonLevel   string   // e.g., "Carbon Footprint Low"
	Allergens     []string // e.g., "milk", "eggs"
	IsHalal       bool
	IsKosher      bool
	IsVegan       bool
	IsVegetarian  bool
	IsGlutenFree  bool
	IsSpicy       bool
}

// NewEatery creates a new Eatery with initialized fields
func NewEatery(name, address, phone, email string) *Eatery {
	return &Eatery{
		Name:           name,
		Address:        address,
		Phone:          phone,
		Email:          email,
		PaymentMethods: []string{},
		Hours:          []Hours{},
		Menus:          []Menu{},
	}
}

// AddPaymentMethod adds a payment method to the eatery
func (e *Eatery) addPaymentMethod(method string) {
	e.PaymentMethods = append(e.PaymentMethods, method)
}

func (e *Eatery) addHours(name string, startTime, endTime time.Time) {
	h := Hours{
		Name:      name,
		StartTime: startTime,
		EndTime:   endTime,
	}
	e.Hours = append(e.Hours, h)
}

func (e *Eatery) addMenu(name string) *Menu {
	m := Menu{
		Name:       name,
		Categories: []MenuCategory{},
	}
	e.Menus = append(e.Menus, m)
	return &e.Menus[len(e.Menus)-1]
}

func (m *Menu) addCategory(title string) *MenuCategory {
	category := MenuCategory{
		Title: title,
		Items: []MenuItem{},
	}
	m.Categories = append(m.Categories, category)
	return &m.Categories[len(m.Categories)-1]
}

// adds a new menu item to a station
func (mc *MenuCategory) AddMenuItem(name string, traits []string) *MenuItem {
	item := MenuItem{
		Name:      name,
		Traits:    traits,
		Allergens: []string{},
	}

	// process traits to fill in appropriate fields
	for _, trait := range traits {
		switch trait {
		case "Halal":
			item.IsHalal = true
		case "Kosher":
			item.IsKosher = true
		case "Vegan":
			item.IsVegan = true
		case "Vegetarian":
			item.IsVegetarian = true
		case "Gluten Free":
			item.IsGlutenFree = true
		case "Spicy":
			item.IsSpicy = true
		case "Nutrient Dense Low":
			item.NutrientLevel = "Low"
		case "Nutrient Dense Low Medium":
			item.NutrientLevel = "Low Medium"
		case "Nutrient Dense Medium":
			item.NutrientLevel = "Medium"
		case "Nutrient Dense Medium High":
			item.NutrientLevel = "Medium High"
		case "Nutrient Dense High":
			item.NutrientLevel = "High"
		case "Carbon Footprint Low":
			item.CarbonLevel = "Low"
		case "Carbon Footprint Medium":
			item.CarbonLevel = "Medium"
		case "Carbon Footprint High":
			item.CarbonLevel = "High"
		}
	}

	mc.Items = append(mc.Items, item)
	return &mc.Items[len(mc.Items)-1]
}

// adds an allergen to a menu item
func (m *MenuItem) AddAllergen(allergen string) {
	m.Allergens = append(m.Allergens, allergen)
}

// FormatMenu returns a formatted string representation of the full menu
func (e *Eatery) FormatMenu() string {
	var result string

	result = "MENU:\n"
	for _, mp := range e.Menus {
		result += fmt.Sprintf("\n=== %s ===\n", mp.Name)

		for _, station := range mp.Categories {
			result += fmt.Sprintf("\n## %s ##\n", station.Title)

			for _, item := range station.Items {
				var traits string
				if len(item.Traits) > 0 {
					traits = fmt.Sprintf(" (%s)", item.Traits[0])
					for i := 1; i < len(item.Traits); i++ {
						traits += fmt.Sprintf(", %s", item.Traits[i])
					}
				}
				result += fmt.Sprintf("* %s%s\n", item.Name, traits)
			}
		}
	}

	return result
}

// Summary prints a summary of the eatery
func (e *Eatery) Summary() string {
	var result string

	result = fmt.Sprintf("==== %s ====\n", e.Name)
	result += fmt.Sprintf("Address: %s\n", e.Address)
	result += fmt.Sprintf("Phone: %s\n", e.Phone)
	result += fmt.Sprintf("Email: %s\n\n", e.Email)

	result += "Payment Methods Accepted:\n"
	for _, method := range e.PaymentMethods {
		result += fmt.Sprintf("- %s\n", method)
	}

	result += "\nMeal Periods:\n"
	for _, h := range e.Hours {
		result += fmt.Sprintf("- %s: %s - %s\n", h.Name, h.StartTime, h.EndTime)
	}

	return result
}
