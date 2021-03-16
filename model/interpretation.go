package model

import (
	"database/sql"
	"licenseplate.wtf/db"
	"log"
	"net/http"
	"time"
)

type Interpretation struct {
	Id        int64
	PlateId   int64
	What      string
	Why       string
	Ip        string
	Username  string
	CreatedAt time.Time
}

func (i *Interpretation) FromRequest(req *http.Request) []string {
	errors := []string{}

	i.What = req.Form["what"][0]
	i.Why = req.Form["why"][0]
	i.Ip = req.RemoteAddr

	if len(i.What) < 5 {
		errors = append(errors, "Your interpretation is too short. It must be at least 5 letters.")
	}

	if len(i.What) > 30 {
		errors = append(errors, "Your interpretation is too long. It must be at most 30 letters.")
	}

	if len(i.Why) < 10 {
		errors = append(errors, "Your explanation is too short. It must be at least 10 letters.")
	}

	if len(i.Why) > 300 {
		errors = append(errors, "Your explanation is too long. It must be at most 300 letters.")
	}

	if len(req.Form["responsibility"]) == 0 {
		errors = append(errors, "You must accept that you will be responsible.")
	}

	return errors
}

func (i *Interpretation) FromRow(row *sql.Rows) {
	var createdAt string
	row.Scan(&i.Id, &i.What, &i.Why, &i.Username, &createdAt)
	i.CreatedAt, _ = time.Parse(db.ISO8601, createdAt)
}

func (i *Interpretation) Create() bool {
	i.Username = NameHash(i.Ip)

	res, err := db.Exec("insert_interpretation", i.What, i.Why, i.Ip, i.Username, i.PlateId)
	if err != nil {
		log.Println(err)
		return false
	}
	i.Id, _ = res.LastInsertId()

	// Little cheat code
	i.CreatedAt = time.Now().UTC().Truncate(time.Second)
	return true
}
