// Package bustime loads transit data from bustime.
package bustime

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

func chunkFetchVehicles(baseUrl, apiKey string, routeIds []string) ([]Vehicle, error) {
	vehicles := []Vehicle{}
	for i, chunk := range utils.ChunkArray(routeIds, 10) {
		chunkVehicles, err := fetchVehicles(baseUrl, apiKey, chunk)
		if err != nil {
			return []Vehicle{}, fmt.Errorf("failed on chunk %d: %w", i, err)
		}
		vehicles = append(vehicles, chunkVehicles...)
	}
	return vehicles, nil
}

func fetchVehicles(baseUrl, apiKey string, routeIds []string) ([]Vehicle, error) {
	if len(routeIds) == 0 {
		return []Vehicle{}, nil
	}
	if len(routeIds) > 10 {
		return nil, fmt.Errorf("too many route ids: %d", len(routeIds))
	}

	routeIdsStr := strings.Join(routeIds, ",")

	fullUrl, err := buildUrl(baseUrl, apiKey, "getvehicles", map[string]string{
		"rt": routeIdsStr,
	})
	if err != nil {
		return []Vehicle{}, fmt.Errorf("failed to build url: %w", err)
	}

	vehicles, err := utils.DoGetRequest[bustimeResponse[vehiclesResponse]](fullUrl, nil)
	if err != nil {
		return []Vehicle{}, fmt.Errorf("failed to fetch vehicles: %w", err)
	}
	return vehicles.Response.Vehicles, nil
}

func fetchRoutes(baseUrl, apiKey string) ([]Route, error) {
	fullUrl, err := buildUrl(baseUrl, apiKey, "getroutes", map[string]string{})
	if err != nil {
		return []Route{}, fmt.Errorf("failed to build url: %w", err)
	}

	routes, err := utils.DoGetRequest[bustimeResponse[routesResponse]](fullUrl, nil)
	if err != nil {
		return []Route{}, fmt.Errorf("failed to fetch routes: %w", err)
	}
	return routes.Response.Routes, nil
}

func buildUrl(baseUrl, apiKey, route string, params map[string]string) (string, error) {
	extendedUrl, err := utils.ExtendUrl(baseUrl, route)
	if err != nil {
		return "", fmt.Errorf("failed to extend url: %w", err)
	}

	parsedUrl, err := url.Parse(*extendedUrl)
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %w", err)
	}

	query := parsedUrl.Query()
	query.Set("format", "json")
	query.Set("key", apiKey)
	for key, value := range params {
		query.Set(key, value)
	}

	parsedUrl.RawQuery = query.Encode()
	return parsedUrl.String(), nil
}
