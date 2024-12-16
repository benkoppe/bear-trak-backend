package external

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type stringifiedGymCapacity struct {
	Name        string
	Count       string
	LastUpdated string
	Percentage  string
}

type GymCapacity struct {
	Name        string
	Count       int64
	LastUpdated time.Time
	Percentage  *int64
}

const capacitiesTimeLayout = "01/02/2006 15:04 PM"

func FetchData(url string) ([]GymCapacity, error) {
	stringifiedCapacities, err := fetchData(url)
	if err != nil {
		return nil, fmt.Errorf("failed to load capacities: %v", err)
	}

	var capacities []GymCapacity
	for _, stringifiedCapacity := range stringifiedCapacities {
		capacity, err := stringifiedCapacity.toGymCapacity()
		if err != nil {
			return nil, fmt.Errorf("failed to parse capacity: %v", err)
		}
		capacities = append(capacities, *capacity)
	}

	return capacities, nil
}

func fetchData(url string) ([]stringifiedGymCapacity, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// load the html document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error loading html: %v", err)
	}

	barCharts := doc.Find(".barChart")
	if barCharts.Length() == 0 {
		return nil, fmt.Errorf("no .barChart elements found")
	}

	var capacities []stringifiedGymCapacity

	barCharts.Each(func(i int, selection *goquery.Selection) {
		fullText := strings.TrimSpace(selection.Text())

		if capacity, err := fromString(fullText); err == nil {
			capacities = append(capacities, *capacity)
		}
	})

	if len(capacities) == 0 {
		return nil, fmt.Errorf("no valid capacities found")
	}

	return capacities, nil
}

func fromString(str string) (*stringifiedGymCapacity, error) {
	pattern := regexp.MustCompile(`^(.*?)(?:Last Count: )(\d+)(?:Updated: )([\d/ :AMPamp]+)\s+(NA|\d+)%?$`)

	match := pattern.FindStringSubmatch(str)
	if match == nil || len(match) < 5 {
		return nil, fmt.Errorf("string does not match the expected pattern")
	}

	return &stringifiedGymCapacity{
		Name:        strings.TrimSpace(match[1]),
		Count:       match[2],
		LastUpdated: strings.TrimSpace(match[3]),
		Percentage:  match[4],
	}, nil
}

func (sgc stringifiedGymCapacity) toGymCapacity() (*GymCapacity, error) {
	count, err := strconv.ParseInt(sgc.Count, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Count: %v", err)
	}

	var percentage *int64
	if sgc.Percentage == "NA" {
		percentage = nil
	} else {
		parsedPercentage, err := strconv.ParseInt(sgc.Percentage, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Percentage: %v", err)
		}
		percentage = &parsedPercentage
	}

	estLocation, err := time.LoadLocation("America/New_York")
	if err != nil {
		return nil, fmt.Errorf("Error loading location: %v", err)
	}

	lastUpdated, err := time.ParseInLocation(capacitiesTimeLayout, sgc.LastUpdated, estLocation)
	if err != nil {
		return nil, fmt.Errorf("failed to parse LastUpdated: %v", err)
	}

	return &GymCapacity{
		Name:        sgc.Name,
		Count:       count,
		LastUpdated: lastUpdated,
		Percentage:  percentage,
	}, nil
}
