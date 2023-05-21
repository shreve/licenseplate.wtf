package db

import (
	"database/sql"
	_ "embed"

	// _ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

const ISO8601 = "2006-01-02 15:04:05"

//go:embed schema.sql
var schemaSql string

func init() {
	var err error
	DB, err = sql.Open("sqlite", "data.sqlite")
	if err != nil {
		panic(err)
	}
	_, err = DB.Exec(schemaSql)
	if err != nil {
		panic(err)
	}
}
