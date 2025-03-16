package umich

import (
	"fmt"
	"log"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/dining/umich/scrape"
	"github.com/benkoppe/bear-trak-backend/go-server/dining/umich/static"
)

type Cache = scrape.Cache

func InitCache(baseUrl string) Cache {
	return scrape.InitCache(baseUrl)
}

func Get(
	scrapeCache Cache,
) ([]api.Eatery, error) {
	cacheResponse, err := scrapeCache.Get()
	if err != nil {
		return nil, fmt.Errorf("error fetching external data: %w", err)
	}

	var eateries []api.Eatery
	for static, scraped := range cacheResponse {
		if len(scraped) == 0 {
			log.Printf("no scraped data for eatery %d", static.ID)
			continue
		}

		newEatery := convertScraped(static, scraped)
		eateries = append(eateries, newEatery)
	}

	return eateries, nil
}

func convertScraped(static static.Eatery, scraped []scrape.Eatery) api.Eatery {
	firstScraped := scraped[0]

	return api.Eatery{
		ID:        static.ID,
		Name:      firstScraped.Name,
		NameShort: firstScraped.Name,
		Hours:     convertHours(scraped),
	}
}

func convertHours(scraped []scrape.Eatery) []api.Hours {
	var hours []api.Hours
	for _, eatery := range scraped {
		for _, scrapeHours := range eatery.Hours {
			h := api.Hours{
				Start: scrapeHours.StartTime,
				End:   scrapeHours.EndTime,
			}
			hours = append(hours, h)
		}
	}
	return hours
}
