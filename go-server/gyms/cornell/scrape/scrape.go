// Package scrape loads all scraped cornell gym content.
package scrape

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

type Cache = *utils.Cache[[]ParsedSchedule]

func InitCache(url string) Cache {
	return utils.NewCache(
		"gymScrape",
		1*time.Minute,
		func() ([]ParsedSchedule, error) {
			return fetchData(url)
		})
}

func fetchData(url string) ([]ParsedSchedule, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching external data: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	tables, err := scrapeTables(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error scraping tables: %w", err)
	}

	var schedules []ParsedSchedule
	for _, table := range tables {
		schedule := parseSchedule(table)
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

func scrapeTables(htmlReader io.Reader) ([]utils.TableData, error) {
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return nil, err
	}

	var results []utils.TableData
	doc.Find("table.striped").Each(func(i int, tableSel *goquery.Selection) {
		tableData := utils.ScrapeTable(tableSel)
		results = append(results, tableData)
	})

	return results, nil
}
