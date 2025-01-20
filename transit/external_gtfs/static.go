package external_gtfs

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/jamespfennell/gtfs"
)

func GetStaticGtfs(url string) *gtfs.Static {
	err := loadDataOnce(url)
	if err != nil {
		fmt.Printf("error getting static gtfs: %v", err)
	}

	return staticGtfs
}

// singleton variables to ensure data is only loaded once.
var (
	staticGtfs *gtfs.Static
	loadOnce   sync.Once
)

func loadDataOnce(url string) error {
	var err error
	loadOnce.Do(func() {
		staticGtfs, err = loadData(url)
	})
	return err
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
