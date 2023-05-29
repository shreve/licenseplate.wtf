package db

import (
	"database/sql"
	"embed"
	"log"
)

//go:embed queries/*.sql
var rawQueries embed.FS
var Queries = make(map[string]*sql.Stmt)

func Query(name string, args ...interface{}) (*sql.Rows, error) {
	log.Println("Query", name, args)
	Lock.RLock()
	defer Lock.RUnlock()
	return GetQuery(name).Query(args...)
}

func Exec(name string, args ...interface{}) (sql.Result, error) {
	log.Println("Exec", name, args)
	Lock.Lock()
	defer Lock.Unlock()
	return GetQuery(name).Exec(args...)
}

func GetQuery(name string) *sql.Stmt {
	// Read and prepare the query if it doesn't exist
	query, ok := Queries[name]
	if !ok {
		data, err := rawQueries.ReadFile("queries/" + name + ".sql")
		if err != nil {
			panic(err)
		}
		query, err = DB.Prepare(string(data))
		if err != nil {
			panic(err)
		}
		Queries[name] = query
	}
	return query
}
