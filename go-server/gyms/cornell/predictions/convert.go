package predictions

import (
	"fmt"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/gyms/cornell/static"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

type Cache = *utils.Cache[[]api.GymCapacityPredictions]

func InitCache(url string) Cache {
	return utils.NewCache(
		"GymCapacityPredictions",
		5*time.Minute,
		func() ([]api.GymCapacityPredictions, error) {
			return GetData(url)
		})
}

func GetData(url string) ([]api.GymCapacityPredictions, error) {
	fetchedPredictions, err := fetchData(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching prediction data: %w", err)
	}

	return convertAllFetched(fetchedPredictions), nil
}

func convertAllFetched(fetched []Prediction) []api.GymCapacityPredictions {
	// group by location name
	byLocation := make(map[string][]api.GymCapacityPredictionPoint)
	for _, prediction := range fetched {
		converted := convertFetched(prediction)
		byLocation[prediction.GymName] = append(byLocation[prediction.GymName], converted)
	}

	// convert location name to IDs
	staticData := static.GetGyms()
	var result []api.GymCapacityPredictions
	for locationName, predictions := range byLocation {
		staticGym := utils.Find(staticData, func(gym static.Gym) bool {
			return gym.PredictionName == locationName
		})

		if staticGym == nil {
			fmt.Printf("couldn't find static ID for location name %s", locationName)
			continue
		}

		result = append(result, api.GymCapacityPredictions{
			LocationId: staticGym.ID,
			Points:     predictions,
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
