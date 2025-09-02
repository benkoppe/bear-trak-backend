// Package bustime loads transit data from bustime.
package bustime

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

func chunkFetchVehicles(baseURL, apiKey string, routeIds []string) ([]Vehicle, error) {
	vehicles := []Vehicle{}
	for i, chunk := range utils.ChunkArray(routeIds, 10) {
		chunkVehicles, err := fetchVehicles(baseURL, apiKey, chunk)
		if err != nil {
			return []Vehicle{}, fmt.Errorf("failed on chunk %d: %w", i, err)
		}
		vehicles = append(vehicles, chunkVehicles...)
	}
	return vehicles, nil
}

func fetchVehicles(baseURL, apiKey string, routeIds []string) ([]Vehicle, error) {
	if len(routeIds) == 0 {
		return []Vehicle{}, nil
	}
	if len(routeIds) > 10 {
		return nil, fmt.Errorf("too many route ids: %d", len(routeIds))
	}

	routeIdsStr := strings.Join(routeIds, ",")

	fullURL, err := buildURL(baseURL, apiKey, "getvehicles", map[string]string{
		"rt": routeIdsStr,
	})
	if err != nil {
		return []Vehicle{}, fmt.Errorf("failed to build url: %w", err)
	}

	vehicles, err := utils.DoGetRequest[bustimeResponse[vehiclesResponse]](fullURL, nil)
	if err != nil {
		return []Vehicle{}, fmt.Errorf("failed to fetch vehicles: %w", err)
	}
	return vehicles.Response.Vehicles, nil
}

func fetchRoutes(baseURL, apiKey string) ([]Route, error) {
	fullURL, err := buildURL(baseURL, apiKey, "getroutes", map[string]string{})
	if err != nil {
		return []Route{}, fmt.Errorf("failed to build url: %w", err)
	}

	routes, err := utils.DoGetRequest[bustimeResponse[routesResponse]](fullURL, nil)
	if err != nil {
		return []Route{}, fmt.Errorf("failed to fetch routes: %w", err)
	}
	return routes.Response.Routes, nil
}

func buildURL(baseURL, apiKey, route string, params map[string]string) (string, error) {
	extendedURL, err := utils.ExtendURL(baseURL, route)
	if err != nil {
		return "", fmt.Errorf("failed to extend url: %w", err)
	}

	parsedURL, err := url.Parse(*extendedURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %w", err)
	}

	query := parsedURL.Query()
	query.Set("format", "json")
	query.Set("key", apiKey)
	for key, value := range params {
		query.Set(key, value)
	}

	parsedURL.RawQuery = query.Encode()
	return parsedURL.String(), nil
}
