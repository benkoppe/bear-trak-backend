package scrape

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
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
	defer resp.Body.Close()

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

func scrapeTables(htmlReader io.Reader) ([]tableData, error) {
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return nil, err
	}

	var results []tableData
	doc.Find("table.striped").Each(func(i int, tableSel *goquery.Selection) {
		tableData := parseTable(tableSel)
		results = append(results, tableData)
	})

	return results, nil
}

// will parse a table and return as TableData
// for rows that span multiple columns, the same value is copied to each column
func parseTable(tableSel *goquery.Selection) tableData {
	var data tableData

	captionSel := tableSel.Find("caption")
	data.Caption = strings.TrimSpace(captionSel.Text())

	headers := make([]string, 0)
	tableSel.Find("thead tr").Each(func(_ int, tr *goquery.Selection) {
		tr.Find("th").Each(func(_ int, th *goquery.Selection) {
			headers = append(headers, strings.TrimSpace(th.Text()))
		})
	})
	data.Headers = headers

	// parse rows
	tableSel.Find("tbody tr").Each(func(_ int, rowSel *goquery.Selection) {
		rowColumns := make([]string, len(headers))
		nextCol := 0

		rowSel.Find("td").Each(func(_ int, td *goquery.Selection) {
			cellText := strings.TrimSpace(td.Text())
			colspanAttr, _ := td.Attr("colspan")
			colspan := 1
			if colspanAttr != "" {
				if c, err := strconv.Atoi(colspanAttr); err == nil {
					colspan = c
				}
			}
			for i := 0; i < colspan; i++ {
				if nextCol >= len(headers) {
					break
				}
				rowColumns[nextCol] = cellText
				nextCol++
			}
		})

		data.Rows = append(data.Rows, rowData{Columns: rowColumns})
	})

	return data
}
