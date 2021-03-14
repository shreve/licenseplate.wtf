package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

const ISO8601 = "2006-01-02 15:04:05"

func init() {
	var err error
	DB, err = sql.Open("sqlite3", "data.sqlite")
	if err != nil {
		panic(err)
	}
}
