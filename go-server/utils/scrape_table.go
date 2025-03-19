package utils

import (
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type TableData struct {
	Caption string
	Headers []string
	Rows    []rowData
}

type rowData struct {
	Columns []string
}

func ScrapeTable(sel *goquery.Selection) TableData {
	var data TableData

	captionSel := sel.Find("caption")
	data.Caption = CleanString(captionSel.Text())

	headers := make([]string, 0)
	sel.Find("thead tr").Each(func(_ int, tr *goquery.Selection) {
		tr.Find("th").Each(func(_ int, th *goquery.Selection) {
			headers = append(headers, CleanString(th.Text()))
		})
	})
	data.Headers = headers

	sel.Find("tbody tr").Each(func(_ int, rowSel *goquery.Selection) {
		rowColumns := make([]string, len(headers))
		nextCol := 0

		rowSel.Find("td").Each(func(_ int, td *goquery.Selection) {
			cellText := CleanString(td.Text())
			// for rows that span multiple columns, the same value is copied to each column
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
