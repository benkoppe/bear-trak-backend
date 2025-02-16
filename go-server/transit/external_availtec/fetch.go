package external_availtec

import (
	"fmt"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

func FetchRoutes(baseUrl string) ([]Route, error) {
	fullUrl, err := utils.ExtendUrl(baseUrl, "Routes/GetVisibleRoutes")
	if fullUrl == nil {
		return []Route{}, fmt.Errorf("failed to extend url: %w", err)
	}

	routes, err := utils.DoGetRequest[[]Route](*fullUrl)
	if routes == nil {
		return []Route{}, err
	}

	return *routes, nil
}
