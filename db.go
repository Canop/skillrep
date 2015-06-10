package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type Saver struct {
	questions chan *Question
}

func NewSaver(bufferSize int) *Saver {
	s := &Saver{}
	s.questions = make(chan *Question, bufferSize)
	return s
}

// die is an acronym for "die if error"
func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Saver) Run() {
	db, err := sql.Open("postgres", config.DB.queryString())
	die(err)
	log.Println(db)
	for {
		q, more := <-s.questions
		if !more {
			return
		}
		log.Printf("Saving %s\n", q.Title)
		var qid int
		sql := "insert into Question(Id, Title, CreationDate, Owner) values ($1, $2, $3, $4) returning Id"
		err = db.QueryRow(sql, q.Id, q.Title, q.CreationDate, q.Owner.Id).Scan(&qid)
		die(err)
		log.Printf("Inserted Question Id: %d\n", qid)
	}
}

func (s *Saver) AddQuestion(q *Question) {
	s.questions <- q
}

func (s *Saver) Done() {
	close(s.questions)
}
