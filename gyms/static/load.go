package static

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

func GetGyms() []Gym {
	// ensure data is loaded
	err := loadData("./gyms/static/gyms.json")
	if err != nil {
		fmt.Printf("error getting gyms: %v", err)
	}

	return gyms
}

// singleton variables to ensure data is only loaded once.
var (
	gyms     []Gym
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
			err = fmt.Errorf("coult not read file: %v", err)
			return
		}

		err = json.Unmarshal(data, &gyms)
		if err != nil {
			err = fmt.Errorf("coult not unmarshal JSON: %v", err)
			return
		}
	})
	return err
}
