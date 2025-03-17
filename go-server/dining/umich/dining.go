package umich

import (
	"fmt"
	"log"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/dining/shared"
	"github.com/benkoppe/bear-trak-backend/go-server/dining/umich/scrape"
	"github.com/benkoppe/bear-trak-backend/go-server/dining/umich/static"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

type Cache = scrape.Cache

func InitCache(baseUrl string) Cache {
	return scrape.InitCache(baseUrl)
}

func Get(
	scrapeCache Cache,
) ([]api.Eatery, error) {
	cacheResponse, err := scrapeCache.Get()
	if err != nil {
		return nil, fmt.Errorf("error fetching external data: %w", err)
	}

	var eateries []api.Eatery
	for static, scraped := range cacheResponse {
		if len(scraped) == 0 {
			log.Printf("no scraped data for eatery %d", static.ID)
			continue
		}

		newEatery := convertScraped(*static, scraped)
		eateries = append(eateries, newEatery)
	}

	return eateries, nil
}

func convertScraped(static static.Eatery, scraped []scrape.Eatery) api.Eatery {
	events := convertScrapedEvents(scraped)

	firstScraped := scraped[0]

	return api.Eatery{
		ID:             static.ID,
		Name:           firstScraped.Name,
		NameShort:      firstScraped.Name,
		ImagePath:      utils.ImageNameToPath("dining/umich", static.ImageName),
		Latitude:       static.Location.Latitude,
		Longitude:      static.Location.Longitude,
		Location:       firstScraped.Address,
		Hours:          convertScrapedHours(scraped),
		Region:         static.Region,
		PayMethods:     convertScrapedPayMethods(firstScraped),
		Categories:     convertStaticCategories(static),
		NextWeekEvents: shared.SelectNextWeekEvents(events),
	}
}

func convertStaticCategories(static static.Eatery) []api.EateryCategoriesItem {
	var categories []api.EateryCategoriesItem

	for _, category := range static.Categories {
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

func convertScrapedPayMethods(scraped scrape.Eatery) []string {
	var payMethods []string

	for _, payMethod := range scraped.PaymentMethods {
		switch payMethod {
		case "Blue Bucks":
			payMethods = append(payMethods, "blueBucks")
		case "Dining Dollars":
			payMethods = append(payMethods, "diningDollars")
		default:
			continue
		}
	}

	return payMethods
}

func convertScrapedHours(scraped []scrape.Eatery) []api.Hours {
	var hours []api.Hours
	for _, eatery := range scraped {
		for _, scrapeHours := range eatery.Hours {
			h := api.Hours{
				Start: scrapeHours.StartTime,
				End:   scrapeHours.EndTime,
			}
			hours = append(hours, h)
		}
	}
	return hours
}

func convertScrapedEvents(scraped []scrape.Eatery) []api.EateryEvent {
	var events []api.EateryEvent

	for _, eatery := range scraped {
		for _, menu := range eatery.Menus {
			hours := getHours(eatery, menu.Name)
			if hours == nil {
				log.Printf("no hours for menu %s", menu.Name)
				continue
			}

			event := convertScrapedEvent(menu, *hours)
			events = append(events, event)
		}
	}
	return events
}

func convertScrapedEvent(menu scrape.Menu, hours scrape.Hours) api.EateryEvent {
	var categories []api.EateryMenuCategory
	for _, category := range menu.Categories {

		var items []api.EateryMenuCategoryItemsItem
		for _, item := range category.Items {

			i := api.EateryMenuCategoryItemsItem{
				Name:    item.Name,
				Healthy: false,
			}
			items = append(items, i)

		}
		cat := api.EateryMenuCategory{
			Name:  category.Title,
			Items: items,
		}
		categories = append(categories, cat)
	}

	return api.EateryEvent{
		Start:          hours.StartTime,
		End:            hours.EndTime,
		MenuCategories: categories,
	}
}

func getHours(eatery scrape.Eatery, name string) *scrape.Hours {
	for _, hours := range eatery.Hours {
		if hours.Name == name {
			return &hours
		}
	}
	return nil
}
