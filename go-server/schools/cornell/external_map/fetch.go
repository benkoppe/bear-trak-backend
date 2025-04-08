package external_map

import (
	"fmt"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

type cache = *utils.Cache[[]itemCategory]

type Cache struct {
	cache cache
}

func (cache Cache) Get(categoryDomId string) ([]Item, error) {
	categories, err := cache.cache.Get()
	if err != nil {
		return nil, err
	}

	category := utils.Find(categories, func(cat itemCategory) bool {
		return cat.DOM_ID == categoryDomId
	})
	if category == nil {
		return nil, fmt.Errorf("no category found with DOM_ID %s", categoryDomId)
	}

	return category.Items, nil
}

func InitCache(url string) Cache {
	return Cache{
		cache: initCache(url),
	}
}

func initCache(url string) cache {
	return utils.NewCache(
		"cornellMapExternal",
		1*time.Hour,
		func() ([]itemCategory, error) {
			return FetchData(url)
		})
}

func FetchData(url string) ([]itemCategory, error) {
	response, err := utils.DoGetRequest[overlayResponse](url, nil)
	if response == nil {
		return []itemCategory{}, err
	}
	return response.Overlays, err
}
