package study

import (
	"fmt"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/study/harvard/libraries"
	"github.com/benkoppe/bear-trak-backend/go-server/study/harvard/libraries/external"
)

type Cache = external.Cache

func InitCache(url string) Cache {
	return external.InitCache(url)
}

func Get(cache Cache) (*api.StudyData, error) {
	libraries, err := libraries.Get(cache)
	if err != nil {
		return nil, fmt.Errorf("error getting libraries: %v", err)
	}

	return &api.StudyData{
		Libraries: libraries,
		Printers:  []api.Printer{},
	}, nil
}
