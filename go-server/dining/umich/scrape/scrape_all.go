package scrape

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/dining/umich/static"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

type Cache = *utils.Cache[map[*static.Eatery][]Eatery]

func InitCache(baseUrl string) Cache {
	return utils.NewCache(
		"diningScrape",
		time.Hour*24,
		func() (map[*static.Eatery][]Eatery, error) {
			staticEateries := static.GetEateries()
			return fetchAll(baseUrl, staticEateries)
		},
	)
}

func fetchAll(baseUrl string, eateries []static.Eatery) (map[*static.Eatery][]Eatery, error) {
	fetchedEateries := make(map[*static.Eatery][]Eatery)
	var mu sync.Mutex
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, 10)

	for _, eatery := range eateries {
		wg.Add(1)
		go func(e static.Eatery) {
			defer wg.Done()

			eateryUrl, err := utils.ExtendUrl(baseUrl, e.ScrapePath)
			if err != nil {
				log.Printf("error extending url for eatery %d: %v", e.ID, err)
				return
			}

			eateryWeek, err := fetchEateryWeekConcurrent(*eateryUrl, semaphore)
			if err != nil {
				log.Printf("error fetching eatery week for eatery %d: %v", e.ID, err)
				return
			}

			mu.Lock()
			fetchedEateries[&e] = eateryWeek
			mu.Unlock()
		}(eatery)
	}

	wg.Wait()
	return fetchedEateries, nil
}

func fetchEateryWeekConcurrent(eateryUrl string, semaphore chan struct{}) ([]Eatery, error) {
	now := time.Now()
	eateryWeek := make([]Eatery, 7)

	var wg sync.WaitGroup
	var errMu sync.Mutex
	var firstErr error

	for i := 0; i < 7; i++ {
		wg.Add(1)
		go func(dayOffset int) {
			defer wg.Done()

			// acquire semaphore slot, or wait
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			date := now.AddDate(0, 0, dayOffset)
			eatery, err := fetchEatery(eateryUrl, date)
			if err != nil {
				errMu.Lock()
				if firstErr == nil {
					firstErr = fmt.Errorf("failed to fetch eatery for date %s: %w", date.Format("2006-01-02"), err)
				}
				errMu.Unlock()
				return
			}

			eateryWeek[dayOffset] = *eatery
		}(i)
	}

	wg.Wait()
	if firstErr != nil {
		return nil, firstErr
	}

	return eateryWeek, nil
}
