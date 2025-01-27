package static

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

func GetAlerts() []Alert {
	// ensure data is loaded
	err := loadData("./alerts/static/alerts.json")
	if err != nil {
		fmt.Printf("error getting eateries: %v", err)
	}

	return alerts
}

// singleton variables to ensure data is only loaded once.
var (
	alerts   []Alert
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

		err = json.Unmarshal(data, &alerts)
		if err != nil {
			err = fmt.Errorf("could not unmarshal JSON: %v", err)
			return
		}
	})
	return err
}
