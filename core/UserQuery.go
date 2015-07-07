// Skill Rep Response
package core

import (
	"log"
)

type UserQuery struct {
	UserId int
}
type UserResponse struct {
	User  RankedUser
	Error string
}

func (q *UserQuery) Answer() UserResponse {
	log.Printf("Processing %#v\n", q)
	var u RankedUser
	sql := `select 
		p.id,
		(select 1 + count(*) from player op where op.skillrep>p.skillrep),
		(select count(*) from answer where owner=$1 and accepted is true),
		p.skillrep,
		p.name,
		p.profile
		from player p where id=$1`
	row := db.QueryRow(sql, q.UserId)
	err := row.Scan(&u.Id, &u.Rank, &u.Accepts, &u.SkillRep, &u.Name, &u.Profile)
	if err != nil {
		return UserResponse{Error: err.Error()}
	} else {
		return UserResponse{User: u}
	}
}
