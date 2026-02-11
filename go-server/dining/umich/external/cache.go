// Package external loads external umich dining content.
package external

import (
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
	"github.com/benkoppe/bear-trak-backend/go-server/utils/timeutils"
)

var (
	DaysToFetch             = 7
	MaxConcurrentLocations  = 6
	LocationDataCacheExpiry = 8 * time.Hour
)

type Cache = *utils.Cache[[]LocationDiningData]

func InitCache(baseURL, apiKey string) Cache {
	return utils.NewCache(
		"diningUMichLocationData",
		LocationDataCacheExpiry,
		func() ([]LocationDiningData, error) {
			return fetchAllLocationDiningData(baseURL, apiKey)
		},
	)
}

func fetchAllLocationDiningData(baseURL, apiKey string) ([]LocationDiningData, error) {
	locations, err := FetchLocations(baseURL, apiKey)
	if err != nil {
		return nil, err
	}

	est := timeutils.LoadEST()
	startDate := time.Now().In(est)
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, est)

	var (
		result []LocationDiningData
		wg     sync.WaitGroup
		mu     sync.Mutex
	)

	semaphore := make(chan struct{}, MaxConcurrentLocations)

	for i := range locations {
		wg.Add(1)
		go func(location Location) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			locationData := fetchLocationDiningData(baseURL, apiKey, location, startDate)

			mu.Lock()
			result = append(result, locationData)
			mu.Unlock()
		}(locations[i])
	}

	wg.Wait()

	sort.Slice(result, func(i, j int) bool {
		return strings.ToLower(result[i].Location.Name) < strings.ToLower(result[j].Location.Name)
	})

	return result, nil
}

func fetchLocationDiningData(baseURL, apiKey string, location Location, startDate time.Time) LocationDiningData {
	locationData := LocationDiningData{
		Location: location,
		Days:     make([]LocationDayData, 0, DaysToFetch),
	}

	for dayOffset := 0; dayOffset < DaysToFetch; dayOffset++ {
		date := startDate.AddDate(0, 0, dayOffset)
		dayData, err := fetchLocationDayData(baseURL, apiKey, location, date)
		if err != nil {
			log.Printf("error fetching umich meal-hours for location=%q date=%s: %v",
				location.Name, date.Format("2006-01-02"), err)
			continue
		}
		locationData.Days = append(locationData.Days, dayData)
	}

	return locationData
}

func fetchLocationDayData(baseURL, apiKey string, location Location, date time.Time) (LocationDayData, error) {
	mealHoursResponse, err := FetchMealHours(baseURL, apiKey, location.Name, date)
	if err != nil {
		return LocationDayData{}, err
	}

	validMeals := filterValidMeals(mealHoursResponse.Meal)
	dayData := LocationDayData{
		Date:  date.Format("02-01-2006"),
		Meals: make([]DayMealData, 0, len(validMeals)),
	}

	for _, meal := range validMeals {
		menuResponse, err := FetchMenu(baseURL, apiKey, location.Name, date, meal.Name)
		if err != nil {
			log.Printf("error fetching umich menu for location=%q date=%s meal=%q: %v",
				location.Name, date.Format("2006-01-02"), meal.Name, err)
			continue
		}

		dayData.Meals = append(dayData.Meals, DayMealData{
			MealName: meal.Name,
			Meal:     meal,
			Hours:    selectHoursForMeal(meal.Name, mealHoursResponse.Hours, len(validMeals)),
			Menu:     menuResponse.Menu,
		})
	}

	return dayData, nil
}

func filterValidMeals(meals []MealItem) []MealItem {
	valid := make([]MealItem, 0, len(meals))
	for _, meal := range meals {
		if meal.HasMenu {
			valid = append(valid, meal)
		}
	}
	return valid
}

func selectHoursForMeal(mealName string, hours []EventHour, validMealCount int) []EventHour {
	target := normalizeMealName(mealName)
	matched := make([]EventHour, 0, len(hours))

	for _, hour := range hours {
		if normalizeMealName(hour.EventTitle) == target {
			matched = append(matched, hour)
		}
	}

	if len(matched) > 0 {
		return matched
	}

	// If there's only one valid meal today and no explicit hour titles,
	// assign the location's provided hours (often "Open") to that meal.
	if validMealCount == 1 {
		return append([]EventHour(nil), hours...)
	}

	return []EventHour{}
}

func normalizeMealName(name string) string {
	return strings.Join(strings.Fields(strings.ToLower(name)), " ")
}
