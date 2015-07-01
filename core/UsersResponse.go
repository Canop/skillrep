// Skill Rep Response
package core

import (
	"log"
)

type UsersQuery struct {
	Page   int
	Search string
}
type UsersResponse struct {
	Users []RankedUser
	Error string
}
type RankedUser struct {
	Id       int
	Rank     int
	SkillRep int
	Accepts  int
	Name     string
	Profile  string
}

func (q *UsersQuery) Fix() {
	if q.Page < 0 {
		q.Page = 0
	}
}

func (q *UsersQuery) Answer() UsersResponse {
	log.Printf("Processing %#v\n", q)
	pageSize := 20
	var r UsersResponse
	r.Users = make([]RankedUser, 0, pageSize)
	if q.Search == "" {
		// when there's no filtering, we don't have to query the rank,
		// which allows for a less expensive query
		i := (q.Page * pageSize)
		sql := `select 
			p.id,
			p.skillrep,
			p.name,
			p.profile
			from player p
			order by p.skillrep desc limit $1 offset $2`
		rows, err := db.TimedQuery(sql, pageSize, i)
		if err != nil {
			r.Error = err.Error()
			return r
		}
		for rows.Next() {
			var u RankedUser
			err = rows.Scan(&u.Id, &u.SkillRep, &u.Name, &u.Profile)
			i++
			u.Rank = i
			die(err)
			r.Users = append(r.Users, u)
		}
	} else {
		// we must query the rank
		sql := `select 
			p.id,
			(select 1 + count(*) from player op where op.skillrep>p.skillrep),
			p.skillrep,
			p.name,
			p.profile
			from player p where p.name ~* $3
			order by p.skillrep desc limit $1 offset $2`
		rows, err := db.TimedQuery(sql, pageSize, (q.Page * pageSize), q.Search)
		if err != nil {
			r.Error = err.Error()
			return r
		}
		for rows.Next() {
			var u RankedUser
			err = rows.Scan(&u.Id, &u.Rank, &u.SkillRep, &u.Name, &u.Profile)
			die(err)
			r.Users = append(r.Users, u)
		}
	}
	return r
}
