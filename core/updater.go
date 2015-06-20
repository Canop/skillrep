package core

import (
	"log"
	"runtime"
	"time"
)

func fetchSomeQuestions(startDate, endDate int64, saver *Saver) {
	runtime.GOMAXPROCS(4)
	log.Println("max proc:", runtime.GOMAXPROCS(0))
	fromDate := startDate
	nbQueriesMax := 50000
	maxCreationDate := fromDate
	page := 1
	for nbQueries := 1; nbQueries <= nbQueriesMax; nbQueries++ {
		log.Printf("Query %d / %d - page=%d \n", nbQueries, nbQueriesMax, page)
		r, err := GetQuestions("stackoverflow", fromDate, endDate, page, 3)
		die(err)
		for _, q := range r.Questions {
			saver.AddQuestion(q)
			if q.CreationDate > fromDate {
				maxCreationDate = q.CreationDate
			}
		}
		if r.HasMore {
			page++
		} else {
			page = 1
			fromDate = maxCreationDate
		}
		log.Println("Quota remaining:", r.QuotaRemaining)
		var wait time.Duration
		if r.QuotaRemaining < 25 {
			log.Println("Let's sleep an hour...")
			wait = time.Hour
		} else if r.Backoff > 0 {
			log.Printf("BACKOFF: %d\n", r.Backoff)
			wait = time.Duration(r.Backoff+2) * time.Second
		} else {
			wait = 200 * time.Millisecond // let's be gentle...
		}
		time.Sleep(wait)
	}
	saver.Done()
}

func Update() {
	log.Printf("Config: %#v\n", config)
	if config.DB.Name == "" {
		log.Fatal("bad config")
	}
	saver := NewSaver(500)
	endDate := int64(time.Now().Add(-7 * 24 * time.Hour).Unix())
	go fetchSomeQuestions(saver.MostRecentQuestionDate(), endDate, saver)
	saver.Run()
}
