package core

import (
	"log"
	"strings"
)

type Saver struct {
	questions chan *Question
}

func NewSaver(bufferSize int) *Saver {
	s := &Saver{}
	s.questions = make(chan *Question, bufferSize)
	return s
}

func (s *Saver) Run() {
	for {
		q, more := <-s.questions
		if !more {
			return
		}
		log.Printf("->Q:\"%.60s\"\n", q.Title)
		tx, err := db.Begin()
		die(err)
		sql := "delete from Answer where question=$1"
		_, err = tx.Exec(sql, q.Id)
		die(err)
		sql = "delete from Question where id=$1" // waiting for PG's upsert... soon...
		_, err = tx.Exec(sql, q.Id)
		die(err)
		sql = "insert into Question(Id, Title, CreationDate, Owner, Tags) values ($1, $2, $3, $4, $5)"
		_, err = tx.Exec(sql, q.Id, q.Title, q.CreationDate, q.Owner.Id, strings.Join(q.Tags, " "))
		die(err)
		for _, a := range q.Answers {
			if a.Owner.Id != 0 {
				sql = "update Player set Name=$1 where id=$2"
				r, err := tx.Exec(sql, a.Owner.Name, a.Owner.Id)
				die(err)
				n, _ := r.RowsAffected()
				if n == 0 { // yes, there's no upsert yet
					sql = "insert into Player(Id, Name) values($1,$2)"
					_, err = tx.Exec(sql, a.Owner.Id, a.Owner.Name)
					die(err)
					// log.Println("new player: ", a.Owner.Name)
				} else {
					// log.Println("Updated player: ", a.Owner.Name)
				}
			}
			sql = "insert into Answer(Id, Owner, Question, CreationDate, Accepted, Score) values ($1, $2, $3, $4, $5, $6)"
			_, err = tx.Exec(sql, a.Id, a.Owner.Id, q.Id, a.CreationDate, a.IsAccepted, a.Score)
			die(err)
		}
		tx.Commit()
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
	db.QueryRow(sql).Scan(&date)
	return date
}
