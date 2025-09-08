// Package harvard loads harvard dining content.
package harvard

import (
	"fmt"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/dining/harvard/external"
	"github.com/benkoppe/bear-trak-backend/go-server/dining/harvard/static"
	"github.com/benkoppe/bear-trak-backend/go-server/dining/shared"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
	"github.com/benkoppe/bear-trak-backend/go-server/utils/timeutils"
)

type Cache = external.Caches

func InitCache(baseURL, apiKey string) Cache {
	return external.InitCaches(baseURL, apiKey)
}

func Get(
	externalCache Cache,
) ([]api.Eatery, error) {
	externalRecipes, err := externalCache.RecipesCache.Get()
	if err != nil {
		return nil, fmt.Errorf("error fetching external recipes: %w", err)
	}

	staticEateries := static.GetEateries()

	var eateries []api.Eatery

	for _, staticEatery := range staticEateries {
		newEatery := convertStatic(staticEatery, externalRecipes)
		eateries = append(eateries, newEatery)
	}

	return eateries, nil
}

func convertStatic(static static.Eatery, externalRecipes []external.Recipe) api.Eatery {
	eatery := api.Eatery{
		ID:         static.ID,
		Name:       static.Name,
		NameShort:  static.Name,
		ImagePath:  utils.ImageNameToPath("dining/harvard", static.ImageName),
		Latitude:   static.Location.Latitude,
		Longitude:  static.Location.Longitude,
		Region:     static.Region,
		PayMethods: []string{},
		Categories: shared.ConvertStaticCategories(static.Categories),
		Hours:      static.WeekHours.CreateFutureHours(),
	}

	if static.APINumber != nil {
		eatery.NextWeekEvents = shared.SelectNextWeekEvents(collectFutureEvents(static, externalRecipes))
	}

	if static.AllWeekMenu != nil {
		eatery.AllWeekMenu = convertStaticMenu(*static.AllWeekMenu)
	}

	return eatery
}

func collectFutureEvents(static static.Eatery, externalRecipes []external.Recipe) []api.EateryEvent {
	if static.APINumber == nil {
		return nil
	}
	locationNumber := *static.APINumber
	weekHours := static.WeekHours

	est := timeutils.LoadEST()
	now := time.Now().In(est)

	var events []api.EateryEvent

	for i := range [7]int{} {
		date := now.AddDate(0, 0, i)
		hours := weekHours.GetHours(date)

		for _, hour := range hours {
			if hour.MealNumber != nil {
				targetServeDate := date.Format("01/02/2006")
				var recipes []external.Recipe

				for _, recipe := range externalRecipes {
					if recipe.LocationNumber == locationNumber && recipe.MealNumber == *hour.MealNumber && recipe.ServeDate == targetServeDate {
						recipes = append(recipes, recipe)
					}
				}

				menu := convertExternalRecipesToMenu(recipes)
				start, err := hour.Open.ToDate(date)
				if err != nil {
					fmt.Printf("error converting open time: %v\n", err)
					continue
				}
				end, err := hour.Close.ToDate(date)
				if err != nil {
					fmt.Printf("error converting close time: %v\n", err)
					continue
				}
				event := api.EateryEvent{
					Start:          start,
					End:            end,
					MenuCategories: menu,
				}
				events = append(events, event)
			}
		}
	}

	return events
}

func convertExternalRecipesToMenu(recipes []external.Recipe) []api.EateryMenuCategory {
	grouped := make(map[string][]external.Recipe)

	for _, recipe := range recipes {
		name := recipe.MenuCategoryName
		if _, ok := grouped[name]; !ok {
			grouped[name] = []external.Recipe{}
		}
		grouped[name] = append(grouped[name], recipe)
	}

	var categories []api.EateryMenuCategory

	for name, recipes := range grouped {
		var items []api.EateryMenuCategoryItemsItem

		for _, recipe := range recipes {
			item := api.EateryMenuCategoryItemsItem{
				Name:    recipe.RecipeName,
				Healthy: false,
			}
			items = append(items, item)
		}

		category := api.EateryMenuCategory{
			Name:  name,
			Items: items,
		}
		categories = append(categories, category)
	}

	return categories
}

func convertStaticMenu(staticCategories []static.MenuCategory) []api.EateryMenuCategory {
	var categories []api.EateryMenuCategory

	for _, staticCategory := range staticCategories {
		var items []api.EateryMenuCategoryItemsItem

		for _, staticItem := range staticCategory.Items {
			items = append(items, api.EateryMenuCategoryItemsItem{
				Name:    staticItem.Item,
				Healthy: staticItem.Healthy,
			})
		}

		categories = append(categories, api.EateryMenuCategory{
			Name:  staticCategory.Category,
			Items: items,
		})
	}

	return categories
}
