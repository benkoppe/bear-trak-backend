package gyms

import (
	"fmt"
	"time"

	backend "github.com/benkoppe/bear-trak-backend/backend"
	"github.com/benkoppe/bear-trak-backend/gyms/external"
	"github.com/benkoppe/bear-trak-backend/gyms/static"
	"github.com/benkoppe/bear-trak-backend/utils"
)

func Get(url string) ([]backend.Gym, error) {
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

	var gyms []backend.Gym

	for _, staticGym := range staticData {
		gym := convertStatic(staticGym)

		capacityData := findCapacityData(staticGym, externalData)
		if capacityData != nil {
			capacity := convertExternalGymCapacity(*capacityData)
			gym.Capacity = backend.NewNilGymCapacity(capacity)
		}

		gyms = append(gyms, gym)
	}

	return gyms, nil
}

func convertStatic(static static.Gym) backend.Gym {
	return backend.Gym{
		ID:                  static.ID,
		Name:                static.Name,
		ImagePath:           utils.ImageNameToPath("gyms", static.ImageName),
		Latitude:            static.Location.Latitude,
		Longitude:           static.Location.Longitude,
		Hours:               createFutureHours(static.WeekHours),
		Facilities:          convertStaticFacilities(static),
		EquipmentCategories: convertStaticEquipmentCategories(static),
		Capacity:            backend.NilGymCapacity{Null: true},
	}
}

func findCapacityData(static static.Gym, externalData []external.GymCapacity) *external.GymCapacity {
	for _, capacityData := range externalData {
		if capacityData.Name == static.ScrapeName {
			return &capacityData
		}
	}
	return nil
}

func createFutureHours(staticHours static.WeekHours) []backend.Hours {
	now := time.Now()
	var futureHours []backend.Hours

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

			futureHours = append(futureHours, backend.Hours{
				Start: start,
				End:   end,
			})
		}
	}

	return futureHours
}

func convertStaticFacilities(static static.Gym) []backend.GymFacilitiesItem {
	var facilities []backend.GymFacilitiesItem

	for _, facility := range static.Facilities {
		facilities = append(facilities, backend.GymFacilitiesItem{
			FacilityType: convertStaticFacilityType(facility),
			Name:         facility.Name,
		})
	}

	return facilities
}

func convertStaticFacilityType(facility static.Facility) backend.GymFacilitiesItemFacilityType {
	switch facility.Type {
	case "pool":
		return backend.GymFacilitiesItemFacilityTypePool
	case "basketball":
		return backend.GymFacilitiesItemFacilityTypeBasketball
	case "bowling":
		return backend.GymFacilitiesItemFacilityTypeBowling
	default:
		return backend.GymFacilitiesItemFacilityTypeUnknown
	}
}

func convertStaticEquipmentCategories(static static.Gym) []backend.GymEquipmentCategoriesItem {
	var categories []backend.GymEquipmentCategoriesItem

	for _, category := range static.Equipment {
		categories = append(categories, backend.GymEquipmentCategoriesItem{
			CategoryType: convertStaticGymCategoryType(category),
			Items:        category.Items,
		})
	}

	return categories
}

func convertStaticGymCategoryType(category static.Equipment) backend.GymEquipmentCategoriesItemCategoryType {
	switch category.Type {
	case "treadmills":
		return backend.GymEquipmentCategoriesItemCategoryTypeTreadmills
	case "ellipticals":
		return backend.GymEquipmentCategoriesItemCategoryTypeEllipticals
	case "rowing":
		return backend.GymEquipmentCategoriesItemCategoryTypeRowing
	case "bike":
		return backend.GymEquipmentCategoriesItemCategoryTypeBike
	case "lifting":
		return backend.GymEquipmentCategoriesItemCategoryTypeLifting
	case "machines":
		return backend.GymEquipmentCategoriesItemCategoryTypeMachines
	case "free weights":
		return backend.GymEquipmentCategoriesItemCategoryTypeFreeWeights
	default:
		return backend.GymEquipmentCategoriesItemCategoryTypeMisc
	}
}

func convertExternalGymCapacity(capacity external.GymCapacity) backend.GymCapacity {
	var percentage backend.NilInt

	if capacity.Percentage != nil {
		percentage = backend.NewNilInt(int(*capacity.Percentage))
	} else {
		percentage = backend.NilInt{Null: true}
	}

	return backend.GymCapacity{
		Count:       int(capacity.Count),
		Percentage:  percentage,
		LastUpdated: capacity.LastUpdated,
	}
}
