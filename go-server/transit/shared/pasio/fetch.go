package pasio

import (
	"fmt"
	"net/url"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

type Cache = *utils.Cache[[]Route]

func InitCache(baseUrl string, systemId string) Cache {
	return utils.NewCache(
		"transitExternalPasio",
		time.Hour*24,
		func() ([]Route, error) {
			return fetchRoutes(baseUrl, systemId)
		},
	)
}

func fetchRoutes(baseUrl string, systemId string) ([]Route, error) {
	parsedUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	// query parameters
	query := parsedUrl.Query()
	query.Set("getRoutes", "1")

	parsedUrl.RawQuery = query.Encode()
	fullUrl := parsedUrl.String()

	requestBody := map[string]string{
		"systemSelected0": systemId,
		"amount":          "1",
	}

	routes, err := utils.DoPostRequest[[]Route](fullUrl, requestBody)
	if err != nil {
		return nil, fmt.Errorf("error fetching routes: %w", err)
	}

	return *routes, nil
}
