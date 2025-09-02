// Package scrape will scrape umich dining content.
package scrape

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/benkoppe/bear-trak-backend/go-server/utils/timeutils"
)

// unused -- a concurrent verion is in scrape_all
// func fetchEateryWeek(eateryURL string) ([]Eatery, error) {
// 	now := time.Now()
// 	var eateryWeek []Eatery
//
// 	for i := range [7]int{} {
// 		date := now.AddDate(0, 0, i)
// 		eatery, err := fetchEatery(eateryURL, date)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to fetch eatery for date %s: %w", date.Format("2006-01-02"), err)
// 		}
// 		eateryWeek = append(eateryWeek, *eatery)
// 	}
//
// 	return eateryWeek, nil
// }

func fetchEatery(eateryURL string, date time.Time) (*Eatery, error) {
	fullURL, err := appendDateSearchParam(eateryURL, date)
	if err != nil {
		return nil, fmt.Errorf("failed to append date search param: %w", err)
	}

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching external data: %w", err)
	}
	defer resp.Body.Close()

	eatery, err := scrape(resp.Body, date)
	if err != nil {
		return nil, fmt.Errorf("error scraping page: %w", err)
	}

	return eatery, nil
}

func appendDateSearchParam(eateryURL string, date time.Time) (string, error) {
	parsedURL, err := url.Parse(eateryURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse base URL: %w", err)
	}

	formattedDate := date.Format("2006-01-02")

	// query parameters
	query := parsedURL.Query()
	query.Set("menuDate", formattedDate)

	parsedURL.RawQuery = query.Encode()
	return parsedURL.String(), nil
}

func scrape(htmlReader io.Reader, date time.Time) (*Eatery, error) {
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return nil, err
	}

	eatery := extractEateryInfo(doc)

	eatery = extractHours(doc, eatery, date)

	eatery = extractMenu(doc, eatery)

	return eatery, nil
}

// removes extra whitespace and normalizes text
func cleanText(s string) string {
	// Replace all whitespace (including non-breaking space) with a single space
	whitespaceRegex := regexp.MustCompile(`[\s\p{Zs}]+`)
	return strings.TrimSpace(whitespaceRegex.ReplaceAllString(s, " "))
}

func extractEateryInfo(doc *goquery.Document) *Eatery {
	name := doc.Find("h2.postTitle .titleText").Text()
	address := doc.Find(".location-details .address").Text()
	phone := doc.Find(".location-details .phone").Text()
	email := doc.Find(".location-details .email").Text()

	eatery := NewEatery(cleanText(name), cleanText(address), cleanText(phone), cleanText(email))

	// Parse payment methods
	doc.Find(".payment-methods li").Each(func(i int, s *goquery.Selection) {
		paymentMethod := cleanText(s.Text())
		eatery.addPaymentMethod(paymentMethod)
	})

	return eatery
}

func extractHours(doc *goquery.Document, eatery *Eatery, date time.Time) *Eatery {
	// Find all meal periods and their times
	doc.Find(".calhours li").Each(func(i int, s *goquery.Selection) {
		mealName := cleanText(s.Find(".calhours-title").Text())
		mealTime := cleanText(s.Find(".calhours-times").Text())
		// Clean up the string
		mealTime = strings.ReplaceAll(mealTime, "‑", "-")      // Replace unicode dash with regular dash
		mealTime = strings.ReplaceAll(mealTime, "&nbsp;", " ") // Replace HTML spaces
		if mealTime == "" {
			return
		}

		parts := strings.Split(mealTime, "-")
		if len(parts) != 2 {
			log.Printf("invalid meal time format: %s", mealTime)
			return
		}

		start, err := timeutils.TimeString(cleanText(parts[0])).ToDate(date)
		if err != nil {
			log.Printf("error parsing start time: %v", err)
			return
		}
		end, err := timeutils.TimeString(cleanText(parts[1])).ToDate(date)
		if err != nil {
			log.Printf("error parsing end time: %v", err)
		}

		if mealName == "24/7 Kiosk" {
			return
		}

		eatery.addHours(mealName, start, end)
	})

	return eatery
}

func extractMenu(doc *goquery.Document, eatery *Eatery) *Eatery {
	container := doc.Find("#mdining-items")

	mealHeaders := container.Find("h3")
	if mealHeaders.Length() == 0 {
		parseMenu(container, "Menu", eatery)
		return eatery
	}

	mealHeaders.Each(func(_ int, hdr *goquery.Selection) {
		name := cleanText(
			strings.NewReplacer("+", "", "-", "").Replace(hdr.Text()),
		)
		parseMenu(hdr.Next(), name, eatery)
	})

	return eatery
}

func parseMenu(contentSel *goquery.Selection, mealName string, eatery *Eatery) {
	if contentSel == nil || contentSel.Length() == 0 {
		return
	}
	mealMenu := eatery.addMenu(mealName)

	contentSel.Find("li > h4").Each(func(_ int, station *goquery.Selection) {
		category := mealMenu.addCategory(cleanText(station.Text()))

		station.Parent().Find("ul.items > li").Each(func(_ int, item *goquery.Selection) {
			itemName := cleanText(item.Find(".item-name").Text())

			var traits []string
			item.Find("ul.traits li").Each(func(_ int, t *goquery.Selection) {
				if title := t.AttrOr("title", ""); title != "" {
					traits = append(traits, title)
				}
			})

			menuItem := category.AddMenuItem(itemName, traits)

			item.Find(".allergens li").Each(func(_ int, a *goquery.Selection) {
				menuItem.AddAllergen(cleanText(a.Text()))
			})
		})
	})
}
