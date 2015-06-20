// Skill Rep Response
package core

import (
	"database/sql"
	"log"
)

type SRQuery struct {
	Page   int
	Search string
}
type SRResponse struct {
	Users   []RankedUser
	Error   string
	DBStats DBStats
}
type RankedUser struct {
	Id       int
	Rank     int
	Accepts  int
	SkillRep int
	Name     string
	Profile  string
}
type DBStats struct {
	MaxQuestionCreationDate int64
	NbQuestions             int
	NbAnswers               int
}

func (q *SRQuery) Fix() {
	if q.Page < 0 {
		q.Page = 0
	}
}

func (q *SRQuery) Answer() SRResponse {
	log.Printf("Processing %#v\n", q)
	pageSize := 20
	db, err := sql.Open("postgres", config.DB.queryString())
	die(err)
	defer db.Close()
	var r SRResponse
	r.Users = make([]RankedUser, 0, pageSize)
	sql := `select 
		p.id,
		(select 1 + count(*) from player op where op.skillrep>p.skillrep),
		(select count(*) from answer where owner=p.id and accepted is true),
		p.skillrep,
		p.name,
		p.profile
		from player p`
	args := []interface{}{pageSize, (q.Page * pageSize)}
	if q.Search != "" {
		sql += ` where p.name ~* $3`
		args = append(args, q.Search)
	}
	sql += ` order by p.skillrep desc limit $1 offset $2`
	log.Println(sql)
	log.Printf("%#v\n", args)
	rows, err := db.Query(sql, args...)
	if err != nil {
		r.Error = err.Error()
		return r
	}
	for rows.Next() {
		var u RankedUser
		err = rows.Scan(&u.Id, &u.Rank, &u.Accepts, &u.SkillRep, &u.Name, &u.Profile)
		die(err)
		r.Users = append(r.Users, u)
	}
	sql = "select max(CreationDate) from Question"
	db.QueryRow(sql).Scan(&r.DBStats.MaxQuestionCreationDate)
	sql = "select count(*) from question"
	db.QueryRow(sql).Scan(&r.DBStats.NbQuestions)
	sql = "select count(*) from answer"
	db.QueryRow(sql).Scan(&r.DBStats.NbAnswers)
	return r
}
