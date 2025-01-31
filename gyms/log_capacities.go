package gyms

import (
	"fmt"
	"log"
)

func LogCapacities(url string) error {
	gyms, err := Get(url)
	if err != nil {
		return fmt.Errorf("error fetching gyms: %v", err)
	}

	for _, gym := range gyms {
		if gym.Capacity.IsNull() {
			log.Printf("No capacity data for gym %s", gym.Name)
			continue
		}

		percentage := gym.Capacity.Value.GetPercentage()
		if percentage.IsNull() {
			log.Printf("No percentage data for gym %s", gym.Name)
			continue
		}

		log.Printf("Gym %s can hold %d people, is %d%% full", gym.Name, gym.Capacity.Value.Total, percentage.Value)
	}

	return nil
}
