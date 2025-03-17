package scrape

import (
	"fmt"
	"testing"

	"github.com/benkoppe/bear-trak-backend/go-server/dining/umich/static"
)

func TestScrape(t *testing.T) {
	static := static.GetEateries()

	allEateries, _ := fetchAll("https://dining.umich.edu/menus-locations/", static)

	for key, val := range allEateries {
		fmt.Printf("Eatery: %d\n", key.ID)
		for _, eatery := range val {
			fmt.Print(eatery.Summary())
		}
	}
}
