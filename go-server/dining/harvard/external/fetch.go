// Package external loads external harvard dining content.
package external

import (
	"fmt"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

type Caches struct {
	LocationsCache *utils.Cache[[]Location]
	RecipesCache   *utils.Cache[[]Recipe]
}

func InitCaches(baseURL, apiKey string) Caches {
	return Caches{
		LocationsCache: initLocationsCache(baseURL, apiKey),
		RecipesCache:   initRecipesCache(baseURL, apiKey),
	}
}

func initLocationsCache(baseURL, apiKey string) *utils.Cache[[]Location] {
	return utils.NewCache(
		"diningLocations",
		24*time.Hour,
		func() ([]Location, error) {
			return fetchLocations(baseURL, apiKey)
		})
}

func initRecipesCache(baseURL, apiKey string) *utils.Cache[[]Recipe] {
	return utils.NewCache(
		"diningRecipes",
		24*time.Hour,
		func() ([]Recipe, error) {
			return fetchRecipes(baseURL, apiKey)
		})
}

func fetchLocations(baseURL, apiKey string) ([]Location, error) {
	fullURL, err := utils.ExtendURL(baseURL, "locations")
	if err != nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	locations, err := utils.DoGetRequest[[]Location](*fullURL, getRequestHeaders(apiKey))
	if locations == nil {
		return []Location{}, err
	}
	return *locations, err
}

func fetchRecipes(baseURL, apiKey string) ([]Recipe, error) {
	fullURL, err := utils.ExtendURL(baseURL, "recipes")
	if err != nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	recipes, err := utils.DoGetRequest[[]Recipe](*fullURL, getRequestHeaders(apiKey))
	if recipes == nil {
		return []Recipe{}, err
	}
	return *recipes, err
}

func getRequestHeaders(apiKey string) map[string]string {
	return map[string]string{
		"X-Api-Key": apiKey,
	}
}
