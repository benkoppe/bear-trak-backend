// Package umich loads umich dining content.
package umich

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/dining/shared"
	"github.com/benkoppe/bear-trak-backend/go-server/dining/umich/external"
	"github.com/benkoppe/bear-trak-backend/go-server/dining/umich/static"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
	"github.com/benkoppe/bear-trak-backend/go-server/utils/timeutils"
	"golang.org/x/text/unicode/norm"
)

type Cache = external.Cache

func InitCache(baseURL, apiKey string) Cache {
	return external.InitCache(baseURL, apiKey)
}

func Get(cache Cache) ([]api.Eatery, error) {
	locationData, err := cache.Get()
	if err != nil {
		return nil, fmt.Errorf("error fetching external data: %w", err)
	}

	staticEateries := static.GetEateries()
	eateries := make([]api.Eatery, 0, len(staticEateries))

	for _, staticEatery := range staticEateries {
		location := matchingLocationData(staticEatery, locationData)
		if location == nil {
			log.Printf("no external data matched static eatery id=%d name=%q buildingId=%d",
				staticEatery.ID, staticEatery.LocationDisplayName, staticEatery.OfficialBuildingID)
			continue
		}

		events := convertLocationEvents(*location)
		hours := hoursFromEvents(events)

		eateryName := strings.TrimSpace(location.Location.DisplayName)
		if eateryName == "" {
			eateryName = strings.TrimSpace(location.Location.Name)
		}
		if eateryName == "" {
			eateryName = staticEatery.LocationDisplayName
		}

		eateries = append(eateries, api.Eatery{
			ID:             staticEatery.ID,
			Name:           eateryName,
			NameShort:      eateryName,
			ImagePath:      utils.ImageNameToPath("dining/umich", staticEatery.ImageName),
			Latitude:       location.Location.Lat,
			Longitude:      location.Location.Lng,
			Location:       formatLocationAddress(location.Location),
			Hours:          hours,
			Region:         staticEatery.Region,
			PayMethods:     append([]string{}, staticEatery.PayMethods...),
			Categories:     shared.ConvertStaticCategories(staticEatery.Categories),
			NextWeekEvents: shared.SelectNextWeekEvents(events),
		})
	}

	sort.Slice(eateries, func(i, j int) bool {
		return eateries[i].ID < eateries[j].ID
	})

	return eateries, nil
}

func matchingLocationData(staticEatery static.Eatery, locationData []external.LocationDiningData) *external.LocationDiningData {
	normalizedDisplayName := normalizeName(staticEatery.LocationDisplayName)
	var buildingMatches []int

	for i := range locationData {
		if locationData[i].Location.OfficialBuildingID == staticEatery.OfficialBuildingID {
			buildingMatches = append(buildingMatches, i)
		}
	}

	if len(buildingMatches) > 0 {
		for _, idx := range buildingMatches {
			location := locationData[idx].Location
			if normalizeName(location.DisplayName) == normalizedDisplayName ||
				normalizeName(location.Name) == normalizedDisplayName {
				return &locationData[idx]
			}
		}

		for _, idx := range buildingMatches {
			location := locationData[idx].Location
			locationName := normalizeName(location.DisplayName)
			if strings.Contains(locationName, normalizedDisplayName) ||
				strings.Contains(normalizedDisplayName, locationName) {
				return &locationData[idx]
			}
		}

		if len(buildingMatches) == 1 {
			return &locationData[buildingMatches[0]]
		}
	}

	for i := range locationData {
		location := locationData[i].Location
		if normalizeName(location.DisplayName) == normalizedDisplayName ||
			normalizeName(location.Name) == normalizedDisplayName {
			return &locationData[i]
		}
	}

	return nil
}

func convertLocationEvents(locationData external.LocationDiningData) []api.EateryEvent {
	events := []api.EateryEvent{}

	for _, day := range locationData.Days {
		for _, meal := range day.Meals {
			menu := convertMenuCategories(meal.Menu.Category)

			hours := meal.Hours
			if len(hours) == 0 {
				fallback, ok := parseMealHours(day.Date, meal.Meal.Hours)
				if ok {
					hours = []external.EventHour{fallback}
				}
			}

			for _, hour := range hours {
				if !hour.EventTimeEnd.After(hour.EventTimeStart) {
					continue
				}
				events = append(events, api.EateryEvent{
					Start:          hour.EventTimeStart,
					End:            hour.EventTimeEnd,
					MenuCategories: menu,
				})
			}
		}
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Start.Before(events[j].Start)
	})

	return events
}

