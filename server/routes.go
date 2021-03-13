package server

import (
	"github.com/gorilla/mux"
)

func (s *server) routes() {
	s.router = mux.NewRouter().StrictSlash(true)

	s.router.HandleFunc("/", home)
	s.router.HandleFunc("/plates", plateList)
	s.router.HandleFunc("/plates/{code}", plateShow)

	s.router.PathPrefix("/static/").Handler(StaticFiles)
}
