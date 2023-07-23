package server

import (
	"github.com/gorilla/mux"
)

func (s *server) routes() {
	s.router = mux.NewRouter().StrictSlash(true)

	s.router.Use(s.logging)

	s.router.HandleFunc("/", s.home)
	s.router.HandleFunc("/plates", s.plateList)
	s.router.HandleFunc("/plates/{code}", s.plateShow)

	s.router.HandleFunc("/plates/{code}/interpretations", s.interpretationCreate)
	s.router.PathPrefix("/").Handler(StaticFiles)
}
