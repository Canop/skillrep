package main

import (
	"log"
	"time"
)

// die is an acronym for "die if error"
func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// fetches one question
func fetchOneQuestion(qid int) {
	saver := NewSaver(50)
	r, err := GetQuestion("stackoverflow", 11353679)
	die(err)
	log.Printf("BACKOFF: %#v\n", r.Backoff)
	for _, q := range r.Questions {
		saver.AddQuestion(q)
	}
	saver.Done()
	saver.Run()
}

func fetchSomeQuestions(startDate, endDate int64, saver *Saver) {
	fromDate := startDate
	nbQueriesMax := 100
	for nbQueries := 0; nbQueries < nbQueriesMax; nbQueries++ {
		log.Printf("Query %d / %d \n", nbQueries, nbQueriesMax)
		r, err := GetQuestions("stackoverflow", fromDate, endDate, 0)
		die(err)
		for _, q := range r.Questions {
			saver.AddQuestion(q)
			if q.CreationDate > fromDate {
				fromDate = q.CreationDate
			}
		}
		if r.Backoff > 0 {
			log.Printf("BACKOFF: %d\n", r.Backoff)
			time.Sleep(time.Duration(r.Backoff) * time.Second)
		}
	}
	saver.Done()
}
func main() {
	ReadConfig()
	log.Printf("Config: %#v\n", config)
	if config.DB.Name == "" {
		log.Fatal("bad config")
	}

	saver := NewSaver(500)
	endDate := int64(time.Now().Add(-3 * 24 * time.Hour).Unix())
	go fetchSomeQuestions(saver.MostRecentQuestionDate(), endDate, saver)
	saver.Run()

}
