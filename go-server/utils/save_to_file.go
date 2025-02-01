package utils

import (
	"fmt"
	"os"
)

func SaveToFile(data []byte, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to create file %s: %w", path, err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %w", path, err)
	}

	err = file.Sync()
	if err != nil {
		return fmt.Errorf("error syncing file %s: %w", path, err)
	}

	return nil
}
