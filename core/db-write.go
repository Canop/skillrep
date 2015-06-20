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
		// log.Printf("->Q:\"%.60s\"\n", q.Title)
		tx, err := db.Begin()
		die(err)
		sql := "delete from Answer where question=$1"
		_, err = tx.Exec(sql, q.Id)
		die(err)
		sql = "delete from Question where id=$1" // waiting for PG's upsert... soon...
		_, err = tx.Exec(sql, q.Id)
		die(err)
		sql = "insert into Question(Id, Title, CreationDate, ClosedDate, Owner, Tags) values ($1,$2,$3,$4,$5,$6)"
		_, err = tx.Exec(sql, q.Id, q.Title, q.CreationDate, q.ClosedDate, q.Owner.Id, strings.Join(q.Tags, " "))
		die(err)
		for _, a := range q.Answers {
			answerSkillRep := 0
			if a.IsAccepted {
				answerSkillRep += 15
			}
			if a.Score >= 10 {
				answerSkillRep += 100
			} else {
				answerSkillRep += 10 * a.Score
			}
			if a.Owner.Id == q.Owner.Id {
				answerSkillRep = 0
			}
			if q.ClosedDate != 0 {
				answerSkillRep = 0
			}
			if a.Owner.Id != 0 {
				sql = `update Player set Name=$1,
					SkillRep=coalesce((select $3+sum(answer.SkillRep) from answer where owner=$2),0)
					where Player.Id=$2`
				r, err := tx.Exec(sql, a.Owner.Name, a.Owner.Id, answerSkillRep)
				if err != nil {
					log.Println(sql, a.Owner.Name, a.Owner.Id, answerSkillRep)
				}
				die(err)
				n, _ := r.RowsAffected()
				if n == 0 { // yes, there's no upsert yet
					sql = "insert into Player(Id, Name, Profile, SkillRep) values($1,$2,$3,$4)"
					_, err = tx.Exec(sql, a.Owner.Id, a.Owner.Name, a.Owner.Profile, answerSkillRep)
					die(err)
				}
			}
			sql = "insert into Answer(Id, Owner, Question, CreationDate, Accepted, Score, SkillRep) values ($1,$2,$3,$4,$5,$6,$7)"
			_, err = tx.Exec(sql, a.Id, a.Owner.Id, q.Id, a.CreationDate, a.IsAccepted, a.Score, answerSkillRep)
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
