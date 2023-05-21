package server

import (
	"io"
	"log"
	"net/http"

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

	if r.URL.Path != plate.URL() {
		http.Redirect(w, r, plate.URL(), http.StatusMovedPermanently)
		return nil, false
	}

	return plate, true
}

func (s *server) home(w http.ResponseWriter, r *http.Request) {
	html.Home(w)
}

func (s *server) plateList(w http.ResponseWriter, r *http.Request) {
	plates := model.AllPlates()
	log.Println("Found list of plates", plates)

	html.PlateList(w, html.ParamsMap{
		"Plates": plates,
	})
}

func (s *server) plateShow(w http.ResponseWriter, r *http.Request) {
	plate, found := currentPlate(w, r)

	if !found {
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
		return
	}

	err := r.ParseForm()
	if err != nil {
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

	http.Redirect(w, r, "/plates/"+plate.Code, http.StatusSeeOther)
}
