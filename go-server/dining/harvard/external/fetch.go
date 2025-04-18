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

func InitCaches(baseUrl, apiKey string) Caches {
	return Caches{
		LocationsCache: initLocationsCache(baseUrl, apiKey),
		RecipesCache:   initRecipesCache(baseUrl, apiKey),
	}
}

func initLocationsCache(baseUrl, apiKey string) *utils.Cache[[]Location] {
	return utils.NewCache(
		"diningLocations",
		24*time.Hour,
		func() ([]Location, error) {
			return fetchLocations(baseUrl, apiKey)
		})
}

func initRecipesCache(baseUrl, apiKey string) *utils.Cache[[]Recipe] {
	return utils.NewCache(
		"diningRecipes",
		24*time.Hour,
		func() ([]Recipe, error) {
			return fetchRecipes(baseUrl, apiKey)
		})
}

func fetchLocations(baseUrl, apiKey string) ([]Location, error) {
	fullUrl, err := utils.ExtendUrl(baseUrl, "locations")
	if err != nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	locations, err := utils.DoGetRequest[[]Location](*fullUrl, getRequestHeaders(apiKey))
	if locations == nil {
		return []Location{}, err
	}
	return *locations, err
}

func fetchRecipes(baseUrl, apiKey string) ([]Recipe, error) {
	fullUrl, err := utils.ExtendUrl(baseUrl, "recipes")
	if err != nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	recipes, err := utils.DoGetRequest[[]Recipe](*fullUrl, getRequestHeaders(apiKey))
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
