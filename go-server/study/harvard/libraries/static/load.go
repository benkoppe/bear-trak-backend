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
		fmt.Printf("error unmarshalling static libraries: %v", err)
	}

	return libraryData
}
