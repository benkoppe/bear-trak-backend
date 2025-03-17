package dining

import (
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/dining/cornell/external"
	"github.com/benkoppe/bear-trak-backend/go-server/dining/shared"

	"github.com/benkoppe/bear-trak-backend/go-server/dining/cornell/static"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
	"golang.org/x/text/unicode/norm"
)

type Cache = external.Cache

func InitCache(url string) Cache {
	return external.InitCache(url)
}

func Get(
	externalCache Cache,
) ([]api.Eatery, error) {
	externalResponse, err := externalCache.Get()
	if err != nil {
		return nil, fmt.Errorf("error fetching external data: %w", err)
	}

	if externalResponse == nil {
		return nil, fmt.Errorf("fetched nil external data")
	}

	externalEateries := externalResponse.Data.Eateries

	eateries := make([]api.Eatery, len(externalEateries))

	for i, externalEatery := range externalEateries {
		eateries[i] = convertExternal(externalEatery)
	}

	staticEateries := static.GetEateries()
	eateries = appendStaticMenus(eateries, staticEateries)

	return eateries, nil
}

func convertExternal(external external.Eatery) api.Eatery {
	events := convertExternalEvents(external)

	return api.Eatery{
		ID:             external.ID,
		Name:           external.Name,
		NameShort:      external.NameShort,
		ImagePath:      getImagePath(external),
		Latitude:       external.Latitude,
		Longitude:      external.Longitude,
		Location:       external.Location,
		Hours:          hoursFromEvents(events),
		Region:         convertExternalRegion(external),
		PayMethods:     convertExternalPayMethods(external),
		Categories:     convertExternalCategories(external),
		NextWeekEvents: shared.SelectNextWeekEvents(events),
	}
}

func convertExternalEvents(external external.Eatery) []api.EateryEvent {
	var events []api.EateryEvent

	for _, operatingHours := range external.OperatingHours {
		for _, event := range operatingHours.Events {
			events = append(events, api.EateryEvent{
				Start:          event.StartTimestamp.ToTime(),
				End:            event.EndTimestamp.ToTime(),
				MenuCategories: convertExternalMenu(event),
			})
		}
	}

	// sort events by Start time
	sort.Slice(events, func(i, j int) bool {
		return events[i].Start.Before(events[j].Start)
	})

	return events
}

func convertExternalMenu(externalEvent external.Event) []api.EateryMenuCategory {
	sortMenuCategories(externalEvent.Menu)
	sortMenuItems(externalEvent.Menu)

	var categories []api.EateryMenuCategory

	for _, category := range externalEvent.Menu {
		var items []api.EateryMenuCategoryItemsItem

		for _, item := range category.Items {
			items = append(items, api.EateryMenuCategoryItemsItem{
				Name:    item.Item,
				Healthy: item.Healthy,
			})
		}

		categories = append(categories, api.EateryMenuCategory{
			Name:  category.Category,
			Items: items,
		})
	}

	return categories
}

func convertExternalRegion(external external.Eatery) string {
	switch external.CampusArea.Descrshort {
	case "Central":
		return "central"
	case "West":
		return "west"
	case "North":
		return "north"
	case "East":
		return "central"
	default:
		return "unknown"
	}
}

func convertExternalPayMethods(external external.Eatery) []string {
	var payMethods []string

	for _, payMethod := range external.PayMethods {
		switch payMethod.DescrShort {
		case "Meal Plan - Swipe":
			payMethods = append(payMethods, "swipes")
		case "Meal Plan - Debit":
			payMethods = append(payMethods, "bigRedBucks")
		case "Cash":
			payMethods = append(payMethods, "cash")
		case "Mobile Payments":
			payMethods = append(payMethods, "digitalWallet")
		case "Major Credit Cards":
			payMethods = append(payMethods, "creditCard")
		case "Cornell Card":
			payMethods = append(payMethods, "cornellCard")
		default:
			continue
		}
	}

	return payMethods
}

