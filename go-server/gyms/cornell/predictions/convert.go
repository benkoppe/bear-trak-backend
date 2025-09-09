package predictions

import (
	"fmt"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/gyms/cornell/scrape"
	"github.com/benkoppe/bear-trak-backend/go-server/gyms/cornell/shared"
	"github.com/benkoppe/bear-trak-backend/go-server/gyms/cornell/static"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
	"github.com/benkoppe/bear-trak-backend/go-server/utils/timeutils"
)

type Cache = *utils.Cache[[]api.GymCapacityPredictions]

func InitCache(url string, hoursCache scrape.Cache) Cache {
	return utils.NewCache(
		"GymCapacityPredictions",
		5*time.Minute,
		func() ([]api.GymCapacityPredictions, error) {
			return GetData(url, hoursCache)
		})
}

func GetData(url string, hoursCache scrape.Cache) ([]api.GymCapacityPredictions, error) {
	fetchedPredictions, err := fetchData(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching prediction data: %w", err)
	}

	return convertAllFetched(fetchedPredictions, hoursCache), nil
}

func convertAllFetched(fetched []Prediction, hoursCache scrape.Cache) []api.GymCapacityPredictions {
	// group by location name
	byLocation := make(map[string][]api.GymCapacityPredictionPoint)
	for _, prediction := range fetched {
		converted := convertFetched(prediction)
		byLocation[prediction.GymName] = append(byLocation[prediction.GymName], converted)
	}

	// convert location name to IDs
	staticData := static.GetGyms()
	scrapedSchedules, err := hoursCache.Get()
	if err != nil {
		fmt.Printf("error fetching scraped schedules: %v\n", err)
	}

	est := timeutils.LoadEST()
	now := time.Now().In(est)

	var result []api.GymCapacityPredictions
	for locationName, predictions := range byLocation {
		staticGym := utils.Find(staticData, func(gym static.Gym) bool {
			return gym.PredictionName == locationName
		})

		if staticGym == nil {
			fmt.Printf("couldn't find static ID for location name %s\n", locationName)
			continue
		}

		hours := shared.CreateFutureHours(*staticGym, scrapedSchedules)
		firstOpen, lastClose := timeutils.FirstOpenAndLastClose(hours, now)
		filteredPredictions := shared.FilterByTimeRange(
			predictions,
			func(p api.GymCapacityPredictionPoint) time.Time { return p.Timestamp },
			firstOpen,
			lastClose,
		)

		result = append(result, api.GymCapacityPredictions{
			LocationId: staticGym.ID,
			Points:     filteredPredictions,
		})
	}

	return result
}

func convertFetched(fetched Prediction) api.GymCapacityPredictionPoint {
	return api.GymCapacityPredictionPoint{
		Timestamp:        fetched.Timestamp,
		PredictionMadeAt: fetched.PredictionMade,
		Count:            fetched.Predicted,
	}
}
