package study

import (
	"fmt"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/schools/cornell/external_map"
	"github.com/benkoppe/bear-trak-backend/go-server/study/cornell/libraries"
	"github.com/benkoppe/bear-trak-backend/go-server/study/cornell/libraries/external"
	"github.com/benkoppe/bear-trak-backend/go-server/study/cornell/printers"
)

type Cache = external.Cache

func InitCache(url string) Cache {
	return external.InitCache(url)
}

func Get(externalCache Cache, mapCache external_map.Cache) (*api.StudyData, error) {
	libraries, err := libraries.Get(externalCache, mapCache)
	if err != nil {
		return nil, fmt.Errorf("error getting libraries: %v", err)
	}

	printers, err := printers.Get(mapCache)
	if err != nil {
		return nil, fmt.Errorf("error getting printers: %v", err)
	}

	return &api.StudyData{
		Libraries: libraries,
		Printers:  printers,
	}, nil
}
