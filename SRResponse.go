// Skill Rep Response
package core

import (
	"database/sql"
)

type SRQuery struct {
	Page int
}
type SRResponse struct {
	Users []RankedUser
}
type RankedUser struct {
	Rank     int
	SkillRep int
	Id       int
	Name     string
}

func (q *SRQuery) Fix() {
	if q.Page < 0 {
		q.Page = 0
	}
}

func (q *SRQuery) Answer() SRResponse {
	pageSize := 50
	db, err := sql.Open("postgres", config.DB.queryString())
	die(err)
	defer db.Close()
	var r SRResponse
	r.Users = make([]RankedUser, 0, pageSize)
	sql := `select owner, sum(least(10,score)) s,
	 (select name from player where id=owner)
	 from answer where owner!=0
	 group by owner order by s desc limit $1 offset $2`
	rows, err := db.Query(sql, pageSize, (q.Page * pageSize))
	die(err)
	i := pageSize * q.Page
	for rows.Next() {
		var u RankedUser
		i++
		u.Rank = i
		err = rows.Scan(&u.Id, &u.SkillRep, &u.Name)
		r.Users = append(r.Users, u)
	}
	return r
}
