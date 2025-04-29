package libraries

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/study/harvard/libraries/external"
	"github.com/benkoppe/bear-trak-backend/go-server/study/harvard/libraries/static"
	"github.com/benkoppe/bear-trak-backend/go-server/study/shared/libcal"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

func Get(cache external.Cache) ([]api.Library, error) {
	staticData := static.GetLibraryData()

	externalData, err := cache.Get()
	if err != nil {
		return nil, fmt.Errorf("error loading external data: %v", err)
	}

	var libraries []api.Library
	for _, externalLibrary := range externalData {
		library, err := convertExternalLibrary(staticData, externalLibrary)
		if err != nil {
			log.Printf("error converting external library: %v", err)
			continue
		}

		if library == nil {
			continue
		}

		libraries = append(libraries, *library)
	}

	return libraries, nil
}

func convertExternalLibrary(static static.LibraryData, external external.Library) (*api.Library, error) {
	for _, excludeId := range static.ExclusionIDs {
		if external.ID == excludeId {
			// skip detected
			return nil, nil
		}
	}

	library := api.Library{
		Name:             external.Name,
		Latitude:         external.Coordinates.Latitude,
		Longitude:        external.Coordinates.Longitude,
		ImagePath:        getImagePath(external),
		PrinterLocations: []string{},
	}

	id, err := strconv.Atoi(external.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert ID to int: %v", err)
	}
	library.ID = id

	if len(external.WeeksHours.Locations) != 1 {
		return nil, fmt.Errorf("expected 1 weeks_hours location, got=%d", len(external.WeeksHours.Locations))
	}
	details := external.WeeksHours.Locations[0]

	cardAccess := false
	for _, cardAccessId := range static.CardAccessIDs {
		if external.ID == cardAccessId {
			cardAccess = true
			break
		}
	}

	hours, err := libcal.ConvertToHours(details.Weeks)
	if err != nil {
		return nil, fmt.Errorf("failed to convert libcal hours: %v", err)
	}
	if !cardAccess {
		library.Hours = hours
	} else {
		library.CardAccessHours = hours
	}

	return &library, nil
}

func getImagePath(external external.Library) string {
	name := external.Name

	lowercased := strings.ToLower(name)

	// filter characters
	var builder strings.Builder
	for _, r := range lowercased {
		// filter marks, and only let letters, numbers, and whitespace through
		if !unicode.IsMark(r) && (unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r)) {
			builder.WriteRune(r)
		}
	}
	stripped := builder.String()

	// regex pattern to match with whitespace
	re := regexp.MustCompile(`\s+`)

	// replace all whitespace with underscores
	imageName := re.ReplaceAllString(stripped, "_")
	return utils.ImageNameToPath("study/harvard", imageName)
}
