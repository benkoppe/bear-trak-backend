package external_gtfs

import (
	"fmt"
	"io"
	"net/http"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
	"github.com/jamespfennell/gtfs"
)

type Cache = *utils.Cache[*gtfs.Static]

func InitCache(url string) Cache {
	return utils.NewCache(
		"transitExternalGtfs",
		utils.NoExpiration,
		func() (*gtfs.Static, error) {
			return loadData(url)
		},
	)
}

func loadData(url string) (*gtfs.Static, error) {
	tcatGtfsData, err := loadTcatGtfs(url)
	if err != nil {
		return nil, fmt.Errorf("error loading tcat data: %v", err)
	}

	staticData, err := gtfs.ParseStatic(tcatGtfsData, gtfs.ParseStaticOptions{})
	if err != nil {
		return nil, fmt.Errorf("error parsing tcat data: %v", err)
	}

	return staticData, nil
}

func loadTcatGtfs(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching gtfs ZIP: %v", err)
	}
	defer resp.Body.Close()

	originalGtfsData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading gtfs ZIP: %v", err)
	}

	return originalGtfsData, nil
}
