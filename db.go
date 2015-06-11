package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type Saver struct {
	questions chan *Question
	db        *sql.DB
}

func NewSaver(bufferSize int) *Saver {
	s := &Saver{}
	s.questions = make(chan *Question, bufferSize)
	var err error
	s.db, err = sql.Open("postgres", config.DB.queryString())
	die(err)
	return s
}

func (s *Saver) Run() {
	for {
		q, more := <-s.questions
		if !more {
			return
		}
		log.Printf("->Q:\"%.60s\"\n", q.Title)
		sql := "delete from Answer where question=$1"
		_, err := s.db.Exec(sql, q.Id)
		die(err)
		sql = "delete from Question where id=$1" // waiting for PG's upsert... soon...
		_, err = s.db.Exec(sql, q.Id)
		die(err)
		sql = "insert into Question(Id, Title, CreationDate, Owner) values ($1, $2, $3, $4)"
		_, err = s.db.Exec(sql, q.Id, q.Title, q.CreationDate, q.Owner.Id)
		die(err)
		for _, a := range q.Answers {
			sql = "insert into Answer(Id, Owner, Question, CreationDate, Score) values ($1, $2, $3, $4, $5)"
			_, err = s.db.Exec(sql, a.Id, a.Owner.Id, q.Id, a.CreationDate, a.Score)
			die(err)
		}
	}
}

func (s *Saver) AddQuestion(q *Question) {
	s.questions <- q
}

func (s *Saver) Done() {
	close(s.questions)
}

func (s *Saver) MostRecentQuestionDate() int64 {
	sql := "select max(CreationDate) from Question"
	var date int64
	s.db.QueryRow(sql).Scan(&date)
	return date
}
