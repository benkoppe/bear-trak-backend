// Package printers loads cornell printer data.
package printers

import (
	"fmt"
	"strings"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/schools/cornell/external_map"
)

func Get(mapCache external_map.Cache) ([]api.Printer, error) {
	mapItems, err := mapCache.Get("CUPrint")
	if err != nil {
		return nil, fmt.Errorf("error loading map data: %v", err)
	}

	var printers []api.Printer
	for _, mapItem := range mapItems {
		printers = append(printers, convertExternalPrinter(mapItem))
	}

	return printers, nil
}

func convertExternalPrinter(mapItem external_map.Item) api.Printer {
	printer := api.Printer{
		Latitude:  mapItem.LatLng.Latitude,
		Longitude: mapItem.LatLng.Longitude,
	}

	if mapItem.Location != "" {
		printer.Location = api.NewNilString(mapItem.Location)
	} else {
		printer.Location = api.NilString{Null: true}
	}

	typeDescription, room := splitAtFirstDash(mapItem.Description)
	switch typeDescription {
	case "Black & White":
		printer.Type = api.PrinterTypeBlackAndWhite
	case "Color":
		printer.Type = api.PrinterTypeColor
	case "Color, Scan, & Copy":
		printer.Type = api.PrinterTypeColorScanCopy
	default:
		printer.Type = api.PrinterTypeUnknown
	}

	if room != "" {
		printer.Room = api.NewNilString(room)
	} else {
		printer.Room = api.NilString{Null: true}
	}

	return printer
}

func splitAtFirstDash(s string) (string, string) {
	index := strings.Index(s, "-")
	if index == -1 {
		// No dash found, return the original string as the first part and an empty second part
		return s, ""
	}
	// Trim spaces around the split parts
	return strings.TrimSpace(s[:index]), strings.TrimSpace(s[index+1:])
}
