package server

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"licenseplate.wtf/html"
	"licenseplate.wtf/model"
)

func currentPlate(w http.ResponseWriter, r *http.Request) (*model.Plate, bool) {
	vars := mux.Vars(r)
	plate := model.NewPlate(vars["code"])

	if !plate.Valid() {
		w.WriteHeader(http.StatusNotFound)
		return nil, false
	}

	success := plate.FindOrCreate()

	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, false
	}

	return plate, true
}

func (s *server) home(w http.ResponseWriter, r *http.Request) {
	html.Home(w)
}

func (s *server) plateList(w http.ResponseWriter, r *http.Request) {
	plates := model.AllPlates()
	html.PlateList(w, html.ParamsMap{
		"Plates": plates,
	})
}

func (s *server) plateShow(w http.ResponseWriter, r *http.Request) {
	plate, found := currentPlate(w, r)

	if !found {
		log.Printf("Couldn't find the plate: %v", plate)
	}

	if r.Method == "GET" && !strings.HasPrefix(plate.URL(), r.URL.Path) {
		log.Printf("Redirecting to canonical path: %s", plate.URL())
		http.Redirect(w, r, plate.URL(), http.StatusMovedPermanently)
		return
	}

	s.cache(
		w,
		[]string{"v1/plate", plate.Code, plate.UpdatedAt.String()},
		func(out io.Writer) {
			plate.LoadInterpretations()
			html.PlateShow(out, html.ParamsMap{
				"Plate": plate,
			})
		},
	)
}

func (s *server) interpretationCreate(w http.ResponseWriter, r *http.Request) {
	plate, found := currentPlate(w, r)

	if !found {
		log.Printf("Couldn't find the plate: %v", plate)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("Couldn't parse the form data %s", err)
		panic(err) // TODO
	}

	interpretation := model.Interpretation{PlateId: plate.Id}
	errors := interpretation.FromRequest(r)

	if len(errors) > 0 {
		log.Println("Rendering with errors", errors)
		plate.LoadInterpretations()
		html.PlateShow(w, html.ParamsMap{
			"Plate":  plate,
			"Errors": errors,
		})
		return
	}

	created := interpretation.Create()
	if !created {
		log.Println("Problem creating")
		errors = []string{"There was a problem saving your response. Please try again."}
		plate.LoadInterpretations()
		html.PlateShow(w, html.ParamsMap{
			"Plate":  plate,
			"Errors": errors,
		})
		return
	}

	// See other because we are redirecting to the plate rather than the interpretation
	http.Redirect(w, r, "/plates/"+plate.Code, http.StatusSeeOther)
}
