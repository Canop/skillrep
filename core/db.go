package core

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", Config().DB.queryString())
	die(err)
}
