package core

import (
	"log"
	"time"
)

func fetchSomeQuestions(startDate, endDate int64, saver *Saver) {
	fromDate := startDate
	nbQueriesMax := 10000
	maxCreationDate := fromDate
	page := 1
	for nbQueries := 1; nbQueries <= nbQueriesMax; nbQueries++ {
		log.Printf("Query %d / %d - page=%d \n", nbQueries, nbQueriesMax, page)
		r, err := GetQuestions("stackoverflow", fromDate, endDate, page)
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
		if r.Backoff > 0 {
			log.Printf("BACKOFF: %d\n", r.Backoff)
			time.Sleep(time.Duration(r.Backoff+2) * time.Second)
		}
	}
	saver.Done()
}

func Update() {
	log.Printf("Config: %#v\n", config)
	if config.DB.Name == "" {
		log.Fatal("bad config")
	}
	saver := NewSaver(500)
	endDate := int64(time.Now().Add(-3 * 24 * time.Hour).Unix())
	go fetchSomeQuestions(saver.MostRecentQuestionDate(), endDate, saver)
	saver.Run()
}
