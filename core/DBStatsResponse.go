// Skill Rep Response
package core

import ()

type DBStatsQuery struct {
}
type DBStatsResponse struct {
	Error   string
	DBStats DBStats
}
type DBStats struct {
	MaxQuestionCreationDate int64
	NbQuestions             int
	NbAnswers               int
}

func (q *DBStatsQuery) Answer() DBStatsResponse {
	var r DBStatsResponse
	sql := "select CreationDate from question order by id desc limit 1"
	db.QueryRow(sql).Scan(&r.DBStats.MaxQuestionCreationDate)
	sql = "select count(*) from question"
	db.QueryRow(sql).Scan(&r.DBStats.NbQuestions)
	sql = "select count(*) from answer"
	db.QueryRow(sql).Scan(&r.DBStats.NbAnswers)
	return r
}
