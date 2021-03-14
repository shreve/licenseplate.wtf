package model

import (
	"licenseplate.wtf/db"
	"log"
	"strings"
	"time"
)

type Plate struct {
	id        int64
	Code      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPlate(code string) *Plate {
	return &Plate{Code: strings.ToUpper(code)}
}

func (p *Plate) Valid() bool {
	if len(p.Code) > 10 || len(p.Code) < 3 {
		return false
	}
	return true
}

func (p *Plate) FindOrCreate() bool {
	rows, err := db.Query("find_plate", p.Code)
	if err != nil || !rows.Next() {
		return p.Create()
	}
	var createdAt string
	var updatedAt string
	err = rows.Scan(&p.id, &p.Code, &createdAt, &updatedAt)
	if err != nil {
		log.Println("Failed", err)
		return false
	}
	p.CreatedAt, _ = time.Parse(db.ISO8601, createdAt)
	p.UpdatedAt, _ = time.Parse(db.ISO8601, updatedAt)
	return true
}

func (p *Plate) Create() bool {
	res, err := db.Exec("insert_plate", p.Code)
	if err != nil {
		return false
	}
	p.id, _ = res.LastInsertId()

	// Little cheat code
	p.CreatedAt = time.Now().UTC().Truncate(time.Second)
	p.UpdatedAt = p.CreatedAt
	return true
}

func (p Plate) FindUpdateTimestamp() bool {
	rows, err := db.Query("find_plate_update_timestamp", p.Code)
	if err != nil {
		return false
	}
	var timestamp string
	rows.Next()
	rows.Scan(&timestamp)
	parsed, err := time.Parse(db.ISO8601, timestamp)
	if err != nil {
		return false
	}
	p.UpdatedAt = parsed
	return true
}
