package external_gtfs

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/jamespfennell/gtfs"
)

func GetStaticGtfs(url string) *gtfs.Static {
	err := loadDataOnce(url)
	if err != nil {
		fmt.Printf("error getting static gtfs: %v", err)
	}

	return staticGtfs
}

// singleton variables to ensure data is only loaded once.
var (
	staticGtfs *gtfs.Static
	loadOnce   sync.Once
)

func loadDataOnce(url string) error {
	var err error
	loadOnce.Do(func() {
		staticGtfs, err = loadData(url)
	})
	return err
}

func loadData(url string) (*gtfs.Static, error) {
	tcatGtfsData, err := loadTcatGtfs(url)
	if err != nil {
		return nil, fmt.Errorf("error loading tcat data: %v", err)
	}

	staticData, err := gtfs.ParseStatic(tcatGtfsData, gtfs.ParseStaticOptions{})
	if err != nil {
		return nil, fmt.Errorf("error parsing tcat data: %v", err)
	}

	return staticData, nil
}

func loadTcatGtfs(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching gtfs ZIP: %v", err)
	}
	defer resp.Body.Close()

	originalGtfsData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading gtfs ZIP: %v", err)
	}

	convertedData, err := convertTcatGtfs(originalGtfsData)
	if err != nil {
		return nil, fmt.Errorf("error converting original gtfs ZIP: %v", err)
	}

	return convertedData, nil
}

// tcat ZIP is technically in the wrong gtfs format.
// it contains a folder zipped containing files, but the zip should just contain the files and no folder
// this function converts a loaded tcat zip to the proper format
func convertTcatGtfs(originalZipData []byte) ([]byte, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(originalZipData), int64(len(originalZipData)))
	if err != nil {
		return nil, fmt.Errorf("error creating ZIP reader: %v", err)
	}

	var newZipBuffer bytes.Buffer
	zipWriter := zip.NewWriter(&newZipBuffer)

	for _, file := range zipReader.File {
		newName := stripRootFolder(file.Name)
		if newName == "" {
			continue
		}

		originalFile, err := file.Open()
		if err != nil {
			zipWriter.Close()
			return nil, fmt.Errorf("Error opening file %s: %v", file.Name, err)
		}

		newFileHeader := &zip.FileHeader{
			Name:     newName,
			Method:   file.Method,
			Modified: file.Modified,
		}
		newFileHeader.SetMode(file.Mode())

		newFileWriter, err := zipWriter.CreateHeader(newFileHeader)
		if err != nil {
			originalFile.Close()
			zipWriter.Close()
			return nil, fmt.Errorf("error creating file %s in new ZIP: %v", newName, err)
		}

		_, err = io.Copy(newFileWriter, originalFile)
		if err != nil {
			originalFile.Close()
			zipWriter.Close()
			return nil, fmt.Errorf("error copying file %s: %v", file.Name, err)
		}

		originalFile.Close()
	}

	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("error closing new ZIP writer: %v", err)
	}

	return newZipBuffer.Bytes(), nil
}

func stripRootFolder(name string) string {
	parts := strings.SplitN(name, "/", 2)
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}
