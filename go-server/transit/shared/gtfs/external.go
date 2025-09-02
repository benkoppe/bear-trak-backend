// Package gtfs handles gtfs data.
package gtfs

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
	"github.com/jamespfennell/gtfs"
)

type Cache = *utils.Cache[*gtfs.Static]

func InitCache(source string) Cache {
	return utils.NewCache(
		"transitExternalGtfs",
		utils.NoExpiration,
		func() (*gtfs.Static, error) {
			return loadData(source)
		},
	)
}

func loadData(source string) (*gtfs.Static, error) {
	tcatGtfsData, err := loadTcatGtfs(source)
	if err != nil {
		return nil, fmt.Errorf("error loading tcat data: %v", err)
	}

	staticData, err := gtfs.ParseStatic(tcatGtfsData, gtfs.ParseStaticOptions{})
	if err != nil {
		return nil, fmt.Errorf("error parsing tcat data: %v", err)
	}

	return staticData, nil
}

func loadTcatGtfs(source string) ([]byte, error) {
	var reader io.ReadCloser
	var err error

	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		// Load from HTTP
		resp, err := http.Get(source)
		if err != nil {
			return nil, fmt.Errorf("error fetching gtfs ZIP: %v", err)
		}
		reader = resp.Body
	} else {
		// Load from file
		reader, err = os.Open(source)
		if err != nil {
			return nil, fmt.Errorf("error opening gtfs ZIP file: %v", err)
		}
	}

	defer reader.Close()

	originalGtfsData, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading gtfs ZIP: %v", err)
	}

	return originalGtfsData, nil
}
