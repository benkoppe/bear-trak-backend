// Package convex loads additional cornell dining content from the Convex admin panel.
package convex

import (
	"fmt"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
	"github.com/benkoppe/bear-trak-backend/go-server/utils/timeutils"
)

// Cache is a TTL cache holding the list of Convex-managed eateries.
type Cache = *utils.Cache[[]Eatery]

// InitCache creates and pre-warms a cache for the Convex eateries endpoint.
// baseURL is the Convex deployment URL (e.g. "https://xxx.convex.cloud").
// token is the bearer token used for authentication (may be empty for unauthenticated queries).
func InitCache(baseURL, token string) Cache {
	return utils.NewCache(
		"diningConvex",
		1*time.Minute,
		func() ([]Eatery, error) {
			return fetchEateries(baseURL, token)
		},
	)
}

func fetchEateries(baseURL, token string) ([]Eatery, error) {
	url := baseURL + "/api/query"
	todayKey := time.Now().In(timeutils.LoadEST()).Format("2006-01-02")

	body := QueryRequest{
		Path: "eateries:getAll",
		Args: map[string]string{
			"token":    token,
			"todayKey": todayKey,
		},
		Format: "json",
	}

	resp, err := utils.DoPostRequest[QueryResponse](url, body)
	if err != nil {
		return nil, fmt.Errorf("convex request failed: %w", err)
	}

	if resp.Status != "success" {
		return nil, fmt.Errorf("convex returned error: %s", resp.ErrorMessage)
	}

	return resp.Value, nil
}
