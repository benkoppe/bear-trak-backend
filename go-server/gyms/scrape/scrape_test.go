package scrape

import "testing"

func TestTableScrape(t *testing.T) {
	FetchData("https://web.archive.org/web/20250211203409/https://scl.cornell.edu/recreation/cornell-fitness-centers")
}
