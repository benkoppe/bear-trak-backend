package dining

import (
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"
	"unicode"

	backend "github.com/benkoppe/bear-trak-backend/backend"
	"github.com/benkoppe/bear-trak-backend/dining/external"
	"github.com/benkoppe/bear-trak-backend/utils"
	"golang.org/x/text/unicode/norm"
)

func Get(url string) ([]backend.Eatery, error) {
	externalResponse, err := external.FetchData(url)
	if err != nil {
		return nil, fmt.Errorf("Error fetching external data: %w", err)
	}

	if externalResponse == nil {
		return nil, fmt.Errorf("Fetched nil external data.")
	}

	externalEateries := externalResponse.Data.Eateries

	eateries := make([]backend.Eatery, len(externalEateries))

	for i, externalEatery := range externalEateries {
		eateries[i] = convertExternal(externalEatery)
	}

	return eateries, nil
}

func convertExternal(external external.Eatery) backend.Eatery {
	events := convertExternalEvents(external)

	return backend.Eatery{
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
		NextWeekEvents: selectNextWeekEvents(events),
	}
}

func convertExternalEvents(external external.Eatery) []backend.EateryEvent {
	var events []backend.EateryEvent

	for _, operatingHours := range external.OperatingHours {
		for _, event := range operatingHours.Events {
			events = append(events, backend.EateryEvent{
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

func convertExternalMenu(externalEvent external.Event) []backend.EateryMenuCategory {
	sortMenuCategories(externalEvent.Menu)
	sortMenuItems(externalEvent.Menu)

	var categories []backend.EateryMenuCategory

	for _, category := range externalEvent.Menu {
		var items []backend.EateryMenuCategoryItemsItem

		for _, item := range category.Items {
			items = append(items, backend.EateryMenuCategoryItemsItem{
				Name:    item.Item,
				Healthy: item.Healthy,
			})
		}

		categories = append(categories, backend.EateryMenuCategory{
			Name:  category.Category,
			Items: items,
		})
	}

	return categories
}

func convertExternalRegion(external external.Eatery) backend.EateryRegion {
	switch external.CampusArea.Descrshort {
	case "Central":
		return backend.EateryRegionCentral
	case "West":
		return backend.EateryRegionWest
	case "North":
		return backend.EateryRegionNorth
	default:
		return backend.EateryRegionUnknown
	}
}

func convertExternalPayMethods(external external.Eatery) []backend.EateryPayMethodsItem {
	var payMethods []backend.EateryPayMethodsItem

	for _, payMethod := range external.PayMethods {
		switch payMethod.DescrShort {
		case "Meal Plan - Swipe":
			payMethods = append(payMethods, backend.EateryPayMethodsItemSwipes)
		case "Meal Plan - Debit":
			payMethods = append(payMethods, backend.EateryPayMethodsItemBigRedBucks)
		case "Cash":
			payMethods = append(payMethods, backend.EateryPayMethodsItemCash)
		case "Mobile Payments":
			payMethods = append(payMethods, backend.EateryPayMethodsItemDigitalWallet)
		case "Major Credit Cards":
			payMethods = append(payMethods, backend.EateryPayMethodsItemCreditCard)
		case "Cornell Card":
			payMethods = append(payMethods, backend.EateryPayMethodsItemCornellCard)
		default:
			continue
		}
	}

	return payMethods
}

func convertExternalCategories(external external.Eatery) []backend.EateryCategoriesItem {
	var categories []backend.EateryCategoriesItem

	for _, eateryType := range external.EateryTypes {
		switch eateryType.Descr {
		case "Convenience Store":
			categories = append(categories, backend.EateryCategoriesItemConvenienceStore)
		case "Cafe":
			categories = append(categories, backend.EateryCategoriesItemCafe)
		case "Dining Room":
			categories = append(categories, backend.EateryCategoriesItemDiningRoom)
		case "Coffee Shop":
			categories = append(categories, backend.EateryCategoriesItemCoffeeShop)
		case "Cart":
			categories = append(categories, backend.EateryCategoriesItemCart)
		case "Food Court":
			categories = append(categories, backend.EateryCategoriesItemFoodCourt)
		default:
			continue
		}
	}

	return categories
}

func hoursFromEvents(events []backend.EateryEvent) []backend.Hours {
	var hours []backend.Hours

	// convert to hours objects
	for _, event := range events {
		hours = append(hours, backend.Hours{
			Start: event.Start,
			End:   event.End,
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
	var merged []backend.Hours
	var currentStart time.Time
	var currentEnd time.Time

	for _, hour := range hours {
		if currentStart.IsZero() {
			currentStart = hour.Start
			currentEnd = hour.End
			continue
		}

		diff := hour.Start.Sub(currentEnd)
		if diff < 0 {
			diff = -diff
		}
		if diff <= 5*time.Minute {
			currentEnd = hour.End
			continue
		} else {
			merged = append(merged, backend.Hours{
				Start: currentStart,
				End:   currentEnd,
			})
			currentStart = time.Time{}
			currentEnd = time.Time{}
		}
	}

	// append the final values if necessary
	if !(currentStart.IsZero()) {
		merged = append(merged, backend.Hours{
			Start: currentStart,
			End:   currentEnd,
		})
	}

	return merged
}

func selectNextWeekEvents(events []backend.EateryEvent) backend.EateryNextWeekEvents {
	now := time.Now()
	weekStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekEnd := weekStart.AddDate(0, 0, 7)

	result := backend.EateryNextWeekEvents{}

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
	return utils.ImageNameToPath("dining", imageName)
}
