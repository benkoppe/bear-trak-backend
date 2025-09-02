// Package shared includes all shared dining methods.
package shared

import (
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/utils/time_utils"
)

func SelectNextWeekEvents(events []api.EateryEvent) api.EateryNextWeekEvents {
	est := time_utils.LoadEST()
	now := time.Now().In(est)
	weekStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekEnd := weekStart.AddDate(0, 0, 7)

	result := api.EateryNextWeekEvents{}

	for _, event := range events {
		if event.Start.After(weekStart) && event.Start.Before(weekEnd) {
			switch event.Start.Weekday() {
			case time.Monday:
				result.Monday = append(result.Monday, event)
			case time.Tuesday:
				result.Tuesday = append(result.Tuesday, event)
			case time.Wednesday:
				result.Wednesday = append(result.Wednesday, event)
			case time.Thursday:
				result.Thursday = append(result.Thursday, event)
			case time.Friday:
				result.Friday = append(result.Friday, event)
			case time.Saturday:
				result.Saturday = append(result.Saturday, event)
			case time.Sunday:
				result.Sunday = append(result.Sunday, event)
			}
		}
	}

	return result
}

func ConvertStaticCategories(staticCategories []string) []api.EateryCategoriesItem {
	var categories []api.EateryCategoriesItem

	for _, category := range staticCategories {
		switch category {
		case "Convenience Store":
			categories = append(categories, api.EateryCategoriesItemConvenienceStore)
		case "Cafe":
			categories = append(categories, api.EateryCategoriesItemCafe)
		case "Dining Room":
			categories = append(categories, api.EateryCategoriesItemDiningRoom)
		case "Coffee Shop":
			categories = append(categories, api.EateryCategoriesItemCoffeeShop)
		case "Cart":
			categories = append(categories, api.EateryCategoriesItemCart)
		case "Food Court":
			categories = append(categories, api.EateryCategoriesItemFoodCourt)
		default:
			continue
		}
	}

	return categories
}
