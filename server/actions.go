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

	if vars["code"] != plate.Code {
		newUrl := strings.ReplaceAll(r.URL.Path, vars["code"], plate.Code)
		http.Redirect(w, r, newUrl, http.StatusMovedPermanently)
		return nil, false
	}

	return plate, true
}

func (s *server) home(w http.ResponseWriter, r *http.Request) {
	html.Home(w)
}

func (s *server) plateList(w http.ResponseWriter, r *http.Request) {
	html.PlateList(w)
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
