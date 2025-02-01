package gyms

import (
	"fmt"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/gyms/external"
	"github.com/benkoppe/bear-trak-backend/go-server/gyms/static"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

func Get(url string) ([]api.Gym, error) {
	staticData := static.GetGyms()

	if len(staticData) == 0 {
		return nil, fmt.Errorf("loaded empty static gyms")
	}

	externalData, err := external.FetchData(url)
	if err != nil {
		// don't break here - if capacities doesn't work, we still want to provide static data.
		// instead, simply print an error.
		fmt.Printf("error fetching external data: %v", err)
	}

	var gyms []api.Gym

	for _, staticGym := range staticData {
		gym := convertStatic(staticGym)

		capacityData := findCapacityData(staticGym, externalData)
		if capacityData != nil {
			capacity := convertExternalGymCapacity(staticGym, *capacityData)
			gym.Capacity = api.NewNilGymCapacity(capacity)
		}

		gyms = append(gyms, gym)
	}

	return gyms, nil
}

func convertStatic(static static.Gym) api.Gym {
	return api.Gym{
		ID:                  static.ID,
		Name:                static.Name,
		ImagePath:           utils.ImageNameToPath("gyms", static.ImageName),
		Latitude:            static.Location.Latitude,
		Longitude:           static.Location.Longitude,
		Hours:               createFutureHours(static.WeekHours),
		Facilities:          convertStaticFacilities(static),
		EquipmentCategories: convertStaticEquipmentCategories(static),
		Capacity:            api.NilGymCapacity{Null: true},
	}
}

func findCapacityData(static static.Gym, externalData []external.Gym) *external.Gym {
	for _, capacityData := range externalData {
		if capacityData.LocationName == static.ScrapeName {
			return &capacityData
		}
	}
	return nil
}

func createFutureHours(staticHours static.WeekHours) []api.Hours {
	est := utils.LoadEST()
	now := time.Now().In(est)
	var futureHours []api.Hours

	for i := 0; i < 7; i++ {
		date := now.AddDate(0, 0, i)
		dayHours := staticHours.GetHours(date)

		for _, hours := range dayHours {
			start, e1 := hours.Open.ToDate(date)
			end, e2 := hours.Close.ToDate(date)

			if e1 != nil {
				fmt.Printf("error parsing hours: %v", e1)
				continue
			}
			if e2 != nil {
				fmt.Printf("error parsing hours: %v", e2)
				continue
			}

			futureHours = append(futureHours, api.Hours{
				Start: start,
				End:   end,
			})
		}
	}

	return futureHours
}

func convertStaticFacilities(static static.Gym) []api.GymFacilitiesItem {
	var facilities []api.GymFacilitiesItem

	for _, facility := range static.Facilities {
		facilities = append(facilities, api.GymFacilitiesItem{
			FacilityType: convertStaticFacilityType(facility),
			Name:         facility.Name,
		})
	}

	return facilities
}

func convertStaticFacilityType(facility static.Facility) api.GymFacilitiesItemFacilityType {
	switch facility.Type {
	case "pool":
		return api.GymFacilitiesItemFacilityTypePool
	case "basketball":
		return api.GymFacilitiesItemFacilityTypeBasketball
	case "bowling":
		return api.GymFacilitiesItemFacilityTypeBowling
	default:
		return api.GymFacilitiesItemFacilityTypeUnknown
	}
}

func convertStaticEquipmentCategories(static static.Gym) []api.GymEquipmentCategoriesItem {
	var categories []api.GymEquipmentCategoriesItem

	for _, category := range static.Equipment {
		categories = append(categories, api.GymEquipmentCategoriesItem{
			CategoryType: convertStaticGymCategoryType(category),
			Items:        category.Items,
		})
	}

	return categories
}

func convertStaticGymCategoryType(category static.Equipment) api.GymEquipmentCategoriesItemCategoryType {
	switch category.Type {
	case "treadmills":
		return api.GymEquipmentCategoriesItemCategoryTypeTreadmills
	case "ellipticals":
		return api.GymEquipmentCategoriesItemCategoryTypeEllipticals
	case "rowing":
		return api.GymEquipmentCategoriesItemCategoryTypeRowing
	case "bike":
		return api.GymEquipmentCategoriesItemCategoryTypeBike
	case "lifting":
		return api.GymEquipmentCategoriesItemCategoryTypeLifting
	case "machines":
		return api.GymEquipmentCategoriesItemCategoryTypeMachines
	case "free weights":
		return api.GymEquipmentCategoriesItemCategoryTypeFreeWeights
	default:
		return api.GymEquipmentCategoriesItemCategoryTypeMisc
	}
}

func convertExternalGymCapacity(gym static.Gym, capacity external.Gym) api.GymCapacity {
	percentage := api.NewNilInt(capacity.GetPercentage())

	// if gym is closed, set percentage to null
	est := utils.LoadEST()
	if !gym.WeekHours.IsOpen(time.Now().In(est)) {
		percentage = api.NilInt{Null: true}
	}

	return api.GymCapacity{
		Total:       capacity.TotalCapacity,
		Percentage:  percentage,
		LastUpdated: capacity.LastUpdatedDateAndTime.ToTime(),
	}
}
