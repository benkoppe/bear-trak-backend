// Package external loads external umich dining content.
package external

import (
	"fmt"
	"net/url"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

func FetchLocations(baseURL, apiKey string) ([]Location, error) {
	response, err := doGet[[]Location](baseURL, "locations", url.Values{
		"key": {apiKey},
	})
	if response == nil {
		return []Location{}, err
	}
	return *response, err
}

func FetchMealHours(baseURL, apiKey, location string, date time.Time) (*MealHoursResponse, error) {
	return doGet[MealHoursResponse](baseURL, "meal-hours", url.Values{
		"key":      {apiKey},
		"location": {location},
		"date":     {date.Format("02-01-2006")},
	})
}

func FetchMenu(baseURL, apiKey, location string, date time.Time, meal string) (*MenuResponse, error) {
	return doGet[MenuResponse](baseURL, "menu", url.Values{
		"key":      {apiKey},
		"location": {location},
		"date":     {date.Format("02-01-2006")},
		"meal":     {meal},
	})
}

func doGet[T any](baseURL, endpoint string, params url.Values) (*T, error) {
	requestURL, err := buildRequestURL(baseURL, endpoint, params)
	if err != nil {
		return nil, err
	}

	return utils.DoGetRequest[T](requestURL, nil)
}

func buildRequestURL(baseURL, endpoint string, params url.Values) (string, error) {
	fullURL, err := utils.ExtendURL(baseURL, endpoint)
	if err != nil {
		return "", fmt.Errorf("failed to extend url for endpoint %q: %w", endpoint, err)
	}

	parsedURL, err := url.Parse(*fullURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse url for endpoint %q: %w", endpoint, err)
	}

	query := parsedURL.Query()
	for key, values := range params {
		for _, value := range values {
			query.Add(key, value)
		}
	}
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}
