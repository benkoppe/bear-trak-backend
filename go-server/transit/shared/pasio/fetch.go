// Package pasio loads transit data from pasio go.
package pasio

import (
	"fmt"
	"net/url"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

type Cache = *utils.Cache[[]Route]

func InitCache(baseURL string, systemID string) Cache {
	return utils.NewCache(
		"transitExternalPasio",
		time.Hour*24,
		func() ([]Route, error) {
			return fetchRoutes(baseURL, systemID)
		},
	)
}

func fetchRoutes(baseURL string, systemID string) ([]Route, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	// query parameters
	query := parsedURL.Query()
	query.Set("getRoutes", "1")

	parsedURL.RawQuery = query.Encode()
	fullURL := parsedURL.String()

	requestBody := map[string]string{
		"systemSelected0": systemID,
		"amount":          "1",
	}

	routes, err := utils.DoPostRequest[[]Route](fullURL, requestBody)
	if err != nil {
		return nil, fmt.Errorf("error fetching routes: %w", err)
	}

	return *routes, nil
}