func convertExternalCategories(external external.Eatery) []api.EateryCategoriesItem {
	var categories []api.EateryCategoriesItem

	for _, eateryType := range external.EateryTypes {
		switch eateryType.Descr {
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

func hoursFromEvents(events []api.EateryEvent) []api.Hours {
	var hours []api.Hours

	est := utils.LoadEST()

	// convert to hours objects
	for _, event := range events {
		hours = append(hours, api.Hours{
			Start: event.Start.In(est),
			End:   event.End.In(est),
		})
	}

	// catch empty hours
	if len(hours) == 0 {
		return hours
	}

	// sort hours by Start time
	sort.Slice(hours, func(i, j int) bool {
		return hours[i].Start.Before(hours[j].Start)
	})

	// merge close start and end times
	var merged []api.Hours
	currentStart := hours[0].Start
	currentEnd := hours[0].End

	for _, hour := range hours {
		diff := hour.Start.Sub(currentEnd)
		if diff < 0 {
			diff = -diff
		}
		if diff <= 5*time.Minute {
			currentEnd = hour.End
			continue
		} else {
			merged = append(merged, api.Hours{
				Start: currentStart,
				End:   currentEnd,
			})
			currentStart = hour.Start
			currentEnd = hour.End
		}
	}

	// append the final values
	merged = append(merged, api.Hours{
		Start: currentStart,
		End:   currentEnd,
	})

	return merged
}

func sortMenuCategories(categories []external.MenuCategory) {
	priorityCategories := []string{
		"Chef's Table",
		"Chef's Table - Sides",
		"Grill",
		"Wok/Asian Station",
	}

	sort.Slice(categories, func(i, j int) bool {
		lhs := categories[i]
		rhs := categories[j]

		containsLeft := slices.Contains(priorityCategories, lhs.Category)
		containsRight := slices.Contains(priorityCategories, rhs.Category)

		if containsLeft && containsRight {
			return slices.Index(priorityCategories, lhs.Category) < slices.Index(priorityCategories, rhs.Category)
		} else if containsLeft {
			return true
		} else if containsRight {
			return false
		} else {
			return lhs.SortIdx < rhs.SortIdx
		}
	})
}

func sortMenuItems(categories []external.MenuCategory) {
	for _, category := range categories {
		items := category.Items
		sort.Slice(items, func(i, j int) bool {
			return items[i].SortIdx < items[j].SortIdx
		})
	}
}

func getImagePath(external external.Eatery) string {
	name := external.Name
	// normalize to decomposed form (NFD)
	// this helps remove things like marks (accents)
	decomposed := norm.NFD.String(name)

	// convert to lowercase
	lowercased := strings.ToLower(decomposed)

	// filter characters
	var builder strings.Builder
	for _, r := range lowercased {
		// filter marks, and only let letters, numbers, and whitespace through
		if !unicode.IsMark(r) && (unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r)) {
			builder.WriteRune(r)
		}
	}
	stripped := builder.String()

	// regex pattern to match with whitespace
	re := regexp.MustCompile(`\s+`)

	// replace all whitespace with underscores
	imageName := re.ReplaceAllString(stripped, "_")
	return utils.ImageNameToPath("dining/cornell", imageName)
}

func appendStaticMenus(eateries []api.Eatery, staticEateries []static.Eatery) []api.Eatery {
	var converted []api.Eatery

	for _, eatery := range eateries {
		staticEatery := matchingStaticEatery(eatery, staticEateries)

		if staticEatery == nil {
			converted = append(converted, eatery)
			continue
		}

		if staticEatery.AllWeekMenu != nil {
			eatery.AllWeekMenu = convertStaticMenu(*staticEatery.AllWeekMenu)
		}

		converted = append(converted, eatery)
	}

	return converted
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

func matchingStaticEatery(eatery api.Eatery, staticEateries []static.Eatery) *static.Eatery {
	for _, staticEatery := range staticEateries {
		if staticEatery.ID == eatery.ID {
			return &staticEatery
		}
	}
	return nil
}
