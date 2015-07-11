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
	NbPlayers               int
}

func (q *DBStatsQuery) Answer() DBStatsResponse {
	var r DBStatsResponse
	sql := "select CreationDate from question order by id desc limit 1"
	db.QueryRow(sql).Scan(&r.DBStats.MaxQuestionCreationDate)
	sql = "select reltuples::int from pg_class where relname='player'"
	db.QueryRow(sql).Scan(&r.DBStats.NbPlayers)
	sql = "select reltuples::int from pg_class where relname='question'"
	db.QueryRow(sql).Scan(&r.DBStats.NbQuestions)
	sql = "select reltuples::int from pg_class where relname='answer'"
	db.QueryRow(sql).Scan(&r.DBStats.NbAnswers)
	return r
}
