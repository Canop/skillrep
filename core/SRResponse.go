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
	Rank    int
	Upvotes int
	Accepts int
	Score   int
	Id      int
	Name    string
	Profile string
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
	sql := `select g.Owner, g.u, g.a, g.s, player.name, player.profile from
	 (select answer.owner,
	  sum(least(10,answer.score)) u,
	  sum(Accepted::int) a,
	  10*sum(least(10,answer.score)) + 15*sum(Accepted::int) s
	  from answer join question on answer.question=question.id 
	  where Accepted is true and answer.owner!=question.owner
	  and question.closeddate=0
	  group by answer.owner
         ) as g
	 left join player on player.id=g.owner
	 where owner!=0`
	args := []interface{}{pageSize, (q.Page * pageSize)}
	if q.Search != "" {
		sql += `and name ~* $3`
		args = append(args, q.Search)
	}
	sql += ` order by s desc limit $1 offset $2`
	log.Println(sql)
	log.Printf("%#v\n", args)
	rows, err := db.Query(sql, args...)
	if err != nil {
		r.Error = err.Error()
		return r
	}
	i := pageSize * q.Page
	for rows.Next() {
		var u RankedUser
		i++
		u.Rank = i
		err = rows.Scan(&u.Id, &u.Upvotes, &u.Accepts, &u.Score, &u.Name, &u.Profile)
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
