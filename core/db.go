package core

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var db TimedDB

type TimedDB struct {
	*sql.DB
}

func (tdb TimedDB) TimedQuery(query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	r, err := tdb.Query(query, args...)
	duration := time.Since(start)
	if duration.Seconds() > .1 {
		log.Println("Long Query duration:", duration)
		log.Println(query)
		log.Println(args...)
	}
	return r, err
}

func init() {
	sqldb, err := sql.Open("postgres", Config().DB.queryString())
	die(err)
	db = TimedDB{sqldb}
}
