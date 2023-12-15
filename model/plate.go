package model

import (
	"fmt"
	"log"
	"strings"
	"time"

	"database/sql"

	"licenseplate.wtf/db"
)

type Plate struct {
	Id              int64
	Code            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Interpretations []Interpretation
}

func NewPlate(code string) *Plate {
	return &Plate{Code: strings.ReplaceAll(strings.ToUpper(code), "+", " ")}
}

func AllPlates() []Plate {
	rows, err := db.Query("all_plates")
	defer rows.Close()
	if err != nil {
		log.Println("Failed", err)
		return []Plate{}
	}

	plates := make([]Plate, 0)
	for rows.Next() {
		plate := Plate{}
		plate.FromRow(rows)
		plates = append(plates, plate)
	}
	log.Printf("Loaded %d plates", len(plates))
	return plates
}

func (p *Plate) FromRow(row *sql.Rows) {
	created_at, updated_at := "", ""

	if err := row.Scan(&p.Id, &p.Code, &created_at, &updated_at); err != nil {
		log.Println("Failed", err)
	}

	p.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", created_at)
	p.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updated_at)
}

func (p *Plate) URL() string {
	slug := strings.ReplaceAll(p.Code, " ", "+")
	return fmt.Sprintf("/plates/%s", slug)
}

func (p *Plate) Valid() bool {
	if len(p.Code) > 10 || len(p.Code) < 3 {
		return false
	}
	if !isAlphaNum(p.Code) {
		return false
	}
	if isNorty(p.Code) {
		return false
	}
	return true
}

func (p *Plate) FindOrCreate() bool {
	rows, err := db.Query("find_plate", p.Code)
	if err != nil || !rows.Next() {
		return p.Create()
	}
	defer rows.Close()
	var createdAt string
	var updatedAt string
	err = rows.Scan(&p.Id, &p.Code, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("Failed to load plate: %v", err)
		return false
	}
	p.CreatedAt, _ = time.Parse(db.ISO8601, createdAt)
	p.UpdatedAt, _ = time.Parse(db.ISO8601, updatedAt)
	return true
}

func (p *Plate) Create() bool {
	res, err := db.Exec("insert_plate", p.Code)
	if err != nil {
		log.Printf("Failed to create plate: %v", err)
		return false
	}
	p.Id, _ = res.LastInsertId()

	// Little cheat code
	p.CreatedAt = time.Now().UTC().Truncate(time.Second)
	p.UpdatedAt = p.CreatedAt
	return true
}

func (p *Plate) LoadInterpretations() {
	rows, err := db.Query("find_interpretations", p.Id)
	defer rows.Close()
	if err != nil {
		return
	}

	p.Interpretations = make([]Interpretation, 0)
	for rows.Next() {
		interp := Interpretation{}
		interp.FromRow(rows)
		p.Interpretations = append(p.Interpretations, interp)
	}
}
