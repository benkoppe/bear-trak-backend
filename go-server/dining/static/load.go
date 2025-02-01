package static

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

func GetEateries() []Eatery {
	// ensure data is loaded
	err := loadData("./dining/static/eateries.json")
	if err != nil {
		fmt.Printf("error getting eateries: %v", err)
	}

	return eateries
}

// singleton variables to ensure data is only loaded once.
var (
	eateries []Eatery
	loadOnce sync.Once
)

func loadData(filePath string) error {
	var err error
	loadOnce.Do(func() {
		file, e := os.Open(filePath)
		if e != nil {
			err = fmt.Errorf("could not open file: %v", err)
			return
		}
		defer file.Close()

		data, e := io.ReadAll(file)
		if e != nil {
			err = fmt.Errorf("could not read file: %v", err)
			return
		}

		err = json.Unmarshal(data, &eateries)
		if err != nil {
			err = fmt.Errorf("could not unmarshal JSON: %v", err)
			return
		}
	})
	return err
}
