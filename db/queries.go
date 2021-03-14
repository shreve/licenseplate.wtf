package db

import (
	"database/sql"
	"embed"
	"sync"
)

//go:embed queries/*.sql
var rawQueries embed.FS
var Queries map[string]*sql.Stmt
var writeLock sync.Mutex

func Query(name string, args ...interface{}) (*sql.Rows, error) {
	return Queries[name].Query(args...)
}

func Exec(name string, args ...interface{}) (sql.Result, error) {
	writeLock.Lock()
	defer writeLock.Unlock()
	return Queries[name].Exec(args...)
}

func init() {
	Queries = make(map[string]*sql.Stmt)

	dirs, err := rawQueries.ReadDir("queries")
	if err != nil {
		panic(err)
	}
	for _, file := range dirs {
		name := file.Name()
		data, err := rawQueries.ReadFile("queries/" + name)
		if err != nil {
			panic(err)
		}
		name = name[0 : len(name)-4]
		stmt, err := DB.Prepare(string(data))
		if err != nil {
			panic(err)
		}
		Queries[name] = stmt
	}
}
