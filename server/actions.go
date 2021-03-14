package server

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"licenseplate.wtf/html"
	"licenseplate.wtf/model"
)

func (s *server) home(w http.ResponseWriter, r *http.Request) {
	html.Home(w)
}

func (s *server) plateList(w http.ResponseWriter, r *http.Request) {
	html.PlateList(w)
}

func (s *server) plateShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	plate := model.NewPlate(vars["code"])

	if !plate.Valid() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	success := plate.FindOrCreate()

	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if vars["code"] != plate.Code {
		http.Redirect(w, r, "/plates/"+plate.Code, http.StatusMovedPermanently)
	}

	s.cache(
		w,
		[]string{"v1/plate", plate.Code, plate.UpdatedAt.String()},
		func() io.Reader {
			var buffer bytes.Buffer
			html.PlateShow(&buffer, *plate)
			return &buffer
		},
	)
}
