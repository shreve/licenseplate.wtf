package db

import (
	"database/sql"
	_ "embed"
	"sync"

	// _ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

var DB *sql.DB
var Lock sync.RWMutex

const LOCATION = "data.sqlite"
const ISO8601 = "2006-01-02 15:04:05"
const KEYFMT = "backup/data-2006-01-02-15-04-05.sqlite"

//go:embed schema.sql
var schemaSql string

func init() {
	var err error
	DB, err = sql.Open("sqlite", LOCATION)
	if err != nil {
		panic(err)
	}
}
