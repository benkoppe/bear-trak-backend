// Package static loads static harvard library content.
package static

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed library-rules.json
var libraryRuleBytes []byte

func GetLibraryData() LibraryData {
	var libraryData LibraryData

	err := json.Unmarshal(libraryRuleBytes, &libraryData)
	if err != nil {
		fmt.Printf("error unmarshalling static libraries: %v\n", err)
	}

	return libraryData
}
