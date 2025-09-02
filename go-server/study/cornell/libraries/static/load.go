// Package static loads static cornell library data.
package static

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed libraries.json
var libraryBytes []byte

func GetLibraries() []Library {
	var libraries []Library

	err := json.Unmarshal(libraryBytes, &libraries)
	if err != nil {
		fmt.Printf("error unmarshalling static libraries: %v", err)
	}

	return libraries
}
