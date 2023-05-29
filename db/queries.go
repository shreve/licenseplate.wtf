package db

import (
	"database/sql"
	"embed"
	"log"
)

//go:embed queries/*.sql
var rawQueries embed.FS
var Queries map[string]*sql.Stmt

func Query(name string, args ...interface{}) (*sql.Rows, error) {
	log.Println("Query", name, args)
	Lock.RLock()
	defer Lock.RUnlock()
	return Queries[name].Query(args...)
}

func Exec(name string, args ...interface{}) (sql.Result, error) {
	Lock.Lock()
	defer Lock.Unlock()
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