func convertMenuCategories(categories []external.MenuCategory) []api.EateryMenuCategory {
	converted := make([]api.EateryMenuCategory, 0, len(categories))
	for _, category := range categories {
		items := make([]api.EateryMenuCategoryItemsItem, 0, len(category.MenuItem))
		for _, item := range category.MenuItem {
			items = append(items, api.EateryMenuCategoryItemsItem{
				Name:    item.Name,
				Healthy: hasHealthyAttribute(item.Attribute),
			})
		}
		converted = append(converted, api.EateryMenuCategory{
			Name:  category.Name,
			Items: items,
		})
	}
	return converted
}

func hasHealthyAttribute(attributes []string) bool {
	for _, attribute := range attributes {
		if strings.HasPrefix(strings.ToLower(attribute), "mhealthy") {
			return true
		}
	}
	return false
}

func hoursFromEvents(events []api.EateryEvent) []api.Hours {
	hours := make([]api.Hours, 0, len(events))
	for _, event := range events {
		hours = append(hours, api.Hours{
			Start: event.Start,
			End:   event.End,
		})
	}
	return hours
}

func formatLocationAddress(location external.Location) string {
	street1 := strings.TrimSpace(location.Address.Street1)
	street2 := strings.TrimSpace(location.Address.Street2)
	if street1 != "" && street2 != "" {
		return street1 + ", " + street2
	}
	if street1 != "" {
		return street1
	}
	if street2 != "" {
		return street2
	}
	building := strings.TrimSpace(location.BuildingPreferredName)
	if building != "" {
		return building
	}
	return location.DisplayName
}

func normalizeName(value string) string {
	decomposed := norm.NFD.String(value)
	lower := strings.ToLower(decomposed)

	var b strings.Builder
	for _, r := range lower {
		if unicode.IsMark(r) {
			continue
		}

		switch {
		case unicode.IsLetter(r), unicode.IsNumber(r):
			b.WriteRune(r)
		case unicode.IsSpace(r), strings.ContainsRune("-_'’`/&", r):
			b.WriteRune(' ')
		}
	}

	return strings.Join(strings.Fields(b.String()), " ")
}

func parseMealHours(day, hours string) (external.EventHour, bool) {
	if strings.TrimSpace(hours) == "" {
		return external.EventHour{}, false
	}

	est := timeutils.LoadEST()
	date, err := time.ParseInLocation("02-01-2006", day, est)
	if err != nil {
		return external.EventHour{}, false
	}

	cleaned := strings.NewReplacer("–", "-", "—", "-").Replace(hours)
	parts := strings.Split(cleaned, "-")
	if len(parts) != 2 {
		return external.EventHour{}, false
	}

	start, err := parseClockTime(parts[0], date)
	if err != nil {
		return external.EventHour{}, false
	}

	end, err := parseClockTime(parts[1], date)
	if err != nil {
		return external.EventHour{}, false
	}

	if end.Before(start) {
		end = end.Add(24 * time.Hour)
	}

	return external.EventHour{
		EventTimeStart: start,
		EventTimeEnd:   end,
		EventTitle:     "Open",
	}, true
}

func parseClockTime(value string, date time.Time) (time.Time, error) {
	est := timeutils.LoadEST()
	trimmed := strings.TrimSpace(value)

	layouts := []string{
		"3:04 PM",
		"3 PM",
		"3:04PM",
		"3PM",
		"3:04 pm",
		"3 pm",
		"3:04pm",
		"3pm",
	}

	for _, layout := range layouts {
		parsed, err := time.ParseInLocation(layout, trimmed, est)
		if err == nil {
			return time.Date(date.Year(), date.Month(), date.Day(), parsed.Hour(), parsed.Minute(), 0, 0, est), nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse clock time %q", value)
}
