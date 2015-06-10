package main

import (
	"log"
)

func main() {
	ReadConfig()

	log.Printf("Config: %#v\n", config)

	if config.DB.Name == "" {
		log.Fatal("bad config")
	}

	saver := NewSaver(500)

	r, err := GetQuestion("stackoverflow", 11353679)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("BACKOFF: %#v\n", r.Backoff)
	for _, q := range r.Questions {
		saver.AddQuestion(q)
	}
	saver.Done()
	saver.Run()

}